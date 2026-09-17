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
            <a-menu-divider />
            <a-menu-item @click="$emit('reset-groups')" danger>
              🔄 重置设备组配置
            </a-menu-item>
          </a-menu>
        </template>
      </a-dropdown-button>

      <a-select
        placeholder="🎯 选择预置命令"
        :options="options"
        v-model:value="selectedPreset"
        style="min-width: 200px;"
        @select="onPresetSelect"
      />

      <a-input
        v-model:value="freeCommandText"
        placeholder="💻 输入 ADB 自由命令 (例如: devices, shell getprop ro.product.model)"
        style="flex: 1; min-width: 320px;"
        allow-clear
        @pressEnter="handleExec"
      >
        <template #addonBefore>
          <a-select
            v-model:value="currentTarget"
            :options="targetOptions"
            placeholder="🎯 目标设备"
            :dropdown-match-select-width="false"
            class="target-addon-select"
          />
        </template>
      </a-input>

      <a-tooltip :title="execTooltip">
        <a-button type="primary" @click="handleExec">
          ▶️ 执行
        </a-button>
      </a-tooltip>
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
    onlineDevices: {
      type: Array,
      default: () => []
    },
    commandOptions: {
      type: Array,
      default: () => COMMAND_OPTIONS
    },
    target: {
      type: String,
      default: ''
    }
  },
  emits: [
    'add-item',
    'fill-ip',
    'refresh-devices',
    'clear-logs',
    'reset-groups',
    'exec-command',
    'exec-free',
    'exec-direct',
    'update:target'
  ],
  data() {
    return {
      selectedPreset: null,
      freeCommandText: '',
      internalTarget: this.target || ''
    };
  },
  computed: {
    currentTarget: {
      get() {
        return this.target !== undefined ? this.target : this.internalTarget;
      },
      set(val) {
        this.internalTarget = val;
        this.$emit('update:target', val);
      }
    },
    options() {
      return this.commandOptions || COMMAND_OPTIONS;
    },
    targetOptions() {
      const groups = [
        {
          label: '🌐 全局模式',
          options: [
            { label: '🌐 全局 / 默认 (不限设备)', value: '' }
          ]
        }
      ];

      // 1. 实时检测到的在线设备 (USB / 模拟器 / 局域网)
      if (this.onlineDevices && this.onlineDevices.length > 0) {
        groups.push({
          label: '📱 实时在线设备',
          options: this.onlineDevices.map(dev => ({
            label: `🟢 ${dev.id} (${dev.type || 'device'})`,
            value: `dev:${dev.id}`
          }))
        });
      }

      // 2. 下方配置的设备分组卡片
      if (this.items && this.items.length > 0) {
        groups.push({
          label: '📁 配置设备组',
          options: this.items.map((item, index) => {
            const ip = item.commandForm?.ip ? item.commandForm.ip.trim() : '';
            const port = item.commandForm?.port ? item.commandForm.port.trim() : '5555';
            const target = ip ? (ip.includes(':') ? ip : `${ip}:${port}`) : '';
            return {
              label: `📂 ${item.title || '设备 ' + (index + 1)} (${target || '未配置IP'})`,
              value: `item:${index}`,
              disabled: !ip
            };
          })
        });
      }

      return groups;
    },
    execTooltip() {
      if (!this.currentTarget) {
        return '全局执行 (不带 -s 参数)';
      }
      if (this.currentTarget.startsWith('dev:')) {
        return `定向至在线设备: ${this.currentTarget.slice(4)}`;
      }
      if (this.currentTarget.startsWith('item:')) {
        const idx = parseInt(this.currentTarget.slice(5), 10);
        const item = this.items[idx];
        return `定向至配置组: ${item ? item.title : '未知'} (${item?.commandForm?.ip || ''})`;
      }
      return '执行当前命令';
    }
  },
  watch: {
    target(newVal) {
      this.internalTarget = newVal;
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

      let target = '';
      let targetLabel = '';

      if (this.currentTarget.startsWith('dev:')) {
        target = this.currentTarget.slice(4);
        targetLabel = target;
      } else if (this.currentTarget.startsWith('item:')) {
        const idx = parseInt(this.currentTarget.slice(5), 10);
        const item = this.items[idx];
        if (!item || !item.commandForm?.ip) {
          message.error('所选设备组未配置有效 IP');
          return;
        }
        const ip = item.commandForm.ip.trim();
        const port = (item.commandForm.port || '5555').trim();
        target = ip.includes(':') ? ip : `${ip}:${port}`;
        targetLabel = item.title || target;
      }

      let fullCmd = trimmed;
      // 若选定了目标设备，且命令本身未包含 -s，则自动拼接 -s
      if (target) {
        if (!fullCmd.startsWith('-s ') && !fullCmd.includes(' -s ')) {
          fullCmd = `-s ${target} ${fullCmd}`;
        }
      }

      // 触发最新统一事件
      this.$emit('exec-command', {
        cmd: fullCmd,
        rawCmd: trimmed,
        target,
        targetLabel
      });

      // 兼容旧事件
      this.$emit('exec-free', fullCmd);
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

.target-addon-select {
  width: 175px;
  text-align: left;
}

:deep(.target-addon-select .ant-select-selector) {
  border: none !important;
  box-shadow: none !important;
  background-color: transparent !important;
}
</style>
