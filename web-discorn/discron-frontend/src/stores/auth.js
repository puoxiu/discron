import { defineStore } from 'pinia';
import { getToken, setToken, removeToken } from '../utils/auth';
import { login as loginApi, register as registerApi } from '../api/auth';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: getToken() || '', // 从本地存储初始化token
    userInfo: null // 用户信息
  }),
  actions: {
    // 登录
    async login(userData) {
      const res = await loginApi(userData);
      // 存储token和用户信息（根据后端返回结构调整）
      this.token = res.data.token;
      this.userInfo = res.data.user;
      console.log("登陆成功，token:", res.data.token);
      console.log("登陆成功，userInfo:", res.data.user);
      setToken(this.token);
      return res;
    },
    // 注册
    async register(userData) {
      const res = await registerApi(userData);
      return res;
    },
    // 退出登录
    logout() {
      this.token = '';
      this.userInfo = null;
      removeToken();
    }
  }
});