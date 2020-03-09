package main

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func createAWSSession(AWSRegion, AWSKeyID, AWSKeySecret string) (*session.Session, error) {
	var awsCredentials *credentials.Credentials
	var awsSession *session.Session

	if (AWSKeyID != "") && (AWSKeySecret != "") {
		awsCredentials = credentials.NewStaticCredentials(AWSKeyID, AWSKeySecret, "")
	}

	awsSession, err := session.NewSession(&aws.Config{
		Region:      aws.String(AWSRegion),
		Credentials: awsCredentials,
	})

	return awsSession, err
}

func getAWSSecret(secretName, AWSRegion, AWSKeyID, AWSKeySecret string) (map[string]interface{}, error) {
	var secrets map[string]interface{}
	awsSession, err := createAWSSession(AWSRegion, AWSKeyID, AWSKeySecret)
	if err != nil {
		return secrets, err
	}

	awsService := secretsmanager.New(awsSession)
	awsRequest := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"),
	}

	awsResponse, err := awsService.GetSecretValue(awsRequest)
	if err != nil {
		return secrets, err
	}

	if err := json.Unmarshal([]byte(*awsResponse.SecretString), &secrets); err != nil {
		return secrets, err
	}

	return secrets, err
}
