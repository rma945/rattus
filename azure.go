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
	var azureSecrets string

	azureBasicClient := keyvault.New()
	azureAuthorizer, err := auth.NewAuthorizerFromEnvironmentWithResource(azureKeyVaultResourceURL)
	if err == nil {
		azureBasicClient.Authorizer = azureAuthorizer
	}

	vaultSecretsList, err := azureBasicClient.GetSecrets(context.Background(), azureVault, nil)
	if err != nil {
		return azureSecrets, err
	}

	secretsList := make(map[string]string)
	for _, secret := range vaultSecretsList.Values() {
		secretValue, err := azureBasicClient.GetSecret(context.Background(), azureVault, path.Base(*secret.ID), "")
		if err == nil {
			secretsList[path.Base(*secret.ID)] = *secretValue.Value
		}
	}

	JSON, err := json.Marshal(secretsList)
	if err != nil {
		return azureSecrets, err
	}

	azureSecrets = string(JSON)

	return azureSecrets, err
}
