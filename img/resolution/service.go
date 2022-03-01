package resolution

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
	"strconv"
	"strings"
)

func init() {
	cmd.RegisterActor(&cmd.Actor{
		Id:    "img_resolution",
		Index: 3,
		Name:  "图片分辨率修改",
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
	cfg, err := loadConfigs(root)
	if err != nil {
		return nil, err
	}

	builder := &img.JobsBuilder{
		Root:    root,
		Workdir: "resolution",
		JobCreator: func(root string, path string) img.Job {
			return &jobImpl{
				root:    root,
				path:    path,
				configs: cfg,
			}
		},
	}
	return builder.Build()
}

func loadConfigs(root string) ([]*config, error) {
	cfgFile := filepath.Join(root, "resolution.txt")
	data, err := os.ReadFile(cfgFile)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("请检查 %s 是否存在", cfgFile))
	}
	str := strings.TrimSpace(string(data))
	str = strings.ReplaceAll(str, "\n\r", "\n")
	str = strings.ReplaceAll(str, "\r\n", "\n")
	str = strings.ReplaceAll(str, "\r", "\n")
	arr := strings.Split(str, "\n")
	ret := make([]*config, 0)
	for _, s := range arr {
		cfg := parse(s)
		if cfg != nil {
			ret = append(ret, cfg)
		}
	}
	if len(ret) == 0 {
		return nil, errors.New("请在resolution.txt中配置正确的分辨率输出")
	}
	return ret, nil
}

func parse(s string) *config {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return nil
	}
	if strings.HasPrefix(s, "#") {
		return nil
	}
	arr := strings.Split(s, "*")
	if len(arr) != 2 {
		return nil
	}
	width, err := strconv.Atoi(strings.TrimSpace(arr[0]))
	if err != nil {
		return nil
	}
	height, err := strconv.Atoi(strings.TrimSpace(arr[1]))
	if err != nil {
		return nil
	}
	return &config{
		width:  width,
		height: height,
	}
}

type config struct {
	width  int
	height int
}

type jobImpl struct {
	root    string
	path    string
	configs []*config
}

func (job jobImpl) Do() error {
	//make dir
	dest := job.getDest()
	dir := filepath.Dir(dest)
	_ = os.MkdirAll(dir, fs.ModePerm)
	//chg resolution
	imageType := img.GetImageType(dest)
	switch imageType {
	case img.IMAGE_TYPE_JPEG, img.IMAGE_TYPE_PNG:
		background, err := gg.LoadImage(job.path)
		if err != nil {
			return err
		}
		for _, cfg := range job.configs {
			ni, er := job.chg(background, cfg)
			if er != nil {
				return er
			}
			str := imgName(cfg, dest)
			er = job.save(ni, imageType, str)
			if er != nil {
				return er
			}
		}
		return nil
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

func imgName(cfg *config, dest string) string {
	arr := strings.Split(dest, ".")
	return fmt.Sprintf("%s-%dx%d.%s", arr[0], cfg.width, cfg.height, arr[1])
}

func (job jobImpl) Info() string {
	return job.getDest()
}

func (job *jobImpl) getDest() string {
	return strings.Replace(job.path, filepath.Join(job.root, "resolution"), filepath.Join(job.root, "resolution_new"), 1)
}

func (job *jobImpl) chg(background image.Image, cfg *config) (image.Image, error) {
	b := background.Bounds()
	bw := b.Dx()
	bh := b.Dy()
	if float64(cfg.width) > (float64(bw)*4) || float64(cfg.height) > (float64(bh)*4) {
		return nil, errors.New(fmt.Sprintf("图片大小%d*%d，尺寸远小于%d*%d,跳过", bw, bh, cfg.width, cfg.height))
	}
	rw := cfg.width
	rh := cfg.height
	scale := 1.0
	if cfg.width > 0 && cfg.height <= 0 {
		scale = float64(cfg.width) / float64(bw)
		rw = cfg.width
		rh = int(float64(bh) * scale)
	} else if cfg.width <= 0 && cfg.height > 0 {
		scale = float64(cfg.height) / float64(bh)
		rw = int(float64(cfg.width) * scale)
		rh = cfg.height
	} else if cfg.width > 0 && cfg.height > 0 {
		scaleX := float64(cfg.width) / float64(bw)
		scaleY := float64(cfg.height) / float64(bh)
		scale = math.Max(scaleX, scaleY)
		rw = cfg.width
		rh = cfg.height
	} else {
		rw = bw
		rh = bh
		scale = 1.0
	}
	dc := gg.NewContext(rw, rh)
	dc.Scale(scale, scale)
	dc.DrawImage(background, 0, 0)
	return dc.Image(), nil
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
