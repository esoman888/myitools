import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useDeviceStore = defineStore('device', () => {
  // 状态
  const devices = ref([])
  const currentDevice = ref(null)
  const isLoading = ref(false)
  const error = ref(null)

  // 计算属性
  const connectedDevices = computed(() => {
    return devices.value.filter(device => device.status === 'connected')
  })

  const deviceCount = computed(() => {
    return connectedDevices.value.length
  })

  // 动作
  async function fetchDevices() {
    isLoading.value = true
    error.value = null
    
    try {
      // 检查window.go是否存在
      if (!window.go || !window.go.main || !window.go.main.App) {
        // 开发环境下使用模拟数据
        console.log('使用模拟设备列表 - 开发环境')
        const mockDevices = [
          {
            udid: 'device1',
            name: 'iPhone 13',
            model: 'iPhone13,1',
            status: 'connected',
            'iOS Version': '15.4.1'
          },
          {
            udid: 'device2',
            name: 'iPad Pro',
            model: 'iPad8,9',
            status: 'connected',
            'iOS Version': '16.0'
          }
        ]
        devices.value = mockDevices
        return mockDevices
      }
      
      // 调用后端API获取设备列表
      const deviceList = await window.go.main.App.GetDevices()
      devices.value = deviceList
      return deviceList
    } catch (err) {
      error.value = '获取设备列表失败: ' + err.message
      console.error('获取设备列表失败:', err)
      return []
    } finally {
      isLoading.value = false
    }
  }

  async function fetchDeviceInfo(udid, silent = false) {
    if (!silent) {
      isLoading.value = true
    }
    error.value = null;
    
    try {
      // 检查window.go是否存在
      if (!window.go || !window.go.main || !window.go.main.App) {
        // 开发环境下使用模拟数据
        console.log('使用模拟数据 - 开发环境')
        const mockDeviceInfo = {
          '设备名称': 'iPhone 模拟设备',
          '设备类型': 'iPhone',
          '销售型号': 'iPhone 13',
          'IMEI': '123456789012345',
          '五码匹配': '是',
          '销售地区': '中国大陆',
          '外壳颜色': '午夜色',
          '充电次数': '125',
          '编译版本': '19A346',
          '蓝牙地址': '00:11:22:33:44:55',
          '硬盘类型': 'NAND',
          'WiFi序列号': 'C0FFEE',
          '屏幕分辨率': '2532 x 1170',
          '红外摄像头': '是',
          '设备标识': 'iPhone14,5',
          '盒命编号': 'ABCDEF',
          '振动器编号': 'XYZ123',
          '越狱状态': '未越狱',
          '设备型号': 'A2482',
          '序列号': 'ABCDEFGHIJK',
          'IMEI2': '987654321098765',
          '芯片型号': 'A15 Bionic',
          '生产日期': '2022年1月',
          '正在充电': '否',
          '电池寿命': '95%',
          '基带版本': '1.59.00',
          '蜂窝地址': '01:23:45:67:89:AB',
          '主板序列号': 'J123S456V789',
          'WiFi模块': 'Broadcom',
          '强制重启': '音量上+音量下+侧边按钮',
          '芯牌': 'TSMC',
          '激活状态': '已激活',
          '产品类型': '零售版',
          '卡槽类型': '双卡',
          'ECID': '0x1A2B3C4D',
          '硬件模型': 'D10AP',
          '面板颜色': '黑色',
          '剩余电量': '87%',
          '固件版本': 'iBoot-7429.61.2',
          'WiFi地址': '11:22:33:44:55:66',
          'CPU架构': 'arm64e',
          '销售类型': '国行',
          '电话号码': '+86 138****1234',
          '距离传感器': '正常',
          '环境光': '正常'
        };
        
        // 更新当前设备信息
        const device = devices.value.find(d => d.udid === udid);
        if (device) {
          currentDevice.value = { ...device, ...mockDeviceInfo };
        } else {
          currentDevice.value = mockDeviceInfo;
        }
        
        return mockDeviceInfo;
      }
      
      // 正常环境下调用后端API
      const deviceInfo = await window.go.main.App.GetDeviceInfo(udid);
      
      // 更新当前设备信息
      const device = devices.value.find(d => d.udid === udid);
      if (device) {
        currentDevice.value = { ...device, ...deviceInfo };
      } else {
        currentDevice.value = deviceInfo;
      }
      
      return deviceInfo;
    } catch (err) {
      error.value = '获取设备信息失败: ' + err.message;
      console.error('获取设备信息失败:', err);
      return null;
    } finally {
      if (!silent) {
        isLoading.value = false;
      }
    }
  }

  // 检查设备是否仍然连接
  async function checkDeviceConnection(udid) {
    try {
      // 检查window.go是否存在
      if (!window.go || !window.go.main || !window.go.main.App) {
        // 开发环境下，始终返回已连接
        console.log('开发环境 - 设备始终连接')
        return true
      }
      
      // 获取最新的设备列表
      await fetchDevices()
      
      // 检查设备是否在列表中
      const deviceStillConnected = devices.value.some(device => device.udid === udid)
      return deviceStillConnected
    } catch (err) {
      console.error('检查设备连接状态失败:', err)
      return false
    }
  }

  // 设置当前设备
  function setCurrentDevice(device) {
    currentDevice.value = device
  }

  return {
    devices,
    currentDevice,
    isLoading,
    error,
    connectedDevices,
    deviceCount,
    fetchDevices,
    fetchDeviceInfo,
    checkDeviceConnection,
    setCurrentDevice
  }
})