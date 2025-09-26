<template>
  <div>
    <div class="flex items-center mb-4">
      <el-button @click="goBack" type="default" size="small" class="mr-2">
        <el-icon><ArrowLeft /></el-icon>
        返回
      </el-button>
      <h2 class="text-2xl font-bold">备份与恢复</h2>
    </div>
    
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- 备份区域 -->
      <div class="bg-white rounded-lg shadow-md p-6">
        <h3 class="text-xl font-semibold mb-4 flex items-center">
          <el-icon class="mr-2"><Download /></el-icon>
          设备备份
        </h3>
        
        <el-form :model="backupForm" label-width="100px">
          <el-form-item label="备份路径">
            <el-input 
              v-model="backupForm.path" 
              placeholder="选择备份保存路径"
              readonly
            >
              <template #append>
                <el-button @click="selectBackupPath">选择路径</el-button>
              </template>
            </el-input>
          </el-form-item>
          
          <el-form-item label="加密备份">
            <el-switch v-model="backupForm.encrypt" />
          </el-form-item>
          
          <el-form-item label="备份密码" v-if="backupForm.encrypt">
            <el-input 
              v-model="backupForm.password" 
              type="password" 
              placeholder="输入备份密码"
              show-password
            />
          </el-form-item>
          
          <el-form-item>
            <el-button 
              type="primary" 
              @click="startBackup"
              :loading="backupLoading"
              :disabled="!backupForm.path"
            >
              开始备份
            </el-button>
          </el-form-item>
        </el-form>
        
        <el-progress 
          v-if="backupProgress > 0"
          :percentage="backupProgress"
          :status="backupStatus"
          class="mt-4"
        />
      </div>
      
      <!-- 恢复区域 -->
      <div class="bg-white rounded-lg shadow-md p-6">
        <h3 class="text-xl font-semibold mb-4 flex items-center">
          <el-icon class="mr-2"><Upload /></el-icon>
          设备恢复
        </h3>
        
        <el-form :model="restoreForm" label-width="100px">
          <el-form-item label="备份路径">
            <el-input 
              v-model="restoreForm.path" 
              placeholder="选择备份文件路径"
              readonly
            >
              <template #append>
                <el-button @click="selectRestorePath">选择路径</el-button>
              </template>
            </el-input>
          </el-form-item>
          
          <el-form-item label="备份密码">
            <el-input 
              v-model="restoreForm.password" 
              type="password" 
              placeholder="输入备份密码（如果有）"
              show-password
            />
          </el-form-item>
          
          <el-form-item>
            <el-button 
              type="danger" 
              @click="startRestore"
              :loading="restoreLoading"
              :disabled="!restoreForm.path"
            >
              开始恢复
            </el-button>
          </el-form-item>
        </el-form>
        
        <el-progress 
          v-if="restoreProgress > 0"
          :percentage="restoreProgress"
          :status="restoreStatus"
          class="mt-4"
        />
      </div>
    </div>
    
    <!-- 备份历史 -->
    <div class="bg-white rounded-lg shadow-md p-6 mt-6">
      <h3 class="text-xl font-semibold mb-4">备份历史</h3>
      
      <el-table :data="backupHistory" stripe style="width: 100%">
        <el-table-column prop="date" label="备份时间" width="180" />
        <el-table-column prop="path" label="备份路径" />
        <el-table-column prop="size" label="大小" width="120" />
        <el-table-column prop="encrypted" label="是否加密" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.encrypted ? 'success' : 'info'">
              {{ scope.row.encrypted ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150">
          <template #default="scope">
            <el-button size="small" @click="restoreFromHistory(scope.row)">
              恢复
            </el-button>
            <el-button size="small" type="danger" @click="deleteBackup(scope.row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft, Download, Upload } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()

const backupForm = ref({
  path: '',
  encrypt: false,
  password: ''
})

const restoreForm = ref({
  path: '',
  password: ''
})

const backupLoading = ref(false)
const restoreLoading = ref(false)
const backupProgress = ref(0)
const restoreProgress = ref(0)
const backupStatus = ref('')
const restoreStatus = ref('')

const backupHistory = ref([
  {
    date: '2023-12-01 14:30:00',
    path: '/Users/user/Backups/iPhone_20231201_143000',
    size: '2.5GB',
    encrypted: true
  },
  {
    date: '2023-11-28 10:15:00',
    path: '/Users/user/Backups/iPhone_20231128_101500',
    size: '2.3GB',
    encrypted: false
  }
])

const selectBackupPath = () => {
  // 这里应该调用系统文件选择对话框
  ElMessage.info('请选择备份保存路径')
  // 模拟选择路径
  backupForm.value.path = '/Users/user/Backups/iPhone_' + new Date().toISOString().slice(0, 19).replace(/[-:]/g, '').replace('T', '_')
}

const selectRestorePath = () => {
  // 这里应该调用系统文件选择对话框
  ElMessage.info('请选择备份文件路径')
  // 模拟选择路径
  restoreForm.value.path = '/Users/user/Backups/iPhone_20231201_143000'
}

const startBackup = async () => {
  if (!backupForm.value.path) {
    ElMessage.error('请选择备份路径')
    return
  }
  
  if (backupForm.value.encrypt && !backupForm.value.password) {
    ElMessage.error('请输入备份密码')
    return
  }
  
  try {
    backupLoading.value = true
    backupProgress.value = 0
    backupStatus.value = ''
    
    // 模拟备份进度
    const interval = setInterval(() => {
      backupProgress.value += 10
      if (backupProgress.value >= 100) {
        clearInterval(interval)
        backupStatus.value = 'success'
        ElMessage.success('备份完成')
        backupLoading.value = false
        
        // 添加到备份历史
        backupHistory.value.unshift({
          date: new Date().toLocaleString(),
          path: backupForm.value.path,
          size: '2.4GB',
          encrypted: backupForm.value.encrypt
        })
      }
    }, 500)
    
    // 这里应该调用后端API
    // await window.go.main.App.BackupDevice(route.params.id, backupForm.value.path, backupForm.value.encrypt, backupForm.value.password)
    
  } catch (error) {
    backupLoading.value = false
    backupStatus.value = 'exception'
    ElMessage.error('备份失败: ' + error.message)
  }
}

const startRestore = async () => {
  if (!restoreForm.value.path) {
    ElMessage.error('请选择备份路径')
    return
  }
  
  const result = await ElMessageBox.confirm(
    '恢复操作将覆盖设备上的所有数据，是否继续？',
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).catch(() => false)
  
  if (!result) return
  
  try {
    restoreLoading.value = true
    restoreProgress.value = 0
    restoreStatus.value = ''
    
    // 模拟恢复进度
    const interval = setInterval(() => {
      restoreProgress.value += 8
      if (restoreProgress.value >= 100) {
        clearInterval(interval)
        restoreStatus.value = 'success'
        ElMessage.success('恢复完成')
        restoreLoading.value = false
      }
    }, 600)
    
    // 这里应该调用后端API
    // await window.go.main.App.RestoreDevice(route.params.id, restoreForm.value.path, restoreForm.value.password)
    
  } catch (error) {
    restoreLoading.value = false
    restoreStatus.value = 'exception'
    ElMessage.error('恢复失败: ' + error.message)
  }
}

const restoreFromHistory = (backup) => {
  restoreForm.value.path = backup.path
  if (backup.encrypted) {
    ElMessage.info('该备份已加密，请输入密码')
  }
}

const deleteBackup = async (backup) => {
  const result = await ElMessageBox.confirm(
    `确定要删除备份 "${backup.path}" 吗？`,
    '确认删除',
    {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).catch(() => false)
  
  if (result) {
    const index = backupHistory.value.indexOf(backup)
    if (index > -1) {
      backupHistory.value.splice(index, 1)
      ElMessage.success('备份已删除')
    }
  }
}

const goBack = () => {
  router.push('/')
}

onMounted(() => {
  // 初始化页面数据
})
</script>