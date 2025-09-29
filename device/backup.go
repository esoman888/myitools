package device

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// BackupInfo 备份信息结构体
type BackupInfo struct {
	ID          string    `json:"id"`           // 备份ID
	DeviceUDID  string    `json:"device_udid"`  // 设备UDID
	DeviceName  string    `json:"device_name"`  // 设备名称
	BackupPath  string    `json:"backup_path"`  // 备份路径
	CreatedAt   time.Time `json:"created_at"`   // 创建时间
	Size        int64     `json:"size"`         // 备份大小(字节)
	IsEncrypted bool      `json:"is_encrypted"` // 是否加密
	IOSVersion  string    `json:"ios_version"`  // iOS版本
}

// BackupProgress 备份进度结构体
type BackupProgress struct {
	Status      string  `json:"status"`       // 状态: preparing, backing_up, finishing, completed, failed
	Progress    float64 `json:"progress"`     // 进度百分比 (0-100)
	CurrentFile string  `json:"current_file"` // 当前正在备份的文件
	Error       string  `json:"error"`        // 错误信息(如果有)
}

// 全局变量，用于跟踪备份进度
var backupProgressMap = make(map[string]*BackupProgress)

// GetBackupProgress 获取备份进度
func GetBackupProgress(backupID string) *BackupProgress {
	progress, exists := backupProgressMap[backupID]
	if !exists {
		return &BackupProgress{
			Status:   "unknown",
			Progress: 0,
		}
	}
	return progress
}

// CheckBackupEncryptionStatus 检查设备备份加密状态
func CheckBackupEncryptionStatus(udid string) (bool, error) {
	if !IsDeviceConnected(udid) {
		return false, errors.New("设备未连接")
	}

	// 使用 ideviceinfo 命令查询备份加密状态
	args := []string{"-u", udid, "-q", "com.apple.mobile.backup"}
	
	fmt.Printf("执行命令检测加密状态: ideviceinfo %s\n", strings.Join(args, " "))
	
	// 创建命令
	cmd := exec.Command("ideviceinfo", args...)
	
	// 执行命令并获取输出
	output, err := cmd.CombinedOutput()
	outputStr := string(output)
	fmt.Printf("命令输出: %s\n", outputStr)
	
	if err != nil {
		fmt.Printf("执行ideviceinfo命令失败: %v\n", err)
		return false, fmt.Errorf("无法获取设备备份信息: %v", err)
	}
	
	// 检查 PasswordProtected 和 WillEncrypt 字段
	// PasswordProtected: true 表示备份需要密码保护
	// WillEncrypt: true 表示备份将被加密
	passwordProtected := strings.Contains(outputStr, "PasswordProtected: true")
	willEncrypt := strings.Contains(outputStr, "WillEncrypt: true")
	
	fmt.Printf("备份加密状态检测结果: PasswordProtected=%v, WillEncrypt=%v\n", passwordProtected, willEncrypt)
	
	// 如果任一字段为true，则认为备份已启用加密
	return passwordProtected || willEncrypt, nil
}

// SetBackupEncryption 设置备份加密状态
func SetBackupEncryption(udid string, enable bool, password string) error {
	if !IsDeviceConnected(udid) {
		return errors.New("设备未连接")
	}

	// 构建命令
	var args []string
	if enable {
		// 使用交互式方式设置加密
		args = []string{"-u", udid, "encryption", "on", "-i"}
	} else {
		// 使用交互式方式关闭加密
		args = []string{"-u", udid, "encryption", "off", "-i"}
	}

	fmt.Printf("执行命令: idevicebackup2 %s\n", strings.Join(args, " "))
	
	// 创建命令
	cmd := exec.Command("idevicebackup2", args...)
	
	// 创建管道用于输入密码
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("创建输入管道失败: %v", err)
	}
	
	// 创建管道用于获取输出
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("创建输出管道失败: %v", err)
	}
	
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("创建错误输出管道失败: %v", err)
	}
	
	// 启动命令
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("启动命令失败: %v", err)
	}
	
	// 创建一个goroutine来读取输出
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stdout.Read(buf)
			if n > 0 {
				fmt.Printf("命令输出: %s", string(buf[:n]))
			}
			if err != nil {
				break
			}
		}
	}()
	
	// 创建一个goroutine来读取错误输出
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stderr.Read(buf)
			if n > 0 {
				fmt.Printf("命令错误输出: %s", string(buf[:n]))
			}
			if err != nil {
				break
			}
		}
	}()
	
	// 等待一段时间，让命令有时间启动并输出提示
	time.Sleep(500 * time.Millisecond)
	
	// 输入密码
	if enable {
		// 启用加密需要输入两次密码
		fmt.Fprintf(stdin, "%s\n%s\n", password, password)
	} else {
		// 禁用加密只需要输入一次密码
		fmt.Fprintf(stdin, "%s\n", password)
	}
	stdin.Close()
	
	// 等待命令完成
	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("设置备份加密状态失败: %v", err)
	}
	
	return nil
}

