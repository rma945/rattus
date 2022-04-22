package main

import (
	"context"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

func getsecretGoogles(secretGoogle string) ([]string, error) {
	var secrets []string

	googleClientContext := context.Background()
	googleClient, err := secretmanager.NewClient(googleClientContext)
	if err != nil {
		return secrets, err
	}

	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretGoogle,
	}

	response, err := googleClient.AccessSecretVersion(googleClientContext, accessRequest)
	if err != nil {
		return secrets, err
	}

	// secretGoogles = fmt.Sprintf("%s", resonse.Payload.Data)

	return secrets, nil
}
