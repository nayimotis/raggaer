package controllers

import (
	"net/http"

	"github.com/nayimotis/raggaer/fibrimad/app/role"

	"golang.org/x/crypto/bcrypt"

	"github.com/nayimotis/raggaer/fibrimad/app/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/justinas/nosurf"
)

// LoginForm shows the login form
func LoginForm(c *gin.Context) {
	// Retrieve user session
	session := sessions.Default(c)

	if _, ok := loggedUser(c); ok {
		c.Redirect(http.StatusSeeOther, "/dashboard")
		return
	}

	// Create XSRF token
	token := nosurf.Token(c.Request)

	flashError := session.Get("error")
	session.Delete("error")
	session.Save()

	c.HTML(http.StatusOK, "login.html", map[string]interface{}{
		"token":    token,
		"error":    flashError,
		"roleList": c.MustGet("roleList"),
	})
}

// LoginProcess executes the login form
func LoginProcess(c *gin.Context) {

	// Retrieve user session
	session := sessions.Default(c)

	// Retrieve context values
	cfg, db, ok := retrieveContextValues(c)
	if !ok {
		session.Set("error", "Algo inesperado ha ocurrido. Inténtalo de nuevo")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	if _, ok := loggedUser(c); ok {
		c.Redirect(http.StatusSeeOther, "/dashboard")
		return
	}

	username := c.PostForm("username")
	password := []byte(c.PostForm("password"))

	// Retrieve user from database
	user, err := models.GetUserByUsername(db, username)
	if err != nil {
		session.Set("error", "Usuario o contraseña incorrectos")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), password); err != nil {
		session.Set("error", "Usuario o contraseña incorrectos")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	if !role.UserHasRole(cfg, user, role.Login) {
		session.Set("error", "Tu cuenta no tiene permisos para iniciar sesión")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	// Set session values
	session.Set("logged_user", user.ID)
	session.Set("is_logged", true)
	if err := session.Save(); err != nil {
		session.Set("error", "Algo inesperado ha ocurrido. Inténtalo de nuevo")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	// Add user login message
	if _, err := user.AddLogMessage(db, "El usuario ha iniciado sesión"); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/dashboard")
}

// Logout closes the user session
func Logout(c *gin.Context) {
	_, db, ok := retrieveContextValues(c)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Check if user is logged
	user, ok := loggedUser(c)
	if !ok {
		c.Redirect(http.StatusSeeOther, "/dashboard")
		return
	}

	// Add user logout message
	if _, err := user.AddLogMessage(db, "El usuario ha cerrado sesión"); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Get session
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.Redirect(http.StatusSeeOther, "/login")
}
