package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

const vaultSkipTLS = true
const requestTimeout = 10

const vaultAuthEndpoint = "/v1/auth/kubernetes/login"

// generate url to k8s login provider from VaultSecretURL
func getVaultLoginURL(URL string) (string, error) {
	var vaultAuthURL string
	parsedURL, err := url.Parse(URL)
	if err != nil {
		log.Fatal(err)
	}
	vaultAuthURL = fmt.Sprintf("%s://%s%s", parsedURL.Scheme, parsedURL.Host, vaultAuthEndpoint)

	return vaultAuthURL, nil
}

// get secret from vault
func getVaultSecret(URL, authToken string) (string, error) {
	var secrets string
	var parsedResponse map[string]interface{}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: vaultSkipTLS}
	client := &http.Client{
		Timeout: time.Second * requestTimeout,
	}

	request, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return secrets, err
	}
	request.Header.Add("X-Vault-Token", authToken)

	response, err := client.Do(request)
	if err != nil {
		return secrets, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return secrets, fmt.Errorf("vault response code: %d", response.StatusCode)
	}

	respBodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return secrets, err
	}

	if err := json.Unmarshal(respBodyBytes, &parsedResponse); err != nil {
		return secrets, err
	}

	if parsedResponse["data"].(map[string]interface{})["metadata"] != nil {
		// vault v2 secret
		secrets, err = mapToJSON(parsedResponse["data"].(map[string]interface{})["data"])
	} else {
		// vault v1 secret
		secrets, err = mapToJSON(parsedResponse["data"].(map[string]interface{}))
	}

	return secrets, nil
}

// get secret from vault
func vaultGetSecret(config applicationConfig) (string, error) {
	vaultToken := config.VaultToken
	var secrets string
	var err error

	// issue new vault token if it was not set from config
	if config.VaultToken == "" {
		vaultToken, err = getK8SVaultToken(config.VaultSecretURL)
		if err != nil {
			return secrets, err
		}
	}

	secrets, err = getVaultSecret(config.VaultSecretURL, vaultToken)
	if err != nil {
		return secrets, fmt.Errorf("failed to retrive vault secrets - %s", err)
	}

	return secrets, nil
}
