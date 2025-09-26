import { createRouter, createWebHashHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/device/:id',
      name: 'device-info',
      component: () => import('../views/DeviceInfoView.vue')
    },
    {
      path: '/backup/:id',
      name: 'backup',
      component: () => import('../views/BackupView.vue')
    },
    {
      path: '/filesystem/:id',
      name: 'filesystem',
      component: () => import('../views/FileSystemView.vue')
    }
  ]
})

export default router