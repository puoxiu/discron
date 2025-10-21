<template>
  <div class="nodes-page">
    <h2 class="page-title">节点管理</h2>
    <el-card>
      <template #header>
        <div class="card-header">
          <span>节点列表</span>
        </div>
      </template>
      
      <div class="search-container">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索节点名称"
          prefix-icon="Search"
          clearable
          class="search-input"
        />
        <el-select v-model="statusFilter" placeholder="节点状态" clearable>
          <el-option label="全部" value="" />
          <el-option label="在线" value="online" />
          <el-option label="离线" value="offline" />
        </el-select>
      </div>
      
      <div class="table-container">
        <el-table :data="filteredNodes" stripe style="width: 100%">
          <el-table-column prop="id" label="节点ID" width="80" />
          <el-table-column prop="nodeName" label="节点名称" />
          <el-table-column prop="ip" label="IP地址" />
          <el-table-column prop="port" label="端口" width="80" />
          <el-table-column prop="status" label="节点状态">
            <template #default="scope">
              <el-tag 
                :type="scope.row.status === 'online' ? 'success' : 'danger'"
              >
                {{ scope.row.status === 'online' ? '在线' : '离线' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="cpuUsage" label="CPU使用率">
            <template #default="scope">
              <div class="progress-container">
                <el-progress 
                  :percentage="scope.row.cpuUsage" 
                  :color="getProgressColor(scope.row.cpuUsage)"
                  :show-text="false"
                />
                <span class="progress-text">{{ scope.row.cpuUsage }}%</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="memoryUsage" label="内存使用率">
            <template #default="scope">
              <div class="progress-container">
                <el-progress 
                  :percentage="scope.row.memoryUsage" 
                  :color="getProgressColor(scope.row.memoryUsage)"
                  :show-text="false"
                />
                <span class="progress-text">{{ scope.row.memoryUsage }}%</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="runningJobs" label="运行任务数" width="120" />
          <el-table-column prop="lastHeartbeat" label="最后心跳时间" />
          <el-table-column label="操作" width="100" fixed="right">
            <template #default="scope">
              <el-button type="primary" link @click="handleViewDetail(scope.row)">
                详情
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { ElMessage } from 'element-plus';
import { getNodes } from '../../api/dashboard';

const nodes = ref([
  {
    id: 1,
    nodeName: 'node-1',
    ip: '192.168.1.101',
    port: 8001,
    status: 'online',
    cpuUsage: 30,
    memoryUsage: 45,
    runningJobs: 3,
    lastHeartbeat: '2024-01-15 10:30:00'
  },
  {
    id: 2,
    nodeName: 'node-2',
    ip: '192.168.1.102',
    port: 8001,
    status: 'online',
    cpuUsage: 60,
    memoryUsage: 75,
    runningJobs: 5,
    lastHeartbeat: '2024-01-15 10:29:00'
  }
]);

const searchKeyword = ref('');
const statusFilter = ref('');

// 筛选后的节点列表
const filteredNodes = computed(() => {
  return nodes.value.filter(node => {
    const keywordMatch = node.nodeName.toLowerCase().includes(searchKeyword.value.toLowerCase());
    const statusMatch = !statusFilter.value || node.status === statusFilter.value;
    return keywordMatch && statusMatch;
  });
});

// 获取进度条颜色
const getProgressColor = (usage) => {
  if (usage < 50) return '#67c23a';
  if (usage < 80) return '#e6a23c';
  return '#f56c6c';
};

// 获取节点列表
const fetchNodes = async () => {
  try {
    const res = await getNodes();
    nodes.value = res.data || [];
  } catch (error) {
    ElMessage.error('获取节点列表失败');
    console.error('获取节点列表失败:', error);
  }
};

// 查看节点详情
const handleViewDetail = (node) => {
  ElMessage.info('查看节点详情功能开发中');
};

// 定时刷新节点状态
let refreshTimer = null;

onMounted(() => {
  fetchNodes();
  // 每10秒刷新一次节点状态
  refreshTimer = setInterval(fetchNodes, 10000);
});

// 清理定时器
const onUnmounted = () => {
  if (refreshTimer) {
    clearInterval(refreshTimer);
  }
};
</script>

<style scoped>
.nodes-page {
  width: 100%;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-container {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
}

.search-input {
  width: 300px;
}

.table-container {
  overflow-x: auto;
}

.progress-container {
  display: flex;
  align-items: center;
  gap: 10px;
}

.progress-text {
  font-size: 12px;
  color: #606266;
  min-width: 40px;
}
</style>