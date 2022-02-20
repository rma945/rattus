package main

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

func getsecretGoogles(secretGoogle string) (string, error) {
	var secretGoogles string

	googleClientContext := context.Background()
	googleClient, err := secretmanager.NewClient(googleClientContext)
	if err != nil {
		return secretGoogles, err
	}

	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretGoogle,
	}

	resonse, err := googleClient.AccessSecretVersion(googleClientContext, accessRequest)
	if err != nil {
		return secretGoogles, err
	}

	secretGoogles = fmt.Sprintf("%s", resonse.Payload.Data)

	return secretGoogles, err
}
