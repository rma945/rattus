package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"text/template"
	"time"
)

// template timestamp function
func templateDatetime() string {
	timeNow := time.Now()
	return timeNow.Format("2006-01-02-15:04:05")
}

// template base64 decode function
func templateBase64Decode(encodedString string) string {
	decodedString, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		return "<failed to decode base64>"
	}
	return string(decodedString)
}

// register custom go template functions
func registerTemplateFunctions() template.FuncMap {
	return template.FuncMap{
		"datetime":     templateDatetime,
		"base64decode": templateBase64Decode,
	}
}

// render template
func generateTemplate(templatePath string, values map[string]interface{}) (string, error) {
	templateFileContent, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return "", err
	}

	templateRender, err := template.New("template").Funcs(registerTemplateFunctions()).Parse(string(templateFileContent))
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

// render template or return plain secrets text
func renderOutput(secrets []string, templatePath string) (string, error) {
	var mergedSecrets map[string]interface{}
	var output string
	var err error

	// merge list of secrets into single one or return them as plain text
	if mergedSecrets, err = mergeSecretListToMap(secrets); err != nil {
		return strings.Join(secrets, "\n"), nil
	}

	// render secrets as template
	if templatePath != "" {
		output, err = generateTemplate(templatePath, mergedSecrets)
		if err != nil {
			return "", fmt.Errorf("failed to render secrets template - %s", err)
		}
	} else {
		secretsJSON, err := json.Marshal(mergedSecrets)
		output = string(secretsJSON)
		if err != nil {
			return "", fmt.Errorf("failed to render map of secrets as json - %s", err)
		}
	}

	return output, nil
}
