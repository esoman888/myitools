<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 顶部导航栏 -->
    <div class="bg-white shadow-sm border-b">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
          <div class="flex items-center">
            <h1 class="text-xl font-semibold text-gray-900">MyiToolsApp</h1>
          </div>
          <div class="flex items-center space-x-4">
            <el-button size="small" @click="refreshDevices" :loading="deviceStore.isLoading">
              <el-icon><Refresh /></el-icon>
              刷新设备
            </el-button>
          </div>
        </div>
      </div>
    </div>

    <!-- 主内容区域 -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- 设备状态提示 -->
      <el-alert
        v-if="deviceStore.error"
        :title="deviceStore.error"
        type="error"
        show-icon
        closable
        class="mb-6"
      />

      <!-- 初始加载状态 -->
      <div v-if="initialLoading" class="flex justify-center items-center py-12">
        <el-skeleton :rows="3" animated />
      </div>

      <!-- 无设备状态 -->
      <div v-else-if="!initialLoading && deviceStore.devices.length === 0" class="text-center py-12">
        <div class="bg-white rounded-lg shadow-sm p-8">
          <el-empty description="未检测到设备">
            <template #description>
              <p class="text-gray-500 mb-4">请连接您的iOS设备并确保已安装必要的驱动程序</p>
            </template>
            <el-button type="primary" @click="refreshDevices" :loading="refreshing">刷新设备</el-button>
          </el-empty>
        </div>
      </div>

      <!-- 设备卡片列表 -->
      <div v-else-if="!initialLoading && deviceStore.devices.length > 0" class="space-y-6">
        <div v-for="device in deviceStore.devices" :key="device.udid" class="bg-white rounded-lg shadow-sm border">
          <div class="p-6">
            <div class="flex items-start justify-between">
              <div class="flex items-center space-x-4">
                <!-- 设备图标 -->
                <div class="w-16 h-16 bg-gradient-to-br from-blue-500 to-purple-600 rounded-xl flex items-center justify-center">
                  <el-icon class="text-white text-2xl"><Iphone /></el-icon>
                </div>
                
                <!-- 设备基本信息 -->
                <div>
                  <h3 class="text-lg font-semibold text-gray-900">{{ device.name || '未命名设备' }}</h3>
                  <p class="text-sm text-gray-500">{{ device.model || '未知型号' }}</p>
                  <div class="flex items-center mt-1">
                    <el-tag type="success" size="small">已连接</el-tag>
                  </div>
                </div>
              </div>

              <!-- 操作按钮 -->
              <div class="flex space-x-2">
                <el-button size="small" @click="viewDeviceInfo(device)">
                  设备详情
                </el-button>
                <el-dropdown>
                  <el-button size="small">
                    更多操作<el-icon class="el-icon--right"><ArrowDown /></el-icon>
                  </el-button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item @click="goToBackup(device)">备份/恢复</el-dropdown-item>
                      <el-dropdown-item @click="goToFileSystem(device)">文件管理</el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </div>
            </div>

            <!-- 设备详细信息 -->
            <div class="mt-6 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
              <div class="bg-gray-50 rounded-lg p-4">
                <div class="text-sm text-gray-500">UDID</div>
                <div class="text-sm font-mono text-gray-900 truncate">{{ device.udid }}</div>
              </div>
              <div class="bg-gray-50 rounded-lg p-4">
                <div class="text-sm text-gray-500">状态</div>
                <div class="text-sm text-green-600 font-medium">{{ device.status === 'connected' ? '已连接' : '未连接' }}</div>
              </div>
              <div class="bg-gray-50 rounded-lg p-4">
                <div class="text-sm text-gray-500">型号</div>
                <div class="text-sm text-gray-900">{{ device.model || '未知' }}</div>
              </div>
              <div class="bg-gray-50 rounded-lg p-4">
                <div class="text-sm text-gray-500">名称</div>
                <div class="text-sm text-gray-900">{{ device.name || '未命名' }}</div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 静默刷新指示器 -->
        <div v-if="silentRefreshing" class="fixed bottom-4 right-4">
          <el-badge is-dot type="info">
            <el-icon class="text-gray-500 animate-spin"><Refresh /></el-icon>
          </el-badge>
        </div>
      </div>

      <!-- 快速操作区域 -->
      <div class="mt-12">
        <h2 class="text-lg font-semibold text-gray-900 mb-6">快速操作</h2>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <div class="bg-white rounded-lg shadow-sm border p-6 hover:shadow-md transition-shadow">
            <div class="flex items-center">
              <div class="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center">
                <el-icon class="text-blue-600 text-xl"><Download /></el-icon>
              </div>
              <div class="ml-4">
                <h3 class="text-lg font-medium text-gray-900">安装应用</h3>
                <p class="text-sm text-gray-500">安装IPA文件到设备</p>
              </div>
            </div>
            <div class="mt-4">
              <el-button type="primary" :disabled="!hasConnectedDevices" class="w-full">
                选择IPA文件
              </el-button>
            </div>
          </div>

          <div class="bg-white rounded-lg shadow-sm border p-6 hover:shadow-md transition-shadow">
            <div class="flex items-center">
              <div class="w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center">
                <el-icon class="text-green-600 text-xl"><FolderOpened /></el-icon>
              </div>
              <div class="ml-4">
                <h3 class="text-lg font-medium text-gray-900">备份设备</h3>
                <p class="text-sm text-gray-500">创建设备数据备份</p>
              </div>
            </div>
            <div class="mt-4">
              <el-button type="success" :disabled="!hasConnectedDevices" class="w-full">
                开始备份
              </el-button>
            </div>
          </div>

          <div class="bg-white rounded-lg shadow-sm border p-6 hover:shadow-md transition-shadow">
            <div class="flex items-center">
              <div class="w-12 h-12 bg-purple-100 rounded-lg flex items-center justify-center">
                <el-icon class="text-purple-600 text-xl"><Document /></el-icon>
              </div>
              <div class="ml-4">
                <h3 class="text-lg font-medium text-gray-900">文件管理</h3>
                <p class="text-sm text-gray-500">浏览设备文件系统</p>
              </div>
            </div>
            <div class="mt-4">
              <el-button type="primary" :disabled="!hasConnectedDevices" class="w-full">
                打开文件管理
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, onBeforeUnmount, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useDeviceStore } from '../stores/device'
import { 
  Refresh, 
  Iphone, 
  ArrowDown, 
  Download, 
  FolderOpened, 
  Document 
} from '@element-plus/icons-vue'

