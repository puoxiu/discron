import request from './request';

// 获取首页统计数据
export const getDashboardStats = () => {
  return request({
    url: '/dashboard/stats',
    method: 'get'
  });
};

// 获取节点列表（用于节点管理页面）
export const getNodes = () => {
  return request({
    url: '/node',
    method: 'get'
  });
};

// 获取任务列表（用于任务管理页面）
export const getJobs = (params) => {
  return request({
    url: '/job',
    method: 'get',
    params
  });
};

// 获取日志列表（用于日志页面）
export const getLogs = (params) => {
  return request({
    url: '/logs',
    method: 'get',
    params
  });
};
