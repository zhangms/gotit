package compress

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"gotit/img"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Jobs struct {
	children []*Job
}

func (j *Jobs) Total() int {
	return len(j.children)
}

func (j *Jobs) Do(index int) error {
	return j.children[index].Do()
}

func (j *Jobs) Info(index int) string {
	return j.children[index].getDest()
}

func (j *Jobs) Summary() string {
	var aIn, aOut float64
	for _, j := range j.children {
		aIn += float64(j.input)
		aOut += float64(j.output)
	}
	imb := aIn / 1024 / 1024
	omb := aOut / 1024 / 1024
	rate := (aIn - aOut) / aIn * 100
	return fmt.Sprintf("原始图片共%.2fMB, 压缩后为%.2fMB，压缩率 %.2f%%", imb, omb, rate)
}

func CreateJobs(root string) (*Jobs, error) {
	apikeyFile := filepath.Join(root, "apikey.txt")
	data, err := os.ReadFile(apikeyFile)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("请检查 %s 是否存在", apikeyFile))
	}
	apikey := strings.TrimSpace(string(data))
	compressDir := filepath.Join(root, "compress")
	dir, err := os.Stat(compressDir)
	if os.IsNotExist(err) || !dir.IsDir() {
		return nil, errors.New(fmt.Sprintf("工作目录 %s 不存在", compressDir))
	}
	ret := make([]*Job, 0)
	_ = filepath.WalkDir(compressDir, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && !strings.HasPrefix(d.Name(), ".") {
			job := CreateJob(apikey, root, path)
			ret = append(ret, job)
		}
		return nil
	})
	return &Jobs{
		children: ret,
	}, nil
}

func CreateJob(apikey string, root string, path string) *Job {
	return &Job{
		root:   root,
		path:   path,
		apikey: apikey,
	}
}

type Job struct {
	root   string
	path   string
	apikey string
	input  int
	output int
}

func (job *Job) Do() error {
	//make dir
	dest := job.getDest()
	dir := filepath.Dir(dest)
	_ = os.MkdirAll(dir, fs.ModePerm)
	//compress
	imageType := img.GetImageType(dest)
	switch imageType {
	case img.IMAGE_TYPE_PNG, img.IMAGE_TYPE_JPEG:
		return job.compress()
	default:
		dt, er := os.ReadFile(job.path)
		if er != nil {
			_ = os.WriteFile(dest, dt, fs.ModePerm)
		}
		fmt.Printf("[WARN ] 非图片直接拷贝：%s\n", job.path)
		return nil
	}
}

func (job *Job) compress() error {
	_, err := os.Stat(job.getDest())
	if !os.IsNotExist(err) {
		fmt.Printf("[WARN ] 已存在跳过：%s\n", job.getDest())
		return nil
	}
	imageData, err := os.ReadFile(job.path)
	if err != nil {
		return err
	}
	resp, err := sendCompressRequest(job.apikey, imageData)
	if err != nil {
		return err
	}
	job.input = resp.Input.Size
	job.output = resp.Output.Size
	if job.output < job.input {
		return downloadAndSave(resp.Output.Url, job.getDest())
	}
	return os.WriteFile(job.getDest(), imageData, fs.ModePerm)
}

func (job *Job) getDest() string {
	return strings.Replace(job.path, filepath.Join(job.root, "compress"), filepath.Join(job.root, "compress_new"), 1)
}

type compressResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Input   struct {
		Size int    `json:"size"`
		Type string `json:"type"`
	}
	Output struct {
		Size   int     `json:"size"`
		Type   string  `json:"type"`
		Width  int     `json:"width"`
		Height int     `json:"height"`
		Ratio  float32 `json:"ratio"`
		Url    string  `json:"url"`
	}
}

func sendCompressRequest(apikey string, data []byte) (*compressResponse, error) {
	request, err := http.NewRequest("POST", "https://api.tinify.com/shrink", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	auth := base64.StdEncoding.EncodeToString([]byte("api:" + apikey))
	request.Header.Add("Authorization", fmt.Sprintf("Basic %s", auth))
	cli := &http.Client{}
	response, err := cli.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	ret := compressResponse{}
	err = json.Unmarshal(responseData, &ret)
	if err != nil {
		return nil, err
	}
	if len(ret.Error) > 0 {
		return nil, errors.New(string(responseData))
	}
	return &ret, nil
}

func downloadAndSave(url string, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return os.WriteFile(dest, data, fs.ModePerm)
}
