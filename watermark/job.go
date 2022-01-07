package watermark

import (
	"errors"
	"fmt"
	"github.com/fogleman/gg"
	"image"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func CreateJobs(root string) ([]*Job, error) {
	markPath := filepath.Join(root, "mark.png")
	watermark, err := gg.LoadImage(markPath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("水印图片加载失败，请检查 %s 是否存在", markPath))
	}
	ret := make([]*Job, 0)
	_ = filepath.WalkDir(filepath.Join(root, "mark"), func(path string, d fs.DirEntry, err error) error {
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
	dc.Push()

	bgw := background.Bounds().Dx()
	bgh := background.Bounds().Dy()
	mw := job.mark.Bounds().Dx()
	mh := job.mark.Bounds().Dy()

	if mw > bgw {
		//水印比背景宽
		scale := 0.5 * float64(bgw) / float64(mw)
		x := float64(bgw / 4)
		y := float64(bgh / 2)
		dc.ScaleAbout(scale, scale, x, y)
		dc.DrawImage(job.mark, 0, 0)
	} else if mh > bgh {
		//水印比背景高
		scale := 0.5 * float64(bgh) / float64(mh)
		x := float64(bgw / 4)
		y := float64(bgh / 2)
		dc.ScaleAbout(scale, scale, x, y)
		dc.DrawImage(job.mark, 0, 0)
	} else {
		x := (bgw - mw) / 2
		y := (bgh - mh) / 2
		dc.DrawImage(job.mark, x, y)
	}
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
