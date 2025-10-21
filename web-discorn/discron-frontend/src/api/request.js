import axios from 'axios';
import { ElMessage } from 'element-plus';
import { getToken, removeToken } from '../utils/auth';

// 创建axios实例
const service = axios.create({
  baseURL: 'http://localhost:8961', 
  timeout: 5000 // 请求超时时间
});

// 请求拦截器（添加token）
service.interceptors.request.use(
  (config) => {
    // 从本地存储获取token，添加到请求头
    const token = getToken();
    if (token) {
      config.headers.Authorization = token;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器（处理错误）
service.interceptors.response.use(
  (response) => {
    const res = response.data;
    // 后端错误码处理（根据实际后端定义调整）
    if (res.code !== 200) {
      ElMessage.error(res.message || '请求失败');
      return Promise.reject(res);
    }
    return res;
  },
  (error) => {
    // 处理401（未登录或token过期）
    if (error.response?.status === 401) {
      ElMessage.error('登录已过期，请重新登录');
      removeToken(); // 清除无效token
      window.location.href = '/login'; // 跳转到登录页
    } else {
      ElMessage.error(error.message || '网络错误');
    }
    return Promise.reject(error);
  }
);

export default service;