package usersvc

import (
	"context"
	"fmt"
	"github.com/sisukasco/commons/crypto"
	"github.com/sisukasco/henki/pkg/db"
	"log"
	"strings"
	"time"
)

func (usvc *UserService) GetAPIKeys(ctx context.Context, userID string) ([]db.GetAPIKeysRow, error) {
	return usvc.svc.DB.Q.GetAPIKeys(ctx, userID)
}

func makeKey(secret string, userID string) string {

	created := time.Date(2020, time.December, 28, 10, 20, 0, 0, time.UTC)
	timestamp := fmt.Sprintf("%v", time.Now().Unix()-created.Unix())

	key, err := crypto.EncryptAES(userID+","+timestamp, secret)
	if err != nil {
		log.Printf("Error encrypting API Key %v ", err)
		return ""
	}
	return key
}
func (usvc *UserService) createAPIKey(ctx context.Context, userID string) string {
	secret := usvc.svc.Konf.String("api.api_key.secret")
	key := makeKey(secret, userID)

	t := 0
	for ; t < 100; t++ {

		if len(key) > 8 {
			exists, err := usvc.svc.DB.Q.DoesAPIKeyExist(ctx, key)
			if exists == false && err == nil {
				break
			}
		}

		key = makeKey(secret, userID)
	}
	return key
}

func (usvc *UserService) NewAPIKey(ctx context.Context, userID string) (string, error) {
	key := usvc.createAPIKey(ctx, userID)

	_, err := usvc.svc.DB.Q.NewApiKey(ctx, db.NewApiKeyParams{Key: key, UserID: userID})

	if err != nil {
		return "", err
	}
	return key, nil
}

func (usvc *UserService) DeleteAPIKey(ctx context.Context, apiKey string, userID string) error {

	return usvc.svc.DB.Q.DeleteAPIKey(ctx, db.DeleteAPIKeyParams{Key: apiKey, UserID: userID})
}

func (usvc *UserService) GetUserFromAPIKey(ctx context.Context, apiKey string) (*db.User, error) {
	user, err := usvc.svc.DB.Q.GetUserFromAPIKey(ctx, apiKey)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (usvc *UserService) GetUserIDFromAPIKey(apiKey string) (string, error) {

	secret := usvc.svc.Konf.String("api.api_key.secret")
	strKey, err := crypto.DecryptAES(apiKey, secret)
	if err != nil {
		return "", err
	}

	log.Printf("Decrypted API Key %v ", strKey)
	parts := strings.Split(strKey, ",")
	if len(parts) == 2 {
		return parts[0], nil
	}

	return "", nil
}
