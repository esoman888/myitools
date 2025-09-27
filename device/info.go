package device

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// isDevicePaired 检查设备是否已配对
func isDevicePaired(udid string) bool {
	cmd := exec.Command("idevicepair", "-u", udid, "validate")
	err := cmd.Run()
	return err == nil
}

// GetDeviceInfo 获取设备详细信息
func GetDeviceInfo(udid string) (map[string]string, error) {
	info := make(map[string]string)

	// 检查设备配对状态
	if !isDevicePaired(udid) {
		info["Status"] = "需要配对"
		info["Name"] = "需要在设备上确认信任"
		info["Model"] = "请在设备上点击\"信任\"按钮"
		info["UDID"] = udid
		return info, nil
	}

	// 调用 ideviceinfo 命令获取设备信息
	cmd := exec.Command("ideviceinfo", "-u", udid)
	output, err := cmd.Output()
	if err != nil {
		// 如果命令执行失败，返回错误信息
		if strings.Contains(err.Error(), "executable file not found") {
			return info, errors.New("需要安装libimobiledevice工具")
		}
		return info, err
	}

	// 解析命令输出
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			info[key] = value
		}
	}

	// 添加一些常用信息的快捷访问
	basicInfo := map[string]string{
		"Name":         getValueOrDefault(info, "DeviceName", "未命名设备"),
		"Model":        getValueOrDefault(info, "ProductType", "未知型号"),
		"iOS Version":  getValueOrDefault(info, "ProductVersion", "未知版本"),
		"Serial":       getValueOrDefault(info, "SerialNumber", "未知序列号"),
		"UDID":         udid,
		"Capacity":     formatCapacity(getValueOrDefault(info, "TotalDiskCapacity", "0")),
		"Battery":      getBatteryInfo(udid),
		"WiFi Address": getValueOrDefault(info, "WiFiAddress", "未知"),
	}

	// 添加爱思助手风格的详细信息
	detailedInfo := map[string]string{
		"设备名称":    getValueOrDefault(info, "DeviceName", "未知"),
		"设备类型":    getValueOrDefault(info, "DeviceClass", "未知"),
		"销售型号":    getValueOrDefault(info, "ModelNumber", "未知"),
		"IMEI":    getValueOrDefault(info, "InternationalMobileEquipmentIdentity", "未知"),
		"外壳颜色":    getDeviceColor(getValueOrDefault(info, "DeviceColor", "未知")),
		"充电次数":    getBatteryCycleCount(udid),
		"销售地区":    getValueOrDefault(info, "RegionInfo", "未知"),
		"编译版本":    getValueOrDefault(info, "BuildVersion", "未知"),
		"蓝牙地址":    getValueOrDefault(info, "BluetoothAddress", "未知"),
		"WiFi序列号": getValueOrDefault(info, "WiFiAddress", "未知"),
		"设备标识":    getValueOrDefault(info, "UniqueDeviceID", "未知"),

		// 硬件信息
		"越狱状态":  checkJailbreak(udid),
		"设备型号":  getValueOrDefault(info, "ProductType", "未知"),
		"序列号":   getValueOrDefault(info, "SerialNumber", "未知"),
		"IMEI2": getValueOrDefault(info, "InternationalMobileEquipmentIdentity2", "未知"),
		"芯片型号":  getValueOrDefault(info, "ChipID", "未知"),
		"生产日期":  ParseProductionDate(getValueOrDefault(info, "SerialNumber", "未知序列号")),
		"基带版本":  getValueOrDefault(info, "BasebandVersion", "未知"),
		"蜂窝地址":  getValueOrDefault(info, "EthernetAddress", "未知"),
		"主板序列号": getValueOrDefault(info, "MLBSerialNumber", "未知"),

		// 状态信息
		"激活状态":   getValueOrDefault(info, "ActivationState", "未知"),
		"产品类型":   getValueOrDefault(info, "ProductType", "未知") + " (" + getValueOrDefault(info, "HardwareModel", "未知") + ")",
		"ECID":   getValueOrDefault(info, "UniqueChipID", "未知"),
		"硬件模型":   getValueOrDefault(info, "HardwareModel", "未知"),
		"固件版本":   getValueOrDefault(info, "FirmwareVersion", "未知"),
		"WiFi地址": getValueOrDefault(info, "WiFiAddress", "未知"),
		"CPU架构":  getValueOrDefault(info, "CPUArchitecture", "未知"),
	}

	// 合并基本信息和详细信息
	for k, v := range basicInfo {
		info[k] = v
	}

	// 合并爱思助手风格的详细信息
	for k, v := range detailedInfo {
		info[k] = v
	}

	return info, nil
}

