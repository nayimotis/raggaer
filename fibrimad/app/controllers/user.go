package controllers

import (
	"net/http"
	"strconv"

	"github.com/nayimotis/raggaer/fibrimad/app/role"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/justinas/nosurf"
	"github.com/nayimotis/raggaer/fibrimad/app/models"
	"golang.org/x/crypto/bcrypt"
)

// CreateUserForm defines an user creation form
type CreateUserForm struct {
	Username string
	Password string
	Role     string
}

// CreateUserProcessForm process the create user form
func CreateUserProcessForm(c *gin.Context) {
	cfg, db, ok := retrieveContextValues(c)
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

	// Create session
	session := sessions.Default(c)

	// Check if user has "create_user" role
	if !role.UserHasRole(cfg, user, role.CreateUser) {
		session.Set("error", "Tu cuenta no tiene permisos para crear un nuevo usuario")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/dashboard")
		return
	}

	username := c.PostForm("username")
	password := c.PostForm("password")
	role := c.PostForm("role")

	// Check for valid password
	if len(password) <= 6 {
		session.Set("error", "La contraseña es demasiado corta, se deben usar como mínimo 6 caracteres. Inténtalo de nuevo")
		session.Set("createUserForm", CreateUserForm{username, password, role})
		session.Save()
		c.Redirect(http.StatusSeeOther, "/admin/user/create")
		return
	}

	// Check for valid role
	if _, ok := cfg.Roles[role]; !ok {
		session.Set("error", "Rol seleccionado desconocido. Inténtalo de nuevo")
		session.Set("createUserForm", CreateUserForm{username, password, role})
		session.Save()
		c.Redirect(http.StatusSeeOther, "/admin/user/create")
		return
	}

	// Check if username exists
	if _, err := models.GetUserByUsername(db, username); err == nil {
		session.Set("error", "El nombre de usuario ya está en uso. Inténtalo de nuevo")
		session.Set("createUserForm", CreateUserForm{username, password, role})
		session.Save()
		c.Redirect(http.StatusSeeOther, "/admin/user/create")
		return
	}

	// Create hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.Error(err)
		session.Set("error", "Algo inesperado ha ocurrido. Inténtalo de nuevo")
		session.Set("createUserForm", CreateUserForm{username, password, role})
		session.Save()
		c.Redirect(http.StatusSeeOther, "/admin/user/create")
		return
	}

	// Create and insert user
	newUser := models.User{
		Username: username,
		Password: string(hashedPassword),
		Role:     role,
	}

	// Insert user into the database
	if _, err := newUser.CreateUser(db); err != nil {
		c.Error(err)
		session.Set("error", "No se ha podido guardar el usuario en la base de datos. Inténtalo de nuevo")
		session.Set("createUserForm", CreateUserForm{username, password, role})
		session.Save()
		c.Redirect(http.StatusSeeOther, "/admin/user/create")
		return
	}

	// Set success message
	session.Set("success", "El usuario "+username+" se ha creado con éxito")
	session.Save()

	c.Redirect(http.StatusSeeOther, "/admin/user/list")
}

// CreateUserShowForm shows the form to create a new user
func CreateUserShowForm(c *gin.Context) {
	cfg, _, ok := retrieveContextValues(c)
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

	// Retrieve default session
	session := sessions.Default(c)

	// Check if user has "create_user" role
	if !role.UserHasRole(cfg, user, role.CreateUser) {
		session.Set("error", "Tu cuenta no tiene permisos para crear un nuevo usuario")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/dashboard")
		return
	}

	// Create XSRF token
	token := nosurf.Token(c.Request)

	// Get flash errors
	flashError := session.Get("error")
	session.Delete("error")

	// Get create user form
	createForm := session.Get("createUserForm")
	session.Delete("createUserForm")
	session.Save()

	c.HTML(http.StatusOK, "create_user.html", map[string]interface{}{
		"user":     user,
		"error":    flashError,
		"roles":    cfg.Roles,
		"token":    token,
		"form":     createForm,
		"roleList": c.MustGet("roleList"),
	})
}

