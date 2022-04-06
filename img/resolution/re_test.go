package resolution

import (
	"fmt"
	"gotit/parallel"
	"regexp"
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

func TestReg(t *testing.T) {
	s := "百度    123*456"
	reg, _ := regexp.Compile("\\s+")
	s = reg.ReplaceAllString(s, " ")
	fmt.Println(s)
}
