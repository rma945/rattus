package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// location of k8s service account token
const tokenFilePath = "/var/run/secrets/kubernetes.io/serviceaccount/token"

// location of k8s namespace
const namespaceFilePath = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

// get token service account token
func getK8SServiceAccountToken() (string, error) {
	var token string

	fileContent, err := ioutil.ReadFile(tokenFilePath)
	if err != nil {
		log.Fatal(err)
	}

	token = string([]byte(fileContent))
	return token, nil
}

func getK8SServiceRole() (string, error) {
	var token string

	fileContent, err := ioutil.ReadFile(namespaceFilePath)
	if err != nil {
		log.Fatal(err)
	}

	token = string([]byte(fileContent))
	return token, nil
}

// getVaultAuthToken - login at vault and retrive vault auth token
func getVaultAuthToken(vaultSecretURL, authToken, authRole string) (string, error) {
	var token string
	var parsedResponse map[string]interface{}
	var requstPayload = []byte(fmt.Sprintf(`{"jwt": "%s", "role": "%s"}`, authToken, authRole))

	vaultLoginURL, err := getVaultLoginURL(vaultSecretURL)
	if err != nil {
		return token, err
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: vaultSkipTLS}
	client := &http.Client{
		Timeout: time.Second * requestTimeout,
	}

	request, err := http.NewRequest("POST", vaultLoginURL, bytes.NewBuffer(requstPayload))
	if err != nil {
		return token, err
	}

	response, err := client.Do(request)
	if err != nil {
		return token, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return token, fmt.Errorf("vault response code: %d", response.StatusCode)
	}

	respBodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return token, err
	}

	if err := json.Unmarshal(respBodyBytes, &parsedResponse); err != nil {
		return token, err
	}

	token = parsedResponse["auth"].(map[string]interface{})["client_token"].(string)

	return token, nil
}

func getK8SVaultToken(vaultSecretURL string) (string, error) {
	var token string
	serviceToken, err := getK8SServiceAccountToken()
	if err != nil {
		return token, fmt.Errorf("failed to get k8s service account token - %s", err)
	}

	serviceRole, err := getK8SServiceRole()
	if err != nil {
		return token, fmt.Errorf("failed to get k8s namespace - %s", err)
	}

	token, err = getVaultAuthToken(vaultSecretURL, serviceToken, serviceRole)
	if err != nil {
		return token, fmt.Errorf("failed to auth at vault - %s", err)
	}

	return token, nil
}
