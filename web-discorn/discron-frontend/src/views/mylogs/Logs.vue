<template>
  <div class="logs-page">
    <h2 class="page-title">日志管理</h2>
    <el-card>
      <template #header>
        <div class="card-header">
          <span>执行日志</span>
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
        <el-select v-model="statusFilter" placeholder="执行状态" clearable>
          <el-option label="全部" value="" />
          <el-option label="成功" value="success" />
          <el-option label="失败" value="failed" />
        </el-select>
        <el-date-picker
          v-model="dateRange"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          clearable
        />
      </div>
      
      <div class="table-container">
        <el-table :data="filteredLogs" stripe style="width: 100%">
          <el-table-column prop="id" label="日志ID" width="80" />
          <el-table-column prop="jobName" label="任务名称" />
          <el-table-column prop="jobId" label="任务ID" width="100" />
          <el-table-column prop="nodeName" label="执行节点" />
          <el-table-column prop="status" label="执行状态">
            <template #default="scope">
              <el-tag 
                :type="scope.row.status === 'success' ? 'success' : 'danger'"
              >
                {{ scope.row.status === 'success' ? '成功' : '失败' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="duration" label="执行时长(ms)" />
          <el-table-column prop="errorMsg" label="错误信息" show-overflow-tooltip>
            <template #default="scope">
              <el-popover
                placement="top"
                title="错误详情"
                trigger="hover"
                :width="400"
                v-if="scope.row.errorMsg"
              >
                <template #reference>
                  <span class="error-preview">{{ scope.row.errorMsg.substring(0, 20) }}...</span>
                </template>
                <pre>{{ scope.row.errorMsg }}</pre>
              </el-popover>
              <span v-else>-</span>
            </template>
          </el-table-column>
          <el-table-column prop="createTime" label="执行时间" />
          <el-table-column label="操作" width="120" fixed="right">
            <template #default="scope">
              <el-button type="primary" link @click="handleViewDetail(scope.row)">
                查看详情
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
      
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="filteredLogs.length"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { ElMessage } from 'element-plus';
import { getLogs } from '../../api/log';

const logs = ref([
  {
    id: 1,
    jobId: 1,
    jobName: '示例任务1',
    nodeName: 'node-1',
    status: 'success',
    duration: 120,
    errorMsg: '',
    createTime: '2024-01-15 10:30:00'
  },
  {
    id: 2,
    jobId: 2,
    jobName: '示例任务2',
    nodeName: 'node-2',
    status: 'failed',
    duration: 50,
    errorMsg: '执行失败：权限不足',
    createTime: '2024-01-15 09:15:00'
  }
]);

const searchKeyword = ref('');
const statusFilter = ref('');
const dateRange = ref(null);
const currentPage = ref(1);
const pageSize = ref(10);

// 筛选后的日志列表
const filteredLogs = computed(() => {
  return logs.value.filter(log => {
    const keywordMatch = log.jobName.toLowerCase().includes(searchKeyword.value.toLowerCase());
    const statusMatch = !statusFilter.value || log.status === statusFilter.value;
    
    // 日期范围筛选
    let dateMatch = true;
    if (dateRange.value && dateRange.value.length === 2) {
      const logDate = new Date(log.createTime);
      dateMatch = logDate >= dateRange.value[0] && logDate <= dateRange.value[1];
    }
    
    return keywordMatch && statusMatch && dateMatch;
  });
});

// 获取日志列表
const fetchLogs = async () => {
  try {
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value,
      keyword: searchKeyword.value,
      status: statusFilter.value
    };
    
    if (dateRange.value && dateRange.value.length === 2) {
      params.startDate = dateRange.value[0];
      params.endDate = dateRange.value[1];
    }
    
    const res = await getLogs(params);
    console.log('获取日志列表成功:', res);
    logs.value = res.data || [];
  } catch (error) {
    ElMessage.error('获取日志列表失败');
    console.error('获取日志列表失败:', error);
  }
};

// 查看详情
const handleViewDetail = (log) => {
  ElMessage.info('查看日志详情功能开发中');
};

// 分页处理
const handleSizeChange = (size) => {
  pageSize.value = size;
  fetchLogs();
};

const handleCurrentChange = (current) => {
  currentPage.value = current;
  fetchLogs();
};

onMounted(() => {
  fetchLogs();
});
</script>

<style scoped>
.logs-page {
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
  flex-wrap: wrap;
}

.search-input {
  width: 300px;
}

.table-container {
  overflow-x: auto;
  margin-bottom: 16px;
}

.error-preview {
  color: #f56c6c;
  cursor: pointer;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>