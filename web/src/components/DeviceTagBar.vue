<template>
  <div v-if="devices.length > 0" class="device-bar">
    <span class="device-bar-title">📱 检测到已连接设备:</span>
    <a-tag
      v-for="dev in devices"
      :key="dev.id"
      :color="dev.status === 'device' ? 'green' : 'orange'"
      class="device-tag"
      @click="$emit('apply-device', dev)"
    >
      {{ dev.id }} ({{ dev.type }}) [{{ dev.status }}]
    </a-tag>
  </div>
</template>

<script>
export default {
  name: 'DeviceTagBar',
  props: {
    devices: {
      type: Array,
      default: () => []
    }
  },
  emits: ['apply-device']
};
</script>

<style scoped>
.device-bar {
  margin-bottom: 16px;
  padding: 8px 12px;
  background: #f6ffed;
  border: 1px solid #b7eb8f;
  border-radius: 6px;
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
}

.device-bar-title {
  font-size: 13px;
  font-weight: 600;
  color: #389e0d;
}

.device-tag {
  cursor: pointer;
  user-select: none;
  transition: all 0.2s ease;
}

.device-tag:hover {
  opacity: 0.85;
  transform: translateY(-1px);
}
</style>
