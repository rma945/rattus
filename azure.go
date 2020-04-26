package main

import (
	"context"
	"encoding/json"
	"os"
	"path"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
)

func getAzureSecrets(azureTenantID, azureClientID, azureClientSecret, azureVault string) (string, error) {
	var azureSecrets string

	os.Setenv("AZURE_TENANT_ID", azureTenantID)
	os.Setenv("AZURE_CLIENT_ID", azureClientID)
	os.Setenv("AZURE_CLIENT_SECRET", azureClientSecret)

	azureAuthorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return azureSecrets, err
	}

	azureBasicClient := keyvault.New()
	azureBasicClient.Authorizer = azureAuthorizer

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
