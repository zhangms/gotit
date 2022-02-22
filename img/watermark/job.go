package watermark

import (
	"errors"
	"fmt"
	"github.com/fogleman/gg"
	"gotit/img"
	"image"
	"io/fs"
	"math"
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
	return "OK"
}

func CreateJobs(root string) (*Jobs, error) {
	markPNG := filepath.Join(root, "mark.png")
	watermark, err := gg.LoadImage(markPNG)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("水印图片加载失败，请检查 %s 是否存在", markPNG))
	}
	markDir := filepath.Join(root, "mark")
	dir, err := os.Stat(markDir)
	if os.IsNotExist(err) || !dir.IsDir() {
		return nil, errors.New(fmt.Sprintf("工作目录 %s 不存在", markDir))
	}
	ret := make([]*Job, 0)
	_ = filepath.WalkDir(markDir, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && !strings.HasPrefix(d.Name(), ".") {
			job := CreateJob(watermark, root, path)
			ret = append(ret, job)
		}
		return nil
	})
	return &Jobs{
		children: ret,
	}, nil
}

func CreateJob(watermark image.Image, root string, path string) *Job {
	return &Job{
		root: root,
		path: path,
		mark: watermark,
	}
}

type Job struct {
	root string
	path string
	mark image.Image
}

func (job *Job) Do() error {
	//make dir
	dest := job.getDest()
	dir := filepath.Dir(dest)
	_ = os.MkdirAll(dir, fs.ModePerm)
	//mark
	imageType := img.GetImageType(dest)
	switch imageType {
	case img.IMAGE_TYPE_JPEG, img.IMAGE_TYPE_PNG:
		background, err := gg.LoadImage(job.path)
		if err != nil {
			return err
		}
		marked := job.drawMark(background)
		return job.save(marked, imageType, dest)
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

func (job *Job) drawMark(background image.Image) image.Image {
	dc := gg.NewContextForImage(background)

	bgw := background.Bounds().Dx()
	bgh := background.Bounds().Dy()
	mw := job.mark.Bounds().Dx()
	mh := job.mark.Bounds().Dy()

	rate := 0.75
	scaleX := float64(bgw) * rate / float64(mw)
	scaleY := float64(bgh) * rate / float64(mh)
	scale := math.Min(scaleX, scaleY)
	mwr := float64(mw) * scale
	mhr := float64(mh) * scale

	dc.Push()
	dc.ScaleAbout(scale, scale, (float64(bgw)-mwr)/2, (float64(bgh)-mhr)/2)
	dc.DrawImage(job.mark, (bgw-int(mwr))/2, (bgh-int(mhr))/2)
	dc.Pop()
	return dc.Image()
}

func (job *Job) save(image image.Image, imageType string, dest string) error {
	switch imageType {
	case img.IMAGE_TYPE_JPEG:
		return gg.SaveJPG(dest, image, 100)
	case img.IMAGE_TYPE_PNG:
		return gg.SavePNG(dest, image)
	default:
		return nil
	}
}

func (job *Job) getDest() string {
	return strings.Replace(job.path, filepath.Join(job.root, "mark"), filepath.Join(job.root, "mark_new"), 1)
}
