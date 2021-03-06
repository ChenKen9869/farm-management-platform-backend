package controller

import (
	"go-backend/dao"
	"go-backend/entity"
	"go-backend/server"
	"go-backend/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateFixedDeviceController(ctx *gin.Context) {
	// company的leader就是这个设备的owner
	userInfo, exists:= ctx.Get("user")
	user := userInfo.(entity.User)
	if !exists {
		panic("error: user information does not exists in application context")
	}
	companyId, err := strconv.Atoi(ctx.PostForm("CompanyId"))
	if err != nil {
		server.Response(ctx, http.StatusInternalServerError, 500, nil, "server inter failed")
		return
	}	
	// 验证company与user权限
	if !service.AuthCompanyUser(user.ID, uint(companyId)) {
		server.Response(ctx, http.StatusUnauthorized, 401, nil, "权限不足")
		return
	}

	company := dao.GetCompanyInfoByID(uint(companyId))

	owner := company.Owner
	deviceId := ctx.PostForm("DeviceId")
	typeId := ctx.PostForm("TypeId")
	// 判断 type 是否在服务器中注册过
	if !dao.ExistFixedDeviceType(typeId) {
		server.Response(ctx, http.StatusBadRequest, 400, nil, "不支持的固定设备类型")
		return
	}

	id := service.CreateFixedDeviceService(deviceId, uint(companyId), typeId, owner)

	server.ResponseSuccess(ctx, gin.H{"Id": id}, server.Success)
}

func DeleteFixedDeviceController(ctx *gin.Context) {
	fixedDeviceIdString := ctx.Query("Id")
	fixedDeviceId, _ := strconv.Atoi(fixedDeviceIdString)
	userInfo, exists:= ctx.Get("user")
	user := userInfo.(entity.User)
	if !exists {
		panic("error: user information does not exists in application context")
	}
	companyId := dao.GetFixedDeviceInfoById(uint(fixedDeviceId)).FarmhouseID
	// 验证权限
	if !service.AuthCompanyUser(user.ID, uint(companyId)) {
		server.Response(ctx, http.StatusUnauthorized, 401, nil, "权限不足")
		return
	}
	
	service.DeleteFixedDeviceService(uint(fixedDeviceId))

	server.ResponseSuccess(ctx, nil, server.Success)
}

func CreatePortableDeviceController(ctx *gin.Context) {
	// id ，biologyId ， typeid，由前端传入, owner = biology.owner
	biologyIdString := ctx.PostForm("BiologyId")
	portableDeviceId := ctx.PostForm("DeviceId")
	typeId := ctx.PostForm("TypeId")
	biologyId, err := strconv.Atoi(biologyIdString)
	if err != nil {
		server.Response(ctx, http.StatusInternalServerError, 500, nil, "server inter failed")
		return
	}
	// 验证用户是否有这个生物的操作权限
	userInfo, exists := ctx.Get("user") 
	if !exists {
		server.Response(ctx, http.StatusInternalServerError, 500, nil, "user information does not exists in application context")
		return
	}
	user := userInfo.(entity.User)
	companyId := dao.GetBiologyInfoById(uint(biologyId)).FarmhouseID
	if !service.AuthCompanyUser(user.ID, uint(companyId)) {
		server.Response(ctx, http.StatusUnauthorized, 401, nil, "权限不足")
		return
	}

	// 判断 type 是否在服务器中注册过
	if !dao.ExistPortableDeviceType(typeId) {
		server.Response(ctx, http.StatusBadRequest, 400, nil, "不支持的携带设备类型")
		return
	}
	id := service.CreatePortableDeviceService(portableDeviceId, uint(biologyId), typeId)

	server.ResponseSuccess(ctx, gin.H{"Id": id}, server.Success)
}

func DeletePortableDeviceController(ctx *gin.Context) {
	portableDeviceIdString := ctx.Query("Id")
	portableDeviceId, _ := strconv.Atoi(portableDeviceIdString) 
	biologyId := dao.GetPortableDeviceInfoById(uint(portableDeviceId)).BiologyID
	companyId := dao.GetBiologyInfoById(biologyId).FarmhouseID

	userInfo, exists := ctx.Get("user") 
	if !exists {
		server.Response(ctx, http.StatusInternalServerError, 500, nil, "user information does not exists in application context")
		return
	}
	user := userInfo.(entity.User)
	if !service.AuthCompanyUser(user.ID, uint(companyId)) {
		server.Response(ctx, http.StatusUnauthorized, 401, nil, "权限不足")
		return
	}
	service.DeletePortableDeviceService(uint(portableDeviceId))

	server.ResponseSuccess(ctx, nil, server.Success)
}

func CreateFixedDeviceTypeController(ctx *gin.Context) {
	fixedDeviceTypeId := ctx.PostForm("FixedDeviceTypeId")
	service.CreateFixedDeviceTypeService(fixedDeviceTypeId)

	server.ResponseSuccess(ctx, nil, server.Success)
}

func DeleteFixedDeviceTypeController(ctx *gin.Context) {
	fixedDeviceTypeId := ctx.Query("FixedDeviceTypeId")
	service.DeleteFixedDeviceTypeService(fixedDeviceTypeId)

	server.ResponseSuccess(ctx, nil, server.Success)
}

