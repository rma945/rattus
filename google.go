package main

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

func getGoogleSecrets(googleSecret string) (string, error) {
	var googleSecrets string

	googleClientContext := context.Background()
	googleClient, err := secretmanager.NewClient(googleClientContext)
	if err != nil {
		return googleSecrets, err
	}

	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: googleSecret,
	}

	resonse, err := googleClient.AccessSecretVersion(googleClientContext, accessRequest)
	if err != nil {
		return googleSecrets, err
	}

	googleSecrets = fmt.Sprintf("%s", resonse.Payload.Data)

	return googleSecrets, err
}
