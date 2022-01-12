package compress

import (
	"gotit/parallel"
	"path/filepath"
)

func Do(args []string) error {
	return doCompress(filepath.Dir(args[0]), 16)
}

func doCompress(workspace string, routine int) error {
	jobs, err := CreateJobs(workspace)
	if err != nil {
		return err
	}
	parallel.Do(jobs, routine)
	return nil
}
