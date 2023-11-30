package dirgroup

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Exec(dir string) error {
	dirName := filepath.Base(dir)
	fmt.Println(dirName)
	psDir := filepath.Join(dir, "PS")
	list, err := os.ReadDir(psDir)
	if err != nil {
		return err
	}
	names := make([]string, 0)
	for _, f := range list {
		name := strings.Split(f.Name(), ".")[0]
		names = append(names, name)
	}

	for _, name := range names {
		f1 := filepath.Join(dir, "PS", name+".psd")
		f11 := filepath.Join(dir, "DEST", dirName, name, "文件提交", name+".psd")
		fileCopy(f1, f11)

		f2 := filepath.Join(dir, "100x100", name+".png")
		f21 := filepath.Join(dir, "DEST", dirName, name, "输出文件2", name+".png")
		fileCopy(f2, f21)

		f3 := filepath.Join(dir, "1000x1000", name+".png")
		f31 := filepath.Join(dir, "DEST", dirName, name, "输出文件1", name+".png")
		fileCopy(f3, f31)
	}

	return nil
}

func fileCopy(src string, dest string) {
	_ = os.MkdirAll(filepath.Dir(dest), os.ModePerm)
	fmt.Println(src)
	srcFile, err := os.Open(src)
	if err != nil {
		panic(err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		panic(err)
	}
	defer destFile.Close()
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		panic(err)
	}
}
