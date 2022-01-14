package res

import (
	"bytes"
	"embed"
	"errors"
	"text/template"
)

//go:embed static
var static embed.FS

func ListDir(dirName string) []string {
	ret := make([]string, 0)
	files, err := static.ReadDir("static/" + dirName)
	if err != nil {
		return ret
	}
	for _, file := range files {
		ret = append(ret, dirName+"/"+file.Name())
	}
	return ret
}

func ReadData(fileName string) ([]byte, error) {
	data, err := static.ReadFile("static/" + fileName)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func TemplateText(name string, params interface{}) ([]byte, error) {
	tplContent, err := ReadData("tpl/" + name)
	if err != nil {
		return nil, errors.New("READ_TEMPLATE_ERROR|" + name + "|" + err.Error())
	}
	if params == nil {
		return tplContent, nil
	}
	t, err := template.New(name).Parse(string(tplContent))
	if err != nil {
		return nil, errors.New("PARSE_TEMPLATE_ERROR|" + name + "|" + err.Error())
	}
	var buffer bytes.Buffer
	err = t.Execute(&buffer, params)
	if err != nil {
		return nil, errors.New("EXECUTE_TEMPLATE_ERROR|" + name + "|" + err.Error())
	}
	return buffer.Bytes(), nil
}