// getValueOrDefault 从map中获取值，如果不存在则返回默认值
func getValueOrDefault(m map[string]string, key string, defaultValue string) string {
	if value, exists := m[key]; exists && value != "" {
		return value
	}
	return defaultValue
}

// getBatteryInfo 获取电池信息
func getBatteryInfo(udid string) string {
	cmd := exec.Command("idevicediagnostics", "-u", udid, "ioregentry", "AppleSmartBattery")
	output, err := cmd.Output()
	if err != nil {
		// 返回命令执行错误的信息
		return fmt.Sprintf("Error: %v", err)
	}

	// 简单解析电池电量
	if strings.Contains(string(output), "\"CurrentCapacity\" = ") {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "\"CurrentCapacity\" = ") {
				// 处理包含 CurrentCapacity 的行
				parts := strings.Split(line, "\"CurrentCapacity\" = ")
				if len(parts) > 1 {
					// 返回去除空白字符的电池容量
					return strings.TrimSpace(parts[1])
				}
			}
		}
	}

	// 如果未能找到信息，则返回"Unknown"
	return "Unknown"
}

// getBatteryCycleCount 获取电池的充电次数
func getBatteryCycleCount(udid string) string {
	cmd := exec.Command("idevicediagnostics", "-u", udid, "ioregentry", "AppleSmartBattery")
	output, err := cmd.Output()
	if err != nil {
		// 返回命令执行错误的信息
		return fmt.Sprintf("Error: %v", err)
	}

	// 解析充电次数
	if strings.Contains(string(output), "\"CycleCount\" = ") {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "\"CycleCount\" = ") {
				// 提取充电次数
				parts := strings.Split(line, "\"CycleCount\" = ")
				if len(parts) > 1 {
					// 返回充电次数
					return strings.TrimSpace(parts[1])
				}
			}
		}
	}

	// 如果未能找到充电次数，则返回"Unknown"
	return "Unknown"
}

// checkJailbreak 检测设备是否越狱
func checkJailbreak(udid string) string {
	// 方法1: 检查是否存在常见越狱应用
	jailbreakApps := []string{
		"/Applications/Cydia.app",
		"/Applications/Sileo.app",
		"/Applications/Zebra.app",
		"/Applications/Installer.app",
	}

	for _, app := range jailbreakApps {
		cmd := exec.Command("ideviceinstaller", "-u", udid, "-l", app)
		output, err := cmd.CombinedOutput()
		if err == nil && !strings.Contains(string(output), "No such file or directory") {
			return "已越狱"
		}
	}

	// 方法2: 检查是否可以访问越狱后才能访问的路径
	jailbreakPaths := []string{
		"/private/var/lib/apt/",
		"/private/var/stash/",
		"/private/var/mobile/Library/SBSettings/",
		"/System/Library/LaunchDaemons/com.saurik.Cydia.Startup.plist",
	}

	for _, path := range jailbreakPaths {
		cmd := exec.Command("idevicefile", "-u", udid, "-l", path)
		output, err := cmd.CombinedOutput()
		if err == nil && !strings.Contains(string(output), "No such file or directory") {
			return "已越狱"
		}
	}

	// 方法3: 检查是否可以执行越狱后才能执行的命令
	cmd := exec.Command("ideviceinfo", "-u", udid, "-k", "PasswordProtected")
	output, err := cmd.Output()
	if err == nil {
		// 如果设备已越狱，通常会显示为false
		if strings.TrimSpace(string(output)) == "false" {
			// 进一步检查其他特征
			cmd = exec.Command("ideviceinfo", "-u", udid, "-k", "ProductVersion")
			output, err = cmd.Output()
			if err == nil {
				version := strings.TrimSpace(string(output))
				// 检查是否为常见的可越狱版本
				if strings.HasPrefix(version, "14.") || strings.HasPrefix(version, "13.") {
					// 这些版本有已知的越狱工具
					return "可能已越狱"
				}
			}
		}
	}

	return "未越狱"
}

