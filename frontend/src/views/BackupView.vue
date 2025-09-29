<template>
  <div class="min-h-screen bg-gray-50 p-4">
    <div class="max-w-7xl mx-auto">
      <!-- 页面标题 -->
      <div class="flex items-center mb-6">
        <el-button @click="goBack" type="default" size="small" class="mr-2">
          <el-icon><ArrowLeft /></el-icon>
          返回
        </el-button>
        <div>
          <h1 class="text-2xl font-bold text-gray-900">设备备份与恢复</h1>
          <p class="text-gray-500 mt-1">管理您的设备备份，创建新备份或恢复现有备份</p>
        </div>
      </div>

      <!-- 设备选择 -->
      <div class="bg-white rounded-lg shadow-sm border p-4 mb-6">
        <h2 class="text-lg font-medium text-gray-900 mb-4">选择设备</h2>
        <div v-if="isLoadingDevices" class="flex justify-center py-4">
          <el-skeleton style="width: 100%" :rows="3" animated />
        </div>
        <div v-else-if="devices.length === 0" class="text-center py-8">
          <div class="text-gray-400 mb-2">
            <i class="el-icon-mobile-phone text-4xl"></i>
          </div>
          <p class="text-gray-500">未检测到已连接的设备</p>
          <el-button type="primary" size="small" class="mt-4" @click="refreshDevices">刷新设备列表</el-button>
        </div>
        <div v-else>
          <el-radio-group v-model="selectedDeviceUDID" class="w-full" @change="onDeviceSelected">
            <div v-for="device in devices" :key="device.udid" class="mb-3 last:mb-0">
              <el-radio :label="device.udid" border class="w-full">
                <div class="flex items-center">
                  <div class="flex-shrink-0">
                    <i class="el-icon-mobile-phone text-2xl text-gray-600"></i>
                  </div>
                  <div class="ml-3">
                    <div class="text-sm font-medium text-gray-900">{{ device.name }}</div>
                    <div class="text-xs text-gray-500">{{ device.model }} | {{ device.udid }}</div>
                  </div>
                </div>
              </el-radio>
            </div>
          </el-radio-group>
          
          <!-- 备份加密状态显示 -->
          <div v-if="selectedDeviceUDID" class="mt-4 p-3 bg-gray-50 rounded-lg">
            <div class="flex items-center justify-between">
              <div class="flex items-center">
                <span class="text-sm font-medium text-gray-700 mr-2">备份加密状态:</span>
                <el-tag 
                  :type="deviceBackupEncrypted ? 'success' : 'info'" 
                  size="small"
                  :loading="isCheckingEncryption"
                >
                  {{ isCheckingEncryption ? '检查中...' : (deviceBackupEncrypted ? '已加密' : '未加密') }}
                </el-tag>
              </div>
              <div class="flex space-x-2">
                <el-button 
                  v-if="deviceBackupEncrypted" 
                  size="small" 
                  type="warning"
                  :disabled="isCheckingEncryption"
                  @click="toggleBackupEncryption(false)"
                >
                  取消加密
                </el-button>
                <el-button 
                  v-else 
                  size="small" 
                  type="primary"
                  :disabled="isCheckingEncryption"
                  @click="toggleBackupEncryption(true)"
                >
                  设置加密
                </el-button>
                <el-button 
                  size="small" 
                  :disabled="isCheckingEncryption"
                  @click="checkBackupEncryptionStatus"
                >
                  刷新状态
                </el-button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 备份操作 -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
        <!-- 创建备份 -->
        <div class="bg-white rounded-lg shadow-sm border p-4">
          <h2 class="text-lg font-medium text-gray-900 mb-4">创建新备份</h2>
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">备份位置</label>
              <div class="flex">
                <el-input v-model="backupPath" placeholder="选择备份保存位置" class="flex-grow" />
                <el-button class="ml-2" @click="selectBackupPath">浏览...</el-button>
              </div>
            </div>
            <div>
              <el-checkbox v-model="encryptBackup">加密备份</el-checkbox>
              <el-input 
                v-if="encryptBackup" 
                v-model="backupPassword" 
                type="password" 
                placeholder="输入加密密码" 
                class="mt-2"
                show-password
              />
            </div>
            <div>
              <el-button 
                type="primary" 
                :disabled="!canCreateBackup" 
                :loading="isCreatingBackup"
                @click="createBackup"
                class="w-full"
              >
                创建备份
              </el-button>
            </div>
          </div>
          
          <!-- 备份进度 -->
          <div v-if="backupProgress" class="mt-4">
            <div class="flex justify-between text-sm mb-1">
              <span>备份进度</span>
              <span>{{ Math.round(backupProgress.progress) }}%</span>
            </div>
            <el-progress 
              :percentage="backupProgress.progress" 
              :status="backupProgressStatus"
            />
            <div class="text-xs text-gray-500 mt-1">
              {{ backupStatusText }}
            </div>
          </div>
        </div>

        <!-- 恢复备份 -->
        <div class="bg-white rounded-lg shadow-sm border p-4">
          <h2 class="text-lg font-medium text-gray-900 mb-4">恢复备份</h2>
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">备份位置</label>
              <div class="flex">
                <el-input v-model="restorePath" placeholder="选择要恢复的备份" class="flex-grow" />
                <el-button class="ml-2" @click="selectRestorePath">浏览...</el-button>
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">密码（如果备份已加密）</label>
              <el-input 
                v-model="restorePassword" 
                type="password" 
                placeholder="输入备份密码" 
                show-password
              />
            </div>
            <div>
              <el-button 
                type="success" 
                :disabled="!canRestoreBackup" 
                :loading="isRestoring"
                @click="restoreBackup"
                class="w-full"
              >
                恢复备份
              </el-button>
            </div>
          </div>
        </div>
      </div>

      <!-- 备份历史 -->
      <div class="bg-white rounded-lg shadow-sm border p-4">
        <div class="flex justify-between items-center mb-4">
          <h2 class="text-lg font-medium text-gray-900">备份历史</h2>
          <el-button size="small" @click="refreshBackups" :loading="isLoadingBackups">刷新</el-button>
        </div>
        
        <div v-if="isLoadingBackups" class="flex justify-center py-8">
          <el-skeleton style="width: 100%" :rows="5" animated />
        </div>
        
        <div v-else-if="backups.length === 0" class="text-center py-8">
          <div class="text-gray-400 mb-2">
            <i class="el-icon-folder-opened text-4xl"></i>
          </div>
          <p class="text-gray-500">暂无备份记录</p>
        </div>
        
        <el-table v-else :data="backups" style="width: 100%">
          <el-table-column prop="deviceName" label="设备名称" />
          <el-table-column prop="iosVersion" label="iOS版本" />
          <el-table-column label="备份大小">
            <template #default="scope">
              {{ formatSize(scope.row.size) }}
            </template>
          </el-table-column>
          <el-table-column label="创建时间">
            <template #default="scope">
              {{ formatDate(scope.row.createdAt) }}
            </template>
          </el-table-column>
          <el-table-column label="加密状态">
            <template #default="scope">
              <el-tag :type="scope.row.isEncrypted ? 'success' : 'info'" size="small">
                {{ scope.row.isEncrypted ? '已加密' : '未加密' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="200">
            <template #default="scope">
              <el-button size="small" @click="selectBackupToRestore(scope.row)">选择恢复</el-button>
              <el-button size="small" type="danger" @click="confirmDeleteBackup(scope.row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'
import { useDeviceStore } from '../stores/device'

const deviceStore = useDeviceStore()
const router = useRouter()

// 设备相关
const devices = computed(() => deviceStore.devices)
const isLoadingDevices = ref(false)
const selectedDeviceUDID = ref('')
const deviceBackupEncrypted = ref(false) // 设备备份加密状态
const isCheckingEncryption = ref(false) // 是否正在检查加密状态

// 备份相关
const backupPath = ref('')
const encryptBackup = ref(false)
const backupPassword = ref('')
const isCreatingBackup = ref(false)
const backupProgress = ref(null)
const currentBackupID = ref(null)
const backupProgressTimer = ref(null)

// 恢复相关
const restorePath = ref('')
const restorePassword = ref('')
const isRestoring = ref(false)

// 备份列表
const backups = ref([])
const isLoadingBackups = ref(false)
const backupBaseDir = ref('') // 初始为空，将从后端获取

// 计算属性
const canCreateBackup = computed(() => {
  return selectedDeviceUDID.value && backupPath.value && (!encryptBackup.value || backupPassword.value)
})

const canRestoreBackup = computed(() => {
  return selectedDeviceUDID.value && restorePath.value
})

const backupProgressStatus = computed(() => {
  if (!backupProgress.value) return ''
  
  switch (backupProgress.value.status) {
    case 'completed':
      return 'success'
    case 'failed':
      return 'exception'
    default:
      return ''
  }
})

const backupStatusText = computed(() => {
  if (!backupProgress.value) return ''
  
  switch (backupProgress.value.status) {
    case 'preparing':
      return '正在准备备份...'
    case 'backing_up':
      return `正在备份: ${backupProgress.value.currentFile || ''}`
    case 'finishing':
      return '正在完成备份...'
    case 'completed':
      return '备份完成'
    case 'failed':
      return `备份失败: ${backupProgress.value.error || '未知错误'}`
    default:
      return '未知状态'
  }
})

// 生命周期钩子
onMounted(async () => {
  // 获取默认备份目录
  try {
    if (window.go && window.go.main && window.go.main.App) {
      const defaultDir = await window.go.main.App.GetDefaultBackupDir()
      if (defaultDir) {
        backupBaseDir.value = defaultDir
        console.log('获取到默认备份目录:', backupBaseDir.value)
      } else {
        // 如果后端返回空，使用临时目录
        backupBaseDir.value = './backups'
        console.warn('后端返回的备份目录为空，使用临时目录')
      }
    } else {
      // 如果 window.go 不存在，使用临时目录
      backupBaseDir.value = './backups'
      console.warn('window.go 不存在，使用临时目录')
    }
  } catch (error) {
    console.error('获取默认备份目录失败:', error)
    // 出错时使用临时目录
    backupBaseDir.value = './backups'
  }
  
  await refreshDevices()
  // 延迟加载备份列表，确保UI先渲染
  setTimeout(async () => {
    await refreshBackups()
  }, 500)
})

onUnmounted(() => {
  if (backupProgressTimer.value) {
    clearInterval(backupProgressTimer.value)
  }
})

// 方法
function goBack() {
  router.push('/')
}

async function refreshDevices() {
  isLoadingDevices.value = true
  try {
    await deviceStore.fetchDevices()
    if (devices.value.length > 0 && !selectedDeviceUDID.value) {
      selectedDeviceUDID.value = devices.value[0].udid
      // 自动检查第一个设备的加密状态
      await checkBackupEncryptionStatus()
    }
  } catch (error) {
    console.error('获取设备列表失败:', error)
    ElMessage.error('获取设备列表失败')
  } finally {
    isLoadingDevices.value = false
  }
}

// 检查设备备份加密状态
async function checkBackupEncryptionStatus() {
  if (!selectedDeviceUDID.value) {
    ElMessage.warning('请先选择设备')
    return
  }
  
  isCheckingEncryption.value = true
  try {
    const result = await window.go.main.App.CheckBackupEncryptionStatus(selectedDeviceUDID.value)
    deviceBackupEncrypted.value = result
    console.log(`设备备份${result ? '已加密' : '未加密'}`)
  } catch (error) {
    console.error('检查备份加密状态失败:', error)
    ElMessage.error('检查备份加密状态失败')
  } finally {
    isCheckingEncryption.value = false
  }
}

// 设备选择后的回调
async function onDeviceSelected() {
  if (selectedDeviceUDID.value) {
    await checkBackupEncryptionStatus()
  }
}

// 切换备份加密状态
async function toggleBackupEncryption(enable) {
  if (!selectedDeviceUDID.value) {
    ElMessage.warning('请先选择设备')
    return
  }
  
  try {
    let password = ''
    if (enable) {
      // 启用加密需要密码
      const { value } = await ElMessageBox.prompt('请输入备份加密密码', '设置备份加密', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        inputType: 'password',
        inputValidator: (value) => {
          if (!value || value.length < 4) {
            return '密码长度至少4位'
          }
          return true
        }
      })
      password = value
    } else {
      // 取消加密需要输入当前密码
      const { value } = await ElMessageBox.prompt('请输入当前的备份加密密码以取消加密', '取消备份加密', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        inputType: 'password',
        inputValidator: (value) => {
          if (!value) {
            return '请输入当前的加密密码'
          }
          return true
        }
      })
      password = value
    }
    
    isCheckingEncryption.value = true
    await window.go.main.App.SetBackupEncryption(selectedDeviceUDID.value, enable, password)
    
    // 重新检查状态
    await checkBackupEncryptionStatus()
    
    ElMessage.success(`备份加密${enable ? '启用' : '取消'}成功`)
  } catch (error) {
    if (error === 'cancel') {
      return // 用户取消操作
    }
    console.error('设置备份加密失败:', error)
    ElMessage.error(`设置备份加密失败: ${error.message || error}`)
  } finally {
    isCheckingEncryption.value = false
  }
}

