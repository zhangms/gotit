package resolution

import (
	"fmt"
	"gotit/parallel"
	"testing"
)

func TestResolution(t *testing.T) {
	jobs, err := newJobs("/Users/zms/Downloads/workspace")
	if err != nil {
		fmt.Println(err)
		return
	}
	parallel.Do(jobs, 1)
}
