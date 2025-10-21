<script setup>
import { useRoute } from 'vue-router';
import MainLayout from './components/layout/MainLayout.vue';

const route = useRoute();

// 判断当前路由是否需要主布局（登录和注册页面不需要）
const needMainLayout = () => {
  const publicRoutes = ['/login', '/register'];
  return !publicRoutes.includes(route.path);
};
</script>

<template>
  <div v-if="needMainLayout()">
    <main-layout />
  </div>
  <router-view v-else />
</template>

<style>
/* 全局样式重置 */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body, #app {
  height: 100%;
  overflow: hidden;
}

/* Element Plus 滚动条样式优化 */
.el-scrollbar__wrap {
  overflow-x: hidden !important;
}
</style>