async function selectBackupPath() {
  try {
    // 检查 window.go 是否存在
    if (!window.go || !window.go.main || !window.go.main.App) {
      console.error('window.go.main.App 不存在，可能是在开发环境中运行');
      // 在开发环境中使用默认路径
      backupPath.value = `${backupBaseDir.value}/${new Date().toISOString().replace(/:/g, '-')}`;
      return;
    }
    
    // 调用后端的文件选择对话框
    const path = await window.go.main.App.OpenDirectoryDialog(
      "选择备份保存位置", 
      backupBaseDir.value
    );
    
    // 如果用户选择了路径（没有取消对话框）
    if (path) {
      backupPath.value = path;
    } else {
      // 如果用户取消了选择，使用默认路径
      const defaultDir = await window.go.main.App.GetDefaultBackupDir();
      backupPath.value = defaultDir || backupBaseDir.value;
    }
  } catch (error) {
    console.error('选择备份路径失败:', error);
    ElMessage.error(`选择备份路径失败: ${error.message || error}`);
    // 出错时使用默认路径
    try {
      const defaultDir = await window.go.main.App.GetDefaultBackupDir();
      backupPath.value = defaultDir || `${backupBaseDir.value}/${new Date().toISOString().replace(/:/g, '-')}`;
    } catch (e) {
      backupPath.value = `${backupBaseDir.value}/${new Date().toISOString().replace(/:/g, '-')}`;
    }
  }
}

