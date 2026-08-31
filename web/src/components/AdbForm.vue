<template>
  <div class="adb-container">
    <!-- 顶部操作栏 -->
    <div class="top-toolbar">
      <a-space size="middle" wrap style="width: 100%;">
        <a-dropdown-button @click="addItem" type="dashed">
          ➕ 新增设备组
          <template #overlay>
            <a-menu>
              <a-menu-item @click="fillLocalIp">
                🌐 填充本机代理 IP
              </a-menu-item>
              <a-menu-item @click="fetchConnectedDevices">
                🔄 刷新已连接设备
              </a-menu-item>
              <a-menu-item @click="clearResultData">
                🧹 清空日志结果区
              </a-menu-item>
            </a-menu>
          </template>
        </a-dropdown-button>

        <a-select
          placeholder="🎯 选择预置命令"
          :options="commandOptions"
          v-model:value="selectCommand"
          style="min-width: 220px;"
          @select="onCommandSelect"
        />

        <a-input
          v-model:value="freeCommand"
          placeholder="💻 输入 ADB 自由命令 (例如: devices, shell getprop ro.product.model)"
          style="flex: 1; min-width: 280px;"
          @pressEnter="execFreeCommand"
        />

        <a-dropdown-button @click="execFreeCommand" type="primary">
          ▶️ 执行
          <template #overlay>
            <a-menu>
              <a-menu-item
                v-for="(item, index) in items"
                :key="index"
                @click="execFreeCommandWithDevice(item)"
              >
                定向至: {{ item.title }} ({{ item.commandForm.ip || '未填IP' }})
              </a-menu-item>
            </a-menu>
          </template>
        </a-dropdown-button>
      </a-space>
    </div>

    <!-- 在线设备快捷栏 (若检测到在线设备) -->
    <div v-if="onlineDevices.length > 0" class="device-bar">
      <span class="device-bar-title">📱 检测到已连接设备:</span>
      <a-tag
        v-for="dev in onlineDevices"
        :key="dev.id"
        :color="dev.status === 'device' ? 'green' : 'orange'"
        class="device-tag"
        @click="applyDeviceToActive(dev)"
      >
        {{ dev.id }} ({{ dev.type }}) [{{ dev.status }}]
      </a-tag>
    </div>

    <!-- 结果展示区 -->
    <div v-show="resultData.length > 0" class="result-section">
      <div class="result-header">
        <span class="result-title">📋 执行日志与输出 ({{ resultData.length }})</span>
        <a-button type="link" size="small" @click="clearResultData">清空日志</a-button>
      </div>
      <a-spin :spinning="spinning">
        <div class="timeline-container">
          <a-timeline>
            <a-timeline-item
              v-for="(item, index) in resultData"
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

    <!-- 设备分组配置面板 -->
    <div class="group-section">
      <a-collapse v-model:activeKey="activeKey">
        <a-collapse-panel
          v-for="(item, index) in items"
          :key="String(index)"
          :header="`📱 ${item.title} ${item.commandForm.ip ? '(' + item.commandForm.ip + ':' + (item.commandForm.port || '5555') + ')' : ''}`"
        >
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
                  @click="sendHttpRequest(item.commandForm, 'setProxy', item.title)"
                >
                  🌐 设置代理
                </a-button>
                <a-button
                  @click="sendHttpRequest(item.commandForm, 'delProxy', item.title)"
                >
                  🚫 清除代理
                </a-button>
                <a-button
                  type="primary"
                  danger
                  @click="sendHttpRequest(item.commandForm, 'stop', item.title)"
                >
                  ⏹️ 停止进程
                </a-button>
                <a-button
                  type="primary"
                  danger
                  @click="sendHttpRequest(item.commandForm, 'clear', item.title)"
                >
                  🧹 清理缓存
                </a-button>
                <a-button
                  type="dashed"
                  danger
                  @click="distoryItem(index)"
                >
                  🗑️ 关闭分组
                </a-button>
              </a-space>
            </a-form-item>
          </a-form>
        </a-collapse-panel>
      </a-collapse>
    </div>
  </div>
</template>

<script>
import { message } from 'ant-design-vue';
import { GET, POST } from '@/utils/http';

const defaultCommandForm = {
  ip: '',
  port: '5555',
  packageName: '',
  proxyAddr: '',
  cmd: ''
};

