<template>
  <div v-show="logs.length > 0" class="result-section">
    <div class="result-header">
      <span class="result-title">📋 执行日志与输出 ({{ logs.length }})</span>
      <a-button type="link" size="small" @click="$emit('clear-logs')">清空日志</a-button>
    </div>
    <a-spin :spinning="spinning">
      <div class="timeline-container">
        <a-timeline>
          <a-timeline-item
            v-for="(item, index) in logs"
            :key="index"
            :color="item.color"
          >
            <div class="log-card">
              <div class="log-meta">
                <span class="log-time">{{ item.time }}</span>
                <span class="log-cmd">{{ item.title }}</span>
                <a-button type="text" size="small" class="copy-btn" @click="copyText(item.data)">复制</a-button>
              </div>
              <pre class="log-output">{{ item.data }}</pre>
            </div>
          </a-timeline-item>
        </a-timeline>
      </div>
    </a-spin>
  </div>
</template>

<script>
import { message } from 'ant-design-vue';

export default {
  name: 'LogConsole',
  props: {
    logs: {
      type: Array,
      default: () => []
    },
    spinning: {
      type: Boolean,
      default: false
    }
  },
  emits: ['clear-logs'],
  methods: {
    copyText(text) {
      if (!text) return;
      navigator.clipboard
        .writeText(text)
        .then(() => {
          message.success('已复制到剪贴板');
        })
        .catch(() => {
          message.error('复制失败，请手动选取复制');
        });
    }
  }
};
</script>

<style scoped>
.result-section {
  margin-bottom: 20px;
  background: #ffffff;
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.result-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  border-bottom: 1px solid #f0f0f0;
  padding-bottom: 8px;
}

.result-title {
  font-weight: 600;
  color: #262626;
}

.timeline-container {
  max-height: 320px;
  overflow-y: auto;
  padding-right: 8px;
}

.log-card {
  background: #fafafa;
  border: 1px solid #f0f0f0;
  border-radius: 6px;
  padding: 8px 12px;
  margin-bottom: 4px;
}

.log-meta {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 6px;
}

.log-time {
  font-size: 12px;
  color: #8c8c8c;
  font-family: monospace;
}

.log-cmd {
  font-size: 13px;
  font-weight: 600;
  color: #1890ff;
  flex: 1;
}

.copy-btn {
  padding: 0;
  height: auto;
}

.log-output {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-all;
  font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, monospace;
  font-size: 12.5px;
  line-height: 1.45;
  color: #262626;
  max-height: 200px;
  overflow-y: auto;
}
</style>
