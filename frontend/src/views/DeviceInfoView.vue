<template>
  <div class="bg-gray-50 min-h-screen">
    <!-- 标题栏 -->
    <div class="bg-white shadow-sm border-b px-6 py-4">
      <div class="flex items-center justify-between">
        <div class="flex items-center">
          <el-button @click="goBack" type="default" size="small" class="mr-3">
            <el-icon><ArrowLeft /></el-icon>
            返回
          </el-button>
          <h1 class="text-xl font-semibold text-gray-800">设备详情</h1>
        </div>
        <div class="flex space-x-2">
          <el-button size="small" @click="loadDeviceInfo">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
          <!-- <el-button :type="autoRefresh ? 'success' : 'default'" size="small" @click="toggleAutoRefresh">
            <el-icon><Timer /></el-icon>
            {{ autoRefresh ? '自动刷新中' : '自动刷新' }}
          </el-button> -->
          <el-button type="primary" size="small">
            数据到剪贴板
          </el-button>
          <el-button size="small">关闭</el-button>
        </div>
      </div>
    </div>

    <div class="p-6">
      <el-alert
        v-if="deviceStore.error"
        :title="deviceStore.error"
        type="error"
        show-icon
        closable
        class="mb-4"
      />
      
      <!-- 初始加载时显示骨架屏 -->
      <div v-if="initialLoading" class="flex justify-center items-center py-8">
        <el-skeleton :rows="10" animated />
      </div>
      
      <!-- 设备信息显示 -->
      <template v-else-if="deviceInfo">
        <!-- 设备信息表格 -->
        <div class="bg-white rounded-lg shadow-sm border">
          <div class="grid text-sm" :class="{'grid-cols-3': hasDetailedInfo, 'grid-cols-1': !hasDetailedInfo}">
            <!-- 左列 -->
            <div :class="{'border-r': hasDetailedInfo}">
              <div class="bg-gray-50 px-4 py-2 border-b font-medium text-gray-700">基本信息</div>
              <div v-for="item in leftColumnData" :key="item.label" class="border-b last:border-b-0">
                <div class="grid grid-cols-2">
                  <div class="px-4 py-3 bg-gray-50 border-r text-gray-600 font-medium">{{ item.label }}</div>
                  <div class="px-4 py-3 text-gray-800">{{ item.value }}</div>
                </div>
              </div>
            </div>
            
            <!-- 中列 -->
            <div v-if="hasDetailedInfo" class="border-r">
              <div class="bg-gray-50 px-4 py-2 border-b font-medium text-gray-700">硬件信息</div>
              <div v-for="item in middleColumnData" :key="item.label" class="border-b last:border-b-0">
                <div class="grid grid-cols-2">
                  <div class="px-4 py-3 bg-gray-50 border-r text-gray-600 font-medium">{{ item.label }}</div>
                  <div class="px-4 py-3 text-gray-800">{{ item.value }}</div>
                </div>
              </div>
            </div>
            
            <!-- 右列 -->
            <div v-if="hasDetailedInfo">
              <div class="bg-gray-50 px-4 py-2 border-b font-medium text-gray-700">状态信息</div>
              <div v-for="item in rightColumnData" :key="item.label" class="border-b last:border-b-0">
                <div class="grid grid-cols-2">
                  <div class="px-4 py-3 bg-gray-50 border-r text-gray-600 font-medium">{{ item.label }}</div>
                  <div class="px-4 py-3 text-gray-800">{{ item.value }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 底部按钮 -->
        <div class="mt-6 flex justify-center space-x-4">
          <el-button @click="copyToClipboard">在记事本中打开</el-button>
          <el-button type="primary">数据到剪贴板</el-button>
          <el-button @click="goBack">关闭</el-button>
        </div>
        
        <!-- 静默刷新时的加载指示器 -->
        <div v-if="silentLoading" class="fixed bottom-4 right-4">
          <el-badge is-dot type="info">
            <el-icon class="text-gray-500 animate-spin"><Refresh /></el-icon>
          </el-badge>
        </div>
      </template>
      
      <div v-else class="text-center py-12">
        <el-empty description="无法获取设备信息">
          <template #description>
            <p>请确保设备已连接并且可以正常通信</p>
          </template>
          <el-button type="primary" @click="loadDeviceInfo">重试</el-button>
        </el-empty>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, onBeforeUnmount } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useDeviceStore } from '../stores/device'
import { ArrowLeft, Refresh, Timer } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const deviceStore = useDeviceStore()