const router = useRouter()
const deviceStore = useDeviceStore()
const refreshInterval = ref(5000) // 刷新间隔，默认5秒
let refreshTimer = null // 刷新定时器
const initialLoading = ref(true) // 初始加载状态
const refreshing = ref(false) // 手动刷新状态
const silentRefreshing = ref(false) // 静默刷新状态

const hasConnectedDevices = computed(() => {
  return deviceStore.connectedDevices.length > 0
})

const refreshDevices = async () => {
  try {
    refreshing.value = true
    await deviceStore.fetchDevices()
  } catch (error) {
    console.error('刷新设备失败:', error)
  } finally {
    refreshing.value = false
    initialLoading.value = false
  }
}

const viewDeviceInfo = (device) => {
  deviceStore.setCurrentDevice(device)
  router.push({ name: 'device-info', params: { id: device.udid } })
}

const goToBackup = (device) => {
  router.push({ name: 'backup', params: { id: device.udid } })
}

const goToFileSystem = (device) => {
  router.push({ name: 'filesystem', params: { id: device.udid } })
}

// 开始自动刷新设备列表
const startAutoRefresh = () => {
  stopAutoRefresh() // 先清除可能存在的定时器
  refreshTimer = setInterval(async () => {
    try {
      silentRefreshing.value = true
      await deviceStore.fetchDevices()
    } catch (error) {
      console.error('自动刷新设备失败:', error)
    } finally {
      silentRefreshing.value = false
    }
  }, refreshInterval.value)
}

// 停止自动刷新
const stopAutoRefresh = () => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
}

onMounted(async () => {
  try {
    initialLoading.value = true
    await refreshDevices()
  } finally {
    initialLoading.value = false
  }
  startAutoRefresh()
})

onBeforeUnmount(() => {
  stopAutoRefresh()
})
</script>