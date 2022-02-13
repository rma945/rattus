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
	VaultSecretURL         string
	AWSSecretName          string
	AWSRegion              string
	AWSKeyID               string
	AWSKeySecret           string
	AWSSessionToken        string
	AzureTenantID          string
	AzureClientID          string
	AzureClientSecret      string
	AzureVault             string
	GoogleSecret           string
	Debug                  *bool
}

func initializeConfiguration() applicationConfig {
	// default configuration
	c := &applicationConfig{}
	c.SecretProvider = "none"
	c.AWSRegion = "us-east-1"

	// cli arguments
	argVaultSecret := flag.String("vault-secret", "", "Vault secret URL - https://vault.example.io/v1/storage/secret\nenv: VAULT_SECRET\n")
	argVaultToken := flag.String("vault-token", "", "Vault authentication token\nenv: VAULT_TOKEN")

	argAWSSecretName := flag.String("aws-secret-name", "", "AWS secret name - example-project-backend\nenv: AWS_SECRET_NAME\n")
	argAWSRegion := flag.String("aws-region", "", "AWS region - us-east-1\nenv: AWS_DEFAULT_REGION or AWS_REGION\n")
	argAWSKeyID := flag.String("aws-key-id", "", "AWS account ID\nenv: AWS_ACCESS_KEY_ID\n")
	argAWSKeySecret := flag.String("aws-key-secret", "", "AWS account secret\nAWS_SECRET_ACCESS_KEY\n")
	argAWSSessionToken := flag.String("aws-session-token", "", "AWS session token\nAWS_SESSION_TOKEN\n")

	argTemplatePath := flag.String("template", "", "Path to template file - /app/config/production.template\nenv: TEMPLATE\n")
	argAzureTenantID := flag.String("azure-tenant-id", "", "Azure tenant ID\nenv: AZURE_TENANT_ID\n")
	argAzureClientID := flag.String("azure-client-id", "", "Azure client ID\nenv: AZURE_CLIENT_ID\n")
	argAzureClientSecret := flag.String("azure-client-secret", "", "Azure client Secret\nenv: AZURE_CLIENT_SECRET\n")
	argAzureVault := flag.String("azure-vault", "", "Azure keyvault storage URL - https://example-key-vault.vault.azure.net/\nenv: AZURE_VAULT\n")

	argGoogleSecret := flag.String("google-secret", "", "Google SecretManager secret - projects/xxxxxxxxxxx/secrets/example-secret/versions/latest \nenv: GOOGLE_SECRET\n")

	c.Debug = flag.Bool("debug", false, "Enable debug information\n")

	flag.Parse()

	// vault secret
	envVaultSecret := os.Getenv("VAULT_SECRET")
	if envVaultSecret != "" {
		c.VaultSecretURL = envVaultSecret
		c.SecretProvider = "vault"
	}
	if *argVaultSecret != "" {
		c.VaultSecretURL = *argVaultSecret
		c.SecretProvider = "vault"
	}

	// vault token
	envVaultToken := os.Getenv("VAULT_TOKEN")
	if envVaultToken != "" {
		c.VaultToken = envVaultToken
	}
	if *argVaultToken != "" {
		c.VaultToken = *argVaultToken
	}

	// aws secret name
	envAWSSecretName := os.Getenv("AWS_SECRET_NAME")
	if envAWSSecretName != "" {
		c.AWSSecretName = envAWSSecretName
		c.SecretProvider = "aws"
	}
	if *argAWSSecretName != "" {
		c.AWSSecretName = *argAWSSecretName
		c.SecretProvider = "aws"
	}

	// aws default region
	envAWSDefaultRegion := os.Getenv("AWS_DEFAULT_REGION")
	if envAWSDefaultRegion != "" {
		c.AWSRegion = envAWSDefaultRegion
	}

	// aws region
	envAWSRegion := os.Getenv("AWS_REGION")
	if envAWSRegion != "" {
		c.AWSRegion = envAWSRegion
	}

	if *argAWSRegion != "" {
		c.AWSRegion = *argAWSRegion
	}

	// aws id
	envAWSKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	if envAWSKeyID != "" {
		c.AWSKeyID = envAWSKeyID
	}
	if *argAWSKeyID != "" {
		c.AWSKeyID = *argAWSKeyID
	}

	// aws secret
	envAWSKeySecret := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if envAWSKeySecret != "" {
		c.AWSKeySecret = envAWSKeySecret
	}
	if *argAWSKeySecret != "" {
		c.AWSKeySecret = *argAWSKeySecret
	}

	// aws session token
	envAWSSessionToken := os.Getenv("AWS_SESSION_TOKEN")
	if envAWSSessionToken != "" {
		c.AWSSessionToken = envAWSSessionToken
	}
	if *argAWSSessionToken != "" {
		c.AWSSessionToken = *argAWSSessionToken
	}

	envAzureTenantID := os.Getenv("AZURE_TENANT_ID")
	if envAzureTenantID != "" {
		c.AzureTenantID = envAzureTenantID
	}
	if *argAWSKeySecret != "" {
		c.AzureTenantID = *argAzureTenantID
	}

	envAzureClientID := os.Getenv("AZURE_CLIENT_ID")
	if envAzureClientID != "" {
		c.AzureClientID = envAzureClientID
	}
	if *argAWSKeySecret != "" {
		c.AzureClientID = *argAzureClientID
	}

	envAzureClientSecret := os.Getenv("AZURE_CLIENT_SECRET")
	if envAzureClientSecret != "" {
		c.AzureClientSecret = envAzureClientSecret
	}
	if *argAWSKeySecret != "" {
		c.AzureClientSecret = *argAzureClientSecret
	}

	envAzureVault := os.Getenv("AZURE_VAULT")
	if envAzureVault != "" {
		c.AzureVault = envAzureVault
		c.SecretProvider = "azure"
	}
	if *argAWSKeySecret != "" {
		c.AzureVault = *argAzureVault
		c.SecretProvider = "azure"
	}

	envGoogleSecret := os.Getenv("GOOGLE_SECRET")
	if envGoogleSecret != "" {
		c.GoogleSecret = envGoogleSecret
		c.SecretProvider = "google"
	}
	if *argGoogleSecret != "" {
		c.GoogleSecret = *argGoogleSecret
		c.SecretProvider = "google"
	}

	// old template path variable for a versions compatability
	envTemplatePath := os.Getenv("TEMPLATE_PATH")
	if envTemplatePath != "" {
		c.TemplatePath = envTemplatePath
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
