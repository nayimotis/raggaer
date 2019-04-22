package main

import "github.com/gin-contrib/sessions"

func createSessionStore(secret string) sessions.CookieStore {
	return sessions.NewCookieStore([]byte(secret))
}