const deviceInfo = computed(() => deviceStore.currentDevice)
const refreshInterval = ref(5000) // 刷新间隔，默认5秒
const autoRefresh = ref(false) // 是否自动刷新
let refreshTimer = null // 刷新定时器
const initialLoading = ref(true) // 初始加载状态
const silentLoading = ref(false) // 静默刷新状态

// 检查是否有详细信息
const hasDetailedInfo = computed(() => {
  if (!deviceInfo.value) return false
  return !!deviceInfo.value['设备名称']
})

// 左列数据 - 基本信息
const leftColumnData = computed(() => {
  if (!deviceInfo.value) return []
  
  // 检查是否有基本信息
  const hasBasicInfo = deviceInfo.value['Name'] || deviceInfo.value['Model'] || deviceInfo.value['iOS Version']
  
  // 如果没有详细信息，显示基本信息
  if (!deviceInfo.value['设备名称'] && hasBasicInfo) {
    return [
      { label: '设备名称', value: deviceInfo.value['Name'] || '未知' },
      { label: '设备型号', value: deviceInfo.value['Model'] || '未知' },
      { label: 'iOS版本', value: deviceInfo.value['iOS Version'] || '未知' },
      { label: '序列号', value: deviceInfo.value['Serial'] || '未知' },
      { label: 'UDID', value: deviceInfo.value['UDID'] || '未知' },
      { label: '容量', value: deviceInfo.value['Capacity'] || '未知' },
      { label: '电池', value: deviceInfo.value['Battery'] || '未知' },
      { label: 'WiFi地址', value: deviceInfo.value['WiFi Address'] || '未知' },
      { label: '状态', value: deviceInfo.value['Status'] || '已连接' }
    ]
  }
  
  // 否则显示详细信息
  return [
    { label: '设备名称', value: deviceInfo.value['设备名称'] || '未知' },
    { label: '设备类型', value: deviceInfo.value['设备类型'] || '未知' },
    { label: '销售型号', value: deviceInfo.value['销售型号'] || '未知' },
    { label: 'IMEI', value: deviceInfo.value['IMEI'] || '未知' },
    { label: '五码匹配', value: deviceInfo.value['五码匹配'] || '未知' },
    { label: '销售地区', value: deviceInfo.value['销售地区'] || '未知' },
    { label: '外壳颜色', value: deviceInfo.value['外壳颜色'] || '未知' },
    { label: '充电次数', value: deviceInfo.value['充电次数'] || '未知' },
    { label: '编译版本', value: deviceInfo.value['编译版本'] || '未知' },
    { label: '蓝牙地址', value: deviceInfo.value['蓝牙地址'] || '未知' },
    { label: '硬盘类型', value: deviceInfo.value['硬盘类型'] || '未知' },
    { label: 'WiFi序列号', value: deviceInfo.value['WiFi序列号'] || '未知' },
    { label: '屏幕分辨率', value: deviceInfo.value['屏幕分辨率'] || '未知' },
    { label: '红外摄像头', value: deviceInfo.value['红外摄像头'] || '未知' },
    { label: '设备标识', value: deviceInfo.value['设备标识'] || '未知' },
    { label: '盒命编号', value: deviceInfo.value['盒命编号'] || '未知' },
    { label: '振动器编号', value: deviceInfo.value['振动器编号'] || '未知' }
  ]
})

// 中列数据 - 硬件信息
const middleColumnData = computed(() => {
  if (!deviceInfo.value) return []
  
  // 检查是否有基本信息
  const hasBasicInfo = deviceInfo.value['Name'] || deviceInfo.value['Model'] || deviceInfo.value['iOS Version']
  
  // 如果没有详细信息，显示空数组或基本信息
  if (!deviceInfo.value['设备名称'] && hasBasicInfo) {
    return []
  }
  
  return [
    { label: '越狱状态', value: deviceInfo.value['越狱状态'] || '未知' },
    { label: '设备型号', value: deviceInfo.value['设备型号'] || '未知' },
    { label: '序列号', value: deviceInfo.value['序列号'] || '未知' },
    { label: 'IMEI2', value: deviceInfo.value['IMEI2'] || '未知' },
    { label: '芯片型号', value: deviceInfo.value['芯片型号'] || '未知' },
    { label: '生产日期', value: deviceInfo.value['生产日期'] || '未知' },
    { label: '正在充电', value: deviceInfo.value['正在充电'] || '未知' },
    { label: '电池寿命', value: deviceInfo.value['电池寿命'] || '未知' },
    { label: '基带版本', value: deviceInfo.value['基带版本'] || '未知' },
    { label: '蜂窝地址', value: deviceInfo.value['蜂窝地址'] || '未知' },
    { label: '主板序列号', value: deviceInfo.value['主板序列号'] || '未知' },
    { label: 'WiFi模块', value: deviceInfo.value['WiFi模块'] || '未知' },
    { label: '强制重启', value: deviceInfo.value['强制重启'] || '未知' },
    { label: '芯牌', value: deviceInfo.value['芯牌'] || '未知' }
  ]
})

