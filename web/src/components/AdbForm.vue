<template>
  <div class="adb-container">
    <!-- 顶部操作栏 -->
    <AdbToolbar
      :items="items"
      @add-item="addItem"
      @fill-ip="fillLocalIp"
      @refresh-devices="fetchConnectedDevices"
      @clear-logs="clearResultData"
      @reset-groups="resetGroups"
      @exec-free="execFreeCommand"
      @exec-direct="execFreeCommandWithDevice"
    />

    <!-- 在线设备快捷栏 -->
    <DeviceTagBar
      :devices="onlineDevices"
      @apply-device="applyDeviceToActive"
    />

    <!-- 结果与日志展示区 -->
    <LogConsole
      :logs="resultData"
      :spinning="spinning"
      @clear-logs="clearResultData"
    />

    <!-- 设备分组配置面板 -->
    <DeviceGroupList
      :items="items"
      v-model:activeKey="activeKey"
      @exec-op="handleGroupOp"
      @rename-group="renameGroup"
      @remove-group="removeGroup"
    />
  </div>
</template>

<script>
import { message, Modal } from 'ant-design-vue';
import { execAdbCommand, fetchDevices } from '@/api/adb';
import { fetchLocalIp } from '@/api/system';
import { DEFAULT_COMMAND_FORM, MAX_GROUP_COUNT } from '@/constants/commands';
import {
  loadStoredDeviceGroups,
  saveStoredDeviceGroups,
  clearStoredDeviceGroups
} from '@/utils/storage';
import AdbToolbar from './AdbToolbar.vue';
import DeviceTagBar from './DeviceTagBar.vue';
import LogConsole from './LogConsole.vue';
import DeviceGroupList from './DeviceGroupList.vue';

export default {
  name: 'AdbForm',
  components: {
    AdbToolbar,
    DeviceTagBar,
    LogConsole,
    DeviceGroupList
  },
  data() {
    const storedItems = loadStoredDeviceGroups();
    return {
      activeKey: ['0'],
      items: storedItems || [
        {
          title: '设备 1',
          commandForm: { ...DEFAULT_COMMAND_FORM }
        }
      ],
      spinning: false,
      resultData: [],
      onlineDevices: []
    };
  },
  watch: {
    items: {
      handler(newItems) {
        saveStoredDeviceGroups(newItems);
      },
      deep: true
    }
  },
  mounted() {
    this.fetchConnectedDevices();
  },
  methods: {
    formatTime() {
      const now = new Date();
      return now.toTimeString().split(' ')[0];
    },
    fetchConnectedDevices() {
      fetchDevices()
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

      execAdbCommand(formPayload)
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
    handleGroupOp({ commandForm, op, title }) {
      this.sendHttpRequest(commandForm, op, title);
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
      if (this.items.length >= MAX_GROUP_COUNT) {
        message.error(`最多添加 ${MAX_GROUP_COUNT} 个设备分组`);
        return;
      }
      const newTitle = prompt('请输入新设备组名称', `设备 ${this.items.length + 1}`);
      if (newTitle) {
        this.items.push({
          title: newTitle,
          commandForm: { ...DEFAULT_COMMAND_FORM }
        });
        this.activeKey = [String(this.items.length - 1)];
      }
    },
    removeGroup(index) {
      if (this.items.length <= 1) {
        message.error('至少保留一个设备分组');
        return;
      }
      this.items.splice(index, 1);
      this.activeKey = ['0'];
    },
    execFreeCommand(cmd) {
      this.sendHttpRequest({ cmd }, 'free');
    },
    execFreeCommandWithDevice({ item, cmd }) {
      const devIp = item.commandForm.ip;
      const port = item.commandForm.port || '5555';
      const target = devIp.includes(':') ? devIp : `${devIp}:${port}`;
      const fullCmd = `-s ${target} ${cmd}`;
      this.sendHttpRequest({ cmd: fullCmd }, 'free', item.title);
    },
    fillLocalIp() {
      const proxyPort = prompt('请输入抓包代理端口号 (例如: 8888, 7890)', '8888');
      if (proxyPort === null) return;

      fetchLocalIp()
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
    renameGroup(index) {
      if (!this.items[index]) return;
      const currentTitle = this.items[index].title || `设备 ${index + 1}`;
      const newTitle = prompt('请输入新的设备组名称', currentTitle);
      if (newTitle !== null) {
        const trimmed = newTitle.trim();
        if (!trimmed) {
          message.warning('设备组名称不能为空');
          return;
        }
        this.items[index].title = trimmed;
        message.success(`已重命名为: ${trimmed}`);
      }
    },
    resetGroups() {
      Modal.confirm({
        title: '重置设备组配置',
        content: '确定要清除本地保存的设备组配置并恢复为默认初始状态吗？',
        okText: '确认重置',
        okType: 'danger',
        cancelText: '取消',
        onOk: () => {
          clearStoredDeviceGroups();
          this.items = [
            {
              title: '设备 1',
              commandForm: { ...DEFAULT_COMMAND_FORM }
            }
          ];
          this.activeKey = ['0'];
          message.success('设备组配置已恢复为默认状态');
        }
      });
    }
  }
};
</script>

<style scoped>
.adb-container {
  text-align: left;
}
</style>