package main

import (
	"flag"
	"fmt"
	"os"
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
	argAWSRegion := flag.String("aws-region", "", "AWS region - us-east-1\nenv: AWS_DEFAULT_REGION\n")
	argAWSKeyID := flag.String("aws-key-id", "", "AWS account ID\nenv: AWS_ACCESS_KEY_ID\n")
	argAWSKeySecret := flag.String("aws-key-secret", "", "AWS account secret\nAWS_SECRET_ACCESS_KEY\n")

	argTemplatePath := flag.String("template", "", "Path to template file - /app/config/production.template\nenv: TEMPLATE_PATH\n")
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

	// aws region
	envAWSRegion := os.Getenv("AWS_DEFAULT_REGION")
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

	// template
	envTemplatePath := os.Getenv("TEMPLATE_PATH")
	if envTemplatePath != "" {
		c.TemplatePath = envTemplatePath
	}
	if *argTemplatePath != "" {
		c.TemplatePath = *argTemplatePath
	}

	return *c
}

func main() {
	// initialize configuration
	var secrets string
	var err error
	config := initializeConfiguration()

	if *config.Debug {
		fmt.Printf("Secret provider: %s\n", config.SecretProvider)
	}

	// get secrets
	switch config.SecretProvider {
	case "vault":
		secrets, err = vaultGetSecret(config)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}

	case "aws":
		secrets, err = getAWSSecretString(config.AWSSecretName, config.AWSRegion, config.AWSKeyID, config.AWSKeySecret)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}

	case "azure":
		secrets, err = getAzureSecrets(config.AzureVault)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}

	case "google":
		secrets, err = getGoogleSecrets(config.GoogleSecret)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
	}

	// render output as template,json or text
	secrets, err = renderOutput(secrets, config.TemplatePath)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	// show secrets output and exit
	fmt.Println(secrets)
	os.Exit(0)
}
