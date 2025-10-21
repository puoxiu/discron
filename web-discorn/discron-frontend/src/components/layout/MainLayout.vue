<template>
  <el-container class="main-layout">
    <side-bar />
    <el-container>
      <el-header class="header">
        <div class="header-right">
          <span class="user-info">欢迎回来{{ userInfo?.userName }}</span>
          <el-button type="danger" @click="handleLogout" size="small">
            <el-icon><circle-close /></el-icon>
            退出登录
          </el-button>
        </div>
      </el-header>
      <el-main class="main-content">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { useAuthStore } from '../../stores/auth';
import { useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';
import { CircleClose } from '@element-plus/icons-vue';
import SideBar from './SideBar.vue';

const authStore = useAuthStore();
const router = useRouter();
const userInfo = authStore.userInfo;

// 退出登录
const handleLogout = () => {
  authStore.logout();
  ElMessage.success('已退出登录');
  router.push('/login');
};
</script>

<style scoped>
.main-layout {
  height: 100vh;
}

.header {
  height: 60px;
  background-color: #fff;
  border-bottom: 1px solid #e4e7ed;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  padding: 0 20px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 20px;
}

.user-info {
  color: #303133;
  font-size: 14px;
}

.main-content {
  padding: 20px;
  overflow-y: auto;
  background-color: #f5f7fa;
}
</style>