async function selectRestorePath() {
  try {
    // 检查 window.go 是否存在
    if (!window.go || !window.go.main || !window.go.main.App) {
      console.error('window.go.main.App 不存在，可能是在开发环境中运行');
      // 在开发环境中使用默认路径
      restorePath.value = backupBaseDir.value;
      return;
    }
    
    // 调用后端的文件选择对话框
    const path = await window.go.main.App.OpenDirectoryDialog(
      "选择要恢复的备份", 
      backupBaseDir.value
    );
    
    // 如果用户选择了路径（没有取消对话框）
    if (path) {
      restorePath.value = path;
    } else {
      // 如果用户取消了选择，使用默认路径
      const defaultDir = await window.go.main.App.GetDefaultBackupDir();
      restorePath.value = defaultDir || backupBaseDir.value;
    }
  } catch (error) {
    console.error('选择恢复路径失败:', error);
    ElMessage.error(`选择恢复路径失败: ${error.message || error}`);
    // 出错时使用默认路径
    try {
      const defaultDir = await window.go.main.App.GetDefaultBackupDir();
      restorePath.value = defaultDir || backupBaseDir.value;
    } catch (e) {
      restorePath.value = backupBaseDir.value;
    }
  }
}

async function createBackup() {
  if (!selectedDeviceUDID.value) {
    ElMessage.warning('请先选择设备')
    return
  }
  
  if (!backupPath.value) {
    ElMessage.warning('请选择备份保存位置')
    return
  }
  
  if (encryptBackup.value && !backupPassword.value) {
    ElMessage.warning('请输入备份密码')
    return
  }
  
  isCreatingBackup.value = true
  backupProgress.value = null
  
  try {
    // 调用后端创建备份
    const backupID = await window.go.main.App.BackupDevice(
      selectedDeviceUDID.value,
      backupPath.value,
      encryptBackup.value,
      backupPassword.value
    )
    
    currentBackupID.value = backupID
    ElMessage.success('备份已开始，请等待完成')
    
    // 开始监控备份进度
    startBackupProgressMonitoring()
    
  } catch (error) {
    console.error('创建备份失败:', error)
    ElMessage.error(`创建备份失败: ${error.message || error}`)
    isCreatingBackup.value = false
  }
}

