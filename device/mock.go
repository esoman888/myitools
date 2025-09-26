package device

import "fmt"

// GetMockDevices 获取模拟设备数据，用于测试
func GetMockDevices() []Device {
	return []Device{
		{
			UDID:   "00008110-001238E23E614015",
			Name:   "我的iPhone",
			Model:  "iPhone14,2",
			Status: "connected",
		},
		{
			UDID:   "00008030-001A35E40C10802E",
			Name:   "测试iPad",
			Model:  "iPad13,1",
			Status: "connected",
		},
	}
}

// GetMockDeviceInfo 获取模拟设备详细信息
func GetMockDeviceInfo(udid string) map[string]string {
	mockInfo := map[string]string{
		"Name":         "iPhone",
		"Model":        "iPhone",
		"iOS Version":  "未越狱",
		"Serial":       "MPU93 CH/A",
		"UDID":         "352340763085655",
		"Capacity":     "128.0 GB",
		"Battery":      "85%",
		"WiFi Address": "aa:bb:cc:dd:ee:ff",
		
		// 爱思助手风格的详细信息
		"设备名称":          "iPhone",
		"设备类型":          "iPhone",
		"销售型号":          "MPU93 CH/A",
		"IMEI":           "352340763085655",
		"五码匹配":          "是",
		"销售地区":          "CH/A(中国)",
		"外壳颜色":          "午夜色",
		"充电次数":          "177次",
		"编译版本":          "20A380",
		"蓝牙地址":          "B8:14:4D:53:4E:3F",
		"硬盘类型":          "TLC",
		"WiFi序列号":        "信息无法读取",
		"屏幕分辨率":         "2532 x 1170",
		"红外摄像头":         "HNQ23853VSV15F7DK",
		"设备标识":          "00008110-001238E23E614015",
		"盒命编号":          "FX92287133C15Y99H+340052910001180015385019833",
		"振动器编号":         "GH92345331G1FVDAR",
		
		// 右侧信息
		"越狱状态":          "未越狱",
		"设备型号":          "iPhone 14",
		"序列号":           "H6JPCQ9W6W",
		"IMEI2":          "352340763157868",
		"芯片型号":          "8110",
		"生产日期":          "2022年10月06日(第41周)",
		"正在充电":          "是",
		"电池寿命":          "100%",
		"基带版本":          "1.00.05",
		"蜂窝地址":          "B8:14:4D:53:70:0D",
		"主板序列号":         "F3Y2414OSY31JRQA",
		"WiFi模块":         "低温WiFi (STATS)",
		"强制重启":          "信息无法读取",
		"芯牌":            "A22F9800033C1CBE24",
		
		// 激活状态
		"激活状态":          "已激活",
		"产品类型":          "iPhone14,7 (A2884)",
		"卡槽类型":          "双卡",
		"ECID":           "001238E23E614015",
		"硬件模型":          "D27AP",
		"面板颜色":          "黑色",
		"剩余电量":          "49%",
		"固件版本":          "16.0.2",
		"WiFi地址":         "B8:14:4D:43:75:6D",
		"CPU架构":          "arm64e",
		"销售类型":          "零售机",
		"电话号码":          "信息无法读取",
		"距离传感器":         "FWP22475B9KUQ948DL",
		"环境光":           "031084941207",
	}
	
	return mockInfo
}

// UseMockData 是否使用模拟数据的标志
var UseMockData = false

// ListDevicesWithMock 列出设备（支持模拟数据）
func ListDevicesWithMock() []Device {
	if UseMockData {
		fmt.Println("使用模拟设备数据")
		return GetMockDevices()
	}
	return ListDevices()
}

// GetDeviceInfoWithMock 获取设备信息（支持模拟数据）
func GetDeviceInfoWithMock(udid string) (map[string]string, error) {
	if UseMockData {
		fmt.Println("使用模拟设备信息")
		return GetMockDeviceInfo(udid), nil
	}
	return GetDeviceInfo(udid)
}