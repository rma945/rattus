package main

import (
	"encoding/json"
	"flag"
	"os"
	"strings"
	"unicode"
)

type applicationConfig struct {
	SecretProvider         string
	K8SServiceAccount      string
	K8SServiceAccountToken string
	TemplatePath           string
	VaultToken             string
	secretVault            string
	secretAWS              string
	secretAzure            string
	secretGoogle           string
	Debug                  *bool
}

func initializeConfiguration() applicationConfig {
	// default configuration
	c := &applicationConfig{}
	c.SecretProvider = "none"

	// cli arguments
	argVaultToken := flag.String("vault-token", "", "Vault authentication token\nenv: VAULT_TOKEN")
	argSecretVault := flag.String("vault-secret", "", "Vault secret URL - https://vault.example.io/v1/storage/secret\nenv: VAULT_SECRET\n")
	argSecretAWS := flag.String("aws-secret", "", "AWS secret name - example-project-backend\nenv: AWS_SECRET\n")
	argSecretAzure := flag.String("azure-secret", "", "Azure keyvault secret URL - https://example-vault.vault.azure.net/secrets/example-secret\nenv: AZURE_SECRET\n")
	argSecretGoogle := flag.String("google-secret", "", "Google SecretManager secret - projects/xxxxxxxxxxx/secrets/example-secret/versions/latest \nenv: GOOGLE_SECRET\n")
	argTemplatePath := flag.String("template", "", "Path to template file - /app/config/production.template\nenv: TEMPLATE\n")
	c.Debug = flag.Bool("debug", false, "Enable debug information\n")

	flag.Parse()

	// vault token
	envVaultToken := os.Getenv("VAULT_TOKEN")
	if envVaultToken != "" {
		c.VaultToken = envVaultToken
	}
	if *argVaultToken != "" {
		c.VaultToken = *argVaultToken
	}

	// vault secret
	envSecretVault := os.Getenv("VAULT_SECRET")
	if envSecretVault != "" {
		c.secretVault = envSecretVault
		c.SecretProvider = "vault"
	}
	if *argSecretVault != "" {
		c.secretVault = *argSecretVault
		c.SecretProvider = "vault"
	}

	// aws secret
	envSecretAWS := os.Getenv("AWS_SECRET")
	if envSecretAWS != "" {
		c.secretAWS = envSecretAWS
		c.SecretProvider = "aws"
	}
	if *argSecretAWS != "" {
		c.secretAWS = *argSecretAWS
		c.SecretProvider = "aws"
	}

	// azure secret
	envAzureVault := os.Getenv("AZURE_SECRET")
	if envAzureVault != "" {
		c.secretAzure = envAzureVault
		c.SecretProvider = "azure"
	}
	if *argSecretAzure != "" {
		c.secretAzure = *argSecretAzure
		c.SecretProvider = "azure"
	}

	// google secret
	envSecretGoogle := os.Getenv("GOOGLE_SECRET")
	if envSecretGoogle != "" {
		c.secretGoogle = envSecretGoogle
		c.SecretProvider = "google"
	}
	if *argSecretGoogle != "" {
		c.secretGoogle = *argSecretGoogle
		c.SecretProvider = "google"
	}

	// template
	envTemplate := os.Getenv("TEMPLATE")
	if envTemplate != "" {
		c.TemplatePath = envTemplate
	}
	if *argTemplatePath != "" {
		c.TemplatePath = *argTemplatePath
	}

	return *c
}

// convert secrets string into list
func stringToList(str string) []string {
	return strings.FieldsFunc(str, func(c rune) bool { return unicode.IsSpace(c) || c == ',' || c == ';' })
}

// merge list of secret strings into single map
func mergeSecretListToMap(secrets []string) (map[string]interface{}, error) {
	var secretJSON map[string]interface{}

	mergedSecrets := make(map[string]interface{})
	for _, secret := range secrets {
		if err := json.Unmarshal([]byte(secret), &secretJSON); err != nil {
			return mergedSecrets, err
		}
		for key := range secretJSON {
			mergedSecrets[key] = secretJSON[key]
		}
	}

	return mergedSecrets, nil
}
