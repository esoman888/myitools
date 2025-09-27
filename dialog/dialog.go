package dialog

import (
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"context"
	"fmt"
)

// DialogManager 对话框管理器
type DialogManager struct {
	ctx context.Context
}

// NewDialogManager 创建新的对话框管理器
func NewDialogManager() *DialogManager {
	return &DialogManager{}
}

// SetContext 设置上下文
func (d *DialogManager) SetContext(ctx context.Context) {
	d.ctx = ctx
}

// OpenDirectoryDialog 打开目录选择对话框
func (d *DialogManager) OpenDirectoryDialog(title string, defaultPath string) (string, error) {
	if d.ctx == nil {
		return "", fmt.Errorf("context not set")
	}
	
	options := runtime.OpenDialogOptions{
		Title:                title,
		DefaultDirectory:     defaultPath,
		CanCreateDirectories: true,
	}
	
	fmt.Printf("Opening directory dialog with title: %s, defaultPath: %s\n", title, defaultPath)
	result, err := runtime.OpenDirectoryDialog(d.ctx, options)
	if err != nil {
		fmt.Printf("Error opening directory dialog: %v\n", err)
		return "", err
	}
	
	fmt.Printf("Selected directory: %s\n", result)
	return result, nil
}

// OpenFileDialog 打开文件选择对话框
func (d *DialogManager) OpenFileDialog(title string, defaultPath string, filters []runtime.FileFilter) (string, error) {
	if d.ctx == nil {
		return "", fmt.Errorf("context not set")
	}
	
	options := runtime.OpenDialogOptions{
		Title:            title,
		DefaultDirectory: defaultPath,
		Filters:          filters,
	}
	
	return runtime.OpenFileDialog(d.ctx, options)
}

// SaveFileDialog 打开文件保存对话框
func (d *DialogManager) SaveFileDialog(title string, defaultPath string, filters []runtime.FileFilter) (string, error) {
	if d.ctx == nil {
		return "", fmt.Errorf("context not set")
	}
	
	options := runtime.SaveDialogOptions{
		Title:            title,
		DefaultDirectory: defaultPath,
		Filters:          filters,
	}
	
	return runtime.SaveFileDialog(d.ctx, options)
}

// MessageDialog 显示消息对话框
func (d *DialogManager) MessageDialog(title string, message string, dialogType string) (string, error) {
	if d.ctx == nil {
		return "", fmt.Errorf("context not set")
	}
	
	var dt runtime.DialogType
	
	switch dialogType {
	case "info":
		dt = runtime.InfoDialog
	case "warning":
		dt = runtime.WarningDialog
	case "error":
		dt = runtime.ErrorDialog
	case "question":
		dt = runtime.QuestionDialog
	default:
		dt = runtime.InfoDialog
	}
	
	options := runtime.MessageDialogOptions{
		Type:    dt,
		Title:   title,
		Message: message,
	}
	
	return runtime.MessageDialog(d.ctx, options)
}