export default {
  name: 'AdbForm',
  data() {
    return {
      activeKey: ['0'],
      items: [
        {
          title: '设备 1',
          commandForm: { ...defaultCommandForm }
        }
      ],
      freeCommand: '',
      selectCommand: null,
      commandOptions: [
        { value: 'devices', label: '📱 adb devices (查看设备)' },
        { value: 'disconnect', label: '🔌 adb disconnect (断开连接)' },
        { value: 'shell pm list packages -3', label: '📦 第三方已安装应用列表' },
        { value: 'shell getprop ro.product.model', label: 'ℹ️ 获取设备型号' },
        { value: 'shell wm size', label: '📐 获取屏幕分辨率' },
        { value: 'shell dumpsys battery', label: '🔋 查看电池状态' },
        { value: 'shell ip route', label: '🌐 查看设备网络路由' }
      ],
      spinning: false,
      resultData: [],
      onlineDevices: []
    };
  },
  mounted() {
    this.fetchConnectedDevices();
  },
  methods: {
    formatTime() {
      const now = new Date();
      return now.toTimeString().split(' ')[0];
    },
    onCommandSelect(value) {
      this.freeCommand = value;
      this.selectCommand = null;
    },
    fetchConnectedDevices() {
      GET('/devices')
        .then(res => {
          if (res.data && res.data.success && Array.isArray(res.data.data)) {
            this.onlineDevices = res.data.data;
          }
        })
        .catch(() => {});
    },
    applyDeviceToActive(dev) {
      const activeIdx = parseInt(this.activeKey[0] || '0', 10);
      if (this.items[activeIdx]) {
        if (dev.id.includes(':')) {
          const parts = dev.id.split(':');
          this.items[activeIdx].commandForm.ip = parts[0];
          this.items[activeIdx].commandForm.port = parts[1] || '5555';
        } else {
          this.items[activeIdx].commandForm.ip = dev.id;
        }
        message.success(`已将设备 ${dev.id} 填充至分组 [${this.items[activeIdx].title}]`);
      }
    },
    sendHttpRequest(commandForm, op, title) {
      const formPayload = { ...commandForm, op };
      this.spinning = true;
      const showCmd = this.dealShowCommandInfo(formPayload, title);
      const currentTime = this.formatTime();

      POST('/adb', formPayload)
        .then(res => {
          const responseData = res.data;
          if (responseData.success) {
            const outputText = responseData.data || '执行成功 (无输出)';
            this.resultData.unshift({
              time: currentTime,
              title: showCmd,
              color: 'green',
              data: outputText
            });
            message.success('执行成功');
          } else {
            this.resultData.unshift({
              time: currentTime,
              title: showCmd,
              color: 'red',
              data: responseData.errMessage || '执行失败'
            });
          }
          if (this.resultData.length > 50) {
            this.resultData.pop();
          }
          this.fetchConnectedDevices();
        })
        .catch(errMsg => {
          this.resultData.unshift({
            time: currentTime,
            title: showCmd,
            color: 'red',
            data: typeof errMsg === 'string' ? errMsg : JSON.stringify(errMsg)
          });
        })
        .finally(() => {
          this.spinning = false;
        });
    },
    dealShowCommandInfo(commandForm, title) {
      let showCmd = commandForm.cmd;
      if (commandForm.op !== 'free') {
        showCmd = commandForm.op;
      }
      if (title) {
        showCmd += ` [${title}]`;
      }
      return '执行: ' + showCmd;
    },
    addItem() {
      if (this.items.length >= 8) {
        message.error('最多添加 8 个设备分组');
        return;
      }
      const newTitle = prompt('请输入新设备组名称', `设备 ${this.items.length + 1}`);
      if (newTitle) {
        this.items.push({
          title: newTitle,
          commandForm: { ...defaultCommandForm }
        });
        this.activeKey = [String(this.items.length - 1)];
      }
    },
    distoryItem(index) {
      if (this.items.length <= 1) {
        message.error('至少保留一个设备分组');
        return;
      }
      this.items.splice(index, 1);
      this.activeKey = ['0'];
    },
    execFreeCommand() {
      if (!this.freeCommand.trim()) {
        message.warning('请输入自由命令');
        return;
      }
      this.sendHttpRequest({ cmd: this.freeCommand }, 'free');
    },
    execFreeCommandWithDevice(item) {
      if (!this.freeCommand.trim()) {
        message.warning('请输入自由命令');
        return;
      }
      const devIp = item.commandForm.ip;
      if (!devIp) {
        message.error('所选分组的设备 IP 不能为空');
        return;
      }
      const port = item.commandForm.port || '5555';
      const target = devIp.includes(':') ? devIp : `${devIp}:${port}`;
      const fullCmd = `-s ${target} ${this.freeCommand}`;
      this.sendHttpRequest({ cmd: fullCmd }, 'free', item.title);
    },
    fillLocalIp() {
      const proxyPort = prompt('请输入抓包代理端口号 (例如: 8888, 7890)', '8888');
      if (proxyPort === null) return;

      GET('/ip')
        .then(res => {
          if (res.data && res.data.data) {
            const localIp = res.data.data;
            const fullProxy = `${localIp}:${proxyPort || '8888'}`;
            this.items.forEach(item => {
              item.commandForm.proxyAddr = fullProxy;
            });
            message.success(`已自动填充代理地址: ${fullProxy}`);
          }
        })
        .catch(err => {
          message.error('获取本机 IP 失败: ' + err);
        });
    },
    clearResultData() {
      this.resultData = [];
      message.info('日志结果区已清空');
    },
    copyText(text) {
      if (!text) return;
      navigator.clipboard.writeText(text).then(() => {
        message.success('已复制到剪贴板');
      }).catch(() => {
        message.error('复制失败，请手动选取复制');
      });
    }
  }
};
</script>

<style scoped>
.adb-container {
  text-align: left;
}

.top-toolbar {
  margin-bottom: 16px;
  padding: 16px;
  background: #ffffff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

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

.group-section {
  background: #ffffff;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}
</style>