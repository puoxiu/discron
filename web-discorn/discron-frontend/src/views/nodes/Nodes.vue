<template>
  <div class="nodes-page">
    <!-- 1. 移除标题《节点信息》 -->

    <!-- 主内容区：节点列表 + 详情大屏 -->
    <div class="main-container">
      <!-- 1. 节点列表区 -->
      <el-card class="nodes-list-card">
        <template #header>
          <div class="card-header">
            <span>节点列表</span>
            <el-button type="primary" @click="fetchNodeList">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
          </div>
        </template>

        <!-- 搜索区域（保留原有） -->
        <div class="search-container">
          <el-row :gutter="20" class="mb-4">
            <el-col :span="6">
              <el-input v-model="searchForm.uuid" placeholder="请输入节点UUID" clearable />
            </el-col>
            <el-col :span="6">
              <el-input v-model="searchForm.ip" placeholder="请输入节点IP地址" clearable />
            </el-col>
            <el-col :span="2">
              <el-button type="primary" @click="handleSearch">搜索</el-button>
            </el-col>
            <el-col :span="2">
              <el-button @click="handleReset">重置</el-button>
            </el-col>
          </el-row>
        </div>

        <!-- 节点卡片网格（删除“最后心跳”“创建时间”，减少空间占用） -->
        <div v-loading="loading" class="nodes-grid">
          <el-card 
            v-for="node in nodes" 
            :key="node.uuid" 
            class="node-card"
            :class="[
              node.status === 1 ? 'status-online' : 'status-offline',
              getJobCountClass(node.jobCount)
            ]"
            shadow="hover"
          >
            <div class="node-header">
              <div class="node-title">
                <el-icon class="node-icon"><Platform /></el-icon>
                <span class="node-name">{{ node.name || node.uuid.substring(0, 8) }}</span>
              </div>
              <el-tag 
                :type="node.status === 1 ? 'success' : 'danger'" 
                size="small"
              >
                {{ node.status === 1 ? '在线' : '离线' }}
              </el-tag>
            </div>

            <div class="node-info">
              <!-- 保留UUID、IP、任务数量，删除最后心跳和创建时间 -->
              <div class="info-item">
                <span class="label">UUID:</span>
                <el-tooltip :content="node.uuid" placement="top">
                  <span class="value">{{ node.uuid.substring(0, 16) }}...</span>
                </el-tooltip>
              </div>
              <div class="info-item">
                <span class="label">IP地址:</span>
                <span class="value">{{ node.ip || 'N/A' }}</span>
              </div>
              <div class="info-item">
                <span class="label">任务数量:</span>
                <span class="value">
                  <el-tag size="small" type="info">{{ node.jobCount || 0 }}</el-tag>
                </span>
              </div>
            </div>

            <div class="node-actions">
              <el-button 
                type="primary" 
                link 
                @click="viewSystemInfo(node)"
                :disabled="node.status !== 1"
              >
                <el-icon><View /></el-icon>
                查看详情
              </el-button>
            </div>
          </el-card>

          <el-empty v-if="!loading && nodes.length === 0" description="暂无节点数据" />
        </div>
      </el-card>

      <!-- 2. 详情大屏（完全保留原有逻辑） -->
      <div v-if="detailVisible" class="detail-fullscreen">
        <div class="detail-header">
          <h3 class="detail-title">
            <el-icon><Monitor /></el-icon>
            节点详情 - {{ currentNode?.name || currentNode?.uuid.substring(0, 8) }}
          </h3>
          <el-button type="primary" @click="closeDetail">
            <el-icon><ArrowLeft /></el-icon>
            返回列表
          </el-button>
          <el-button @click="refreshSystemInfo">
            <el-icon><Refresh /></el-icon>
            刷新数据
          </el-button>
        </div>

        <div v-loading="systemInfoLoading" class="detail-loading">
          <div class="detail-content">
            <div class="chart-row">
              <div class="chart-card">
                <h4>CPU使用率</h4>
                <div class="chart-container" ref="cpuChartRef"></div>
              </div>
              <div class="chart-card">
                <h4>内存使用率</h4>
                <div class="chart-container" ref="ramChartRef"></div>
              </div>
            </div>

            <div class="chart-row">
              <div class="chart-card">
                <h4>磁盘使用率</h4>
                <div class="chart-container" ref="diskChartRef"></div>
              </div>
              <div class="info-card">
                <h4>系统信息</h4>
                <div class="info-table">
                  <div class="info-row">
                    <span class="info-label">操作系统</span>
                    <span class="info-value">{{ systemInfo?.os?.goos || 'N/A' }}</span>
                  </div>
                  <div class="info-row">
                    <span class="info-label">CPU核心数</span>
                    <span class="info-value">{{ systemInfo?.cpu?.cores || 'N/A' }}</span>
                  </div>
                  <div class="info-row">
                    <span class="info-label">总内存</span>
                    <span class="info-value">{{ systemInfo?.ram?.totalMb ? `${systemInfo.ram.totalMb} MB (${formatMBtoGB(systemInfo.ram.totalMb)})` : 'N/A' }}</span>
                  </div>
                  <div class="info-row">
                    <span class="info-label">总磁盘</span>
                    <span class="info-value">{{ systemInfo?.disk?.totalGb ? `${systemInfo.disk.totalGb} GB` : 'N/A' }}</span>
                  </div>
                  <div class="info-row">
                    <span class="info-label">Go版本</span>
                    <span class="info-value">{{ systemInfo?.os?.goVersion || 'N/A' }}</span>
                  </div>
                  <div class="info-row">
                    <span class="info-label">协程数量</span>
                    <span class="info-value">{{ systemInfo?.os?.numGoroutine || 'N/A' }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
// 脚本部分完全保留，未做任何修改
import { ref, onMounted, onUnmounted } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { 
  Platform, Check, Close, Refresh, View, Delete, 
  Monitor, Coin, Files, InfoFilled, ArrowLeft 
} from '@element-plus/icons-vue';
import { getNodeList, deleteNode, getSystemInfo } from '../../api/node';
import * as echarts from 'echarts';

const nodes = ref([]);
const loading = ref(false);
const systemInfoLoading = ref(false);
const detailVisible = ref(false);
const systemInfo = ref(null);
const currentNode = ref(null);

const searchForm = ref({
  uuid: '',
  ip: ''
});

const cpuChartRef = ref(null);
const ramChartRef = ref(null);
const diskChartRef = ref(null);
const cpuChart = ref(null);
const ramChart = ref(null);
const diskChart = ref(null);

const formatTime = (timestamp) => {
  if (!timestamp) return 'N/A';
  const date = new Date(timestamp * 1000);
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  });
};