// formatCapacity 格式化容量显示
func formatCapacity(capacity string) string {
	if capacity == "" || capacity == "0" {
		return "未知"
	}

	// 尝试解析为数字
	if bytes, err := strconv.ParseInt(capacity, 10, 64); err == nil {
		const unit = 1024
		if bytes < unit {
			return capacity + " B"
		}
		div, exp := int64(unit), 0
		for n := bytes / unit; n >= unit; n /= unit {
			div *= unit
			exp++
		}
		return strconv.FormatFloat(float64(bytes)/float64(div), 'f', 1, 64) + " " + []string{"B", "KB", "MB", "GB", "TB"}[exp+1]
	}

	return capacity
}

// 解析苹果序列号推算生产日期
func ParseProductionDate(serial string) string {
	if len(serial) < 5 {
		return "未知"
	}

	// 清理序列号，移除空格和特殊字符
	serial = strings.TrimSpace(serial)

	// 根据序列号长度判断格式
	var year, week int
	var ok bool

	if len(serial) == 10 && serial[0] == 'C' {
		// 新格式序列号 (例如 C1QQG9XD4R)
		// 年份信息在第3位
		yearCode := strings.ToUpper(string(serial[2]))
		yearMap := map[string]int{
			"C": 2010, "D": 2011, "F": 2012, "G": 2013, "H": 2014,
			"J": 2015, "K": 2016, "L": 2017, "M": 2018, "N": 2019,
			"P": 2020, "Q": 2021, "R": 2022, "S": 2023, "T": 2024,
			"V": 2025, "W": 2026, "X": 2027, "Y": 2028, "Z": 2029,
		}

		year, ok = yearMap[yearCode]
		if !ok {
			return "未知"
		}

		// 周数信息在第9位
		weekCode := rune(serial[8])
		weekMap := map[rune]int{
			'1': 1, '2': 5, '3': 9, '4': 13, '5': 18,
			'6': 22, '7': 27, '8': 31, '9': 35, '0': 40,
			'A': 44, 'B': 48, 'C': 52,
		}

		week, ok = weekMap[weekCode]
		if !ok {
			return "未知"
		}
	} else {
		// 旧格式序列号
		yearMap := map[string]int{
			"C": 2010, "D": 2011, "F": 2012, "G": 2013, "H": 2014,
			"J": 2015, "K": 2016, "L": 2017, "M": 2018, "N": 2019,
			"P": 2020, "Q": 2021, "R": 2022, "S": 2023, "T": 2024,
			"V": 2025, "W": 2026, "X": 2027, "Y": 2028, "Z": 2029,
		}

		weekMap := map[rune]int{
			'0': 1, '1': 2, '2': 3, '3': 4, '4': 5, '5': 6, '6': 7, '7': 8, '8': 9, '9': 10,
			'A': 11, 'B': 12, 'C': 13, 'D': 14, 'E': 15, 'F': 16, 'G': 17, 'H': 18, 'J': 19, 'K': 20,
			'L': 21, 'M': 22, 'N': 23, 'P': 24, 'Q': 25, 'R': 26, 'S': 27, 'T': 28, 'V': 29, 'W': 30,
			'X': 31, 'Y': 32, 'Z': 33,
		}

		yearCode := strings.ToUpper(string(serial[3]))
		weekCode := rune(serial[4])

		year, ok = yearMap[yearCode]
		if !ok {
			return "未知"
		}

		week, ok = weekMap[weekCode]
		if !ok {
			return "未知"
		}
	}

	// 计算生产日期（取周三）
	firstDay := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	prodDate := firstDay.AddDate(0, 0, (week-1)*7+2)

	return fmt.Sprintf("%d年%02d月%02d日(第%d周)",
		prodDate.Year(),
		int(prodDate.Month()),
		prodDate.Day(),
		week,
	)
}

// DeviceColor颜色解析函数
func getDeviceColor(colorCode string) string {
	// 尝试将颜色代码转换为整数
	code, err := strconv.Atoi(colorCode)
	if err != nil {
		return "未知颜色"
	}

	colorMap := map[int]string{
		0: "黑色 / 深空灰",
		1: "白色 / 银色",
		2: "红色 (PRODUCT RED)",
		3: "蓝色",
		4: "粉色",
		5: "绿色",
		6: "紫色",
	}

	// 返回对应的颜色，如果没有匹配，返回"未知颜色"
	if color, exists := colorMap[code]; exists {
		return color
	}
	return "未知颜色"
}
