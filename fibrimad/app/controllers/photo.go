package controllers

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/nayimotis/raggaer/fibrimad/app/role"
	"github.com/gin-contrib/sessions"

	"github.com/nayimotis/raggaer/fibrimad/app/models"
	"github.com/gin-gonic/gin"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

// UploadBoxPhoto uploads a file for the given box
func UploadBoxPhoto(c *gin.Context) {
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

	// Retrieve photo identifier
	photoIdentifierStr := c.Param("id")
	photoID, err := strconv.ParseInt(photoIdentifierStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve photo
	photo, err := models.GetBoxPhotoByID(db, photoID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve box
	box, err := models.GetBoxWorkOrderByID(db, photo.BoxID)
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

	// Retrieve default session
	session := sessions.Default(c)

	// Check if user has access to upload files for the box
	if !role.UserHasRole(cfg, user, role.UploadFileBox) {

		// Check if user is assigned to box work order
		assigned, err := workOrder.IsUserAssigned(db, user)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if !assigned {
			session.Set("error", "Tu cuenta no tiene permisos para subir fotos a esta caja")
			session.Save()
			c.Redirect(http.StatusSeeOther, "/dashboard")
			return
		}
	}

	// Retrieve form file
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

	// Close form file
	defer file.Close()

	// Read whole file
	fileBuffer, err := ioutil.ReadAll(file)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Detect content type
	contentType := http.DetectContentType(fileBuffer)

	if contentType != "image/png" && contentType != "image/gif" && contentType != "image/jpeg" {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Create paths
	dirPath := filepath.Join("boxes", strconv.FormatInt(box.ID, 10))

	// Create needed work order directory
	if err := os.MkdirAll(dirPath, os.ModeDir); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Create photo file
	photoFile, err := os.OpenFile(filepath.Join(dirPath, fileHeader.Filename), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	defer photoFile.Close()

	// Write content
	if _, err := photoFile.Write(fileBuffer); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Update photo fields
	photo.Filename = fileHeader.Filename
	photo.Path = filepath.Join(dirPath, fileHeader.Filename)

	// Update photo
	if _, err := photo.UpdateBoxPhoto(db); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Add log message
	user.AddLogMessage(db, "El usuario "+user.Username+" ha subido una foto (Obra "+workOrder.Code+", Caja "+box.Code+", Foto "+photo.Name+")")

	// Set success message
	session.Set("success", "Foto subida con éxito")
	session.Save()

	c.Redirect(http.StatusSeeOther, "/admin/box/view/"+strconv.FormatInt(box.ID, 10))
}

// ViewBoxPhoto views a file for the given box
func ViewBoxPhoto(c *gin.Context) {
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

	// Retrieve photo identifier
	photoIdentifierStr := c.Param("id")
	photoID, err := strconv.ParseInt(photoIdentifierStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve photo
	photo, err := models.GetBoxPhotoByID(db, photoID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve box
	box, err := models.GetBoxWorkOrderByID(db, photo.BoxID)
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

	// Retrieve default session
	session := sessions.Default(c)

	// Check if user has access to view files for the box
	if !role.UserHasRole(cfg, user, role.ViewFileBox) {

		// Check if user is assigned to box work order
		assigned, err := workOrder.IsUserAssigned(db, user)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if !assigned {
			session.Set("error", "Tu cuenta no tiene permisos para visualizar fotos sobre esta caja")
			session.Save()
			c.Redirect(http.StatusSeeOther, "/dashboard")
			return
		}
	}

	// Check view mode
	viewMode := c.Param("mode")

	if viewMode == "download" {
		c.Header("Content-Disposition", "attachment; filename="+photo.Filename)
	}

	c.File(photo.Path)
}
