# MyiToolsApp

一个使用 **Go + Wails + Vue** 构建的桌面应用示例，功能类似 **爱思助手**
设备信息展示页面
增加备份/恢复、安装 IPA 的功能
引入 Tailwind / Element Plus 美化界面
支持多设备并发操作

---

## 📂 项目结构

myitools/
├─ go.mod
├─ go.sum
├─ main.go # Wails 程序入口
├─ app.go # 应用生命周期管理
├─ device/ # 设备相关逻辑
│ ├─ device.go # 设备基础操作
│ ├─ info.go # 设备信息获取
│ └─ backup.go # 备份与恢复功能
├─ installer/ # 应用安装相关
│ ├─ ipa.go # IPA 安装处理
│ └─ utils.go # 安装工具函数
├─ filesystem/ # 文件系统操作
│ ├─ mount.go # 设备挂载
│ └─ browser.go # 文件浏览器
├─ frontend/ # Vue 前端
│ ├─ index.html
│ ├─ package.json
│ ├─ vite.config.js
│ ├─ tailwind.config.js # Tailwind 配置
│ └─ src/
│   ├─ main.js # Vue 入口
│   ├─ App.vue # 根组件
│   ├─ assets/ # 静态资源
│   │ ├─ logo.svg
│   │ └─ styles/ # 样式文件
│   ├─ components/ # 通用组件
│   │ ├─ DeviceCard.vue # 设备卡片组件
│   │ ├─ AppHeader.vue # 应用头部
│   │ └─ FileExplorer.vue # 文件浏览器组件
│   ├─ stores/ # Pinia 状态管理
│   │ ├─ device.js # 设备状态
│   │ └─ app.js # 应用状态
│   └─ views/ # 页面视图
│     ├─ HomeView.vue # 首页
│     ├─ DeviceInfoView.vue # 设备信息页
│     ├─ BackupView.vue # 备份恢复页
│     └─ FileSystemView.vue # 文件系统页

---

## 🚀 功能

### 设备管理
- [ ] 列出已连接 iOS 设备（调用 `idevice_id -l`）
- [ ] 自动检测设备连接/断开并更新界面
- [ ] 支持多设备并行操作
- [ ] 设备分组和标签管理

### 设备信息
- [ ] 获取设备详细信息（`ideviceinfo`）
  - 基本信息：设备名称、型号、iOS 版本、序列号
  - 存储信息：总容量、已用空间、可用空间
  - 电池信息：电量、健康状态
  - 网络信息：WiFi MAC、蓝牙地址

### 备份与恢复
- [ ] 完整设备备份/恢复（`idevicebackup2`）
- [ ] 应用单独备份与恢复
  - 支持选择性备份特定应用数据
  - 支持备份加密与密码管理
- [ ] 备份历史管理与自动备份计划
- [ ] 备份数据浏览与导出

### 应用管理
- [ ] 安装 IPA 应用（`ideviceinstaller`）
  - 支持拖放安装
  - 支持批量安装
- [ ] 应用列表查看与管理
  - 卸载应用
  - 查看应用详情
- [ ] 应用数据导出

### 文件系统
- [ ] 文件浏览（`ifuse` 挂载 AFC 文件系统）
  - 文件/文件夹浏览
  - 文件上传/下载
  - 文件删除/重命名
- [ ] 照片管理
  - 照片浏览与导出
  - 相册管理
- [ ] 文档管理

---

## ⚙️ 环境准备

### 系统要求
- macOS 10.15+ / Windows 10+ / Linux (Ubuntu 20.04+)
- 8GB+ RAM 推荐
- 1GB+ 可用存储空间

### 依赖工具
1. **Go 环境（1.21+）**  
   ```bash
   # 检查是否已安装
   go version
   
   # macOS 安装
   brew install go
   
   # Windows 安装
   # 下载并运行安装程序: https://golang.org/dl/
   ```

2. **Node.js 环境（18+）**
   ```bash
   # 检查是否已安装
   node -v
   
   # macOS 安装
   brew install node
   
   # Windows 安装
   # 下载并运行安装程序: https://nodejs.org/
   ```

3. **Wails CLI**
   ```bash
   # 安装 Wails CLI
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   ```

4. **libimobiledevice 工具集**
   ```bash
   # macOS 安装
   brew install --HEAD libimobiledevice
   brew install --HEAD ideviceinstaller
   brew install --HEAD ifuse
   
   # Ubuntu 安装
   sudo apt-get install libimobiledevice6 libimobiledevice-utils ideviceinstaller ifuse
   
   # Windows 安装
   # 请参考 libimobiledevice-win32 项目: https://github.com/libimobiledevice-win32/imobiledevice-net
   ```

5. **开发工具**
   - 推荐使用 Visual Studio Code
   - 安装 Go、Vue、Tailwind CSS 扩展

---

## 📦 安装与运行

### 从源码构建

1. **克隆仓库**
   ```bash
   git clone https://github.com/yourusername/myitools.git
   cd myitools
   ```

2. **安装前端依赖**
   ```bash
   cd frontend
   npm install
   cd ..
   ```

3. **开发模式运行**
   ```bash
   export PATH=$PATH:$(go env GOPATH)/bin && wails dev
   wails dev
   ```
   这将启动开发服务器并打开应用程序。

4. **构建生产版本**
   ```bash
   # macOS
   wails build -platform darwin/universal
   
   # Windows
   wails build -platform windows/amd64
   
   # Linux
   wails build -platform linux/amd64
   ```

5. **打包应用**
   ```bash
   # macOS 打包为 .app
   wails build -platform darwin/universal -package
   
   # Windows 打包为安装程序
   wails build -platform windows/amd64 -nsis
   ```

### 预编译版本

1. 从 [Releases](https://github.com/yourusername/myitools/releases) 页面下载最新版本
2. 解压或安装下载的文件
3. 运行应用程序

---

## 💻 开发指南

### 项目结构说明

- **Go 后端**：处理与 iOS 设备的通信，调用 libimobiledevice 工具
- **Vue 前端**：提供用户界面，使用 Tailwind CSS 和 Element Plus 组件库

### 添加新功能

1. 在 Go 后端添加新的设备操作函数
   ```go
   // device/example.go
   package device
   
   func NewFeature() string {
       // 实现新功能
       return "New feature result"
   }
   ```

2. 在 app.go 中导出函数供前端调用
   ```go
   // app.go
   func (a *App) DeviceNewFeature() string {
       return device.NewFeature()
   }
   ```

3. 在 Vue 前端调用后端函数
   ```js
   // 在 Vue 组件中
   import { DeviceNewFeature } from '../wailsjs/go/main/App'
   
   // 调用函数
   const result = await DeviceNewFeature()
   ```

### 贡献指南

欢迎提交 Pull Request 或创建 Issue！

1. Fork 本仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

---

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件
