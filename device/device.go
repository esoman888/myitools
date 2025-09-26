package device

import (
	"errors"
	"os/exec"
	"strings"
)

// Device 表示iOS设备
type Device struct {
	UDID   string `json:"udid"`
	Name   string `json:"name"`
	Model  string `json:"model"`
	Status string `json:"status"`
}

// ListDevices 列出所有已连接的iOS设备
func ListDevices() []Device {
	devices := []Device{}

	// 调用 idevice_id -l 命令获取设备列表
	cmd := exec.Command("idevice_id", "-l")
	output, err := cmd.Output()
	if err != nil {
		println("获取设备列表失败:", err.Error())
		return devices
	}

	// 解析命令输出
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		udid := strings.TrimSpace(line)
		if udid != "" {
			// 获取设备名称
			name := getDeviceName(udid)
			// 获取设备型号
			model := getDeviceModel(udid)

			// 添加调试信息
			println("设备UDID:", udid)
			println("设备名称:", name)
			println("设备型号:", model)

			// 检查是否需要配对
			if name == "未命名设备" && model == "未知型号" {
				pairingStatus := checkPairingStatus(udid)
				println("设备配对状态:", pairingStatus)
				
				devices = append(devices, Device{
					UDID:   udid,
					Name:   "需要在设备上确认信任",
					Model:  "请在设备上点击\"信任\"按钮",
					Status: "需要配对",
				})
			} else {
				devices = append(devices, Device{
					UDID:   udid,
					Name:   name,
					Model:  model,
					Status: "connected",
				})
			}
		}
	}
	return devices
}

// checkPairingStatus 检查设备配对状态
func checkPairingStatus(udid string) string {
	cmd := exec.Command("idevicepair", "-u", udid, "validate")
	output, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(output), "trust dialog") {
			return "需要在设备上确认信任"
		}
		return "配对失败: " + string(output)
	}
	return "已配对"
}

// getDeviceName 获取设备名称
func getDeviceName(udid string) string {
	cmd := exec.Command("ideviceinfo", "-u", udid, "-k", "DeviceName")
	output, err := cmd.Output()
	if err != nil {
		// 如果命令执行失败，返回更详细的错误信息
		println("获取设备名称失败:", err.Error())
		if strings.Contains(err.Error(), "executable file not found") {
			return "需要安装libimobiledevice"
		}
		return "未命名设备"
	}
	result := strings.TrimSpace(string(output))
	if result == "" {
		return "未命名设备"
	}
	return result
}

// getDeviceModel 获取设备型号
func getDeviceModel(udid string) string {
	cmd := exec.Command("ideviceinfo", "-u", udid, "-k", "ProductType")
	output, err := cmd.Output()
	if err != nil {
		// 如果命令执行失败，返回更详细的错误信息
		println("获取设备型号失败:", err.Error())
		if strings.Contains(err.Error(), "executable file not found") {
			return "需要安装libimobiledevice"
		}
		return "未知型号"
	}
	result := strings.TrimSpace(string(output))
	if result == "" {
		return "未知型号"
	}
	return result
}

// IsDeviceConnected 检查设备是否已连接
func IsDeviceConnected(udid string) bool {
	cmd := exec.Command("idevice_id", "-l")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.Contains(string(output), udid)
}

// BackupDevice 备份设备数据
func BackupDevice(udid string, backupDir string, encrypt bool, password string) error {
	if !IsDeviceConnected(udid) {
		return errors.New("设备未连接")
	}

	args := []string{"-u", udid, "backup", "--full", backupDir}
	if encrypt {
		args = append(args, "--password", password)
	}

	cmd := exec.Command("idevicebackup2", args...)
	_, err := cmd.Output()
	return err
}

// RestoreDevice 恢复设备数据
func RestoreDevice(udid string, backupDir string, password string) error {
	if !IsDeviceConnected(udid) {
		return errors.New("设备未连接")
	}

	args := []string{"-u", udid, "restore", "--full", backupDir}
	if password != "" {
		args = append(args, "--password", password)
	}

	cmd := exec.Command("idevicebackup2", args...)
	_, err := cmd.Output()
	return err
}
