<template>
  <div class="group-section">
    <a-collapse :activeKey="activeKey" @change="onCollapseChange">
      <a-collapse-panel
        v-for="(item, index) in items"
        :key="String(index)"
      >
        <template #header>
          <span class="panel-header-content">
            <span class="panel-title">📱 {{ item.title }}</span>
            <span v-if="item.commandForm.ip" class="panel-ip-tag">
              ({{ item.commandForm.ip }}:{{ item.commandForm.port || '5555' }})
            </span>
          </span>
        </template>
        <template #extra>
          <a-button
            type="link"
            size="small"
            class="header-rename-btn"
            @click.stop="$emit('rename-group', index)"
          >
            ✏️ 重命名
          </a-button>
        </template>

        <a-form
          :model="item.commandForm"
          layout="horizontal"
          :label-col="{ span: 4 }"
          :wrapper-col="{ span: 18 }"
        >
          <a-row :gutter="16">
            <a-col :span="14">
              <a-form-item label="设备 IP" required>
                <a-input
                  v-model:value="item.commandForm.ip"
                  placeholder="例如: 192.168.1.100"
                />
              </a-form-item>
            </a-col>
            <a-col :span="10">
              <a-form-item label="端口">
                <a-input
                  v-model:value="item.commandForm.port"
                  placeholder="默认 5555"
                />
              </a-form-item>
            </a-col>
          </a-row>

          <a-row :gutter="16">
            <a-col :span="14">
              <a-form-item label="应用包名">
                <a-input
                  v-model:value="item.commandForm.packageName"
                  placeholder="例如: com.example.app"
                />
              </a-form-item>
            </a-col>
            <a-col :span="10">
              <a-form-item label="代理地址">
                <a-input
                  v-model:value="item.commandForm.proxyAddr"
                  placeholder="例如: 192.168.1.50:8888"
                />
              </a-form-item>
            </a-col>
          </a-row>

          <a-form-item :wrapper-col="{ offset: 4, span: 18 }">
            <a-space wrap>
              <a-button
                type="primary"
                @click="$emit('exec-op', { commandForm: item.commandForm, op: 'setProxy', title: item.title })"
              >
                🌐 设置代理
              </a-button>
              <a-button
                @click="$emit('exec-op', { commandForm: item.commandForm, op: 'delProxy', title: item.title })"
              >
                🚫 清除代理
              </a-button>
              <a-button
                type="primary"
                danger
                @click="$emit('exec-op', { commandForm: item.commandForm, op: 'stop', title: item.title })"
              >
                ⏹️ 停止进程
              </a-button>
              <a-button
                type="primary"
                danger
                @click="$emit('exec-op', { commandForm: item.commandForm, op: 'clear', title: item.title })"
              >
                🧹 清理缓存
              </a-button>
              <a-button
                @click="$emit('rename-group', index)"
              >
                ✏️ 重命名
              </a-button>
              <a-button
                type="dashed"
                danger
                @click="$emit('remove-group', index)"
              >
                🗑️ 关闭分组
              </a-button>
            </a-space>
          </a-form-item>
        </a-form>
      </a-collapse-panel>
    </a-collapse>
  </div>
</template>

<script>
export default {
  name: 'DeviceGroupList',
  props: {
    items: {
      type: Array,
      default: () => []
    },
    activeKey: {
      type: Array,
      default: () => ['0']
    }
  },
  emits: ['update:activeKey', 'exec-op', 'rename-group', 'remove-group'],
  methods: {
    onCollapseChange(keys) {
      this.$emit('update:activeKey', keys);
    }
  }
};
</script>

<style scoped>
.group-section {
  background: #ffffff;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.panel-header-content {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.panel-title {
  font-weight: 600;
  color: #1f2937;
}

.panel-ip-tag {
  color: #6b7280;
  font-size: 13px;
  font-weight: normal;
}

.header-rename-btn {
  padding: 0 4px;
  font-size: 13px;
}
</style>
