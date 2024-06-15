package handler

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/yomek33/talki/internal/logger"
	"google.golang.org/api/option"
)

type Firebase struct {
	App        *firebase.App
	AuthClient *auth.Client
}

const serviceAccountJsonPath = "./service-account-credentials.json"

func InitFirebase(ctx context.Context) (*Firebase, error) {
	opt := option.WithCredentialsFile(serviceAccountJsonPath)

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		logger.Errorf("Error initializing app: %v", err)
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	auth, err := app.Auth(ctx)
	if err != nil {
		logger.Errorf("Error getting Auth client :%v", err)
		return nil, fmt.Errorf("error getting Auth client: %v", err)
	}
	logger.Infof("Firebase initialized successfully")
	return &Firebase{App: app, AuthClient: auth}, nil
}