// CreateBackup 创建设备备份
func CreateBackup(udid string, backupDir string, encrypt bool, password string) (string, error) {
	if !IsDeviceConnected(udid) {
		return "", errors.New("设备未连接")
	}

	// 确保备份目录存在
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", fmt.Errorf("创建备份目录失败: %v", err)
	}

	// 获取设备信息
	deviceInfo, err := GetDeviceInfo(udid)
	if err != nil {
		return "", fmt.Errorf("获取设备信息失败: %v", err)
	}

	// 打印设备信息用于调试
	fmt.Printf("设备信息: %+v\n", deviceInfo)
	fmt.Printf("设备名称: %s\n", deviceInfo["DeviceName"])
	fmt.Printf("iOS版本: %s\n", deviceInfo["ProductVersion"])
	fmt.Printf("用户选择加密备份: %v\n", encrypt)

	// 生成唯一的备份ID
	backupID := fmt.Sprintf("backup_%s_%d", udid, time.Now().Unix())

	// 初始化备份进度
	backupProgressMap[backupID] = &BackupProgress{
		Status:   "preparing",
		Progress: 0,
	}

	// 在后台执行备份
	go func() {
		// 使用原始UDID进行备份
		args := []string{"-u", udid, "backup", "--full", backupDir}
		if encrypt && password != "" {
			args = append(args, "--password", password)
			fmt.Printf("添加加密参数: %v\n", args)
		}

		fmt.Printf("备份命令参数: %v\n", args)
		cmd := exec.Command("idevicebackup2", args...)
		output, err := cmd.CombinedOutput()
		
		fmt.Printf("备份命令输出: %s\n", string(output))

		if err != nil {
			backupProgressMap[backupID].Status = "failed"
			backupProgressMap[backupID].Error = fmt.Sprintf("备份失败: %v - %s", err, string(output))
			return
		}

		// 更新进度为完成
		backupProgressMap[backupID].Status = "completed"
		backupProgressMap[backupID].Progress = 100

		// 备份完成后，将目录重命名，添加时间戳
		originalDir := filepath.Join(backupDir, udid)
		timestamp := time.Now().Format("20060102_150405")
		newUDID := udid + "_" + timestamp
		newDir := filepath.Join(backupDir, newUDID)
		
		// 检查原始目录是否存在
		if _, err := os.Stat(originalDir); os.IsNotExist(err) {
			fmt.Printf("备份目录不存在: %s, 错误: %v\n", originalDir, err)
			return
		}
		
		// 重命名目录
		err = os.Rename(originalDir, newDir)
		if err != nil {
			fmt.Printf("重命名备份目录失败: %v\n", err)
			return
		}
		
		fmt.Printf("备份目录已重命名: %s -> %s\n", originalDir, newDir)

		// 从Info.plist文件中读取设备信息
		infoPath := filepath.Join(newDir, "Info.plist")
		
		// 确保设备名称和iOS版本有值
		deviceName := deviceInfo["DeviceName"]
		iosVersion := deviceInfo["ProductVersion"]
		
		// 从Info.plist中提取信息
		cmd = exec.Command("plutil", "-p", infoPath)
		plistOutput, err := cmd.Output()
		if err == nil {
			outputStr := string(plistOutput)
			
			// 提取设备名称
			if idx := strings.Index(outputStr, "\"Device Name\" => "); idx != -1 {
				deviceName = extractQuotedValue(outputStr[idx+16:])
				fmt.Printf("从Info.plist获取到设备名称: %s\n", deviceName)
			}
			
			// 提取iOS版本
			if idx := strings.Index(outputStr, "\"Product Version\" => "); idx != -1 {
				iosVersion = extractQuotedValue(outputStr[idx+20:])
				fmt.Printf("从Info.plist获取到iOS版本: %s\n", iosVersion)
			}
		} else {
			fmt.Printf("读取Info.plist失败: %v\n", err)
		}
		
		// 如果仍然没有值，使用默认值
		if deviceName == "" {
			deviceName = "未知设备"
		}
		
		if iosVersion == "" {
			iosVersion = "未知版本"
		}
		
		// 检查备份是否真的加密了
		isEncrypted := isBackupEncrypted(newDir)
		fmt.Printf("备份加密状态检测结果: %v\n", isEncrypted)
		
		// 创建备份信息文件，保存设备名称和iOS版本
		backupInfoPath := filepath.Join(newDir, "MyiToolsBackupInfo.json")
		backupInfo := BackupInfo{
			ID:          backupID,
			DeviceUDID:  newUDID,
			DeviceName:  deviceName,
			BackupPath:  newDir,
			CreatedAt:   time.Now(),
			IsEncrypted: isEncrypted,
			IOSVersion:  iosVersion,
		}
		
		// 将备份信息序列化为JSON
		infoJSON := fmt.Sprintf(`{
			"id": "%s",
			"device_udid": "%s",
			"device_name": "%s",
			"backup_path": "%s",
			"created_at": "%s",
			"is_encrypted": %v,
			"ios_version": "%s"
		}`, 
		backupInfo.ID, 
		backupInfo.DeviceUDID, 
		backupInfo.DeviceName, 
		backupInfo.BackupPath, 
		backupInfo.CreatedAt.Format(time.RFC3339),
		backupInfo.IsEncrypted,
		backupInfo.IOSVersion)

		// 写入备份信息文件
		os.MkdirAll(filepath.Dir(backupInfoPath), 0755)
		os.WriteFile(backupInfoPath, []byte(infoJSON), 0644)
	}()

	return backupID, nil
}