// ShowUserList shows all the users of the application
func ShowUserList(c *gin.Context) {
	cfg, db, ok := retrieveContextValues(c)
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

	// Retrieve default session
	session := sessions.Default(c)

	// Check if user has "create_user" role
	if !role.UserHasRole(cfg, user, role.UserList) {
		session.Set("error", "Tu cuenta no tiene permisos para ver la lista de usuarios")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/dashboard")
		return
	}

	// Retrieve success flash message
	flashSuccess := session.Get("success")
	session.Delete("success")

	// Retrieve error flash message
	flashError := session.Get("error")
	session.Delete("error")
	session.Save()

	// Retrieve list of users
	userList, err := models.UserList(db)
	if err != nil {
		c.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Create XSRF token
	token := nosurf.Token(c.Request)

	c.HTML(http.StatusOK, "list_user.html", map[string]interface{}{
		"list":     userList,
		"user":     user,
		"token":    token,
		"success":  flashSuccess,
		"error":    flashError,
		"roleList": c.MustGet("roleList"),
	})
}

// EditUserForm shows the edit user form
func EditUserForm(c *gin.Context) {
	cfg, db, ok := retrieveContextValues(c)
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

	// Retrieve default session
	session := sessions.Default(c)

	// Check if user has "create_user" role
	if !role.UserHasRole(cfg, user, role.EditUser) {
		session.Set("error", "Tu cuenta no tiene permisos para ver editar a un usuario")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/dashboard")
		return
	}

	// Retrieve user identifier
	userIdentifierStr := c.Param("id")
	userID, err := strconv.ParseInt(userIdentifierStr, 10, 64)
	if err != nil {
		c.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve user by its identifier
	u, err := models.GetUserByID(db, userID)
	if err != nil {
		c.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Create XSRF token
	token := nosurf.Token(c.Request)

	// Get flash errors
	flashError := session.Get("error")
	session.Delete("error")

	// Get create user form
	createForm := session.Get("editUserForm")
	session.Delete("editUserForm")
	session.Save()

	c.HTML(http.StatusOK, "edit_user.html", map[string]interface{}{
		"user":     u,
		"error":    flashError,
		"roles":    cfg.Roles,
		"token":    token,
		"form":     createForm,
		"roleList": c.MustGet("roleList"),
	})
}

// EditUserProcessForm process the edit user form
func EditUserProcessForm(c *gin.Context) {
	cfg, db, ok := retrieveContextValues(c)
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

	// Create session
	session := sessions.Default(c)

	// Check if user has "create_user" role
	if !role.UserHasRole(cfg, user, role.EditUser) {
		session.Set("error", "Tu cuenta no tiene permisos para editar un usuario")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/dashboard")
		return
	}

	// Retrieve user identifier
	userIdentifierStr := c.Param("id")
	userID, err := strconv.ParseInt(userIdentifierStr, 10, 64)
	if err != nil {
		c.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve user by its identifier
	u, err := models.GetUserByID(db, userID)
	if err != nil {
		c.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	username := c.PostForm("username")
	password := c.PostForm("password")
	role := c.PostForm("role")

	// Check for valid role
	if _, ok := cfg.Roles[role]; !ok {
		session.Set("error", "Rol seleccionado desconocido. Inténtalo de nuevo")
		session.Set("editUserForm", CreateUserForm{username, password, role})
		session.Save()
		c.Redirect(http.StatusSeeOther, "/admin/user/edit/"+strconv.FormatInt(u.ID, 10))
		return
	}

	// Create and insert user
	newUser := models.User{
		Username: username,
		Password: u.Password,
		Role:     role,
		ID:       u.ID,
	}

	if username != u.Username {
		// Check if username exists
		if _, err := models.GetUserByUsername(db, username); err == nil {
			session.Set("error", "El nombre de usuario ya está en uso. Inténtalo de nuevo")
			session.Set("editUserForm", CreateUserForm{username, password, role})
			session.Save()
			c.Redirect(http.StatusSeeOther, "/admin/user/edit/"+strconv.FormatInt(u.ID, 10))
			return
		}
	}

	if password != "" {
		// Check for valid password
		if len(password) <= 6 {
			session.Set("error", "La contraseña es demasiado corta, se deben usar como mínimo 6 caracteres. Inténtalo de nuevo")
			session.Set("editUserForm", CreateUserForm{username, password, role})
			session.Save()
			c.Redirect(http.StatusSeeOther, "/admin/user/edit/"+strconv.FormatInt(u.ID, 10))
			return
		}

		// Create hashed password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			c.Error(err)
			session.Set("error", "Algo inesperado ha ocurrido. Inténtalo de nuevo")
			session.Set("editUserForm", CreateUserForm{username, password, role})
			session.Save()
			c.Redirect(http.StatusSeeOther, "/admin/user/edit/"+strconv.FormatInt(u.ID, 10))
			return
		}

		newUser.Password = string(hashedPassword)
	}

	// Insert user into the database
	if _, err := newUser.EditUser(db); err != nil {
		c.Error(err)
		session.Set("error", "No se ha podido guardar el usuario en la base de datos. Inténtalo de nuevo")
		session.Set("editUserForm", CreateUserForm{username, password, role})
		session.Save()
		c.Redirect(http.StatusSeeOther, "/admin/user/edit/"+strconv.FormatInt(u.ID, 10))
		return
	}

	// Set success message
	session.Set("success", "El usuario ha sido editado con éxito")
	session.Save()

	c.Redirect(http.StatusSeeOther, "/admin/user/list")
}

// DeleteUser deletes the given user
func DeleteUser(c *gin.Context) {
	cfg, db, ok := retrieveContextValues(c)
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

	// Create session
	session := sessions.Default(c)

	// Check if user has role
	if !role.UserHasRole(cfg, user, role.DeleteUser) {
		session.Set("error", "Tu cuenta no tiene permisos para eliminar un usuario")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/dashboard")
		return
	}

	// Retrieve user identifier
	userIdentifierStr := c.Param("id")
	userID, err := strconv.ParseInt(userIdentifierStr, 10, 64)
	if err != nil {
		c.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve user by its identifier
	u, err := models.GetUserByID(db, userID)
	if err != nil {
		c.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Start transaction
	transaction, err := db.Begin()
	if err != nil {
		c.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Delete user using the transaction
	if _, err := u.DeleteUserTransaction(transaction); err != nil {
		c.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Delete user from work orders using the transaction
	if _, err := u.DeleteFromWorkOrderTransaction(transaction); err != nil {
		c.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Commit changes
	if err := transaction.Commit(); err != nil {
		// Rollback changes
		transaction.Rollback()

		c.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Set success message
	session.Set("success", "El usuario ha sido eliminado con éxito")
	session.Save()

	c.Redirect(http.StatusSeeOther, "/admin/user/list")
}

// ViewUser views the given user
func ViewUser(c *gin.Context) {
	cfg, db, ok := retrieveContextValues(c)
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

	// Create session
	session := sessions.Default(c)

	// Check if user has role
	if !role.UserHasRole(cfg, user, role.ViewUser) {
		session.Set("error", "Tu cuenta no tiene permisos para visualizar un usuario")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/dashboard")
		return
	}

	// Retrieve user identifier
	userIdentifierStr := c.Param("id")
	userID, err := strconv.ParseInt(userIdentifierStr, 10, 64)
	if err != nil {
		c.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve user by its identifier
	u, err := models.GetUserByID(db, userID)
	if err != nil {
		c.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve user logs
	logs, err := models.UserLogs(db, u)
	if err != nil {
		c.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve user work orders
	workOrders, err := models.GetUserWorkOrders(db, u)
	if err != nil {
		c.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// XSRF token
	token := nosurf.Token(c.Request)

	c.HTML(http.StatusOK, "view_user.html", map[string]interface{}{
		"user":            u,
		"currentUser":     user,
		"logs":            logs,
		"workOrders":      workOrders,
		"view_work_order": role.UserHasRole(cfg, u, role.ViewWorkOrder),
		"roleList":        c.MustGet("roleList"),
		"token":           token,
	})
}

// RemoveUserFromWorkOrder removes the gein user from the given work order
func RemoveUserFromWorkOrder(c *gin.Context) {
	cfg, db, ok := retrieveContextValues(c)
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

	// Create session
	session := sessions.Default(c)

	// Check if user has role
	if !role.UserHasRole(cfg, user, role.RemoveUser) {
		session.Set("error", "Tu cuenta no tiene permisos para retirar a un usuario de una obra")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/dashboard")
		return
	}

	// Retrieve user identifier
	userIdentifierStr := c.Param("id")
	userID, err := strconv.ParseInt(userIdentifierStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve user by its identifier
	u, err := models.GetUserByID(db, userID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve work order identifier
	workIdentifierStr := c.Param("work")
	workID, err := strconv.ParseInt(workIdentifierStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve work order by its identifier
	workOrder, err := models.GetWorkOrderByID(db, workID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Delete assigned user
	if _, err := workOrder.DeleteAssignedUser(db, u); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Set success message
	session.Set("success", "Usuario eliminado de la obra con éxito")
	session.Save()

	c.Redirect(http.StatusSeeOther, "/admin/work/view/"+workIdentifierStr)
}
