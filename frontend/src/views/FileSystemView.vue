<template>
  <div>
    <div class="flex items-center mb-4">
      <el-button @click="goBack" type="default" size="small" class="mr-2">
        <el-icon><ArrowLeft /></el-icon>
        返回
      </el-button>
      <h2 class="text-2xl font-bold">文件系统</h2>
    </div>
    
    <div class="bg-white rounded-lg shadow-md p-6">
      <!-- 工具栏 -->
      <div class="flex items-center justify-between mb-4">
        <div class="flex items-center">
          <el-breadcrumb separator="/">
            <el-breadcrumb-item @click="navigateTo('/')">根目录</el-breadcrumb-item>
            <el-breadcrumb-item 
              v-for="(path, index) in breadcrumbs" 
              :key="index"
              @click="navigateTo(path.fullPath)"
            >
              {{ path.name }}
            </el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        
        <div class="flex items-center space-x-2">
          <el-button size="small" @click="refreshFiles">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
          <el-button size="small" @click="uploadFile">
            <el-icon><Upload /></el-icon>
            上传
          </el-button>
          <el-button size="small" @click="createFolder">
            <el-icon><FolderAdd /></el-icon>
            新建文件夹
          </el-button>
        </div>
      </div>
      
      <!-- 文件列表 -->
      <el-table 
        :data="files" 
        stripe 
        style="width: 100%"
        @row-dblclick="handleRowDoubleClick"
        v-loading="loading"
      >
        <el-table-column width="50">
          <template #default="scope">
            <el-icon v-if="scope.row.type === 'folder'" class="text-yellow-500">
              <Folder />
            </el-icon>
            <el-icon v-else class="text-blue-500">
              <Document />
            </el-icon>
          </template>
        </el-table-column>
        
        <el-table-column prop="name" label="名称" min-width="200">
          <template #default="scope">
            <span 
              class="cursor-pointer hover:text-blue-600"
              @click="handleFileClick(scope.row)"
            >
              {{ scope.row.name }}
            </span>
          </template>
        </el-table-column>
        
        <el-table-column prop="size" label="大小" width="120">
          <template #default="scope">
            {{ scope.row.type === 'folder' ? '-' : formatFileSize(scope.row.size) }}
          </template>
        </el-table-column>
        
        <el-table-column prop="modifiedTime" label="修改时间" width="180" />
        
        <el-table-column label="操作" width="200">
          <template #default="scope">
            <el-button size="small" @click="downloadFile(scope.row)" v-if="scope.row.type === 'file'">
              下载
            </el-button>
            <el-button size="small" @click="renameFile(scope.row)">
              重命名
            </el-button>
            <el-button size="small" type="danger" @click="deleteFile(scope.row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
    
    <!-- 上传对话框 -->
    <el-dialog v-model="uploadDialogVisible" title="上传文件" width="500px">
      <el-upload
        class="upload-demo"
        drag
        :auto-upload="false"
        :file-list="uploadFileList"
        :on-change="handleFileChange"
        multiple
      >
        <el-icon class="el-icon--upload"><upload-filled /></el-icon>
        <div class="el-upload__text">
          将文件拖到此处，或<em>点击上传</em>
        </div>
      </el-upload>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="uploadDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="confirmUpload" :loading="uploading">
            上传
          </el-button>
        </span>
      </template>
    </el-dialog>
    
    <!-- 重命名对话框 -->
    <el-dialog v-model="renameDialogVisible" title="重命名" width="400px">
      <el-form>
        <el-form-item label="新名称">
          <el-input v-model="newFileName" placeholder="输入新的文件名" />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="renameDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="confirmRename">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  ArrowLeft, 
  Refresh, 
  Upload, 
  FolderAdd, 
  Folder, 
  Document,
  UploadFilled 
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()

const currentPath = ref('/')
const files = ref([])
const loading = ref(false)
const uploadDialogVisible = ref(false)
const renameDialogVisible = ref(false)
const uploadFileList = ref([])
const uploading = ref(false)
const newFileName = ref('')
const currentRenameFile = ref(null)

// 面包屑导航
const breadcrumbs = computed(() => {
  if (currentPath.value === '/') return []
  
  const paths = currentPath.value.split('/').filter(p => p)
  const result = []
  let fullPath = ''
  
  for (const path of paths) {
    fullPath += '/' + path
    result.push({
      name: path,
      fullPath: fullPath
    })
  }
  
  return result
})

