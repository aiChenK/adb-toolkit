import { GET, POST } from '@/utils/http';

/**
 * 执行 ADB 指令
 * @param {Object} payload 命令参数 (ip, port, op, proxyAddr, packageName, cmd)
 */
export function execAdbCommand(payload) {
  return POST('/adb', payload);
}

/**
 * 获取当前已连接的设备列表
 */
export function fetchDevices() {
  return GET('/devices');
}
