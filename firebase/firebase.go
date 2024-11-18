package firebase

import (
	"context"
	"log"
	"os"

	"firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

var FirebaseAuth *auth.Client

func InitFirebase() {
	rawJSON := os.Getenv("FIREBASE_SA_KEY")
    if rawJSON == "" {
        log.Fatalf("FIREBASE_SA_KEY not set in environment variables")
    }

    opt := option.WithCredentialsJSON([]byte(rawJSON))

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v\n", err)
	}

	FirebaseAuth, err = app.Auth(context.Background())
	if err != nil {
		log.Fatalf("Error initializing Firebase Auth client: %v\n", err)
	}
}