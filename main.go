package main

import (
	"fmt"
	"os"
)

func main() {
	var output string
	var secrets []string
	var err error

	// initialize configuration
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
		secrets, err = getAWSSecretString(config.AWSSecretName, config.AWSRegion, config.AWSKeyID, config.AWSKeySecret, config.AWSSessionToken)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}

		// case "azure":
		// 	secrets, err = getAzureSecrets(config.AzureVault)
		// 	if err != nil {
		// 		fmt.Printf("Error: %s\n", err)
		// 		os.Exit(1)
		// 	}

		// case "google":
		// 	secrets, err = getGoogleSecrets(config.GoogleSecret)
		// 	if err != nil {
		// 		fmt.Printf("Error: %s\n", err)
		// 		os.Exit(1)
		// 	}
	}

	// render output as template,json or text
	output, err = renderOutput(secrets, config.TemplatePath)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	// show secrets output and exit
	fmt.Println(output)
	os.Exit(0)
}
