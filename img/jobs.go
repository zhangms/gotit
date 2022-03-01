package img

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Job interface {
	Do() error
	Info() string
}

type Jobs struct {
	Children []Job
	summary  func(jobs *Jobs) string
}

func (j *Jobs) Total() int {
	return len(j.Children)
}

func (j *Jobs) Do(index int) error {
	return j.Children[index].Do()
}

func (j *Jobs) Info(index int) string {
	return j.Children[index].Info()
}

func (j *Jobs) Summary() string {
	if j.summary == nil {
		return "OK"
	}
	return j.summary(j)
}

type JobsBuilder struct {
	Root       string
	Workdir    string
	JobCreator func(root string, path string) Job
	Summary    func(jobs *Jobs) string
}

func (j *JobsBuilder) Build() (*Jobs, error) {
	workspace := filepath.Join(j.Root, j.Workdir)
	dir, err := os.Stat(workspace)
	if os.IsNotExist(err) || !dir.IsDir() {
		return nil, errors.New(fmt.Sprintf("工作目录 %s 不存在", workspace))
	}
	ret := make([]Job, 0)
	_ = filepath.WalkDir(workspace, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && !strings.HasPrefix(d.Name(), ".") {
			job := j.JobCreator(j.Root, path)
			ret = append(ret, job)
		}
		return nil
	})
	return &Jobs{
		Children: ret,
		summary:  j.Summary,
	}, nil
}
