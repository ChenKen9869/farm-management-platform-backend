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

func CreateCompanyController(ctx *gin.Context) {
	name := ctx.PostForm("Name")
	parentId := ctx.PostForm("ParentId")
	location := ctx.PostForm("Location")
	leaderInfo, exists := ctx.Get("user")
	if !exists {
		panic("error: user information does not exists in application context")
	}
	leaderId := leaderInfo.(entity.User).ID
	parent, err := strconv.Atoi(parentId)
	if err != nil {
		server.Response(ctx, http.StatusInternalServerError, 500, nil, "atoi error")
		return
	}
	id, errService := service.CreateCompanyService(uint(parent), name, leaderId, location)
	if errService != nil {
		server.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "param error")
		return
	}
	// 如果是新公司，则需要将新企业与用户名关系写入到关系表
	if parent == 0 {
		companyUser := entity.CompanyUser{
			CompanyID: id,
			UserID: leaderId,
		}
		dao.CreateCompanyUser(companyUser)
	}
	server.ResponseSuccess(ctx, gin.H{"CompanyId": id}, server.Success)
}


func GetCompanyTreeListController(ctx *gin.Context) {
	userId := ctx.Query("UserId")
	userInfo, exists := ctx.Get("user")
	if !exists {
		server.Response(ctx, http.StatusInternalServerError, 500, nil, "user infromation does not exists in application context")
		return
	}
	user := userInfo.(entity.User)
	uid, err := strconv.Atoi(userId)
	if err != nil {
		server.Response(ctx, http.StatusInternalServerError, 500, nil, "atoi error")
		return
	}
	if user.ID != uint(uid) {
		server.Response(ctx, http.StatusUnauthorized, 401, nil, "权限不足")
		return
	}
	treeList := service.GetCompanyTreeListService(uint(uid))
	companyTreeList := gin.H{
		"mechanism": treeList,
	}
	server.ResponseSuccess(ctx, companyTreeList, server.Success)
}

// 删除应该同时删除权限表中的企业
func DeleteCompanyController(ctx *gin.Context) {
	companyIdString := ctx.Query("CompanyId")
	companyId, errAtoi := strconv.Atoi(companyIdString)
	if errAtoi != nil {
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
		server.Response(ctx, http.StatusUnauthorized, 401, nil, "权限不足啦")
		return
	}
	err := service.DeleteCompanyService(uint(companyId))
	if err != nil {
		if msg := err.Error(); msg == server.CompanyNotExist {
			server.Response(ctx, http.StatusBadRequest, 400, nil, msg)
			return
		} else if msg == server.NodeHasSubcompany {
			server.Response(ctx, http.StatusBadRequest, 400, nil, msg)
			return
		}
	}
	if dao.GetCompanyInfoByID(uint(companyId)).ParentID == 0 {
		service.DeleteCompanyUserService(uint(companyId), user.ID)
	}
	server.ResponseSuccess(ctx, nil, server.Success)
}

func CreateCompanyUserController(ctx *gin.Context) {
	companyIdString := ctx.PostForm("CompanyId")
	userIdString := ctx.PostForm("UserId")
	companyId, errAtoiComanyId := strconv.Atoi(companyIdString)
	userId, errAtoiUserId := strconv.Atoi(userIdString)
	if errAtoiComanyId != nil || errAtoiUserId != nil {
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
	err := service.CreateCompanyUserService(uint(companyId), uint(userId))
	if err != nil {
		if msg := err.Error(); msg == server.CompanyNotExist {
			server.Response(ctx, http.StatusBadRequest, 400, nil, msg)
			return
		} else if msg == server.UserNotExist {
			server.Response(ctx, http.StatusBadRequest, 400, nil, msg)
			return
		} else {
			server.Response(ctx, http.StatusBadRequest, 400, nil, msg)
			return
		}
	}
	server.ResponseSuccess(ctx, 
		gin.H{"companyId" : companyId, "userId" : userId,}, 
		server.Success)
}

func DeleteCompanyUserController(ctx *gin.Context) {
	companyIdString := ctx.Query("CompanyId")
	userIdString := ctx.Query("UserId")
	companyId, errAtoiComanyId := strconv.Atoi(companyIdString)
	userId, errAtoiUserId := strconv.Atoi(userIdString)
	if errAtoiComanyId != nil || errAtoiUserId != nil {
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
	service.DeleteCompanyUserService(uint(companyId), uint(userId))
	server.ResponseSuccess(ctx, 
		gin.H{"companyId" : companyId, "userId" : userId,}, 
		server.Success)
}

func GetEmployeeListController(ctx *gin.Context) {
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
	employeeList := service.GetEmployeeListService(uint(companyId))
	
	result := []gin.H{}
	for employee := range employeeList {
		result = append(result, gin.H{
			"id" : employee.ID,
			"name" : employee.Name,
			"authCompany" : employeeList[employee],
		})
	}
	
	server.ResponseSuccess(ctx, gin.H{"employeeList":result}, server.Success)
}

func GetCompanyInfoController(ctx *gin.Context) {
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
	result := service.GetCompanyInfoService(uint(companyId))
	server.ResponseSuccess(ctx, gin.H{"companyInfo": result}, server.Success)
}