<template>
  <div class="home">
    <h1>欢迎回来，{{ userInfo?.userName }}</h1>
    <el-button type="danger" @click="handleLogout">退出登录</el-button>
  </div>
</template>

<script setup>
import { useAuthStore } from '../stores/auth';
import { useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';

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