// RestoreBackup 恢复设备备份
func RestoreBackup(udid string, backupDir string, password string) error {
	if !IsDeviceConnected(udid) {
		return errors.New("设备未连接")
	}

	// 检查备份是否加密
	isEncrypted := isBackupEncrypted(backupDir)
	fmt.Printf("备份加密状态: %v, 用户提供的密码: %v\n", isEncrypted, password != "")

	// 如果备份加密了但没有提供密码
	if isEncrypted && password == "" {
		return errors.New("此备份已加密，请提供密码")
	}

	// 如果备份没有加密但提供了密码
	if !isEncrypted && password != "" {
		fmt.Printf("警告：备份未加密，但提供了密码，将忽略密码\n")
	}

	args := []string{"-u", udid, "restore", "--full", backupDir}
	if isEncrypted && password != "" {
		args = append(args, "--password", password)
	}

	fmt.Printf("恢复命令参数: %v\n", args)
	cmd := exec.Command("idevicebackup2", args...)
	output, err := cmd.CombinedOutput()
	fmt.Printf("恢复命令输出: %s\n", string(output))
	
	if err != nil {
		return fmt.Errorf("恢复失败: %v - %s", err, string(output))
	}

	return nil
}

// ListBackups 列出所有备份
func ListBackups(backupBaseDir string) ([]BackupInfo, error) {
	backups := []BackupInfo{}
	
	// 确保备份基础目录存在
	if _, err := os.Stat(backupBaseDir); os.IsNotExist(err) {
		return backups, nil // 返回空列表而不是错误
	}
	
	// 遍历备份目录
	entries, err := os.ReadDir(backupBaseDir)
	if err != nil {
		return nil, fmt.Errorf("读取备份目录失败: %v", err)
	}
	
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		
		// 尝试读取备份信息
		udidDir := filepath.Join(backupBaseDir, entry.Name())
		infoPath := filepath.Join(udidDir, "Info.plist")
		myiToolsInfoPath := filepath.Join(udidDir, "MyiToolsBackupInfo.json")
		
		// 检查是否是有效的备份目录
		if _, err := os.Stat(infoPath); os.IsNotExist(err) {
			continue
		}
		
		// 获取备份大小
		var size int64
		filepath.Walk(udidDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if !info.IsDir() {
				size += info.Size()
			}
			return nil
		})
		
		// 尝试获取设备名称和iOS版本
		deviceName := "未知设备"
		iosVersion := "未知版本"
		createdAt := time.Now() // 默认为当前时间
		
		// 首先尝试从我们自己的备份信息文件中读取
		if data, err := os.ReadFile(myiToolsInfoPath); err == nil {
			// 解析JSON
			infoStr := string(data)
			
			// 提取设备名称
			if idx := strings.Index(infoStr, "\"device_name\""); idx != -1 {
				deviceName = extractJSONValue(infoStr[idx+14:])
			}
			
			// 提取iOS版本
			if idx := strings.Index(infoStr, "\"ios_version\""); idx != -1 {
				iosVersion = extractJSONValue(infoStr[idx+14:])
			}
			
			// 提取创建时间
			if idx := strings.Index(infoStr, "\"created_at\""); idx != -1 {
				timeStr := extractJSONValue(infoStr[idx+13:])
				if parsedTime, err := time.Parse(time.RFC3339, timeStr); err == nil {
					createdAt = parsedTime
				}
			}
		} else {
			// 如果没有我们的信息文件，则从Info.plist中提取信息
			cmd := exec.Command("plutil", "-p", infoPath)
			output, err := cmd.Output()
			if err == nil {
				outputStr := string(output)
				
				// 提取设备名称
				if idx := strings.Index(outputStr, "\"Device Name\" => "); idx != -1 {
					deviceName = extractQuotedValue(outputStr[idx+16:])
				}
				
				// 提取iOS版本
				if idx := strings.Index(outputStr, "\"Product Version\" => "); idx != -1 {
					iosVersion = extractQuotedValue(outputStr[idx+20:])
				}
				
				// 尝试从文件修改时间获取创建时间
				if fileInfo, err := os.Stat(infoPath); err == nil {
					createdAt = fileInfo.ModTime()
				}
			}
		}
		
		// 创建备份信息
		backup := BackupInfo{
			ID:          entry.Name(),
			DeviceUDID:  entry.Name(),
			DeviceName:  deviceName,
			BackupPath:  udidDir,
			CreatedAt:   createdAt,
			Size:        size,
			IsEncrypted: isBackupEncrypted(udidDir),
			IOSVersion:  iosVersion,
		}
		
		backups = append(backups, backup)
	}
	
	// 打印调试信息
	fmt.Printf("找到 %d 个备份\n", len(backups))
	for i, backup := range backups {
		fmt.Printf("备份 #%d: 设备名=%s, iOS版本=%s, 创建时间=%s\n", 
			i+1, backup.DeviceName, backup.IOSVersion, backup.CreatedAt.Format(time.RFC3339))
	}
	
	return backups, nil
}

