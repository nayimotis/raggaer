package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/justinas/nosurf"
	"github.com/nayimotis/raggaer/fibrimad/app/models"
	"github.com/nayimotis/raggaer/fibrimad/app/role"
)

// CreateBoxForm create box data form
type CreateBoxForm struct {
	Code   string
	Photos []string
}

// ProcessCreateBoxForm process the create box form
func ProcessCreateBoxForm(c *gin.Context) {
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
	if !role.UserHasRole(cfg, user, role.CreateBox) {
		session.Set("error", "Tu cuenta no tiene permisos para crear una nueva caja")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/dashboard")
		return
	}

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

	// Retrieve form values
	code := c.PostForm("code")
	photos := c.PostFormArray("photos[]")

	// Create work order box
	b := models.WorkOrderBox{
		Code:        code,
		WorkOrderID: workOrder.ID,
	}

	// Create box and retrieve inserted ID
	boxID, err := b.CreateBox(db)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve box identifier
	boxIdentifier, err := boxID.LastInsertId()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Set box ID
	b.ID = boxIdentifier

	for _, pic := range photos {
		// Skip empty pic
		if pic == "" {
			continue
		}

		// Create photo entry for the box
		if _, err := b.CreatePhoto(db, pic); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	// Set success message
	session.Set("success", "Caja creada con éxito")
	session.Save()

	c.Redirect(http.StatusSeeOther, "/admin/work/view/"+workOrderIdentifierStr)
}

// ShowCreateBoxForm shows the create box form
func ShowCreateBoxForm(c *gin.Context) {
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
	if !role.UserHasRole(cfg, user, role.CreateBox) {
		session.Set("error", "Tu cuenta no tiene permisos para crear una nueva caja")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/dashboard")
		return
	}

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

	// Create XSRF token
	token := nosurf.Token(c.Request)

	// Get flash errors
	flashError := session.Get("error")
	session.Delete("error")

	// Get create user form
	createForm := session.Get("createBoxForm")
	session.Delete("createBoxForm")
	session.Save()

	c.HTML(http.StatusOK, "create_box.html", map[string]interface{}{
		"errors":    flashError,
		"token":     token,
		"workOrder": workOrder,
		"form":      createForm,
		"roleList":  c.MustGet("roleList"),
	})
}

// ProcessEditBoxForm process the edit box form
func ProcessEditBoxForm(c *gin.Context) {
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
	if !role.UserHasRole(cfg, user, role.EditBox) {
		session.Set("error", "Tu cuenta no tiene permisos para editar una caja")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/dashboard")
		return
	}

	// Retrieve boxr identifier
	boxIdentifierStr := c.Param("id")
	boxID, err := strconv.ParseInt(boxIdentifierStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve box by its identifier
	box, err := models.GetBoxWorkOrderByID(db, boxID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve form values
	code := c.PostForm("code")

	// Start MySQL transaction
	tx, err := db.Begin()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Edit box code
	if _, err := box.EditBoxTransaction(tx, code); err != nil {
		tx.Rollback()
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve form values
	photos := c.PostFormArray("photos[]")
	currentPhotosID := c.PostFormArray("photos_current_id[]")
	currentPhotosName := c.PostFormArray("photos_current[]")

	for _, pic := range photos {
		// Skip empty pic
		if pic == "" {
			continue
		}

		// Create photo entry for the box
		if _, err := box.CreatePhotoTransaction(tx, pic); err != nil {
			tx.Rollback()
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	// Retrieve box photos
	currentPhotos, err := box.GetPhotos(db)
	if err != nil {
		tx.Rollback()
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	for _, pic := range currentPhotos {
		found := false
		foundIndex := 0

		// Range current photos and find deleted ones
		for index, i := range currentPhotosID {
			id, err := strconv.ParseInt(i, 10, 64)
			if err != nil {
				tx.Rollback()
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			if id == pic.ID {
				found = true
				foundIndex = index
				break
			}
		}

		// Delete photo if needed
		if !found {
			if _, err := pic.DeleteBoxPhotoTransaction(tx); err != nil {
				tx.Rollback()
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			continue
		}

		// Retrieve edited photo name
		currentName := currentPhotosName[foundIndex]

		// Edit photo if not deleting
		if _, err := pic.EditBoxPhotoTransaction(tx, currentName); err != nil {
			tx.Rollback()
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	// Commit transaction changes
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Set success message
	session.Set("success", "Caja editada con éxito")
	session.Save()

	c.Redirect(http.StatusSeeOther, "/admin/work/view/"+strconv.FormatInt(box.WorkOrderID, 10))
}

// ShowEditBoxForm shows the edit box form
func ShowEditBoxForm(c *gin.Context) {
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
	if !role.UserHasRole(cfg, user, role.CreateBox) {
		session.Set("error", "Tu cuenta no tiene permisos para crear una nueva caja")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/dashboard")
		return
	}

	// Retrieve box identifier
	boxIdentifierStr := c.Param("id")
	boxID, err := strconv.ParseInt(boxIdentifierStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve box by its identifier
	box, err := models.GetBoxWorkOrderByID(db, boxID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve box photos
	photos, err := box.GetPhotos(db)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Create XSRF token
	token := nosurf.Token(c.Request)

	// Get flash errors
	flashError := session.Get("error")
	session.Delete("error")
	session.Save()

	c.HTML(http.StatusOK, "edit_box.html", map[string]interface{}{
		"errors":   flashError,
		"token":    token,
		"box":      box,
		"photos":   photos,
		"roleList": c.MustGet("roleList"),
	})
}

// ViewWorkOrderBox views the given work order box
func ViewWorkOrderBox(c *gin.Context) {
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

	// Retrieve box identifier
	boxIdentifierStr := c.Param("id")
	boxID, err := strconv.ParseInt(boxIdentifierStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve box
	box, err := models.GetBoxWorkOrderByID(db, boxID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve box work order
	workOrder, err := models.GetWorkOrderByID(db, box.WorkOrderID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve box photos
	boxPhotos, err := box.GetPhotos(db)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve default session
	session := sessions.Default(c)

	// Check if user has role
	if !role.UserHasRole(cfg, user, role.ViewBox) {

		// Check if user is assigned to box work order
		assigned, err := workOrder.IsUserAssigned(db, user)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if !assigned {
			session.Set("error", "Tu cuenta no tiene permisos para visualizar esta caja")
			session.Save()
			c.Redirect(http.StatusSeeOther, "/dashboard")
			return
		}
	}

	// Retrieve flash values
	flashError := session.Get("error")
	session.Delete("error")

	flashSuccess := session.Get("success")
	session.Delete("success")
	session.Save()

	// Generate token XSRF
	token := nosurf.Token(c.Request)

	c.HTML(http.StatusOK, "view_box.html", map[string]interface{}{
		"box":        box,
		"workOrder":  workOrder,
		"user":       user,
		"create_box": role.UserHasRole(cfg, user, role.CreateBox),
		"roleList":   c.MustGet("roleList"),
		"error":      flashError,
		"success":    flashSuccess,
		"photos":     boxPhotos,
		"token":      token,
	})
}

// DeleteBox deletes the given work order box
func DeleteBox(c *gin.Context) {
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

	// Retrieve box identifier
	boxIdentifierStr := c.Param("id")
	boxID, err := strconv.ParseInt(boxIdentifierStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve default session
	session := sessions.Default(c)

	// Retrieve box
	box, err := models.GetBoxWorkOrderByID(db, boxID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Check if user has role
	if !role.UserHasRole(cfg, user, role.DeleteBox) {
		session.Set("error", "Tu cuenta no tiene permisos para eliminar una caja")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/dashboard")
		return
	}

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Delete box with transaction
	if _, err := box.DeleteBoxTransaction(tx); err != nil {
		tx.Rollback()
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Delete all box photos
	if _, err := box.DeleteBoxPhotosTransaction(tx); err != nil {
		tx.Rollback()
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Remove box directory
	os.RemoveAll(filepath.Join("boxes", strconv.FormatInt(box.ID, 10)))

	// Commit transaction changes
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Set success message
	session.Set("success", "Caja eliminada con éxito")
	session.Save()

	c.Redirect(http.StatusSeeOther, "/admin/work/view/"+strconv.FormatInt(box.WorkOrderID, 10))
}
