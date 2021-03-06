package dao

import (
	"go-backend/common"
	"go-backend/entity"
	"time"

	"github.com/jinzhu/gorm"
)

func CreateFenceRecord(fenceRecord entity.FenceRecord) uint {
	common.GetDB().Create(&fenceRecord)
	return fenceRecord.ID
}

func AddAlarmTime(fenceId uint) {
	common.GetDB().Model(&entity.FenceRecord{}).Where("id = ? ", fenceId).Update("alarm_time", gorm.Expr("alarm_time + ?", 1))
}

func ModifyFenceStat(fenceId uint, stat int) {
	common.GetDB().Model(&entity.FenceRecord{}).Where("id = ?", fenceId).Update("stat", stat)
}

func UpdateFenceEndTime(fenceId uint, endTime time.Time) {
	common.GetDB().Model(&entity.FenceRecord{}).Where("id = ?", fenceId).Update("end_time", endTime)
}


func GetFenceRecordById(fenceId uint) entity.FenceRecord {
	var fenceRecord entity.FenceRecord
	common.GetDB().Where("id = ?", fenceId).First(&fenceRecord)
	return fenceRecord
}

// 查看用户所有的活跃围栏
func GetActiveFenceListByUserId(userId uint) []entity.FenceRecord {
	var fenceRecordList []entity.FenceRecord
	common.GetDB().Where("owner = ? and stat = ?", userId, entity.FenceRunning).Find(&fenceRecordList)
	return fenceRecordList
}

// 查看农牧场所有的活跃围栏
func GetActiveFenceListByCompanyId(companyId uint) []entity.FenceRecord {
	var fenceRecordList []entity.FenceRecord
	common.GetDB().Where("parent_id = ? and stat = ?", companyId, entity.FenceRunning).Find(&fenceRecordList)
	return fenceRecordList
}

// 查看农牧场围栏任务历史记录
func GetFenceRecordListByCompanyId(companyId uint) []entity.FenceRecord {
	var fenceRecordList []entity.FenceRecord
	common.GetDB().Where("parent_id = ? and stat = ?", companyId, entity.FenceFinished).Find(&fenceRecordList)
	return fenceRecordList
}