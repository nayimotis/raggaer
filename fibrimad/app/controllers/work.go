package controllers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/justinas/nosurf"
	"github.com/nayimotis/raggaer/fibrimad/app/models"
	"github.com/nayimotis/raggaer/fibrimad/app/role"
)

// CreateWorkForm defines a create work order
type CreateWorkForm struct {
	Code        string
	Description string
	State       string
	StartDate   string
	EndDate     string
}

// WorkOrderList shows the work order list
func WorkOrderList(c *gin.Context) {
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

	// Check if user has role
	if !role.UserHasRole(cfg, user, role.WorkOrderList) {
		session.Set("error", "Tu cuenta no tiene permisos para ver la lista de obras")
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

	// Create XSRF token
	token := nosurf.Token(c.Request)

	// Load work orders
	workOrders, err := models.LoadWorkOrders(db)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.HTML(http.StatusOK, "list_work.html", map[string]interface{}{
		"list":     workOrders,
		"user":     user,
		"error":    flashError,
		"success":  flashSuccess,
		"token":    token,
		"roleList": c.MustGet("roleList"),
	})
}

// CreateWorkOrderProcessForm process the create order form
func CreateWorkOrderProcessForm(c *gin.Context) {
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

	// Check if user has role
	if !role.UserHasRole(cfg, user, role.CreateWorkOrder) {
		session.Set("error", "Tu cuenta no tiene permisos para crear una nueva obra")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/dashboard")
		return
	}

	// Get work order values
	code := c.PostForm("code")
	description := c.PostForm("description")
	users := c.PostFormArray("users[]")
	start := c.PostForm("start-date")
	end := c.PostForm("end-date")
	state := c.PostForm("state")

	// Parse start date
	startDate, err := time.Parse("2006-01-02", start)
	if err != nil {
		c.Error(err)
		session.Set("error", "Fecha de entrada inválida. Inténtalo de nuevo")
		session.Set("createWorkForm", CreateWorkForm{code, description, state, start, end})
		session.Save()
		c.Redirect(http.StatusSeeOther, "/admin/work/create")
		return
	}

	// Parse end date
	endDate, err := time.Parse("2006-01-02", end)
	if err != nil {
		c.Error(err)
		session.Set("error", "Fecha de salida inválida. Inténtalo de nuevo")
		session.Set("createWorkForm", CreateWorkForm{code, description, state, start, end})
		session.Save()
		c.Redirect(http.StatusSeeOther, "/admin/work/create")
		return
	}

	// Create work order
	workOrder := models.WorkOrder{
		Code:        code,
		Description: description,
		State:       state,
		StartDate:   startDate,
		EndDate:     endDate,
	}
	workOrderResult, err := workOrder.CreateWorkOrder(db)
	if err != nil {
		c.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve last inserted work order identifier
	workOrderLastID, err := workOrderResult.LastInsertId()
	if err != nil {
		c.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	workOrder.ID = workOrderLastID

	for _, u := range users {
		// Retrieve user identififer from string
		userID, err := strconv.ParseInt(u, 10, 64)
		if err != nil {
			c.Error(err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// Load user by id
		currentUser, err := models.GetUserByID(db, userID)
		if err != nil {
			c.Error(err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// Add user to work order
		if _, err := workOrder.AddUserToWorkOrder(db, currentUser); err != nil {
			c.Error(err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	// Set success message
	session.Set("success", "Obra creada con éxito")
	session.Save()

	c.Redirect(http.StatusSeeOther, "/admin/work/list")
}

// CreateWorkOrderForm shows the create work order form
func CreateWorkOrderForm(c *gin.Context) {
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
	if !role.UserHasRole(cfg, user, role.CreateWorkOrder) {
		session.Set("error", "Tu cuenta no tiene permisos para crear una nueva obra")
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
	createForm := session.Get("createWorkForm")
	session.Delete("createWorkForm")
	session.Save()

	// Load user list for the work order
	users, err := models.UserList(db)
	if err != nil {
		c.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.HTML(http.StatusOK, "create_work.html", map[string]interface{}{
		"user":     user,
		"error":    flashError,
		"users":    users,
		"token":    token,
		"form":     createForm,
		"roleList": c.MustGet("roleList"),
	})
}

// DeleteWorkOrder deletes the given work order
func DeleteWorkOrder(c *gin.Context) {
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
	if !role.UserHasRole(cfg, user, role.DeleteWorkOrder) {
		session.Set("error", "Tu cuenta no tiene permisos para eliminar una obra")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/dashboard")
		return
	}

	// Retrieve work order identifier
	workOrderIdentifierStr := c.Param("id")
	workID, err := strconv.ParseInt(workOrderIdentifierStr, 10, 64)
	if err != nil {
		c.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve order by its identifier
	order, err := models.GetWorkOrderByID(db, workID)
	if err != nil {
		c.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Begin MySQL transaction
	transaction, err := db.Begin()
	if err != nil {
		c.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Delete all assigned users
	if _, err := order.DeleteAssignedUsersTransaction(transaction); err != nil {
		// Rollback transaction
		transaction.Rollback()

		c.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Delete work order
	if _, err := order.DeleteWorkOrderTransaction(transaction); err != nil {
		// Rollback transaction
		transaction.Rollback()

		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Delete work order files
	if _, err := order.DeleteWorkOrderFilesTransaction(transaction); err != nil {
		// Rollback transaction
		transaction.Rollback()

		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Delete work order files directory
	os.RemoveAll(filepath.Join("documents", strconv.FormatInt(order.ID, 10)))

	// Commit transaction changes
	if err := transaction.Commit(); err != nil {
		// Rollback transaction
		transaction.Rollback()

		c.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Set success message
	session.Set("success", "La obra ha sido eliminada con éxito")
	session.Save()

	c.Redirect(http.StatusSeeOther, "/admin/work/list")
}

// EditWorkOrderForm shows the edit work order form
func EditWorkOrderForm(c *gin.Context) {
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
	if !role.UserHasRole(cfg, user, role.EditWorkOrder) {
		session.Set("error", "Tu cuenta no tiene permisos para editar una obra")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/dashboard")
		return
	}

	// Retrieve work order identifier
	workOrderIdentifierStr := c.Param("id")
	workID, err := strconv.ParseInt(workOrderIdentifierStr, 10, 64)
	if err != nil {
		c.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve all users
	users, err := models.UserList(db)
	if err != nil {
		c.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve order by its identifier
	order, err := models.GetWorkOrderByID(db, workID)
	if err != nil {
		c.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve order users
	workUsers, err := order.GetUsers(db)
	if err != nil {
		c.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Set assign field of users
	for i, u := range users {
		for _, assigned := range workUsers {
			if u.ID == assigned.ID {
				users[i].IsAssigned = true
			}
		}
	}

	// Create XSRF token
	token := nosurf.Token(c.Request)

	// Get flash errors
	flashError := session.Get("error")
	session.Delete("error")

	// Get create work form
	createForm := session.Get("editWorkForm")
	session.Delete("editWorkForm")
	session.Save()

	c.HTML(http.StatusOK, "edit_work.html", map[string]interface{}{
		"user":     user,
		"users":    users,
		"order":    order,
		"error":    flashError,
		"token":    token,
		"form":     createForm,
		"roleList": c.MustGet("roleList"),
	})
}

// EditWorkOrderProcessForm process the work order edit form
func EditWorkOrderProcessForm(c *gin.Context) {
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
	if !role.UserHasRole(cfg, user, role.EditWorkOrder) {
		session.Set("error", "Tu cuenta no tiene permisos para editar una obra")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/dashboard")
		return
	}

	// Retrieve work order identifier
	workOrderIdentifierStr := c.Param("id")
	workID, err := strconv.ParseInt(workOrderIdentifierStr, 10, 64)
	if err != nil {
		c.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve order by its identifier
	order, err := models.GetWorkOrderByID(db, workID)
	if err != nil {
		c.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Get work order values
	code := c.PostForm("code")
	description := c.PostForm("description")
	users := c.PostFormArray("users[]")
	start := c.PostForm("start-date")
	end := c.PostForm("end-date")
	state := c.PostForm("state")

	// Parse start date
	startDate, err := time.Parse("2006-01-02", start)
	if err != nil {
		c.Error(err)
		session.Set("error", "Fecha de entrada inválida. Inténtalo de nuevo")
		session.Set("editWorkForm", CreateWorkForm{code, description, state, start, end})
		session.Save()
		c.Redirect(http.StatusSeeOther, "/admin/work/edit/"+strconv.FormatInt(order.ID, 10))
		return
	}

	// Parse end date
	endDate, err := time.Parse("2006-01-02", end)
	if err != nil {
		c.Error(err)
		session.Set("error", "Fecha de salida inválida. Inténtalo de nuevo")
		session.Set("editWorkForm", CreateWorkForm{code, description, state, start, end})
		session.Save()
		c.Redirect(http.StatusSeeOther, "/admin/work/edit/"+strconv.FormatInt(order.ID, 10))
		return
	}

	workOrder := &models.WorkOrder{
		ID:          order.ID,
		StartDate:   startDate,
		EndDate:     endDate,
		State:       state,
		Code:        code,
		Description: description,
	}

	// Edit work order
	if _, err := workOrder.EditWorkOrder(db); err != nil {
		c.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	for _, u := range users {
		// Retrieve user identififer from string
		userID, err := strconv.ParseInt(u, 10, 64)
		if err != nil {
			c.Error(err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// Load user by id
		currentUser, err := models.GetUserByID(db, userID)
		if err != nil {
			c.Error(err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// Check if user is already assigned to work order
		exists := false
		row := db.QueryRow(
			"SELECT EXISTS (SELECT 1 FROM work_order_users WHERE user_id = ? AND work_order_id = ?)",
			currentUser.ID,
			workOrder.ID,
		)
		if err := row.Scan(&exists); err != nil {
			c.Error(err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if exists {
			continue
		}

		// Add user to work order
		if _, err := workOrder.AddUserToWorkOrder(db, currentUser); err != nil {
			c.Error(err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	session.Set("success", "Obra editada con éxito")
	session.Save()

	c.Redirect(http.StatusSeeOther, "/admin/work/list")
}

// ViewWorkOrder views the given work order
func ViewWorkOrder(c *gin.Context) {
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

	// Get default session
	session := sessions.Default(c)

	// Get session success messages
	flashSuccess := session.Get("success")
	session.Delete("success")

	// Get session error messages
	flashError := session.Get("error")
	session.Delete("error")
	session.Save()

	// Retrieve work order identifier
	workOrderIdentifierStr := c.Param("id")
	workID, err := strconv.ParseInt(workOrderIdentifierStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve order by its identifier
	workOrder, err := models.GetWorkOrderByID(db, workID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve work order files
	workOrderFiles, err := workOrder.GetFiles(db)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve box list
	workOrderBoxes, err := workOrder.GetBoxes(db)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Create XSRF token
	token := nosurf.Token(c.Request)

	// Check if user has role
	if !role.UserHasRole(cfg, user, role.ViewWorkOrder) {

		// Check if user is not assigned
		isAssigned, err := workOrder.IsUserAssigned(db, user)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if !isAssigned {
			session.Set("error", "Tu cuenta no esta asignada a la obra "+workOrder.Code)
			session.Save()
			c.Redirect(http.StatusSeeOther, "/dashboard")
			return
		}

		c.HTML(http.StatusOK, "view_work.html", map[string]interface{}{
			"workOrder":   workOrder,
			"error":       flashError,
			"success":     flashSuccess,
			"files":       workOrderFiles,
			"fileDelete":  role.UserHasRole(cfg, user, role.DeleteFile),
			"userRemove":  role.UserHasRole(cfg, user, role.RemoveUser),
			"fileEdit":    role.UserHasRole(cfg, user, role.EditFile),
			"boxRemove":   role.UserHasRole(cfg, user, role.DeleteBox),
			"boxView":     role.UserHasRole(cfg, user, role.ViewBox),
			"boxEdit":     role.UserHasRole(cfg, user, role.EditBox),
			"token":       token,
			"roleList":    c.MustGet("roleList"),
			"osSeparator": string(os.PathSeparator),
			"boxList":     workOrderBoxes,
		})
		return
	}

	// Retrieve users assigned to work order
	workOrderUsers, err := workOrder.GetUsers(db)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.HTML(http.StatusOK, "view_work.html", map[string]interface{}{
		"workOrder":   workOrder,
		"users":       workOrderUsers,
		"files":       workOrderFiles,
		"error":       flashError,
		"success":     flashSuccess,
		"fileDelete":  role.UserHasRole(cfg, user, role.DeleteFile),
		"userRemove":  role.UserHasRole(cfg, user, role.RemoveUser),
		"fileEdit":    role.UserHasRole(cfg, user, role.EditFile),
		"boxRemove":   role.UserHasRole(cfg, user, role.DeleteBox),
		"boxView":     role.UserHasRole(cfg, user, role.ViewBox),
		"boxEdit":     role.UserHasRole(cfg, user, role.EditBox),
		"token":       token,
		"osSeparator": string(os.PathSeparator),
		"roleList":    c.MustGet("roleList"),
		"boxList":     workOrderBoxes,
	})
}

// UploadFileWorkOrder uploads a file for the given work order
func UploadFileWorkOrder(c *gin.Context) {
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
	if !role.UserHasRole(cfg, user, role.UploadFileWorkOrder) {
		session.Set("error", "Tu cuenta no tiene permisos para subir archivos a una obra")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/dashboard")
		return
	}

	// Retrieve work order identifier
	workOrderIdentifierStr := c.Param("id")
	workID, err := strconv.ParseInt(workOrderIdentifierStr, 10, 64)
	if err != nil {
		c.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve order by its identifier
	order, err := models.GetWorkOrderByID(db, workID)
	if err != nil {
		c.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Get form values
	fileName := c.PostForm("name")

	// Get fileheader from the form
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Open form file
	file, err := fileHeader.Open()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	defer file.Close()

	// Create paths
	dirPath := filepath.Join("documents", strconv.FormatInt(order.ID, 10))

	// Create needed work order directory
	if err := os.MkdirAll(dirPath, os.ModeDir); err != nil {
		c.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Create document file
	documentFile, err := os.OpenFile(filepath.Join(dirPath, fileHeader.Filename), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		c.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	defer documentFile.Close()

	// Copy contents of form file to document file
	if _, err := io.Copy(documentFile, file); err != nil {
		c.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Create document entry in database
	workOrderDocument := &models.WorkOrderFile{
		WorkOrderID: order.ID,
		Name:        fileName,
		Path:        filepath.Join(dirPath, fileHeader.Filename),
		Filename:    fileHeader.Filename,
	}

	// Insert document entry into database
	if _, err := workOrderDocument.CreateWorkOrderDocument(db); err != nil {
		c.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	session.Set("success", "Archivo subido con éxito")
	session.Save()

	c.Redirect(http.StatusSeeOther, "/admin/work/view/"+strconv.FormatInt(order.ID, 10))
}
