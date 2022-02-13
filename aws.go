package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func createAWSSession(AWSRegion, AWSKeyID, AWSKeySecret, AWSSessionToken string) (*session.Session, error) {
	var awsCredentials *credentials.Credentials
	var awsSession *session.Session

	if (AWSKeyID != "") && (AWSKeySecret != "") {
		awsCredentials = credentials.NewStaticCredentials(AWSKeyID, AWSKeySecret, AWSSessionToken)
	}

	awsSession, err := session.NewSession(&aws.Config{
		Region:      aws.String(AWSRegion),
		Credentials: awsCredentials,
	})

	return awsSession, err
}

func getAWSSecretString(AWSSecrets, AWSRegion, AWSKeyID, AWSKeySecret, AWSSessionToken string) ([]string, error) {
	var secrets []string

	awsSession, err := createAWSSession(AWSRegion, AWSKeyID, AWSKeySecret, AWSSessionToken)
	if err != nil {
		return secrets, err
	}

	awsService := secretsmanager.New(awsSession)

	for _, s := range stringToList(AWSSecrets) {
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

	return secrets, err
}
