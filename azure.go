package main

import (
	"context"
	"encoding/json"
	"path"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

const azureKeyVaultResourceURL = "https://vault.azure.net"

func getAzureSecrets(azureVault string) (string, error) {
	var azureSecrets string = ""

	azureBasicClient := keyvault.New()
	azureAuthorizer, err := auth.NewAuthorizerFromEnvironmentWithResource(azureKeyVaultResourceURL)
	if err == nil {
		azureBasicClient.Authorizer = azureAuthorizer
	}

	vaultSecretsList, err := azureBasicClient.GetSecrets(context.Background(), azureVault, nil)
	if err != nil {
		return azureSecrets, err
	}

	// retrive all secrets from keyVault
	secretsList := make(map[string]string)
	for {
		// retrive values from keyVault secrets
		for _, secret := range vaultSecretsList.Values() {
			secretValue, err := azureBasicClient.GetSecret(context.Background(), azureVault, path.Base(*secret.ID), "")
			if err == nil {
				secretsList[path.Base(*secret.ID)] = *secretValue.Value
			}
		}
		// retrive next page of secrets from keyVault
		if err := vaultSecretsList.NextWithContext(context.Background()); err != nil {
			return azureSecrets, err
		}
		// check that all secrets was already retrieved
		if len(vaultSecretsList.Values()) == 0 {
			break
		}
	}

	// convert secrets map to json
	JSON, err := json.Marshal(secretsList)
	if err != nil {
		return azureSecrets, err
	}

	azureSecrets = string(JSON)

	return azureSecrets, err
}
