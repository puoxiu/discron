<template>
  <div class="logs-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>执行日志</span>
        </div>
      </template>
      
      <div class="search-container">
        <el-form :inline="true" :model="searchForm" class="search-form">
          <el-form-item label="任务名称">
            <el-input
              v-model="searchKeyword"
              placeholder="请输入任务名称"
              prefix-icon="Search"
              clearable
              class="search-input"
              @clear="handleSearch"
              @keyup.enter="handleSearch"
            />
          </el-form-item>
          
          <el-form-item label="执行状态">
            <el-select 
              v-model="statusFilter" 
              placeholder="请选择执行状态" 
              clearable
              class="search-select"
              @change="handleSearch"
            >
              <el-option label="全部" value="" />
              <el-option label="成功" value="success" />
              <el-option label="失败" value="failed" />
            </el-select>
          </el-form-item>
          
          <el-form-item label="执行时间">
            <el-date-picker
              v-model="dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              clearable
              class="search-date"
              @change="handleSearch"
            />
          </el-form-item>
          
          <el-form-item class="button-group-item">
            <el-button type="primary" icon="Search" @click="handleSearch">搜索</el-button>
            <el-button icon="Refresh" @click="handleReset">重置</el-button>
          </el-form-item>
        </el-form>
      </div>
      
      <div class="table-container">
        <el-table 
          :data="logs" 
          stripe 
          style="width: 100%"
          row-key="id"
        >
          <!-- 展开列 -->
          <el-table-column type="expand">
            <template #default="scope">
              <div class="expand-content">
                <el-descriptions :column="2" border size="default">
                  <el-descriptions-item label="日志ID">
                    {{ scope.row.id }}
                  </el-descriptions-item>
                  <el-descriptions-item label="任务ID">
                    {{ scope.row.job_id }}
                  </el-descriptions-item>
                  <el-descriptions-item label="执行命令" :span="2">
                    <el-text class="command-text">{{ scope.row.command || '-' }}</el-text>
                  </el-descriptions-item>
                  <el-descriptions-item label="节点IP">
                    {{ scope.row.ip || '-' }}
                  </el-descriptions-item>
                  <el-descriptions-item label="节点UUID">
                    {{ scope.row.node_uuid || '-' }}
                  </el-descriptions-item>
                  <el-descriptions-item label="执行时长">
                    {{ calculateDuration(scope.row.start_time, scope.row.end_time) }}
                  </el-descriptions-item>
                  <el-descriptions-item label="执行时间">
                    {{ formatTime(scope.row.start_time) }}
                  </el-descriptions-item>
                  <el-descriptions-item label="执行结果" :span="2">
                    <el-tag :type="scope.row.success ? 'success' : 'danger'">
                      {{ scope.row.success ? '执行成功' : '执行失败' }}
                    </el-tag>
                  </el-descriptions-item>
                  <el-descriptions-item label="重试次数">
                    {{ scope.row.retry_times || 0 }} 次
                  </el-descriptions-item>
                  <el-descriptions-item label="主机名">
                    {{ scope.row.hostname || '-' }}
                  </el-descriptions-item>
                  <el-descriptions-item label="调度规则" :span="2">
                    {{ scope.row.spec || '-' }}
                  </el-descriptions-item>
                  <!-- 新增：输出信息展示 -->
                  <el-descriptions-item label="输出信息" :span="2">
                    <el-text class="output-text">{{ scope.row.output || '-' }}</el-text>
                  </el-descriptions-item>
                </el-descriptions>
              </div>
            </template>
          </el-table-column>

          <!-- 主表格列：显示关键信息 -->
          <el-table-column prop="job_id" label="任务ID" width="100" />
          <el-table-column prop="name" label="任务名称" min-width="150" show-overflow-tooltip />
          <el-table-column prop="success" label="执行状态" width="100">
            <template #default="scope">
              <el-tag :type="scope.row.success ? 'success' : 'danger'">
                {{ scope.row.success ? '成功' : '失败' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="重试次数" width="100">
            <template #default="scope">
              {{ scope.row.retry_times || 0 }}
            </template>
          </el-table-column>
          <el-table-column label="执行时间" width="180">
            <template #default="scope">
              {{ formatTime(scope.row.start_time) }}
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
          :total="total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
// 脚本部分保持不变
import { ref, onMounted } from 'vue';
import { ElMessage } from 'element-plus';
import { getLogs } from '../../api/job';

const logs = ref([]);
const total = ref(0);
const searchKeyword = ref('');
const statusFilter = ref('');
const dateRange = ref(null);
const currentPage = ref(1);
const pageSize = ref(10);
const searchForm = ref({});

const fetchLogs = async () => {
  try {
    const params = { page: currentPage.value, page_size: pageSize.value };
    if (searchKeyword.value) params.name = searchKeyword.value;
    if (statusFilter.value) params.success = statusFilter.value === 'success';
    if (dateRange.value && dateRange.value.length === 2) {
      params.start_time = Math.floor(dateRange.value[0].getTime() / 1000);
      params.end_time = Math.floor(dateRange.value[1].getTime() / 1000);
    }

    const res = await getLogs(params);
    if (res.code === 200 && res.data) {
      logs.value = res.data.list || [];
      total.value = res.data.total || 0;
      ElMessage.success(`加载成功，共 ${total.value} 条日志`);
    } else {
      ElMessage.error(res.msg || '获取日志列表失败');
    }
  } catch (error) {
    ElMessage.error('获取日志列表失败');
    console.error(error);
  }
};

const handleSearch = () => { currentPage.value = 1; fetchLogs(); };
const handleReset = () => {
  searchKeyword.value = '';
  statusFilter.value = '';
  dateRange.value = null;
  currentPage.value = 1;
  fetchLogs();
};
const handleSizeChange = (size) => { pageSize.value = size; currentPage.value = 1; fetchLogs(); };
const handleCurrentChange = (current) => { currentPage.value = current; fetchLogs(); };

const formatTime = (timestamp) => {
  if (!timestamp || timestamp === 0) return '-';
  const ms = typeof timestamp === 'number' ? timestamp * 1000 : Number(timestamp) * 1000;
  const date = new Date(ms);
  if (isNaN(date.getTime())) return '-';
  return `${date.getFullYear()}-${String(date.getMonth()+1).padStart(2,0)}-${String(date.getDate()).padStart(2,0)} ${String(date.getHours()).padStart(2,0)}:${String(date.getMinutes()).padStart(2,0)}:${String(date.getSeconds()).padStart(2,0)}`;
};

const calculateDuration = (startTime, endTime) => {
  if (startTime === 0 || endTime === 0) return '执行中';
  const sec = endTime - startTime;
  if (sec < 0) return '00:00';
  return `${String(Math.floor(sec/60)).padStart(2,0)}:${String(sec%60).padStart(2,0)}`;
};

onMounted(() => fetchLogs());
</script>

<style scoped>
.logs-page { width: 100%; }

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-container {
  margin-bottom: 20px;
  width: 100%;
  overflow-x: auto;
}

.search-form {
  display: flex;
  flex-wrap: nowrap;
  align-items: center;
  gap: 12px;
  width: 100%;
  padding: 4px 0;
}

.search-form :deep(.el-form-item) {
  margin: 0;
  white-space: nowrap;
}

.search-form :deep(.button-group-item) {
  padding: 0;
  margin-left: 8px;
}

.search-form :deep(.el-form-item__label) {
  font-weight: 500;
  color: #606266;
  padding-right: 6px;
}

.search-input { width: 180px; }
.search-select { width: 140px; }
.search-date { width: 280px; }

.search-form :deep(.el-button) {
  vertical-align: middle;
  margin-right: 8px;
}
.search-form :deep(.el-button:last-child) { margin-right: 0; }

.table-container {
  overflow-x: auto;
  margin-bottom: 16px;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.expand-content {
  padding: 20px 20px 20px 48px;
  background-color: #f5f7fa;
}

.command-text {
  word-break: break-all;
  color: #409eff;
  font-family: 'Courier New', monospace;
}

/* 输出信息样式：支持多行显示并保留换行符 */
.output-text {
  word-break: break-all;
  white-space: pre-wrap; /* 保留文本中的换行符 */
  font-family: 'Courier New', monospace;
  line-height: 1.6; /* 增加行高提升可读性 */
  color: #303133;
  max-height: 300px; /* 限制最大高度 */
  overflow-y: auto; /* 内容过长时显示滚动条 */
  padding: 8px 12px;
  background-color: #fff;
  border-radius: 4px;
  border: 1px solid #eee;
}

:deep(.el-table__expand-icon) { font-size: 16px; }
:deep(.el-descriptions__label) {
  font-weight: 600;
  background-color: #fafafa;
}
:deep(.el-table__expanded-cell) { padding: 0; }
/* 调整描述项间距 */
:deep(.el-descriptions-item) {
  padding: 12px 16px;
}
</style>