package main

import (
	"context"
	"myitools/device"
)

// App 应用结构体
type App struct {
	ctx context.Context
}

// NewApp 创建一个新的App应用实例
func NewApp() *App {
	return &App{}
}

// startup 在应用启动时调用
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	// 使用真实设备数据
	device.UseMockData = false
}

// shutdown 在应用关闭时调用
func (a *App) shutdown(ctx context.Context) {
	// 执行清理操作
}

// GetDevices 获取已连接的iOS设备列表
func (a *App) GetDevices() []device.Device {
	return device.ListDevicesWithMock()
}

// GetDeviceInfo 获取设备详细信息
func (a *App) GetDeviceInfo(udid string) (map[string]string, error) {
	return device.GetDeviceInfoWithMock(udid)
}

// BackupDevice 备份设备数据
func (a *App) BackupDevice(udid string, backupDir string, encrypt bool, password string) error {
	return device.BackupDevice(udid, backupDir, encrypt, password)
}

// RestoreDevice 恢复设备数据
func (a *App) RestoreDevice(udid string, backupDir string, password string) error {
	return device.RestoreDevice(udid, backupDir, password)
}