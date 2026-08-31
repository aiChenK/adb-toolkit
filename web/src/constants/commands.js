/**
 * 预设常用 ADB 快捷命令配置
 */
export const COMMAND_OPTIONS = [
  { value: 'devices', label: '📱 adb devices (查看设备)' },
  { value: 'disconnect', label: '🔌 adb disconnect (断开连接)' },
  { value: 'shell pm list packages -3', label: '📦 第三方已安装应用列表' },
  { value: 'shell getprop ro.product.model', label: 'ℹ️ 获取设备型号' },
  { value: 'shell wm size', label: '📐 获取屏幕分辨率' },
  { value: 'shell dumpsys battery', label: '🔋 查看电池状态' },
  { value: 'shell ip route', label: '🌐 查看设备网络路由' }
];

/**
 * 默认设备表单初始数据模板
 */
export const DEFAULT_COMMAND_FORM = {
  ip: '',
  port: '5555',
  packageName: '',
  proxyAddr: '',
  cmd: ''
};

/**
 * 最大支持设备分组数
 */
export const MAX_GROUP_COUNT = 8;
