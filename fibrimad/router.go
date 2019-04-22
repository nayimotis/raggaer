package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nayimotis/raggaer/fibrimad/app/controllers"
)

func setRoutes(router *gin.Engine) {
	router.GET("/", controllers.Index)
	router.GET("/file/view/:id/:mode", controllers.ViewFile)
	router.POST("/admin/file/delete/:id", controllers.DeleteFile)
	router.POST("/admin/file/edit/:id", controllers.EditFile)
	router.GET("/login", controllers.LoginForm)
	router.POST("/login", controllers.LoginProcess)
	router.GET("/dashboard", controllers.Dashboard)
	router.GET("/admin/user/list", controllers.ShowUserList)
	router.GET("/admin/user/create", controllers.CreateUserShowForm)
	router.POST("/admin/user/create", controllers.CreateUserProcessForm)
	router.GET("/admin/user/edit/:id", controllers.EditUserForm)
	router.POST("/admin/user/edit/:id", controllers.EditUserProcessForm)
	router.GET("/admin/user/view/:id", controllers.ViewUser)
	router.POST("/admin/user/delete/:id", controllers.DeleteUser)
	router.POST("/admin/user/remove/:id/:work", controllers.RemoveUserFromWorkOrder)
	router.POST("/logout", controllers.Logout)
	router.GET("/admin/work/create", controllers.CreateWorkOrderForm)
	router.POST("/admin/work/create", controllers.CreateWorkOrderProcessForm)
	router.POST("/admin/work/delete/:id", controllers.DeleteWorkOrder)
	router.GET("/admin/work/list", controllers.WorkOrderList)
	router.GET("/admin/work/edit/:id", controllers.EditWorkOrderForm)
	router.POST("/admin/work/edit/:id", controllers.EditWorkOrderProcessForm)
	router.POST("/admin/work/upload/:id", controllers.UploadFileWorkOrder)
	router.GET("/admin/work/view/:id", controllers.ViewWorkOrder)
	router.GET("/admin/box/edit/:id", controllers.ShowEditBoxForm)
	router.POST("/admin/box/edit/:id", controllers.ProcessEditBoxForm)
	router.GET("/admin/box/create/:id", controllers.ShowCreateBoxForm)
	router.POST("/admin/box/create/:id", controllers.ProcessCreateBoxForm)
	router.GET("/admin/box/view/:id", controllers.ViewWorkOrderBox)
	router.POST("/admin/box/delete/:id", controllers.DeleteBox)
	router.POST("/admin/photo/upload/:id", controllers.UploadBoxPhoto)
	router.GET("/admin/photo/view/:id/:mode", controllers.ViewBoxPhoto)
}
