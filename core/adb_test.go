package core

import (
	"net"
	"reflect"
	"testing"
)

func TestParseCommandArgs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "简单命令参数",
			input:    "shell getprop ro.product.model",
			expected: []string{"shell", "getprop", "ro.product.model"},
		},
		{
			name:     "包含单引号的长参数",
			input:    "shell 'settings put global http_proxy 192.168.1.100:8888'",
			expected: []string{"shell", "settings put global http_proxy 192.168.1.100:8888"},
		},
		{
			name:     "包含双引号的管道命令",
			input:    `shell "dumpsys window | grep -E 'mCurrentFocus'"`,
			expected: []string{"shell", `dumpsys window | grep -E 'mCurrentFocus'`},
		},
		{
			name:     "多个空格与制表符处理",
			input:    "   shell   pm   clear   com.example.app\t",
			expected: []string{"shell", "pm", "clear", "com.example.app"},
		},
		{
			name:     "空字符串输入",
			input:    "",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := parseCommandArgs(tt.input)
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("parseCommandArgs(%q) = %v, expected %v", tt.input, actual, tt.expected)
			}
		})
	}
}

func TestExtractPackageName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "dumpsys window 焦点窗口格式",
			input:    "mCurrentFocus=Window{5a91720 u0 com.android.settings/com.android.settings.SubSettings}",
			expected: "com.android.settings",
		},
		{
			name:     "mFocusedApp 格式",
			input:    "mFocusedApp=AppWindowToken{379d2b2 token=Token{ActivityRecord{75ba023 u0 com.tencent.mm/.ui.LauncherUI t123}}}",
			expected: "com.tencent.mm",
		},
		{
			name:     "dumpsys activity activities 中的 topResumedActivity",
			input:    "topResumedActivity=ActivityRecord{2c76a59 u0 com.google.android.apps.photos/.home.HomeActivity t1}",
			expected: "com.google.android.apps.photos",
		},
		{
			name:     "未包含有效包名格式",
			input:    "Nothing focused here",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := extractPackageName(tt.input)
			if actual != tt.expected {
				t.Errorf("extractPackageName(%q) = %q, expected %q", tt.input, actual, tt.expected)
			}
		})
	}
}

func TestIsPrivate172(t *testing.T) {
	tests := []struct {
		ip       string
		expected bool
	}{
		{"172.16.0.1", true},
		{"172.31.255.255", true},
		{"172.20.10.5", true},
		{"172.15.255.255", false},
		{"172.32.0.1", false},
		{"192.168.1.1", false},
		{"10.0.0.1", false},
	}

	for _, tt := range tests {
		t.Run(tt.ip, func(t *testing.T) {
			parsedIP := net.ParseIP(tt.ip)
			if parsedIP == nil {
				t.Fatalf("无法解析测试 IP: %s", tt.ip)
			}
			actual := isPrivate172(parsedIP)
			if actual != tt.expected {
				t.Errorf("isPrivate172(%s) = %v, expected %v", tt.ip, actual, tt.expected)
			}
		})
	}
}

func TestGetLocalIp(t *testing.T) {
	ip, err := GetLocalIp()
	if err != nil {
		t.Fatalf("GetLocalIp() 发生错误: %v", err)
	}
	if ip == "" {
		t.Errorf("GetLocalIp() 返回的 IP 字符串为空")
	}
	if parsed := net.ParseIP(ip); parsed == nil {
		t.Errorf("GetLocalIp() 返回的不是合法的 IP: %s", ip)
	}
}

