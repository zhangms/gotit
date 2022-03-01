package compress

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"gotit/cmd"
	"gotit/img"
	"gotit/parallel"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	cmd.RegisterActor(&cmd.Actor{
		Id:    "img_compressor",
		Index: img.Compress,
		Name:  "图片压缩",
		Action: func(args []string) error {
			jobs, err := newJobs(filepath.Dir(args[0]))
			if err != nil {
				return err
			}
			parallel.Do(jobs, 16)
			return nil
		},
	})
}

func newJobs(root string) (*img.Jobs, error) {
	apikey, err := readApikey(root)
	if err != nil {
		return nil, err
	}
	builder := &img.JobsBuilder{
		Root:    root,
		Workdir: "compress",
		JobCreator: func(root string, path string) img.Job {
			return newJob(apikey, root, path)
		},
		Summary: func(jobs *img.Jobs) string {
			var aIn, aOut float64
			for _, job := range jobs.Children {
				j := job.(*jobImpl)
				aIn += float64(j.input)
				aOut += float64(j.output)
			}
			imb := aIn / 1024 / 1024
			omb := aOut / 1024 / 1024
			rate := (aIn - aOut) / aIn * 100
			return fmt.Sprintf("原始图片共%.2fMB, 压缩后为%.2fMB，压缩率 %.2f%%", imb, omb, rate)
		},
	}
	return builder.Build()
}

func readApikey(root string) ([]string, error) {
	apikeyFile := filepath.Join(root, "apikey.txt")
	data, err := os.ReadFile(apikeyFile)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("请检查 %s 是否存在", apikeyFile))
	}
	keys := strings.TrimSpace(string(data))
	keys = strings.ReplaceAll(keys, "\n\r", "\n")
	keys = strings.ReplaceAll(keys, "\r\n", "\n")
	keys = strings.ReplaceAll(keys, "\r", "\n")
	apiKeys := strings.Split(keys, "\n")
	return apiKeys, nil
}

func newJob(apikey []string, root string, path string) *jobImpl {
	return &jobImpl{
		root:   root,
		path:   path,
		apikey: apikey,
	}
}

type jobImpl struct {
	root   string
	path   string
	apikey []string
	input  int64
	output int64
}

func (job *jobImpl) Info() string {
	return job.getDest()
}

func (job *jobImpl) Do() error {
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
			return er
		}
		_ = os.WriteFile(dest, dt, fs.ModePerm)
		fmt.Printf("[WARN ] 非图片直接拷贝：%s\n", job.path)
		return nil
	}
}

func (job *jobImpl) compress() error {
	imageData, err := os.ReadFile(job.path)
	if err != nil {
		return err
	}
	inf, err := os.Stat(job.getDest())
	if !os.IsNotExist(err) {
		job.output = inf.Size()
		job.input = int64(len(imageData))
		fmt.Printf("[WARN ] 已存在跳过：%s\n", job.getDest())
		return nil
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

func (job *jobImpl) getDest() string {
	return strings.Replace(job.path, filepath.Join(job.root, "compress"), filepath.Join(job.root, "compress_new"), 1)
}

type compressResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Input   struct {
		Size int64  `json:"size"`
		Type string `json:"type"`
	}
	Output struct {
		Size   int64   `json:"size"`
		Type   string  `json:"type"`
		Width  int     `json:"width"`
		Height int     `json:"height"`
		Ratio  float32 `json:"ratio"`
		Url    string  `json:"url"`
	}
}

func sendCompressRequest(apiKeys []string, data []byte) (*compressResponse, error) {
	var globalErr error
	for _, apikey := range apiKeys {
		apikey = strings.TrimSpace(apikey)
		if len(apikey) != 32 {
			continue
		}
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
		responseData, err := ioutil.ReadAll(response.Body)
		_ = response.Body.Close()
		if err != nil {
			return nil, err
		}
		ret := compressResponse{}
		err = json.Unmarshal(responseData, &ret)
		if err != nil {
			globalErr = err
			continue
		}
		if len(ret.Error) > 0 {
			globalErr = errors.New(string(responseData))
			continue
		}
		return &ret, nil
	}
	return nil, globalErr
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
