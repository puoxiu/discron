<template>
  <div class="jobs-page">
    <h2 class="page-title">任务管理</h2>
    <el-card>
      <template #header>
        <div class="card-header">
          <span>任务列表</span>
          <el-button type="primary" @click="handleAddJob">
            <el-icon><plus /></el-icon>
            创建任务
          </el-button>
        </div>
      </template>
      
      <div class="search-container">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索任务名称"
          prefix-icon="Search"
          clearable
          class="search-input"
        />
        <el-select v-model="statusFilter" placeholder="任务状态" clearable>
          <el-option label="全部" value="" />
          <el-option label="运行中" value="running" />
          <el-option label="已停止" value="stopped" />
          <el-option label="已完成" value="completed" />
        </el-select>
      </div>
      
      <div class="table-container">
        <el-table :data="filteredJobs" stripe style="width: 100%">
          <el-table-column prop="id" label="任务ID" width="80" />
          <el-table-column prop="jobName" label="任务名称" />
          <el-table-column prop="cronExpr" label="Cron表达式" />
          <el-table-column prop="command" label="执行命令" />
          <el-table-column prop="status" label="状态">
            <template #default="scope">
              <el-tag 
                :type="getTagType(scope.row.status)"
              >
                {{ getStatusText(scope.row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="createTime" label="创建时间" />
          <el-table-column prop="lastRunTime" label="上次执行时间" />
          <el-table-column label="操作" width="180" fixed="right">
            <template #default="scope">
              <el-button 
                :type="scope.row.status === 'running' ? 'warning' : 'success'"
                @click="handleToggleJob(scope.row)"
                size="small"
              >
                {{ scope.row.status === 'running' ? '停止' : '启动' }}
              </el-button>
              <el-button type="primary" link @click="handleEditJob(scope.row)">
                编辑
              </el-button>
              <el-button type="danger" link @click="handleDeleteJob(scope.row.id)">
                删除
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
import { Plus } from '@element-plus/icons-vue';
import { getJobs } from '../../api/dashboard';

const jobs = ref([
  {
    id: 1,
    jobName: '示例任务1',
    cronExpr: '0/5 * * * * ?',
    command: 'echo hello world',
    status: 'running',
    createTime: '2024-01-01 10:00:00',
    lastRunTime: '2024-01-15 10:30:00'
  }
]);

const searchKeyword = ref('');
const statusFilter = ref('');

// 筛选后的任务列表
const filteredJobs = computed(() => {
  return jobs.value.filter(job => {
    const keywordMatch = job.jobName.toLowerCase().includes(searchKeyword.value.toLowerCase());
    const statusMatch = !statusFilter.value || job.status === statusFilter.value;
    return keywordMatch && statusMatch;
  });
});

// 获取状态标签类型
const getTagType = (status) => {
  switch (status) {
    case 'running':
      return 'success';
    case 'stopped':
      return 'warning';
    case 'completed':
      return 'info';
    default:
      return 'primary';
  }
};

// 获取状态文本
const getStatusText = (status) => {
  switch (status) {
    case 'running':
      return '运行中';
    case 'stopped':
      return '已停止';
    case 'completed':
      return '已完成';
    default:
      return status;
  }
};

// 获取任务列表
const fetchJobs = async () => {
  try {
    const res = await getJobs();
    jobs.value = res.data || [];
  } catch (error) {
    ElMessage.error('获取任务列表失败');
    console.error('获取任务列表失败:', error);
  }
};

// 处理添加任务
const handleAddJob = () => {
  ElMessage.info('创建任务功能开发中');
};

// 处理编辑任务
const handleEditJob = (job) => {
  ElMessage.info('编辑任务功能开发中');
};

// 处理删除任务
const handleDeleteJob = (jobId) => {
  ElMessage.info('删除任务功能开发中');
};

// 处理启动/停止任务
const handleToggleJob = (job) => {
  ElMessage.info(`${job.status === 'running' ? '停止' : '启动'}任务功能开发中`);
};

onMounted(() => {
  fetchJobs();
});
</script>

<style scoped>
.jobs-page {
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
</style>