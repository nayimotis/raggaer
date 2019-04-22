package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/nayimotis/raggaer/fibrimad/app/config"
	"github.com/nayimotis/raggaer/fibrimad/app/controllers"
	"github.com/nayimotis/raggaer/fibrimad/app/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/justinas/nosurf"
)

func main() {
	// Register gob types
	gob.Register(&models.User{})
	gob.Register(&controllers.CreateUserForm{})
	gob.Register(&controllers.CreateWorkForm{})

	// Load config file
	configFile, err := config.LoadConfigFile("app/config/config.toml")
	if err != nil {
		log.Fatal(err)
	}

	// Set gin release mode
	switch configFile.Mode {
	case "debug":
		gin.SetMode(gin.DebugMode)
	case "release":
		gin.SetMode(gin.ReleaseMode)
	}

	// Load database connection
	database, err := openDatabaseConnection(
		fmt.Sprintf(
			"%s:%s@/%s?parseTime=true",
			configFile.Database.Username,
			configFile.Database.Password,
			configFile.Database.Database,
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create session store
	sessionStore := createSessionStore(configFile.SessionSecret)

	// Create new gin router
	router := gin.Default()
	router.Static("/public", filepath.Join("app", "public"))
	router.SetFuncMap(setTemplateFuncMap())
	router.LoadHTMLGlob(filepath.Join("app", "views", "*.html"))

	// Set middlewares
	router.Use(config.ConfigUseContext(configFile))
	router.Use(databaseUseContext(database))
	router.Use(sessions.Sessions("mysession", sessionStore))
	router.Use(controllers.LoggedUserContext)

	// Set routes
	setRoutes(router)

	// Run gin http server
	http.ListenAndServe(configFile.Address, nosurf.NewPure(router))
}
