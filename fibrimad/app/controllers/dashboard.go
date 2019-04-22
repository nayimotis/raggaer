package controllers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/justinas/nosurf"

	"github.com/nayimotis/raggaer/fibrimad/app/models"
	"github.com/nayimotis/raggaer/fibrimad/app/role"

	"github.com/gin-gonic/gin"
)

// Dashboard shows the application dashboard
func Dashboard(c *gin.Context) {
	cfg, db, ok := retrieveContextValues(c)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Check if user is logged in
	user, ok := loggedUser(c)
	if !ok {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	// Get user work orders
	workOrders, err := models.GetUserWorkOrders(db, user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Get session
	session := sessions.Default(c)

	// Get flash errors
	flashError := session.Get("error")
	session.Delete("error")
	session.Save()

	// Create XSRF token
	token := nosurf.Token(c.Request)

	c.HTML(http.StatusOK, "dashboard.html", map[string]interface{}{
		"user":              user,
		"user_list":         role.UserHasRole(cfg, user, role.UserList),
		"create_user":       role.UserHasRole(cfg, user, role.CreateUser),
		"create_work_order": role.UserHasRole(cfg, user, role.CreateWorkOrder),
		"work_order_list":   role.UserHasRole(cfg, user, role.WorkOrderList),
		"error":             flashError,
		"roleList":          c.MustGet("roleList"),
		"token":             token,
		"workOrders":        workOrders,
	})
}
