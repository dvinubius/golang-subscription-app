package main

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
)

var routes = []string{
	"/",
	"/login",
	"/logout",
	"/register",
	"/activate",
	"/members/plans",
	"/members/subscribe",
}

func Test_Routes_Exist(t *testing.T) {
	chiRoutes := testApp.routes().(chi.Router)

	for _, route := range routes {
		routeExists(t, chiRoutes, route)
	}
}

func routeExists(t *testing.T, routes chi.Router, route string) {
	found := false

	walkFn := func(method string, foundRoute string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if route == foundRoute {
			found = true
		}
		return nil
	}
	_ = chi.Walk(routes, walkFn)

	if !found {
		t.Errorf("Did not find %s in registered routes", route)
	}
}
