<template>
  <div class="app-container">
    <AppHeader />
    <div class="content-container">
      <router-view />
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useDeviceStore } from './stores/device'
import AppHeader from './components/AppHeader.vue'

const deviceStore = useDeviceStore()

onMounted(async () => {
  // 应用启动时获取设备列表
  await deviceStore.fetchDevices()
  
  // 设置定时刷新设备列表
  setInterval(async () => {
    await deviceStore.fetchDevices()
  }, 5000) // 每5秒刷新一次
})
</script>

<style>
/* 全局样式已在main.css中定义 */
</style>