// 右列数据 - 状态信息
const rightColumnData = computed(() => {
  if (!deviceInfo.value) return []
  
  // 检查是否有基本信息
  const hasBasicInfo = deviceInfo.value['Name'] || deviceInfo.value['Model'] || deviceInfo.value['iOS Version']
  
  // 如果没有详细信息，显示空数组或基本信息
  if (!deviceInfo.value['设备名称'] && hasBasicInfo) {
    return []
  }
  
  return [
    { label: '激活状态', value: deviceInfo.value['激活状态'] || '未知' },
    { label: '产品类型', value: deviceInfo.value['产品类型'] || '未知' },
    { label: '卡槽类型', value: deviceInfo.value['卡槽类型'] || '未知' },
    { label: 'ECID', value: deviceInfo.value['ECID'] || '未知' },
    { label: '硬件模型', value: deviceInfo.value['硬件模型'] || '未知' },
    { label: '面板颜色', value: deviceInfo.value['面板颜色'] || '未知' },
    { label: '剩余电量', value: deviceInfo.value['剩余电量'] || '未知' },
    { label: '固件版本', value: deviceInfo.value['固件版本'] || '未知' },
    { label: 'WiFi地址', value: deviceInfo.value['WiFi地址'] || '未知' },
    { label: 'CPU架构', value: deviceInfo.value['CPU架构'] || '未知' },
    { label: '销售类型', value: deviceInfo.value['销售类型'] || '未知' },
    { label: '电话号码', value: deviceInfo.value['电话号码'] || '未知' },
    { label: '距离传感器', value: deviceInfo.value['距离传感器'] || '未知' },
    { label: '环境光', value: deviceInfo.value['环境光'] || '未知' }
  ]
})

const loadDeviceInfo = async (showLoading = true) => {
  const deviceId = route.params.id
  if (deviceId) {
    if (showLoading) {
      // 如果是初始加载，使用initialLoading
      if (initialLoading.value) {
        initialLoading.value = true
      } else {
        // 否则使用静默加载指示器
        silentLoading.value = true
      }
    }
    
    try {
      await deviceStore.fetchDeviceInfo(deviceId, true) // 始终使用静默模式，由我们自己控制UI状态
    } catch (error) {
      console.error('加载设备信息失败:', error)
    } finally {
      if (showLoading) {
        initialLoading.value = false
        silentLoading.value = false
      }
    }
  }
}

// 开始自动刷新
const startAutoRefresh = () => {
  stopAutoRefresh() // 先清除可能存在的定时器
  if (autoRefresh.value) {
    refreshTimer = setInterval(async () => {
      const deviceId = route.params.id
      
      try {
        // 先检查设备是否仍然连接
        const isConnected = await deviceStore.checkDeviceConnection(deviceId)
        
        if (!isConnected) {
          // 设备已断开，返回首页
          ElMessage.warning('设备已断开连接')
          stopAutoRefresh()
          router.push('/')
          return
        }
        
        // 设备仍然连接，静默刷新设备信息
        silentLoading.value = true
        await deviceStore.fetchDeviceInfo(deviceId, true)
        silentLoading.value = false
      } catch (error) {
        console.error('自动刷新失败:', error)
        silentLoading.value = false
      }
    }, refreshInterval.value)
  }
}

// 停止自动刷新
const stopAutoRefresh = () => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
}

// 切换自动刷新状态
const toggleAutoRefresh = () => {
  autoRefresh.value = !autoRefresh.value
  if (autoRefresh.value) {
    startAutoRefresh()
    ElMessage.success('已开启自动刷新')
  } else {
    stopAutoRefresh()
    ElMessage.info('已关闭自动刷新')
  }
}

const copyToClipboard = () => {
  ElMessage.info('功能开发中...')
}

const goBack = () => {
  router.push('/')
}

onMounted(() => {
  loadDeviceInfo(true)
  startAutoRefresh()
})

onBeforeUnmount(() => {
  stopAutoRefresh()
})
</script>