// GetBackupInfo 从备份目录读取备份信息
func GetBackupInfo(backupPath string) (BackupInfo, error) {
	// 打印调试信息
	fmt.Printf("GetBackupInfo 被调用，路径: %s\n", backupPath)
	
	// 默认备份信息
	backup := BackupInfo{
		DeviceName: "未知设备",
		IOSVersion: "未知版本",
		CreatedAt:  time.Now(),
	}
	
	// 检查备份目录是否存在
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		fmt.Printf("备份目录不存在: %s, 错误: %v\n", backupPath, err)
		return backup, fmt.Errorf("备份目录不存在: %v", err)
	}
	
	// 尝试读取备份信息
	infoPath := filepath.Join(backupPath, "Info.plist")
	myiToolsInfoPath := filepath.Join(backupPath, "MyiToolsBackupInfo.json")
	
	fmt.Printf("Info.plist 路径: %s\n", infoPath)
	fmt.Printf("MyiToolsBackupInfo.json 路径: %s\n", myiToolsInfoPath)
	
	// 检查Info.plist是否存在
	if _, err := os.Stat(infoPath); os.IsNotExist(err) {
		fmt.Printf("Info.plist文件不存在: %s, 错误: %v\n", infoPath, err)
		return backup, fmt.Errorf("Info.plist文件不存在: %v", err)
	}
	
	// 获取备份大小
	var size int64
	filepath.Walk(backupPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	
	backup.Size = size
	backup.BackupPath = backupPath
	backup.ID = filepath.Base(backupPath)
	backup.DeviceUDID = backup.ID
	
	// 首先尝试从我们自己的备份信息文件中读取
	if data, err := os.ReadFile(myiToolsInfoPath); err == nil {
		// 解析JSON
		infoStr := string(data)
		
		// 提取设备名称
		if idx := strings.Index(infoStr, "\"device_name\""); idx != -1 {
			backup.DeviceName = extractJSONValue(infoStr[idx+14:])
		}
		
		// 提取iOS版本
		if idx := strings.Index(infoStr, "\"ios_version\""); idx != -1 {
			backup.IOSVersion = extractJSONValue(infoStr[idx+14:])
		}
		
		// 提取创建时间
		if idx := strings.Index(infoStr, "\"created_at\""); idx != -1 {
			timeStr := extractJSONValue(infoStr[idx+13:])
			if parsedTime, err := time.Parse(time.RFC3339, timeStr); err == nil {
				backup.CreatedAt = parsedTime
			}
		}
		
		// 提取加密状态
		if idx := strings.Index(infoStr, "\"is_encrypted\""); idx != -1 {
			encryptedStr := strings.TrimSpace(infoStr[idx+15:])
			if strings.HasPrefix(encryptedStr, "true") {
				backup.IsEncrypted = true
			}
		}
	} else {
		// 如果没有我们的信息文件，则从Info.plist中提取信息
		cmd := exec.Command("plutil", "-p", infoPath)
		output, err := cmd.Output()
		if err == nil {
			outputStr := string(output)
			
			// 提取设备名称
			if idx := strings.Index(outputStr, "\"Device Name\" => "); idx != -1 {
				backup.DeviceName = extractQuotedValue(outputStr[idx+16:])
			}
			
			// 提取iOS版本
			if idx := strings.Index(outputStr, "\"Product Version\" => "); idx != -1 {
				backup.IOSVersion = extractQuotedValue(outputStr[idx+20:])
			}
			
			// 尝试从文件修改时间获取创建时间
			if fileInfo, err := os.Stat(infoPath); err == nil {
				backup.CreatedAt = fileInfo.ModTime()
			}
		}
		
		// 检查备份是否加密
		backup.IsEncrypted = isBackupEncrypted(backupPath)
	}
	
	// 打印调试信息
	fmt.Printf("读取备份信息: 设备名=%s, iOS版本=%s, 创建时间=%s\n", 
		backup.DeviceName, backup.IOSVersion, backup.CreatedAt.Format(time.RFC3339))
	
	return backup, nil
}