func CreatePortableDeviceTypeController(ctx *gin.Context) {
	portableDeviceTypeId := ctx.PostForm("PortableDeviceTypeId")
	service.CreatePortableDeviceTypeService(portableDeviceTypeId)

	server.ResponseSuccess(ctx, nil, server.Success)
}

func DeletePortableDeviceTypeController(ctx *gin.Context) {
	portableDeviceTypeId := ctx.Query("PortableDeviceTypeId")
	service.DeletePortableDeviceTypeService(portableDeviceTypeId)

	server.ResponseSuccess(ctx, nil, server.Success)
}

func GetMonitorStreamController(ctx *gin.Context) {
	deviceId := ctx.Query("Id")
	deviceIdInt, _ := strconv.Atoi(deviceId)
	if len(deviceId) < 1 {
		server.Response(ctx, http.StatusBadRequest, 400, nil, "monitorId is empty")
		return
	}
	// 权限验证
	userInfo, exists := ctx.Get("user") 
	if !exists {
		server.Response(ctx, http.StatusInternalServerError, 500, nil, "user information does not exists in application context")
		return
	}
	user := userInfo.(entity.User)
	companyId := dao.GetFixedDeviceInfoById(uint(deviceIdInt)).FarmhouseID
	if !service.AuthCompanyUser(user.ID, uint(companyId)) {
		server.Response(ctx, http.StatusUnauthorized, 401, nil, "权限不足")
		return
	}
	resultUrl, resultExpireTime, id:= service.GetMonitorStreamByDeviceIdService(uint(deviceIdInt))
	payload := gin.H{
		"id" : id,
		"url" : resultUrl,
		"expireTime" : resultExpireTime,
	}
	server.ResponseSuccess(ctx, payload, server.Success)
}

func GetLatestFioController(ctx *gin.Context) {
	fioIdString := ctx.Query("Id")
	fioId, _ := strconv.Atoi(fioIdString)
	// 权限验证
	userInfo, exists := ctx.Get("user") 
	if !exists {
		server.Response(ctx, http.StatusInternalServerError, 500, nil, "user information does not exists in application context")
		return
	}
	user := userInfo.(entity.User)
	companyId := dao.GetFixedDeviceInfoById(uint(fioId)).FarmhouseID
	if !service.AuthCompanyUser(user.ID, uint(companyId)) {
		server.Response(ctx, http.StatusUnauthorized, 401, nil, "权限不足")
		return
	}
	data := service.GetLatestFioService(uint(fioId))
	server.ResponseSuccess(ctx, gin.H{"latest": data}, server.Success)
}

func GetFioListByTime(ctx *gin.Context) {
	fioIdString := ctx.Query("Id")
	startTime := ctx.Query("StartTime")
	endTime := ctx.Query("EndTime")
	fioId, _ := strconv.Atoi(fioIdString)
	// 权限验证
	userInfo, exists := ctx.Get("user") 
	if !exists {
		server.Response(ctx, http.StatusInternalServerError, 500, nil, "user information does not exists in application context")
		return
	}
	user := userInfo.(entity.User)
	companyId := dao.GetFixedDeviceInfoById(uint(fioId)).FarmhouseID
	if !service.AuthCompanyUser(user.ID, uint(companyId)) {
		server.Response(ctx, http.StatusUnauthorized, 401, nil, "权限不足")
		return
	}
	result := service.GetFioRecordListByTimeService(uint(fioId), startTime, endTime)
	server.ResponseSuccess(ctx, gin.H{"recordList": result}, server.Success)
}

func GetFixedDeviceListByFarmhouseController(ctx *gin.Context) {
	companyIdString := ctx.Query("CompanyId")
	companyId, errAtoiComanyId := strconv.Atoi(companyIdString)
	if errAtoiComanyId != nil {
		server.Response(ctx, http.StatusInternalServerError, 500, nil, "服务器内部错误")
	}
	userInfo, exists := ctx.Get("user")
	if !exists {
		server.Response(ctx, http.StatusInternalServerError, 500, nil, "user infromation does not exists in application context")
		return
	}
	user := userInfo.(entity.User)
	// 权限验证
	if !service.AuthCompanyUser(user.ID, uint(companyId)) {
		server.Response(ctx, http.StatusUnauthorized, 401, nil, "权限不足")
		return
	}
	result := service.GetFixedDeviceListByFarmhouseService(uint(companyId))
	server.ResponseSuccess(ctx, gin.H{"fixedDeviceList": result}, server.Success)
}

func GetPortableDeviceListByFarmhouseController(ctx *gin.Context) {
	companyIdString := ctx.Query("CompanyId")
	companyId, errAtoiComanyId := strconv.Atoi(companyIdString)
	if errAtoiComanyId != nil {
		server.Response(ctx, http.StatusInternalServerError, 500, nil, "服务器内部错误")
	}
	userInfo, exists := ctx.Get("user")
	if !exists {
		server.Response(ctx, http.StatusInternalServerError, 500, nil, "user infromation does not exists in application context")
		return
	}
	user := userInfo.(entity.User)
	// 权限验证
	if !service.AuthCompanyUser(user.ID, uint(companyId)) {
		server.Response(ctx, http.StatusUnauthorized, 401, nil, "权限不足")
		return
	}
	result := service.GetPortableDeviceListByFarmhouseService(uint(companyId))
	server.ResponseSuccess(ctx, gin.H{"portableDeviceList": result}, server.Success)
}

