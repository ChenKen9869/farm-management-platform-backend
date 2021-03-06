package controller

import (
	"go-backend/geocontainer"
	"go-backend/dao"
	"go-backend/entity"
	"go-backend/monitor"
	"go-backend/server"
	"go-backend/service"
	"go-backend/vo"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func CreateFenceController(ctx *gin.Context) {
	// 读取参数
	position := ctx.PostForm("Position")
	deviceList := ctx.PostForm("DeviceList")
	duraString := ctx.PostForm("Duration")
	coordinate := ctx.PostForm("Coordinate")
	name := ctx.PostForm("Name")
	parentIdString := ctx.PostForm("ParentId")
	// 验证坐标系是否合法
	allow, ok := geocontainer.Coordinates[coordinate]
	if !ok {
		server.Response(ctx, http.StatusBadRequest, 400, nil, "coordinate does not exists")
		return
	}
	if !allow {
		server.Response(ctx, http.StatusBadRequest, 400, nil, "coordiante does not allow to use")
		return
	}
	duration, errA := strconv.Atoi(duraString)
	if errA != nil {
		panic("atoi error")
	}  
	parentId, errAtoi := strconv.Atoi(parentIdString)
	if errAtoi != nil {
		panic("atoi error")
	}
	company := dao.GetCompanyInfoByID(uint(parentId))
	userInfo, exists := ctx.Get("user")
	if !exists {
		server.Response(ctx, http.StatusInternalServerError, 500, nil, "user info does not exists in application context")
		return
	}
	user := userInfo.(entity.User)
	// 验证 token 中的 user 是否有操作权限
	if !service.AuthCompanyUser(user.ID, company.ID) {
		server.Response(ctx, http.StatusUnauthorized, 401, nil, "权限不足")
		return
	}
	// 验证 deviceList 权限
	if !service.AuthFenceDeviceList(company.ID, deviceList) {
		server.Response(ctx, http.StatusUnauthorized, 401, nil, "权限不足")
		return
	}
	fenceId := monitor.StartFenceJob(user.ID, position, deviceList, duration, uint(parentId), name, coordinate)
	server.ResponseSuccess(ctx, gin.H{"id":fenceId}, server.Success)
}

func StopFenceController(ctx *gin.Context) {
	fenceIdString := ctx.Query("FenceId")
	fenceId, errAtoi := strconv.Atoi(fenceIdString)
	if errAtoi != nil {
		server.Response(ctx, http.StatusInternalServerError, 500, nil, "atoi error")
		return
	}
	userInfo, exists := ctx.Get("user")
	if !exists {
		server.Response(ctx, http.StatusInternalServerError, 500, nil, "user info does not exists in application context")
		return
	}
	user := userInfo.(entity.User)
	// 权限验证
	if !service.AuthCompanyUser(user.ID, dao.GetFenceRecordById(uint(fenceId)).ParentId) {
		server.Response(ctx, http.StatusUnauthorized, 401, nil, "权限不足")
		return
	}
	monitor.StopFenceJob(uint(fenceId))
	server.ResponseSuccess(ctx, nil, server.Success)
}

func GetActiveFenceByCompanyIdController(ctx *gin.Context) {
	companyIdString := ctx.Query("CompanyId")
	companyId, errAtoi := strconv.Atoi(companyIdString)
	if errAtoi != nil {
		server.Response(ctx, http.StatusInternalServerError, 500, nil, "atoi error")
		return
	}
	userInfo, exists := ctx.Get("user")
	if !exists {
		server.Response(ctx, http.StatusInternalServerError, 500, nil, "user info does not exists in application context")
		return
	}
	user := userInfo.(entity.User)
	// 权限验证
	if !service.AuthCompanyUser(user.ID, uint(companyId)) {
		server.Response(ctx, http.StatusUnauthorized, 401, nil, "权限不足")
		return
	}
	activeList := service.GetActiveFenceListByCompanyService(uint(companyId))
	server.ResponseSuccess(ctx, gin.H{"activeList" : activeList}, server.Success)
}

func GetActiveFenceStat(ctx *gin.Context) {
	fenceIdString := ctx.Query("FenceId")
	fenceId, errAtoi := strconv.Atoi(fenceIdString)
	if errAtoi != nil {
		server.Response(ctx, http.StatusInternalServerError, 500, nil, "atoi error")
		return
	}
	userInfo, exists := ctx.Get("user")
	if !exists {
		server.Response(ctx, http.StatusInternalServerError, 500, nil, "user info does not exists in application context")
		return
	}
	user := userInfo.(entity.User)
	// 权限验证
	if !service.AuthCompanyUser(user.ID, dao.GetFenceRecordById(uint(fenceId)).ParentId) {
		server.Response(ctx, http.StatusUnauthorized, 401, nil, "权限不足")
		return
	}
	coordinate, position, deviceList, alarmTime, remain := service.GetFenceStatService(uint(fenceId))
	result := vo.FenceRunningStatus{
		Coordinate: coordinate,
		Position: position,
		DeviceList: deviceList,
		AlarmTime: alarmTime,
		Remain: remain,
	}
	server.ResponseSuccess(ctx, gin.H{"status" : result}, server.Success)
}
