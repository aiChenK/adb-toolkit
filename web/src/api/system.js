import { GET } from '@/utils/http';

/**
 * 获取宿主机局域网 IP
 */
export function fetchLocalIp() {
  return GET('/ip');
}
