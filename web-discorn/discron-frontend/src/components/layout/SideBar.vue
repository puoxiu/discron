<template>
  <el-aside width="240px" class="sidebar">
    <div class="logo">
      <h1>DistCron</h1>
    </div>
    <el-menu
      :default-active="activeMenu"
      class="sidebar-menu"
      background-color="#001529"
      text-color="#fff"
      active-text-color="#409EFF"
      router
    >
      <el-menu-item index="/">
        <el-icon><home-filled /></el-icon>
        <span>首页</span>
      </el-menu-item>
      <el-menu-item index="/users">
        <el-icon><user /></el-icon>
        <span>用户管理</span>
      </el-menu-item>
      <el-menu-item index="/jobs">
        <el-icon><document /></el-icon>
        <span>任务管理</span>
      </el-menu-item>
      <el-menu-item index="/logs">
        <el-icon><message /></el-icon>
        <span>日志</span>
      </el-menu-item>
      <el-menu-item index="/nodes">
        <el-icon><link /></el-icon>
        <span>节点管理</span>
      </el-menu-item>
    </el-menu>
  </el-aside>
</template>

<script setup>
import { computed } from 'vue';
import { useRoute } from 'vue-router';
import { HomeFilled, User, Document, Message, Link } from '@element-plus/icons-vue';

const route = useRoute();

// 计算当前活动菜单
const activeMenu = computed(() => {
  return route.path;
});
</script>

<style scoped>
.sidebar {
  height: 100vh;
  background-color: #001529;
  overflow-y: auto;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  border-bottom: 1px solid #1f2937;
}

.logo h1 {
  font-size: 20px;
  font-weight: 600;
}

.sidebar-menu {
  border-right: none;
  height: calc(100vh - 60px);
}

.el-menu-item {
  height: 60px;
  line-height: 60px;
}
</style>