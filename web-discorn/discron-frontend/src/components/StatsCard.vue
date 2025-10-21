<template>
  <el-card :body-style="{ padding: '20px' }" class="stats-card">
    <div class="stats-content">
      <div class="stats-icon" :class="iconClass">
        <component :is="icon" />
      </div>
      <div class="stats-info">
        <div class="stats-value">{{ value }}</div>
        <div class="stats-label">{{ label }}</div>
      </div>
    </div>
  </el-card>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  label: {
    type: String,
    required: true
  },
  value: {
    type: [String, Number],
    required: true
  },
  icon: {
    type: Object,
    required: true
  },
  type: {
    type: String,
    default: 'primary'
  }
});

// 根据类型计算图标样式类
const iconClass = computed(() => {
  return `icon-${props.type}`;
});
</script>

<style scoped>
.stats-card {
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  transition: all 0.3s;
}

.stats-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px 0 rgba(0, 0, 0, 0.12);
}

.stats-content {
  display: flex;
  align-items: center;
  gap: 20px;
}

.stats-icon {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32px;
}

.icon-primary {
  background-color: #ecf5ff;
  color: #409eff;
}

.icon-success {
  background-color: #f0f9eb;
  color: #67c23a;
}

.icon-warning {
  background-color: #fdf6ec;
  color: #e6a23c;
}

.icon-danger {
  background-color: #fef0f0;
  color: #f56c6c;
}

.stats-info {
  flex: 1;
}

.stats-value {
  font-size: 32px;
  font-weight: 600;
  color: #303133;
  line-height: 1;
  margin-bottom: 8px;
}

.stats-label {
  font-size: 14px;
  color: #909399;
}
</style>