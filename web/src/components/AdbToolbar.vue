<template>
  <div class="top-toolbar">
    <a-space size="middle" wrap style="width: 100%;">
      <a-dropdown-button @click="$emit('add-item')" type="dashed">
        ➕ 新增设备组
        <template #overlay>
          <a-menu>
            <a-menu-item @click="$emit('fill-ip')">
              🌐 填充本机代理 IP
            </a-menu-item>
            <a-menu-item @click="$emit('refresh-devices')">
              🔄 刷新已连接设备
            </a-menu-item>
            <a-menu-item @click="$emit('clear-logs')">
              🧹 清空日志结果区
            </a-menu-item>
          </a-menu>
        </template>
      </a-dropdown-button>

      <a-select
        placeholder="🎯 选择预置命令"
        :options="options"
        v-model:value="selectedPreset"
        style="min-width: 220px;"
        @select="onPresetSelect"
      />

      <a-input
        v-model:value="freeCommandText"
        placeholder="💻 输入 ADB 自由命令 (例如: devices, shell getprop ro.product.model)"
        style="flex: 1; min-width: 280px;"
        @pressEnter="handleExec"
      />

      <a-dropdown-button @click="handleExec" type="primary">
        ▶️ 执行
        <template #overlay>
          <a-menu>
            <a-menu-item
              v-for="(item, index) in items"
              :key="index"
              @click="handleDirectExec(item)"
            >
              定向至: {{ item.title }} ({{ item.commandForm.ip || '未填IP' }})
            </a-menu-item>
          </a-menu>
        </template>
      </a-dropdown-button>
    </a-space>
  </div>
</template>

<script>
import { message } from 'ant-design-vue';
import { COMMAND_OPTIONS } from '@/constants/commands';

export default {
  name: 'AdbToolbar',
  props: {
    items: {
      type: Array,
      default: () => []
    },
    commandOptions: {
      type: Array,
      default: () => COMMAND_OPTIONS
    }
  },
  emits: ['add-item', 'fill-ip', 'refresh-devices', 'clear-logs', 'exec-free', 'exec-direct'],
  data() {
    return {
      selectedPreset: null,
      freeCommandText: ''
    };
  },
  computed: {
    options() {
      return this.commandOptions || COMMAND_OPTIONS;
    }
  },
  methods: {
    onPresetSelect(value) {
      this.freeCommandText = value;
      this.selectedPreset = null;
    },
    handleExec() {
      const trimmed = (this.freeCommandText || '').trim();
      if (!trimmed) {
        message.warning('请输入自由命令');
        return;
      }
      this.$emit('exec-free', trimmed);
    },
    handleDirectExec(item) {
      const trimmed = (this.freeCommandText || '').trim();
      if (!trimmed) {
        message.warning('请输入自由命令');
        return;
      }
      if (!item.commandForm.ip) {
        message.error('所选分组的设备 IP 不能为空');
        return;
      }
      this.$emit('exec-direct', { item, cmd: trimmed });
    }
  }
};
</script>

<style scoped>
.top-toolbar {
  margin-bottom: 16px;
  padding: 16px;
  background: #ffffff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}
</style>
