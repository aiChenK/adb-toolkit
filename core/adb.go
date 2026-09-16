package core

import (
	"adb-toolkit/dto/request"
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

type DeviceInfo struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Type   string `json:"type"` // "usb" 或 "tcpip"
}

// Exec 执行指定 ADB 命令并返回输出和可能的错误
func Exec(commandForm request.CommandForm) (output string, err error) {
	var args []string

	targetAddr := ""
	if commandForm.Ip != "" {
		port := commandForm.Port
		if port == "" {
			port = "5555"
		}
		targetAddr = fmt.Sprintf("%s:%s", commandForm.Ip, port)
	}

	needConnectDevice := false

	switch commandForm.Op {
	case "setProxy":
		needConnectDevice = true
		args = []string{"-s", targetAddr, "shell", "settings", "put", "global", "http_proxy", commandForm.ProxyAddr}
	case "delProxy":
		needConnectDevice = true
		args = []string{"-s", targetAddr, "shell", "settings", "put", "global", "http_proxy", ":0"}
	case "clear":
		needConnectDevice = true
		args = []string{"-s", targetAddr, "shell", "pm", "clear", commandForm.PackageName}
	case "stop":
		needConnectDevice = true
		args = []string{"-s", targetAddr, "shell", "am", "force-stop", commandForm.PackageName}
	case "getPackage":
		if targetAddr == "" {
			return "设备 IP 不能为空", fmt.Errorf("设备 IP 不能为空")
		}
		if connErr := deviceConnect(targetAddr, false); connErr != nil {
			return fmt.Sprintf("连接设备 %s 失败: %v", targetAddr, connErr), connErr
		}
		pkgName, pkgErr := getActivePackage(targetAddr)
		if pkgErr != nil {
			return pkgErr.Error(), pkgErr
		}
		fmt.Printf("📱 获取设备 [%s] 当前运行包名: %s\n", targetAddr, pkgName)
		return pkgName, nil
	default:
		// 自由命令处理
		trimmedCmd := strings.TrimSpace(commandForm.Cmd)
		if trimmedCmd == "" {
			return "命令内容不可为空", fmt.Errorf("命令内容不可为空")
		}
		// 去除前导 "adb "
		if strings.HasPrefix(trimmedCmd, "adb ") {
			trimmedCmd = strings.TrimPrefix(trimmedCmd, "adb ")
		}
		args = parseCommandArgs(trimmedCmd)
	}

	if needConnectDevice && targetAddr != "" {
		if connErr := deviceConnect(targetAddr, false); connErr != nil {
			return fmt.Sprintf("连接设备 %s 失败: %v", targetAddr, connErr), connErr
		}
	}

	cmdStr := "adb " + strings.Join(args, " ")
	fmt.Println("🚀 执行命令:", cmdStr)

	cmd := exec.Command("adb", args...)
	var stdoutStderr bytes.Buffer
	cmd.Stdout = &stdoutStderr
	cmd.Stderr = &stdoutStderr

	runErr := cmd.Run()
	output = strings.TrimSpace(stdoutStderr.String())

	if runErr != nil {
		if output == "" {
			output = runErr.Error()
		}
		return output, runErr
	}

	return output, nil
}

// parseCommandArgs 解析带引号的命令行参数
func parseCommandArgs(cmdStr string) []string {
	var args []string
	var current strings.Builder
	inQuote := false
	var quoteChar rune

	for _, r := range cmdStr {
		switch {
		case r == '\'' || r == '"':
			if inQuote && r == quoteChar {
				inQuote = false
			} else if !inQuote {
				inQuote = true
				quoteChar = r
			} else {
				current.WriteRune(r)
			}
		case r == ' ' || r == '\t':
			if inQuote {
				current.WriteRune(r)
			} else if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}
	}
	if current.Len() > 0 {
		args = append(args, current.String())
	}
	return args
}

// GetDevices 获取当前已连接的所有设备列表
func GetDevices() ([]DeviceInfo, error) {
	cmd := exec.Command("adb", "devices")
	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("执行 adb devices 失败: %w, 输出: %s", err, string(outputBytes))
	}

	lines := strings.Split(string(outputBytes), "\n")
	if len(lines) <= 1 {
		return []DeviceInfo{}, nil
	}

	var devices []DeviceInfo
	// 过滤第一行 "List of devices attached"
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			devID := fields[0]
			devStatus := fields[1]
			devType := "usb"
			if strings.Contains(devID, ":") {
				devType = "tcpip"
			}
			devices = append(devices, DeviceInfo{
				ID:     devID,
				Status: devStatus,
				Type:   devType,
			})
		}
	}
	return devices, nil
}

// deviceConnect 自动检测并连接指定 IP:Port 设备
func deviceConnect(targetAddr string, isRetry bool) error {
	devices, err := GetDevices()
	if err != nil {
		return err
	}

	deviceStatus := ""
	found := false
	for _, dev := range devices {
		if dev.ID == targetAddr {
			found = true
			deviceStatus = dev.Status
			break
		}
	}

	// 如果未连接或状态非 device，尝试重连
	if !found || deviceStatus != "device" {
		if !isRetry {
			if found && deviceStatus != "device" {
				// 定向断开异常设备，避免影响全局
				_, _ = exec.Command("adb", "disconnect", targetAddr).CombinedOutput()
			}
			// 执行重连
			connectOut, connErr := exec.Command("adb", "connect", targetAddr).CombinedOutput()
			fmt.Printf("连接设备 [%s] 结果: %s\n", targetAddr, strings.TrimSpace(string(connectOut)))
			if connErr != nil {
				return fmt.Errorf("adb connect %s 失败: %w", targetAddr, connErr)
			}
			// 二次校验
			return deviceConnect(targetAddr, true)
		} else {
			return fmt.Errorf("设备 [%s] 连接后状态仍为异常 (%s)", targetAddr, deviceStatus)
		}
	}

	return nil
}

// extractPackageName 从 dumpsys 输出中提取前台应用包名
func extractPackageName(output string) string {
	pkgRegex := regexp.MustCompile(`([a-zA-Z][a-zA-Z0-9_]*(?:\.[a-zA-Z0-9_]+)+)/`)
	if match := pkgRegex.FindStringSubmatch(output); len(match) > 1 {
		return match[1]
	}
	return ""
}

// getActivePackage 获取设备当前运行的前台包名 (跨 Android 版本与主流模拟器兼容)
func getActivePackage(targetAddr string) (string, error) {
	// 策略 1: dumpsys window (优先获取当前窗口焦点，适用于大多数原生系统与高版本 Android)
	out, _ := exec.Command("adb", "-s", targetAddr, "shell", "dumpsys window | grep -E 'mCurrentFocus|mFocusedApp'").CombinedOutput()
	if pkg := extractPackageName(string(out)); pkg != "" {
		return pkg, nil
	}

	// 策略 2: dumpsys activity activities (兼容 MuMu、雷电等模拟器以及部分定制 ROM)
	out, _ = exec.Command("adb", "-s", targetAddr, "shell", "dumpsys activity activities | grep -E 'mResumedActivity|topResumedActivity|ResumedActivity'").CombinedOutput()
	if pkg := extractPackageName(string(out)); pkg != "" {
		return pkg, nil
	}

	// 策略 3: dumpsys activity top (备用后备方案)
	out, _ = exec.Command("adb", "-s", targetAddr, "shell", "dumpsys activity top").CombinedOutput()
	if pkg := extractPackageName(string(out)); pkg != "" {
		return pkg, nil
	}

	return "", fmt.Errorf("未检测到当前正在运行的应用包名（请确认屏幕已解锁且应用处于前台）")
}