const formatMBtoGB = (mb) => {
  if (!mb) return '0 GB';
  return (mb / 1024).toFixed(2) + ' GB';
};

const getJobCountClass = (jobCount = 0) => {
  if (jobCount === 0) return 'job-none';
  if (jobCount <= 5) return 'job-few';
  if (jobCount <= 10) return 'job-many';
  return 'job-much';
};

const fetchNodeList = async () => {
  loading.value = true;
  try {
    const params = {
      uuid: searchForm.value.uuid,
      ip: searchForm.value.ip
    };

    const res = await getNodeList(params);
    if (res.code === 200) {
      nodes.value = res.data.list || [];
    } else {
      ElMessage.error(res.msg || '获取节点列表失败');
    }
  } catch (error) {
    ElMessage.error('获取节点列表失败');
    console.error('获取节点列表失败:', error);
  } finally {
    loading.value = false;
  }
};

const handleSearch = () => {
  fetchNodeList();
};

const handleReset = () => {
  searchForm.value = {
    uuid: '',
    ip: ''
  };
  fetchNodeList();
};

const viewSystemInfo = async (node) => {
  currentNode.value = node;
  detailVisible.value = true;
  await loadSystemInfo(node.uuid);
  setTimeout(() => {
    initCharts();
    updateCharts();
  }, 100);
};

const closeDetail = () => {
  detailVisible.value = false;
  systemInfo.value = null;
  currentNode.value = null;
  destroyCharts();
};

