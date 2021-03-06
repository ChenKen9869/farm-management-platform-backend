package dao

import (
	"go-backend/common"
	"go-backend/entity"
)

// 创建固定设备
func CreateFixedDevice(fixedDevice entity.FixedDevice) uint {
	common.GetDB().Create(&fixedDevice)
	return fixedDevice.ID
}
// 创建携带设备
func CreatePortableDevice(portableDevice entity.PortableDevice) uint {
	common.GetDB().Create(&portableDevice)
	return portableDevice.ID
}

// 新增便携式设备类型
func CreateFixedDeviceType(fixedDeviceType entity.FixedDeviceType) {
	common.GetDB().Create(fixedDeviceType)
}

func CreatePortableDeviceType(portableDeviceType entity.PortableDeviceType) {
	common.GetDB().Create(&portableDeviceType)
}
// 删除便携式设备与类型
func DeletePortableDevice(portableDeviceId uint) entity.PortableDevice {
	var portableDevice entity.PortableDevice
	common.GetDB().Where("id = ?", portableDeviceId).First(&portableDevice)
	common.GetDB().Delete(&portableDevice)
	return portableDevice
}

func DeletePortableDeviceType(portableDeviceTypeId string) entity.PortableDeviceType {
	var portableDeviceType entity.PortableDeviceType
	common.GetDB().Where("portable_device_type_id = ?", portableDeviceTypeId).First(&portableDeviceType)
	common.GetDB().Delete(&portableDeviceType)
	return portableDeviceType
}

// 删除固定式设备与类型
func DeleteFixedDevice(fixedDeviceId uint) entity.FixedDevice {
	var fixedDevice entity.FixedDevice
	common.GetDB().Where("id = ?", fixedDeviceId).First(&fixedDevice)
	common.GetDB().Delete(&fixedDevice)
	return fixedDevice
}

func DeleteFixedDeviceType(fixedDeviceTypeId string) entity.FixedDeviceType {
	var fixedDeviceType entity.FixedDeviceType
	common.GetDB().Where("fixed_device_type_id = ?", fixedDeviceTypeId).First(&fixedDeviceType)
	common.GetDB().Delete(&fixedDeviceType)
	return fixedDeviceType
}

// 修改携带设备（修改绑定的生物 / 属主公司等）

// 修改固定设备（修改绑定的农舍 / 属主等）


func ExistFixedDeviceType(fixedDeviceTypeId string) bool {
	var fixedType entity.FixedDeviceType
	common.GetDB().Table("fixed_device_types").Where("fixed_device_type_id = ?", fixedDeviceTypeId).First(&fixedType)
	return len(fixedType.FixedDeviceTypeID) != 0
}

func ExistPortableDeviceType(portableDeviceTypeId string) bool {
	var portableType entity.PortableDeviceType
	common.GetDB().Table("portable_device_types").Where("portable_device_type_id = ?", portableDeviceTypeId).First(&portableType)
	return len(portableType.PortableDeviceTypeID) != 0
}

func GetPortableDeviceInfoByBiology(biologyId uint) entity.PortableDevice {
	var portableDevice entity.PortableDevice
	common.GetDB().Where("biology_id = ?", biologyId).First(&portableDevice)
	return portableDevice
}

func GetFixedDeviceListByFarmhouse(farmhouseId uint) []entity.FixedDevice {
	var fixedDeviceList []entity.FixedDevice
	common.GetDB().Where("farmhouse_id = ?", farmhouseId).Find(&fixedDeviceList)
	return fixedDeviceList
}

// 通过 id 查看设备
func GetFixedDeviceInfoById(fixedDeviceId uint) entity.FixedDevice {
	var fixedDevice entity.FixedDevice
	common.GetDB().Where("id = ?", fixedDeviceId).First(&fixedDevice)
	return fixedDevice
}

func GetPortableDeviceInfoById(portableDeviceId uint) entity.PortableDevice {
	var portableDevice entity.PortableDevice
	common.GetDB().Where("id = ?", portableDeviceId).First(&portableDevice)
	return portableDevice
}

