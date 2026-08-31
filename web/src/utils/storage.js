import { DEFAULT_COMMAND_FORM, MAX_GROUP_COUNT } from '@/constants/commands';

export const STORAGE_KEY_DEVICE_GROUPS = 'adb_toolkit_device_groups';

/**
 * 从 localStorage 读取已保存的设备分组配置
 * @returns {Array|null} 格式化后的设备组列表，读取失败或无数据时返回 null
 */
export function loadStoredDeviceGroups() {
  try {
    const rawData = localStorage.getItem(STORAGE_KEY_DEVICE_GROUPS);
    if (!rawData) {
      return null;
    }
    const parsed = JSON.parse(rawData);
    if (!Array.isArray(parsed) || parsed.length === 0) {
      return null;
    }

    // 限制最大分组数并规范化字段
    const validItems = parsed.slice(0, MAX_GROUP_COUNT).map((item, index) => {
      const itemTitle = item && typeof item.title === 'string' && item.title.trim()
        ? item.title.trim()
        : `设备 ${index + 1}`;
      const itemForm = item && typeof item.commandForm === 'object' && item.commandForm !== null
        ? item.commandForm
        : {};

      return {
        title: itemTitle,
        commandForm: {
          ...DEFAULT_COMMAND_FORM,
          ...itemForm
        }
      };
    });

    return validItems.length > 0 ? validItems : null;
  } catch (err) {
    console.warn('[Storage] 读取设备分组缓存失败:', err);
    return null;
  }
}

/**
 * 将设备分组配置写入 localStorage
 * @param {Array} items 设备分组列表
 */
export function saveStoredDeviceGroups(items) {
  try {
    if (!Array.isArray(items)) {
      return;
    }
    localStorage.setItem(STORAGE_KEY_DEVICE_GROUPS, JSON.stringify(items));
  } catch (err) {
    console.warn('[Storage] 保存设备分组缓存失败:', err);
  }
}

/**
 * 清除设备分组的 localStorage 缓存
 */
export function clearStoredDeviceGroups() {
  try {
    localStorage.removeItem(STORAGE_KEY_DEVICE_GROUPS);
  } catch (err) {
    console.warn('[Storage] 清除设备分组缓存失败:', err);
  }
}
