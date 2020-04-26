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

// register custom go template functions
func registerTemplateFunctions() template.FuncMap {
	return template.FuncMap{
		"datetime": templateDatetime,
	}
}

// render template
func generateTemplate(templatePath string, values map[string]interface{}) (string, error) {
	var renderedTemplate string

	templateFileContent, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return renderedTemplate, err
	}

	templateFunctions := registerTemplateFunctions()
	templateRender, err := template.New("template").Funcs(templateFunctions).Parse(string(templateFileContent))
	if err != nil {
		return renderedTemplate, err
	}

	templateRederBuffer := &bytes.Buffer{}
	err = templateRender.Execute(templateRederBuffer, values)
	if err != nil {
		return renderedTemplate, err
	}
	renderedTemplate = templateRederBuffer.String()

	return renderedTemplate, nil
}

// convert map of interface to JSON
func mapToJSON(values interface{}) (string, error) {
	JSON, err := json.Marshal(values)
	if err != nil {
		return "", err
	}

	return string(JSON), nil
}

// render template or return plain secrets text
func renderOutput(secretsString, templatePath string) (string, error) {
	var secretsMap map[string]interface{}
	var secretsOutput string
	var err error

	secretsOutput = secretsString

	// try to convert json secrets to map interfaces or return plan secret value
	err = json.Unmarshal([]byte(secretsString), &secretsMap)
	if err != nil {
		return secretsOutput, nil
	}

	// render secrets as template
	if templatePath != "" {
		secretsOutput, err = generateTemplate(templatePath, secretsMap)
		if err != nil {
			return secretsOutput, fmt.Errorf("failed to render secrets template - %s", err)
		}
	}

	// return plain json
	return secretsOutput, nil
}
