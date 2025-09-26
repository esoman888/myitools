<template>
  <div class="bg-white rounded-lg shadow-md p-4 hover:shadow-lg transition-shadow">
    <div class="flex items-center justify-between mb-3">
      <h3 class="text-lg font-semibold text-gray-800">{{ device.name || '未命名设备' }}</h3>
      <el-tag type="success" v-if="device.status === 'connected'">已连接</el-tag>
      <el-tag type="danger" v-else>已断开</el-tag>
    </div>
    
    <div class="text-sm text-gray-600 mb-4">
      <div class="mb-1">
        <span class="font-medium">型号:</span> {{ device.model || '未知' }}
      </div>
      <div class="mb-1">
        <span class="font-medium">UDID:</span> 
        <span class="text-xs font-mono">{{ truncateUDID(device.udid) }}</span>
      </div>
    </div>
    
    <div class="flex justify-between mt-4">
      <el-button type="primary" size="small" @click="viewDeviceInfo">
        设备信息
      </el-button>
      <el-dropdown>
        <el-button type="default" size="small">
          更多操作<el-icon class="el-icon--right"><arrow-down /></el-icon>
        </el-button>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item @click="goToBackup">备份/恢复</el-dropdown-item>
            <el-dropdown-item @click="goToFileSystem">文件系统</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useDeviceStore } from '../stores/device'
import { ArrowDown } from '@element-plus/icons-vue'

const props = defineProps({
  device: {
    type: Object,
    required: true
  }
})

const router = useRouter()
const deviceStore = useDeviceStore()

const truncateUDID = (udid) => {
  if (!udid) return '未知'
  return udid.length > 12 ? `${udid.substring(0, 8)}...` : udid
}

const viewDeviceInfo = () => {
  deviceStore.setCurrentDevice(props.device)
  router.push({ name: 'device-info', params: { id: props.device.udid } })
}

const goToBackup = () => {
  router.push({ name: 'backup', params: { id: props.device.udid } })
}

const goToFileSystem = () => {
  router.push({ name: 'filesystem', params: { id: props.device.udid } })
}
</script>