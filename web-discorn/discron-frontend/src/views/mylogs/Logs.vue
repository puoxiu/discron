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
          
          <el-form-item>
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
                    <!-- 用开始/结束时间计算时长（后端未返回time_cost） -->
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
          <!-- 修复：重试次数字段适配后端retry_times -->
          <el-table-column label="重试次数" width="100">
            <template #default="scope">
              {{ scope.row.retry_times || 0 }}
            </template>
          </el-table-column>
          <!-- 修复：时间字段改用后端start_time -->
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
import { ref, onMounted } from 'vue';
import { ElMessage } from 'element-plus';
import { getLogs } from '../../api/log';

// 数据定义
const logs = ref([]);
const total = ref(0);
const searchKeyword = ref('');
const statusFilter = ref('');
const dateRange = ref(null);
const currentPage = ref(1);
const pageSize = ref(10);
const searchForm = ref({}); // 用于el-form的model

// 获取日志列表
const fetchLogs = async () => {
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
    };
    
    // 添加任务名称筛选
    if (searchKeyword.value) {
      params.name = searchKeyword.value;
    }
    
    // 添加状态筛选
    if (statusFilter.value) {
      params.success = statusFilter.value === 'success';
    }
    
    // 添加日期范围筛选（适配后端start_time字段，转秒级时间戳）
    if (dateRange.value && dateRange.value.length === 2) {
      params.start_time = Math.floor(dateRange.value[0].getTime() / 1000);
      params.end_time = Math.floor(dateRange.value[1].getTime() / 1000);
    }
    
    console.log('请求参数:', params);
    const res = await getLogs(params);
    console.log('接口返回结果:', res);
    
    if (res.code === 200 && res.data) {
      logs.value = res.data.list || [];
      total.value = res.data.total || 0;
      
      // 调试：确认关键字段
      if (logs.value.length > 0) {
        console.log('时间字段:', logs.value[0].start_time, '类型:', typeof logs.value[0].start_time);
        console.log('重试次数字段:', logs.value[0].retry_times);
      }
      
      ElMessage.success(`加载成功，共 ${total.value} 条日志`);
    } else {
      ElMessage.error(res.msg || '获取日志列表失败');
    }
  } catch (error) {
    ElMessage.error('获取日志列表失败');
    console.error('获取日志列表失败:', error);
  }
};

// 搜索处理
const handleSearch = () => {
  currentPage.value = 1; // 搜索时重置到第一页
  fetchLogs();
};

// 重置搜索条件
const handleReset = () => {
  searchKeyword.value = '';
  statusFilter.value = '';
  dateRange.value = null;
  currentPage.value = 1;
  fetchLogs();
};

// 分页处理
const handleSizeChange = (size) => {
  pageSize.value = size;
  currentPage.value = 1;
  fetchLogs();
};

const handleCurrentChange = (current) => {
  currentPage.value = current;
  fetchLogs();
};

// 格式化时间（适配后端秒级时间戳）
const formatTime = (timestamp) => {
  // 后端返回0表示无时间，直接显示"-"
  if (!timestamp || timestamp === 0) return '-';
  
  // 确认是秒级时间戳（10位），转成毫秒给Date
  const msTimestamp = typeof timestamp === 'number' ? timestamp * 1000 : Number(timestamp) * 1000;
  const date = new Date(msTimestamp);
  
  // 检查时间有效性
  if (isNaN(date.getTime())) return '-';
  
  // 格式化输出
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  const seconds = String(date.getSeconds()).padStart(2, '0');
  
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
};

// 计算执行时长（用start_time和end_time计算，适配后端字段）
const calculateDuration = (startTime, endTime) => {
  // 后端end_time为0表示未结束，显示"执行中"
  if (startTime === 0 || endTime === 0) return '执行中';
  
  // 计算秒数差，转成"分:秒"格式
  const durationSec = endTime - startTime;
  if (durationSec < 0) return '00:00';
  
  const minutes = Math.floor(durationSec / 60);
  const seconds = durationSec % 60;
  
  return `${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`;
};

// 初始化加载
onMounted(() => {
  fetchLogs();
});
</script>

<style scoped>
/* 样式保持不变，沿用之前的 */
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
  margin-bottom: 20px;
}

.search-form {
  display: flex;
  flex-wrap: wrap;
  gap: 0;
}

.search-form :deep(.el-form-item) {
  margin-bottom: 16px;
  margin-right: 24px;
}

.search-form :deep(.el-form-item__label) {
  font-weight: 500;
  color: #606266;
}

.search-input {
  width: 240px;
}

.search-select {
  width: 180px;
}

.search-date {
  width: 360px;
}

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

:deep(.el-table__expand-icon) {
  font-size: 16px;
}

:deep(.el-descriptions__label) {
  font-weight: 600;
  background-color: #fafafa;
}

:deep(.el-table__expanded-cell) {
  padding: 0;
}
</style>