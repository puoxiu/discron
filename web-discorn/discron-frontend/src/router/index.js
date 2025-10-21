import { createRouter, createWebHistory } from 'vue-router';
import { getToken } from '../utils/auth';

// 路由规则
const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('../views/Home.vue'),
    meta: { requiresAuth: true } // 需要登录才能访问
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/access/Login.vue'),
    meta: { requiresAuth: false } // 公开访问
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../views/access/Register.vue'),
    meta: { requiresAuth: false }
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

// 路由守卫（验证登录状态）
router.beforeEach((to, from, next) => {
  // 需要登录的路由，若未登录则跳转到登录页
  if (to.meta.requiresAuth && !getToken()) {
    next('/login');
  } else {
    next();
  }
});

export default router;