const loadSystemInfo = async (uuid) => {
  systemInfoLoading.value = true;
  try {
    const res = await getSystemInfo(uuid);
    if (res.code === 200) {
      systemInfo.value = res.data;
      updateCharts();
    } else {
      ElMessage.error(res.msg || '获取系统信息失败');
    }
  } catch (error) {
    ElMessage.error('获取系统信息失败');
    console.error('获取系统信息失败:', error);
  } finally {
    systemInfoLoading.value = false;
  }
};

const refreshSystemInfo = async () => {
  if (currentNode.value) {
    await loadSystemInfo(currentNode.value.uuid);
    ElMessage.success('刷新成功');
  }
};

const initCharts = () => {
  if (cpuChartRef.value && !cpuChart.value) {
    cpuChart.value = echarts.init(cpuChartRef.value);
    cpuChart.value.setOption({
      tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
      grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
      xAxis: { type: 'category', data: [] },
      yAxis: { type: 'value', max: 100, name: '使用率(%)' },
      series: [{
        name: 'CPU核心使用率',
        type: 'bar',
        data: [],
        itemStyle: { color: '#409eff' },
        barWidth: '60%'
      }]
    });
  }

  if (ramChartRef.value && !ramChart.value) {
    ramChart.value = echarts.init(ramChartRef.value);
    ramChart.value.setOption({
      tooltip: { formatter: '{a} <br/>{b}: {c}%' },
      series: [{
        name: '内存使用率',
        type: 'gauge',
        radius: '80%',
        startAngle: 90,
        endAngle: -270,
        pointer: { show: false },
        progress: {
          show: true,
          overlap: false,
          roundCap: true,
          clip: false,
          itemStyle: { color: '#67c23a' }
        },
        axisLine: {
          lineStyle: { width: 12 }
        },
        splitLine: { show: false },
        axisTick: { show: false },
        axisLabel: { show: false },
        title: { show: false },
        detail: {
          valueAnimation: true,
          fontSize: 20,
          offsetCenter: [0, 0],
          formatter: '{value}%',
          color: '#303133'
        },
        data: [{ value: 0, name: '使用率' }]
      }]
    });
  }

  if (diskChartRef.value && !diskChart.value) {
    diskChart.value = echarts.init(diskChartRef.value);
    diskChart.value.setOption({
      tooltip: { formatter: '{a} <br/>{b}: {c}%' },
      series: [{
        name: '磁盘使用率',
        type: 'gauge',
        radius: '80%',
        startAngle: 90,
        endAngle: -270,
        pointer: { show: false },
        progress: {
          show: true,
          overlap: false,
          roundCap: true,
          clip: false,
          itemStyle: { color: '#e6a23c' }
        },
        axisLine: {
          lineStyle: { width: 12 }
        },
        splitLine: { show: false },
        axisTick: { show: false },
        axisLabel: { show: false },
        title: { show: false },
        detail: {
          valueAnimation: true,
          fontSize: 20,
          offsetCenter: [0, 0],
          formatter: '{value}%',
          color: '#303133'
        },
        data: [{ value: 0, name: '使用率' }]
      }]
    });
  }
};

const updateCharts = () => {
  const info = systemInfo.value;
  if (!info) return;

  if (cpuChart.value && info.cpu?.cpus) {
    const cpuData = info.cpu.cpus.map((p, i) => Number(p.toFixed(2)));
    const xData = info.cpu.cpus.map((_, i) => `核心${i + 1}`);
    cpuChart.value.setOption({
      xAxis: { data: xData },
      series: [{ data: cpuData }]
    });
  }

  if (ramChart.value && info.ram?.usedPercent) {
    const ramPercent = Number(info.ram.usedPercent.toFixed(2));
    const color = ramPercent < 60 ? '#67c23a' : ramPercent < 80 ? '#e6a23c' : '#f56c6c';
    ramChart.value.setOption({
      series: [{
        progress: { itemStyle: { color } },
        data: [{ value: ramPercent }]
      }]
    });
  }

  if (diskChart.value && info.disk?.usedPercent) {
    const diskPercent = Number(info.disk.usedPercent.toFixed(2));
    const color = diskPercent < 60 ? '#67c23a' : diskPercent < 80 ? '#e6a23c' : '#f56c6c';
    diskChart.value.setOption({
      series: [{
        progress: { itemStyle: { color } },
        data: [{ value: diskPercent }]
      }]
    });
  }
};

