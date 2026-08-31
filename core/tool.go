package core

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// GetLocalIp 获取本机局域网 IPv4 地址（跨平台支持 macOS / Windows / Linux）
func GetLocalIp() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("获取网络接口列表失败: %w", err)
	}

	// 候选 IP 列表，优先局域网常用私有网段 (192.168.x.x > 10.x.x.x > 172.16~31.x.x)
	var fallbackIP string

	for _, iface := range interfaces {
		// 跳过未启用、环回以及点对点接口
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagPointToPoint != 0 {
			continue
		}

		nameLower := strings.ToLower(iface.Name)
		// 过滤常见的虚拟网卡接口
		if strings.Contains(nameLower, "docker") ||
			strings.Contains(nameLower, "veth") ||
			strings.Contains(nameLower, "br-") ||
			strings.Contains(nameLower, "vmnet") ||
			strings.Contains(nameLower, "vethernet") ||
			strings.Contains(nameLower, "tailscale") ||
			strings.Contains(nameLower, "virbr") ||
			strings.Contains(nameLower, "tun") ||
			strings.Contains(nameLower, "tap") {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			ip4 := ip.To4()
			if ip4 == nil {
				continue
			}

			ipStr := ip4.String()
			// 优先匹配 192.168.x.x
			if strings.HasPrefix(ipStr, "192.168.") {
				return ipStr, nil
			}
			// 其次匹配 10.x.x.x 或 172.16~31.x.x
			if strings.HasPrefix(ipStr, "10.") || isPrivate172(ip4) {
				fallbackIP = ipStr
				continue
			}

			if fallbackIP == "" {
				fallbackIP = ipStr
			}
		}
	}

	if fallbackIP != "" {
		return fallbackIP, nil
	}

	return "127.0.0.1", nil
}

func isPrivate172(ip net.IP) bool {
	if len(ip) == 4 || len(ip) == 16 {
		ip4 := ip.To4()
		if ip4 != nil && ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31 {
			return true
		}
	}
	return false
}

// CheckAdb 检查 adb 环境变量配置是否正常
func CheckAdb() bool {
	cmd := exec.Command("adb", "version")
	err := cmd.Run()
	if err != nil {
		PrintWarn("⚠️ 未检测到 adb 环境变量或 adb 未安装，部分 ADB 操作可能无法执行。请确保系统已安装 Android Platform Tools 并配置 PATH。")
		return false
	}
	return true
}

// PrintErr 在终端打印红色错误信息
func PrintErr(msg string) {
	fmt.Print("\033[31m")
	fmt.Println(msg)
	fmt.Print("\033[0m")
}

// PrintWarn 在终端打印黄色警告信息
func PrintWarn(msg string) {
	fmt.Print("\033[33m")
	fmt.Println(msg)
	fmt.Print("\033[0m")
}

// PrintSuccess 在终端打印绿色成功信息
func PrintSuccess(msg string) {
	fmt.Print("\033[32m")
	fmt.Println(msg)
	fmt.Print("\033[0m")
}

// OpenBrowser 打开系统默认浏览器并访问指定 url
func OpenBrowser(url string) {
	// 等待服务监听就绪
	time.Sleep(300 * time.Millisecond)

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default: // Linux 等
		cmd = exec.Command("xdg-open", url)
	}

	_ = cmd.Start()
}
