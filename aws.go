package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// get secrets from aws secret storage
func getsecretAWSString(secretAWS string) ([]string, error) {
	var secrets []string

	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	awsService := secretsmanager.New(awsSession)

	for _, s := range stringToList(secretAWS) {
		awsRequest := &secretsmanager.GetSecretValueInput{
			SecretId:     aws.String(s),
			VersionStage: aws.String("AWSCURRENT"),
		}
		awsResponse, err := awsService.GetSecretValue(awsRequest)
		if err != nil {
			return secrets, err
		}

		secrets = append(secrets, *awsResponse.SecretString)
		if err != nil {
			return secrets, fmt.Errorf("failed to retrive vault secrets - %s", err)
		}
	}

	return secrets, nil
}
