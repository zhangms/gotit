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

type Job struct {
	root string
	path string
	mark image.Image
}

func (job *Job) Do() error {
	background, err := gg.LoadImage(job.path)
	dest := job.getDest()
	if err != nil {
		dt, er := os.ReadFile(job.path)
		if er != nil {
			dir := filepath.Dir(dest)
			_ = os.MkdirAll(dir, fs.ModePerm)
			_ = os.WriteFile(dest, dt, fs.ModePerm)
		}
		return errors.New("加载失败，拷贝：" + job.path)
	}
	markImage := job.drawMark(background)
	return job.saveImage(markImage, dest)
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

func (job *Job) saveImage(img image.Image, dest string) error {
	dir := filepath.Dir(dest)
	_ = os.MkdirAll(dir, fs.ModePerm)
	arr := strings.Split(dest, ".")
	suffix := arr[len(arr)-1]
	var err error = nil
	switch suffix {
	case "jpg", "JPG", "jpeg", "JPEG":
		err = gg.SaveJPG(dest, img, 100)
		break
	case "png", "PNG":
		err = gg.SavePNG(dest, img)
		break
	default:
		err = errors.New("不支持的图片类型")
		break
	}
	if err != nil {
		return errors.New("保存失败：" + dest + "," + err.Error())
	}
	return nil
}

func (job *Job) getDest() string {
	return strings.Replace(job.path, filepath.Join(job.root, "mark"), filepath.Join(job.root, "mark_new"), 1)
}
