package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Index shows the index page
func Index(c *gin.Context) {
	if _, ok := loggedUser(c); ok {
		c.Redirect(http.StatusSeeOther, "/dashboard")
		return
	}

	c.Redirect(http.StatusSeeOther, "/login")
}
