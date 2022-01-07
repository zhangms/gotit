package watermark

import (
	"errors"
	"fmt"
	"github.com/fogleman/gg"
	"image"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"strings"
)

func CreateJobs(root string) ([]*Job, error) {
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
		if !d.IsDir() {
			job := CreateJob(watermark, root, path)
			ret = append(ret, job)
		}
		return nil
	})
	return ret, nil
}

func CreateJob(watermark image.Image, root string, path string) *Job {
	return &Job{
		root: root,
		path: path,
		mark: watermark,
	}
}

const (
	imageTypeJPEG = "JPEG"
	imageTypePNG  = "PNG"
)

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
	fileType := job.getFileType()
	switch fileType {
	case imageTypePNG, imageTypeJPEG:
		background, err := gg.LoadImage(job.path)
		if err != nil {
			return err
		}
		marked := job.drawMark(background)
		return job.save(marked, fileType, dest)
	default:
		dt, er := os.ReadFile(job.path)
		if er != nil {
			_ = os.WriteFile(dest, dt, fs.ModePerm)
		}
		return errors.New("未知类型，直接拷贝：" + job.path)
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

func (job *Job) save(img image.Image, fileType string, dest string) error {
	switch fileType {
	case imageTypeJPEG:
		return gg.SaveJPG(dest, img, 100)
	case imageTypePNG:
		return gg.SavePNG(dest, img)
	default:
		return nil
	}
}

func (job *Job) getDest() string {
	return strings.Replace(job.path, filepath.Join(job.root, "mark"), filepath.Join(job.root, "mark_new"), 1)
}

func (job *Job) getFileType() string {
	dest := job.getDest()
	arr := strings.Split(dest, ".")
	suffix := arr[len(arr)-1]
	switch suffix {
	case "jpg", "JPG", "jpeg", "JPEG":
		return imageTypeJPEG
	case "png", "PNG":
		return imageTypePNG
	default:
		return suffix
	}
}