// 模拟文件数据
const mockFiles = {
  '/': [
    { name: 'Documents', type: 'folder', size: 0, modifiedTime: '2023-12-01 10:30:00' },
    { name: 'Photos', type: 'folder', size: 0, modifiedTime: '2023-12-01 09:15:00' },
    { name: 'Music', type: 'folder', size: 0, modifiedTime: '2023-11-30 14:20:00' },
    { name: 'config.plist', type: 'file', size: 2048, modifiedTime: '2023-12-01 08:45:00' }
  ],
  '/Documents': [
    { name: 'readme.txt', type: 'file', size: 1024, modifiedTime: '2023-12-01 10:30:00' },
    { name: 'data.json', type: 'file', size: 4096, modifiedTime: '2023-12-01 09:20:00' }
  ],
  '/Photos': [
    { name: 'IMG_001.jpg', type: 'file', size: 2048576, modifiedTime: '2023-12-01 09:15:00' },
    { name: 'IMG_002.jpg', type: 'file', size: 1536000, modifiedTime: '2023-12-01 09:10:00' }
  ]
}

const loadFiles = async (path = '/') => {
  loading.value = true
  
  try {
    // 模拟API调用延迟
    await new Promise(resolve => setTimeout(resolve, 500))
    
    // 这里应该调用后端API获取文件列表
    // const fileList = await window.go.main.App.ListFiles(route.params.id, path)
    
    files.value = mockFiles[path] || []
    currentPath.value = path
    
  } catch (error) {
    ElMessage.error('加载文件列表失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

const navigateTo = (path) => {
  loadFiles(path)
}

const handleFileClick = (file) => {
  if (file.type === 'folder') {
    const newPath = currentPath.value === '/' ? `/${file.name}` : `${currentPath.value}/${file.name}`
    navigateTo(newPath)
  }
}

const handleRowDoubleClick = (row) => {
  handleFileClick(row)
}

const refreshFiles = () => {
  loadFiles(currentPath.value)
}

const uploadFile = () => {
  uploadDialogVisible.value = true
  uploadFileList.value = []
}

const createFolder = async () => {
  const { value: folderName } = await ElMessageBox.prompt('请输入文件夹名称', '新建文件夹', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
  }).catch(() => ({ value: null }))
  
  if (folderName) {
    // 这里应该调用后端API创建文件夹
    ElMessage.success(`文件夹 "${folderName}" 创建成功`)
    refreshFiles()
  }
}

const handleFileChange = (file, fileList) => {
  uploadFileList.value = fileList
}

const confirmUpload = async () => {
  if (uploadFileList.value.length === 0) {
    ElMessage.warning('请选择要上传的文件')
    return
  }
  
  uploading.value = true
  
  try {
    // 这里应该调用后端API上传文件
    await new Promise(resolve => setTimeout(resolve, 2000))
    
    ElMessage.success('文件上传成功')
    uploadDialogVisible.value = false
    refreshFiles()
    
  } catch (error) {
    ElMessage.error('文件上传失败: ' + error.message)
  } finally {
    uploading.value = false
  }
}

const downloadFile = async (file) => {
  try {
    ElMessage.info(`开始下载 "${file.name}"`)
    // 这里应该调用后端API下载文件
    
  } catch (error) {
    ElMessage.error('文件下载失败: ' + error.message)
  }
}

const renameFile = (file) => {
  currentRenameFile.value = file
  newFileName.value = file.name
  renameDialogVisible.value = true
}

const confirmRename = async () => {
  if (!newFileName.value.trim()) {
    ElMessage.warning('请输入新的文件名')
    return
  }
  
  try {
    // 这里应该调用后端API重命名文件
    ElMessage.success('重命名成功')
    renameDialogVisible.value = false
    refreshFiles()
    
  } catch (error) {
    ElMessage.error('重命名失败: ' + error.message)
  }
}

const deleteFile = async (file) => {
  const result = await ElMessageBox.confirm(
    `确定要删除 "${file.name}" 吗？`,
    '确认删除',
    {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).catch(() => false)
  
  if (result) {
    try {
      // 这里应该调用后端API删除文件
      ElMessage.success('删除成功')
      refreshFiles()
      
    } catch (error) {
      ElMessage.error('删除失败: ' + error.message)
    }
  }
}

const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 B'
  
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const goBack = () => {
  router.push('/')
}

onMounted(() => {
  loadFiles()
})
</script>