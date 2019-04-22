package controllers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/raggaer/fibrimad/app/role"

	"github.com/gin-gonic/gin"
	"github.com/raggaer/fibrimad/app/models"
)

// EditFile edits the given file
func EditFile(c *gin.Context) {
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

	// Retrieve session
	session := sessions.Default(c)

	// Check if user has "create_user" role
	if !role.UserHasRole(cfg, user, role.EditFile) {
		session.Set("error", "Tu cuenta no tiene permisos para editar un fichero")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/dashboard")
		return
	}

	// Retrieve file identifier
	fileIdentifierStr := c.Param("id")
	fileID, err := strconv.ParseInt(fileIdentifierStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve file by its identifier
	currentFile, err := models.GetFileByID(db, fileID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Update file name
	currentFile.Name = c.PostForm("name")

	// Retrieve work order by file identifier
	order, err := models.GetWorkOrderByID(db, currentFile.WorkOrderID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	updateFile := true

	// Get fileheader from the form
	fileHeader, err := c.FormFile("file")
	if err != nil {
		updateFile = false
	}

	if updateFile {
		if fileHeader.Filename == "" || fileHeader.Size == 0 {
			updateFile = false
		}
	}

	if updateFile {
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
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// Create document file
		documentFile, err := os.OpenFile(filepath.Join(dirPath, fileHeader.Filename), os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		defer documentFile.Close()

		// Copy contents of form file to document file
		if _, err := io.Copy(documentFile, file); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// Update file path
		currentFile.Path = filepath.Join(dirPath, fileHeader.Filename)

		// Update file name
		currentFile.Filename = fileHeader.Filename
	}

	// Update file name and path
	if _, err := currentFile.EditFile(db); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Set success message
	session.Set("success", "Fichero editado con éxito")
	session.Save()

	c.Redirect(http.StatusSeeOther, "/admin/work/view/"+strconv.FormatInt(currentFile.WorkOrderID, 10))
}

// DeleteFile deletes the given file
func DeleteFile(c *gin.Context) {
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

	// Retrieve session
	session := sessions.Default(c)

	// Check if user has "create_user" role
	if !role.UserHasRole(cfg, user, role.DeleteFile) {
		session.Set("error", "Tu cuenta no tiene permisos para eliminar un fichero")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/dashboard")
		return
	}

	// Retrieve file identifier
	fileIdentifierStr := c.Param("id")
	fileID, err := strconv.ParseInt(fileIdentifierStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve file by its identifier
	file, err := models.GetFileByID(db, fileID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Delete file from database
	if _, err := file.DeleteWorkOrderFile(db); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Delete file from file system
	if err := os.Remove(file.Path); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	session.Set("success", "Fichero eliminado con éxito")
	session.Save()

	c.Redirect(http.StatusSeeOther, "/admin/work/view/"+strconv.FormatInt(file.WorkOrderID, 10))
}

// ViewFile shows the given file
func ViewFile(c *gin.Context) {
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

	// Retrieve session
	session := sessions.Default(c)

	// Retrieve file identifier
	fileIdentifierStr := c.Param("id")
	fileID, err := strconv.ParseInt(fileIdentifierStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve file by its identifier
	file, err := models.GetFileByID(db, fileID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve work order
	workOrder, err := models.GetWorkOrderByID(db, file.WorkOrderID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Check if user has file view rights
	if role.UserHasRole(cfg, user, role.ViewFile) {
		// Check view mode
		viewMode := c.Param("mode")

		if viewMode == "download" {
			c.Header("Content-Disposition", "attachment; filename="+file.Filename)
		}

		c.File(file.Path)
		return
	}

	// Check if user is assigned to the file work order
	assigned, err := workOrder.IsUserAssigned(db, user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if !assigned {
		session.Set("error", "Tu cuenta no esta asignada a la obra "+workOrder.Code+" no tiene acceso al fichero")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/dashboard")
		return
	}

	// Check view mode
	viewMode := c.Param("mode")

	if viewMode == "download" {
		c.Header("Content-Disposition", "attachment; filename="+file.Filename)
	}

	c.File(file.Path)
}
