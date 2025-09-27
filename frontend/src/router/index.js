import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import DeviceInfoView from '../views/DeviceInfoView.vue'
import BackupView from '../views/BackupView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/device/:id',
      name: 'device-info',
      component: DeviceInfoView
    },
    {
      path: '/backup',
      name: 'backup',
      component: BackupView
    }
  ]
})

export default router