package compress

import (
	"gotit/parallel"
	"testing"
)

func TestCompress(t *testing.T) {
	jobs, _ := CreateJobs("/Users/zms/Downloads/workspace")
	parallel.Do(jobs, 1)
}
