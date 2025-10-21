<template>
  <div class="home">
    <h2 class="page-title">当前状态</h2>
    <div class="stats-grid">
      <stats-card
        v-for="card in statsCards"
        :key="card.label"
        :label="card.label"
        :value="card.value"
        :icon="card.icon"
        :type="card.type"
      />
    </div>
    
    <div class="recent-activities">
      <el-card>
        <template #header>
          <div class="card-header">
            <span>最近活动</span>
          </div>
        </template>
        <el-table :data="recentLogs" stripe style="width: 100%">
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="jobName" label="任务名称" />
          <el-table-column prop="status" label="状态">
            <template #default="scope">
              <el-tag :type="scope.row.status === 'success' ? 'success' : 'danger'">
                {{ scope.row.status === 'success' ? '成功' : '失败' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="nodeName" label="执行节点" />
          <el-table-column prop="createTime" label="执行时间" />
        </el-table>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { ElMessage } from 'element-plus';
import { Setting, Document, Operation, Warning } from '@element-plus/icons-vue';
import StatsCard from '../components/StatsCard.vue';
import { getTodayStatistics } from '../api/home';

// 统计数据卡片配置
const statsCards = ref([
  {
    label: '正常运行节点数',
    value: '0',
    icon: Setting,
    type: 'success'
  },
  {
    label: '异常节点数',
    value: '0',
    icon: Warning,
    type: 'danger'
  },
  {
    label: '正在执行任务数',
    value: '0',
    icon: Operation,
    type: 'warning'
  }
]);

// 最近活动日志
const recentLogs = ref([
  {
    id: 1,
    jobName: '示例任务1',
    status: 'success',
    nodeName: 'node-1',
    createTime: '2024-01-15 10:30:00'
  },
  {
    id: 2,
    jobName: '示例任务2',
    status: 'danger',
    nodeName: 'node-2',
    createTime: '2024-01-15 09:15:00'
  },
  {
    id: 3,
    jobName: '示例任务3',
    status: 'success',
    nodeName: 'node-1',
    createTime: '2024-01-15 08:45:00'
  }
]);

// 获取统计数据
const fetchStats = async () => {
  try {
    const todayStats = await getTodayStatistics();

    console.log('todayStats:', todayStats);

    const data = todayStats.data;
    // 更新统计卡片数据
    statsCards.value[0].value = data.normal_node_count || '0';
    statsCards.value[1].value = data.fail_node_count || '0';
    statsCards.value[2].value = data.job_running_count || '0';
    
    // 如果后端返回了最近日志，更新日志数据
    if (data.recentLogs) {
      recentLogs.value = data.recentLogs;
    }
  } catch (error) {
    ElMessage.error('获取统计数据失败');
    console.error('获取统计数据失败:', error);
  }
};

// 组件挂载时获取数据
onMounted(() => {
  fetchStats();
  // 每分钟刷新一次数据
  setInterval(fetchStats, 60000);
});
</script>

<style scoped>
.home {
  width: 100%;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 20px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 20px;
  margin-bottom: 24px;
}

.recent-activities {
  margin-top: 24px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
}
</style>