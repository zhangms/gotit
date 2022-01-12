package watermark

import (
	"gotit/parallel"
	"path/filepath"
)

func Do(args []string) error {
	err := doMark(filepath.Dir(args[0]), 16)
	return err
}

func doMark(workspace string, routine int) error {
	jobs, err := CreateJobs(workspace)
	if err != nil {
		return err
	}
	parallel.Do(jobs, routine)
	return nil
}
