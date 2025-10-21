import request from './request';

// 登录接口
export function login(data) {
  return request({
    url: '/login', // 后端登录接口地址（对应router.go中的base.POST("login")）
    method: 'post',
    data
  });
}

// 注册接口
export function register(data) {
  return request({
    url: '/register', // 后端注册接口地址
    method: 'post',
    data
  });
}