// 辅助函数：从字符串中提取引号内的值
func extractQuotedValue(s string) string {
	start := strings.Index(s, "\"")
	if start == -1 {
		return ""
	}
	end := strings.Index(s[start+1:], "\"")
	if end == -1 {
		return ""
	}
	return s[start+1 : start+1+end]
}

// 辅助函数：从JSON字符串中提取值
func extractJSONValue(s string) string {
	start := strings.Index(s, "\"")
	if start == -1 {
		return ""
	}
	end := strings.Index(s[start+1:], "\"")
	if end == -1 {
		return ""
	}
	return s[start+1 : start+1+end]
}

// 辅助函数：检查备份是否加密
func isBackupEncrypted(backupDir string) bool {
	manifestPath := filepath.Join(backupDir, "Manifest.plist")
	
	// 检查文件是否存在
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		fmt.Printf("Manifest.plist 文件不存在: %s\n", manifestPath)
		return false
	}
	
	cmd := exec.Command("plutil", "-p", manifestPath)
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("读取 Manifest.plist 失败: %v\n", err)
		return false
	}
	
	outputStr := string(output)
	fmt.Printf("Manifest.plist 内容: %s\n", outputStr)
	
	// 检查加密状态，可能是 "IsEncrypted" => 1 或 "IsEncrypted" => true
	isEncrypted := strings.Contains(outputStr, "\"IsEncrypted\" => 1") || 
		strings.Contains(outputStr, "\"IsEncrypted\" => true")||
		strings.Contains(outputStr, "\"IsEncrypted\" => YES")
	
	// 检查是否明确标记为未加密
	isNotEncrypted := strings.Contains(outputStr, "\"IsEncrypted\" => 0") || 
		strings.Contains(outputStr, "\"IsEncrypted\" => false") ||
		strings.Contains(outputStr, "\"IsEncrypted\" => NO")
	
	fmt.Printf("备份加密状态检测: isEncrypted=%v, isNotEncrypted=%v\n", isEncrypted, isNotEncrypted)
	
	return isEncrypted && !isNotEncrypted
}
