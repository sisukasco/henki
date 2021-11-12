package usersvc_test

import (
	"context"
	stesting "github.com/sisukasco/henki/pkg/testing"
	"github.com/sisukasco/henki/pkg/usersvc"
	"testing"
)

func TestCreatingAPIKey(t *testing.T) {
	user, err := stesting.SignupUser(t, svc)
	if err != nil {
		t.Fatalf("Error creating a new user %v", err)
		return
	}
	userSvcObj := usersvc.New(svc)

	ctx = context.Background()

	apikey1, err := userSvcObj.NewAPIKey(ctx, user.ID)
	if err != nil {
		t.Fatalf("Error creating a new user %v", err)
		return
	}

	if len(apikey1) <= 8 {
		t.Errorf("Error creating API Key generated key %s", apikey1)
	}
	t.Logf("API key 1 created %s", apikey1)

	u1, err := userSvcObj.GetUserFromAPIKey(ctx, apikey1)

	if err != nil {
		t.Fatalf("Error getting User from API Key %v", err)
	}

	if u1.Email != user.Email {
		t.Errorf("The user from API Key does not match %s %s", u1.Email, user.Email)
	}
	t.Logf("Got User back from API Key %s ", u1.Email)

	userID1, err := userSvcObj.GetUserIDFromAPIKey(apikey1)
	if err != nil {
		t.Fatalf("Error getting userID from API Key %v", err)
	}
	t.Logf("Got UserID back from API Key %s ", userID1)
	if userID1 != u1.ID {
		t.Fatalf("Expected UserID %s received %s ", u1.ID, userID1)
	}

	apikey2, err := userSvcObj.NewAPIKey(ctx, user.ID)
	if err != nil {
		t.Fatalf("Error creating a new user %v", err)
		return
	}
	if len(apikey2) <= 8 {
		t.Errorf("Error creating API Key generated key %s", apikey2)
	}
	t.Logf("API key 2 created %s", apikey2)

	keys, err := userSvcObj.GetAPIKeys(ctx, user.ID)
	if err != nil {
		t.Fatalf("Error getting API keys %v", err)
	}
	if keys[0].Key != apikey1 {
		t.Errorf("GetKey received key does not match!  %v vs %s  ", keys[0], apikey1)
	}
	if keys[1].Key != apikey2 {
		t.Errorf("GetKey received key does not match!  %v vs %s  ", keys[1], apikey2)
	}
	t.Logf("API Keys got back for user %s, keys %v %v", user.ID, keys[0], keys[1])
}

func TestDeletingAPIKey(t *testing.T) {
	user, err := stesting.SignupUser(t, svc)
	if err != nil {
		t.Fatalf("Error creating a new user %v", err)
		return
	}
	userSvcObj := usersvc.New(svc)

	ctx = context.Background()

	apikey1, err := userSvcObj.NewAPIKey(ctx, user.ID)
	if err != nil {
		t.Fatalf("Error creating a new user %v", err)
		return
	}
	t.Logf("Created new API Key %s", apikey1)

	err = userSvcObj.DeleteAPIKey(ctx, apikey1, user.ID)
	if err != nil {
		t.Fatalf("Error deleting  apiKey %v", err)
		return
	}

	keys, err := userSvcObj.GetAPIKeys(ctx, user.ID)
	if err != nil {
		t.Fatalf("Error getting API keys %v", err)
	}

	if len(keys) > 0 {
		t.Errorf("Expected to delete the key Keys in db %d", len(keys))
	}
}
