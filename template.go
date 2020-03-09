package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"text/template"
	"time"
)

func templateDatetime() string {
	timeNow := time.Now()
	return timeNow.Format("2006-01-02-15:04:05")
}

func registerTemplateFunctions() template.FuncMap {
	return template.FuncMap{
		"datetime": templateDatetime,
	}
}

func generateTemplate(templatePath string, values map[string]interface{}) (string, error) {

	templateFileContent, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return "", err
	}

	templateFunctions := registerTemplateFunctions()
	templateRender, err := template.New("template").Funcs(templateFunctions).Parse(string(templateFileContent))
	if err != nil {
		return "", err
	}

	templateRederBuffer := &bytes.Buffer{}
	err = templateRender.Execute(templateRederBuffer, values)
	if err != nil {
		return "", err
	}

	return templateRederBuffer.String(), nil
}

// convert map of interface to JSON
func mapToJSON(values map[string]interface{}) (string, error) {
	JSON, err := json.Marshal(values)
	if err != nil {
		return "", err
	}

	return string(JSON), nil
}

func renderOutput(secrets map[string]interface{}, templatePath string) (string, error) {
	var stdout string
	var err error

	if templatePath != "" {
		stdout, err = generateTemplate(templatePath, secrets)
		if err != nil {
			return stdout, fmt.Errorf("failed to render secrets template - %s", err)
		}
	} else {
		stdout, err = mapToJSON(secrets)
		if err != nil {
			return stdout, fmt.Errorf("failed to convert secrets at json - %s", err)
		}
	}

	return stdout, nil
}