function startBackupProgressMonitoring() {
  if (backupProgressTimer.value) {
    clearInterval(backupProgressTimer.value)
  }
  
  backupProgressTimer.value = setInterval(async () => {
    try {
      if (!currentBackupID.value) return
      
      const progress = await window.go.main.App.GetBackupProgress(currentBackupID.value)
      backupProgress.value = progress
      
      if (progress.status === 'completed' || progress.status === 'failed') {
        clearInterval(backupProgressTimer.value)
        isCreatingBackup.value = false
        currentBackupID.value = null
        
        if (progress.status === 'completed') {
          ElMessage.success('备份完成')
          // 刷新备份列表
          await refreshBackups()
        } else {
          ElMessage.error(`备份失败: ${progress.error || '未知错误'}`)
        }
      }
    } catch (error) {
      console.error('获取备份进度失败:', error)
    }
  }, 1000)
}

async function restoreBackup() {
  if (!selectedDeviceUDID.value) {
    ElMessage.warning('请先选择设备')
    return
  }
  
  if (!restorePath.value) {
    ElMessage.warning('请选择要恢复的备份')
    return
  }
  
  try {
    await ElMessageBox.confirm(
      '恢复备份将覆盖设备上的所有数据，确定要继续吗？',
      '确认恢复',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    isRestoring.value = true
    
    await window.go.main.App.RestoreDevice(
      selectedDeviceUDID.value,
      restorePath.value,
      restorePassword.value
    )
    
    ElMessage.success('恢复完成')
    
  } catch (error) {
    if (error === 'cancel') {
      return
    }
    console.error('恢复备份失败:', error)
    ElMessage.error(`恢复备份失败: ${error.message || error}`)
  } finally {
    isRestoring.value = false
  }
}

async function refreshBackups() {
  isLoadingBackups.value = true
  try {
    const result = await window.go.main.App.ListBackups(backupBaseDir.value)
    console.log('获取到的备份列表:', result)
    
    if (Array.isArray(result)) {
      // 延迟处理，避免阻塞UI
      setTimeout(async () => {
        const updatedBackups = []
        
        for (const backup of result) {
          try {
            const backupPath = backup.backupPath || backup.backup_path
            console.log('正在获取备份详细信息，路径:', backupPath);
            
            // 从备份目录直接读取最新的备份信息
            // 注意：这里可能会出现路径问题，确保路径正确
            const fixedPath = backupPath.replace(/\\/g, '/');
            console.log('修正后的路径:', fixedPath);
            
            const backupInfo = await window.go.main.App.GetBackupInfo(fixedPath);
            console.log('获取到备份详细信息:', backupInfo);
            
            // 使用最新的备份信息
            updatedBackups.push({
              ...backup,
              deviceName: backupInfo.deviceName || backup.device_name || '未知设备',
              iosVersion: backupInfo.iosVersion || backup.ios_version || '未知版本',
              createdAt: backupInfo.createdAt || backup.created_at || new Date().toISOString()
            });
          } catch (error) {
            console.error('获取备份详细信息失败:', error);
            // 如果获取详细信息失败，使用原始信息
            updatedBackups.push({
              ...backup,
              deviceName: backup.device_name || '未知设备',
              iosVersion: backup.ios_version || '未知版本',
              createdAt: backup.created_at || new Date().toISOString()
            });
          }
        }
        
        backups.value = updatedBackups;
        console.log('处理后的备份列表:', backups.value);
      }, 100);
    } else {
      console.error('备份列表不是数组:', result);
      backups.value = [];
    }
  } catch (error) {
    console.error('获取备份列表失败:', error);
    ElMessage.error('获取备份列表失败');
    backups.value = [];
  } finally {
    isLoadingBackups.value = false;
  }
}

function selectBackupToRestore(backup) {
  restorePath.value = backup.backupPath || backup.backup_path;
  if (backup.isEncrypted || backup.is_encrypted) {
    ElMessage.warning('此备份已加密，请输入密码');
  }
}

async function confirmDeleteBackup(backup) {
  try {
    const deviceName = backup.deviceName || backup.device_name || '未知设备';
    await ElMessageBox.confirm(
      `确定要删除"${deviceName}"的备份吗？此操作不可逆。`,
      '删除备份',
      {
        confirmButtonText: '删除',
        cancelButtonText: '取消',
        type: 'danger'
      }
    )
    
    await deleteBackup(backup)
  } catch (error) {
    if (error === 'cancel') {
      return
    }
    console.error('删除备份失败:', error)
  }
}

async function deleteBackup(backup) {
  try {
    const backupPath = backup.backupPath || backup.backup_path
    await window.go.main.App.DeleteBackup(backupPath)
    
    ElMessage.success('备份删除成功')
    await refreshBackups()
  } catch (error) {
    console.error('删除备份失败:', error)
    ElMessage.error(`删除备份失败: ${error.message || error}`)
  }
}

function formatSize(bytes) {
  if (!bytes) return '0 B'
  
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(1024))
  
  return Math.round(bytes / Math.pow(1024, i) * 100) / 100 + ' ' + sizes[i]
}

function formatDate(dateString) {
  if (!dateString) return '未知时间'
  
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}
</script>
