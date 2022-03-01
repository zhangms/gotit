package watermark

import (
	"errors"
	"fmt"
	"github.com/fogleman/gg"
	"gotit/cmd"
	"gotit/img"
	"gotit/parallel"
	"image"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	cmd.RegisterActor(&cmd.Actor{
		Id:    "watermarker",
		Index: 0,
		Name:  "图片加水印",
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
	markPNG := filepath.Join(root, "mark.png")
	watermark, err := gg.LoadImage(markPNG)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("水印图片加载失败，请检查 %s 是否存在", markPNG))
	}
	builder := &img.JobsBuilder{
		Root:    root,
		Workdir: "mark",
		JobCreator: func(root string, path string) img.Job {
			return newJob(watermark, root, path)
		},
	}
	return builder.Build()
}

func newJob(watermark image.Image, root string, path string) *jobImpl {
	return &jobImpl{
		root: root,
		path: path,
		mark: watermark,
	}
}

type jobImpl struct {
	root string
	path string
	mark image.Image
}

func (job *jobImpl) Info() string {
	return job.getDest()
}

func (job *jobImpl) Do() error {
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

func (job *jobImpl) drawMark(background image.Image) image.Image {
	dc := gg.NewContextForImage(background)

	bgw := background.Bounds().Dx()
	bgh := background.Bounds().Dy()
	mw := job.mark.Bounds().Dx()
	mh := job.mark.Bounds().Dy()

	rate := 0.4
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

func (job *jobImpl) save(image image.Image, imageType string, dest string) error {
	switch imageType {
	case img.IMAGE_TYPE_JPEG:
		return gg.SaveJPG(dest, image, 100)
	case img.IMAGE_TYPE_PNG:
		return gg.SavePNG(dest, image)
	default:
		return nil
	}
}

func (job *jobImpl) getDest() string {
	return strings.Replace(job.path, filepath.Join(job.root, "mark"), filepath.Join(job.root, "mark_new"), 1)
}