const destroyCharts = () => {
  if (cpuChart.value) {
    cpuChart.value.dispose();
    cpuChart.value = null;
  }
  if (ramChart.value) {
    ramChart.value.dispose();
    ramChart.value = null;
  }
  if (diskChart.value) {
    diskChart.value.dispose();
    diskChart.value = null;
  }
};

const handleResize = () => {
  cpuChart.value && cpuChart.value.resize();
  ramChart.value && ramChart.value.resize();
  diskChart.value && diskChart.value.resize();
};

onMounted(() => {
  fetchNodeList();
  window.addEventListener('resize', handleResize);
});

onUnmounted(() => {
  destroyCharts();
  window.removeEventListener('resize', handleResize);
});
</script>

<style scoped>
/* 样式部分仅删除标题相关样式，其余保留 */
.nodes-page {
  width: 100%;
  padding: 20px;
  box-sizing: border-box;
}

/* 2. 移除标题样式（原.page-title类删除） */

/* 卡片头部样式（保留原有） */
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.search-container {
  margin-bottom: 16px;
}

/* 节点卡片样式（因删除2个info-item，卡片自动缩小，无需额外改样式） */
.nodes-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
  min-height: 200px;
}

.node-card {
  transition: all 0.3s;
  border-radius: 8px;
  border-top-width: 4px;
  border-top-style: solid;
}

.status-online {
  border-top-color: #67c23a;
}
.status-offline {
  border-top-color: #f56c6c;
  opacity: 0.8;
}

.job-none {
  background-color: #f0f9ff;
}
.job-few {
  background-color: #ecfdf5;
}
.job-many {
  background-color: #fffbeb;
}
.job-much {
  background-color: #fff2f0;
}

.node-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.node-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
  padding-bottom: 12px;
  border-bottom: 1px solid #ebeef5;
}

.node-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.node-icon {
  font-size: 20px;
  color: #409eff;
}

.node-name {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.node-info {
  margin-bottom: 15px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  font-size: 13px;
}

.info-item .label {
  color: #909399;
  font-weight: 500;
}

.info-item .value {
  color: #606266;
  text-align: right;
  max-width: 180px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.node-actions {
  display: flex;
  justify-content: space-around;
  padding-top: 12px;
  border-top: 1px solid #ebeef5;
}

/* 详情大屏样式（完全保留） */
.detail-fullscreen {
  width: 100%;
  min-height: calc(100vh - 80px);
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  padding: 20px;
  box-sizing: border-box;
}

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 12px;
  border-bottom: 1px solid #ebeef5;
}

.detail-title {
  font-size: 20px;
  font-weight: 600;
  color: #303133;
  display: flex;
  align-items: center;
  gap: 8px;
}

.detail-loading {
  width: 100%;
  min-height: calc(100vh - 180px);
}

.detail-content {
  width: 100%;
}

.chart-row {
  display: flex;
  gap: 20px;
  margin-bottom: 20px;
}

.chart-card {
  flex: 1;
  background-color: #f5f7fa;
  border-radius: 8px;
  padding: 16px;
  box-sizing: border-box;
}

.chart-card h4 {
  margin: 0 0 16px 0;
  font-size: 16px;
  color: #303133;
}

.chart-container {
  width: 100%;
  height: 300px;
}

.info-card {
  flex: 1;
  background-color: #f5f7fa;
  border-radius: 8px;
  padding: 16px;
  box-sizing: border-box;
}

.info-card h4 {
  margin: 0 0 16px 0;
  font-size: 16px;
  color: #303133;
}

.info-table {
  width: 100%;
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px dashed #ebeef5;
}

.info-row:last-child {
  border-bottom: none;
}

.info-label {
  font-size: 14px;
  color: #909399;
  font-weight: 500;
}

.info-value {
  font-size: 14px;
  color: #303133;
}
</style>