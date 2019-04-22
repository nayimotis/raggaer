package controllers

import (
	"database/sql"

	"github.com/raggaer/fibrimad/app/role"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/raggaer/fibrimad/app/config"
	"github.com/raggaer/fibrimad/app/models"
)

type roleListSidebar struct {
	UserList   bool
	CreateUser bool
	CreateWork bool
	WorkList   bool
}

// LoggedUserContext sets the logged user inside the context
func LoggedUserContext(c *gin.Context) {
	cfg, _, ok := retrieveContextValues(c)
	if !ok {
		c.Set("roleList", roleListSidebar{})
		c.Next()
		return
	}

	loggedUser, isLogged := loggedUser(c)
	if !isLogged {
		c.Set("roleList", roleListSidebar{})
		c.Next()
		return
	}

	roleList := roleListSidebar{
		UserList:   role.UserHasRole(cfg, loggedUser, role.UserList),
		CreateUser: role.UserHasRole(cfg, loggedUser, role.CreateUser),
		WorkList:   role.UserHasRole(cfg, loggedUser, role.WorkOrderList),
		CreateWork: role.UserHasRole(cfg, loggedUser, role.CreateWorkOrder),
	}

	c.Set("roleList", roleList)
	c.Next()
}

func loggedUser(c *gin.Context) (*models.User, bool) {
	session := sessions.Default(c)

	_, db, ok := retrieveContextValues(c)
	if !ok {
		return nil, false
	}

	logged, ok := session.Get("is_logged").(bool)
	if !ok {
		return nil, false
	}

	if !logged {
		return nil, false
	}

	loggedUserID, ok := session.Get("logged_user").(int64)
	if !ok {
		return nil, false
	}

	loggedUser, err := models.GetUserByID(db, loggedUserID)
	if err != nil {
		return nil, false
	}

	return loggedUser, true
}

func retrieveContextValues(c *gin.Context) (*config.Config, *sql.DB, bool) {
	cfg, ok := c.Get("config")
	if !ok {
		return nil, nil, false
	}

	cfgp, ok := cfg.(*config.Config)
	if !ok {
		return nil, nil, false
	}

	db, ok := c.Get("db")
	if !ok {
		return nil, nil, false
	}

	dbp, ok := db.(*sql.DB)
	if !ok {
		return nil, nil, false
	}

	return cfgp, dbp, true
}
