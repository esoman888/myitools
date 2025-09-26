<template>
  <header class="bg-white shadow-md">
    <div class="container mx-auto px-4 py-3 flex items-center justify-between">
      <div class="flex items-center">
        <img src="../assets/logo.svg" alt="MyiToolsApp Logo" class="h-8 w-8 mr-2" />
        <h1 class="text-xl font-bold text-primary-600">MyiToolsApp</h1>
      </div>
      
      <div class="flex items-center">
        <div class="mr-4 text-sm">
          <span class="text-gray-500">已连接设备:</span>
          <span class="ml-1 font-medium" :class="deviceCount > 0 ? 'text-green-600' : 'text-red-500'">
            {{ deviceCount }}
          </span>
        </div>
        
        <el-button type="primary" size="small" @click="refreshDevices" :loading="isLoading">
          <el-icon><Refresh /></el-icon>
          刷新设备
        </el-button>
      </div>
    </div>
  </header>
</template>

<script setup>
import { computed } from 'vue'
import { useDeviceStore } from '../stores/device'
import { Refresh } from '@element-plus/icons-vue'

const deviceStore = useDeviceStore()

const deviceCount = computed(() => deviceStore.deviceCount)
const isLoading = computed(() => deviceStore.isLoading)

const refreshDevices = async () => {
  await deviceStore.fetchDevices()
}
</script>