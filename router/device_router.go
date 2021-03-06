package router

import (
	"go-backend/controller"
	"go-backend/middleware"

	"github.com/gin-gonic/gin"
)

func DeviceRouter(r *gin.Engine) *gin.Engine {
	device := r.Group("/device")
	device.Use(middleware.AuthMiddleware())
	
	fixedDevice := device.Group("/fixed")
	fixedDevice.POST("/create", controller.CreateFixedDeviceController)
	fixedDevice.DELETE("/delete", controller.DeleteFixedDeviceController)
	fixedDevice.POST("/create_type", controller.CreateFixedDeviceTypeController)
	fixedDevice.DELETE("/delete_type", controller.DeleteFixedDeviceTypeController)

	fixedDevice.GET("/get_monitor", controller.GetMonitorStreamController)
	fixedDevice.GET("/get_fio_latest", controller.GetLatestFioController)
	fixedDevice.GET("/get_by_farmhouse", controller.GetFixedDeviceListByFarmhouseController)
	fixedDevice.GET("/get_fio_list_by_time", controller.GetFioListByTime)

	portableDevice := device.Group("/portable")
	portableDevice.POST("/create", controller.CreatePortableDeviceController)
	portableDevice.DELETE("/delete", controller.DeletePortableDeviceController)
	portableDevice.POST("/create_type", controller.CreatePortableDeviceTypeController)
	portableDevice.DELETE("/delete_type", controller.DeletePortableDeviceTypeController)
	
	portableDevice.GET("/get_by_farmhouse", controller.GetPortableDeviceListByFarmhouseController)

	return r
}