package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_AddDefaultData(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	ctx := getCtxWithSession(r)
	r = r.WithContext(ctx)

	testApp.Session.Put(ctx, "flash", "flash123")
	testApp.Session.Put(ctx, "warning", "warning123")
	testApp.Session.Put(ctx, "error", "error123")

	td := testApp.AddDefaultData(&TemplateData{}, r)

	if td.Flash != "flash123" {
		t.Error("Failed to get flash data")
	}
	if td.Warning != "warning123" {
		t.Error("Failed to get warning data")
	}
	if td.Error != "error123" {
		t.Error("Failed to get error data")
	}
}

func Test_render(t *testing.T) {
	pathToTemplates = "./templates"
	rr := httptest.NewRecorder()

	r, _ := http.NewRequest("GET", "/", nil)
	ctx := getCtxWithSession(r)
	r = r.WithContext(ctx)

	testApp.render(rr, r, "home.page.gohtml", &TemplateData{})

	if rr.Code != 200 {
		t.Error("Failed to render page")
	}
}
