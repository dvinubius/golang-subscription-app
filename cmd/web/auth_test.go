package main

import (
	"net/http"
	"testing"
)

func Test_IsAuthenticated(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	ctx := getCtxWithSession(r)
	r = r.WithContext(ctx)

	auth := testApp.IsAuthenticated(r)
	if auth {
		t.Error("says it is authenticated when it shouldn't be")
	}

	testApp.Session.Put(ctx, "userID", 1)

	auth = testApp.IsAuthenticated(r)
	if !auth {
		t.Error("says it is not authenticated when it should be")
	}
}
