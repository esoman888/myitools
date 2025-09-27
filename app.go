package main

import (
	"context"
	"fmt"
	"myitools/device"
	"myitools/dialog"
	"os"
	"path/filepath"
	"runtime"
)

// App 应用结构体
type App struct {
	ctx    context.Context
	dialog *dialog.DialogManager
}

// NewApp 创建一个新的App应用实例
func NewApp() *App {
	return &App{
		dialog: dialog.NewDialogManager(),
	}
}

// startup 在应用启动时调用
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.dialog.SetContext(ctx)
	// 使用真实设备数据
	device.UseMockData = false

	// 确保备份目录存在
	a.ensureBackupDirExists()
}

// ensureBackupDirExists 确保备份目录存在
func (a *App) ensureBackupDirExists() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("获取用户主目录失败: %v\n", err)
		return
	}

	// 创建默认备份目录
	backupDir := filepath.Join(homeDir, "Desktop", "MateTools", "backups")
	err = os.MkdirAll(backupDir, 0755)
	if err != nil {
		fmt.Printf("创建备份目录失败: %v\n", err)
	} else {
		fmt.Printf("备份目录已创建或已存在: %s\n", backupDir)
	}
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
func (a *App) BackupDevice(udid string, backupDir string, encrypt bool, password string) (string, error) {
	return device.CreateBackup(udid, backupDir, encrypt, password)
}

// RestoreDevice 恢复设备数据
func (a *App) RestoreDevice(udid string, backupDir string, password string) error {
	return device.RestoreBackup(udid, backupDir, password)
}

// GetBackupProgress 获取备份进度
func (a *App) GetBackupProgress(backupID string) *device.BackupProgress {
	return device.GetBackupProgress(backupID)
}

// ListBackups 列出所有备份
func (a *App) ListBackups(backupBaseDir string) ([]device.BackupInfo, error) {
	return device.ListBackups(backupBaseDir)
}

// DeleteBackup 删除备份
func (a *App) DeleteBackup(backupPath string) error {
	// return device.DeleteBackup(backupPath)
	return nil
}

// OpenDirectoryDialog 打开目录选择对话框
func (a *App) OpenDirectoryDialog(title string, defaultPath string) (string, error) {
	return a.dialog.OpenDirectoryDialog(title, defaultPath)
}

// GetBackupInfo 获取备份信息
func (a *App) GetBackupInfo(backupPath string) (device.BackupInfo, error) {
	return device.GetBackupInfo(backupPath)
}

// CheckBackupEncryptionStatus 检查设备备份加密状态
func (a *App) CheckBackupEncryptionStatus(udid string) (bool, error) {
	return device.CheckBackupEncryptionStatus(udid)
}

// SetBackupEncryption 设置备份加密状态
func (a *App) SetBackupEncryption(udid string, enable bool, password string) error {
	return device.SetBackupEncryption(udid, enable, password)
}

// OpenFileDialog 打开文件选择对话框
func (a *App) OpenFileDialog(title string, defaultPath string) (string, error) {
	return a.dialog.OpenFileDialog(title, defaultPath, nil)
}

// SaveFileDialog 打开文件保存对话框
func (a *App) SaveFileDialog(title string, defaultPath string) (string, error) {
	return a.dialog.SaveFileDialog(title, defaultPath, nil)
}

// GetDefaultBackupDir 获取默认备份目录
func (a *App) GetDefaultBackupDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// 如果无法获取用户主目录，使用当前目录
		currentDir, _ := os.Getwd()
		backupDir := filepath.Join(currentDir, "backups")
		return backupDir
	}

	// 使用用户文档目录下的路径
	var backupDir string
	if runtime.GOOS == "darwin" {
		// macOS 上使用 Documents 目录
		backupDir = filepath.Join(homeDir, "Documents", "myitools", "backups")
	} else if runtime.GOOS == "windows" {
		// Windows 上使用 Documents 目录
		backupDir = filepath.Join(homeDir, "Documents", "myitools", "backups")
	} else {
		// Linux 和其他系统使用 home 目录下的 .myitools
		backupDir = filepath.Join(homeDir, ".myitools", "backups")
	}

	// 确保目录存在
	err = os.MkdirAll(backupDir, 0755)
	if err != nil {
		fmt.Printf("创建备份目录失败: %v\n", err)
		// 如果创建失败，使用临时目录
		tempDir := filepath.Join(os.TempDir(), "myitools", "backups")
		os.MkdirAll(tempDir, 0755)
		return tempDir
	}

	return backupDir
}
