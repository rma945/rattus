package main

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
)

// get secrets from azure vault
func getsecretAzures(secretAzure string) ([]string, error) {
	var secrets []string

	azureCredentials, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return secrets, err
	}

	for _, s := range stringToList(secretAzure) {
		parsedSecretURL, err := url.Parse(s)
		if err != nil {
			return secrets, err
		}

		azureVaultClient, err := azsecrets.NewClient(fmt.Sprintf("https://%s", parsedSecretURL.Host), azureCredentials, nil)
		if err != nil {
			return secrets, err
		}

		resp, err := azureVaultClient.GetSecret(context.TODO(), strings.TrimPrefix(parsedSecretURL.Path, "/secrets/"), nil)
		if err != nil {
			return secrets, err
		}

		secrets = append(secrets, *resp.Value)
	}

	return secrets, nil
}
