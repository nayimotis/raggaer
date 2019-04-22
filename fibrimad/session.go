package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
)

func createSessionStore(secret string) sessions.Store {
	store := cookie.NewStore([]byte(secret))
	return store
}
