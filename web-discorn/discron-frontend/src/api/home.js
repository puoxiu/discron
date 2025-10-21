import request from './request';


// 获取首页今日统计数据
export const getTodayStatistics = () => {
  return request({
    url: '/statis/today',
    method: 'get'
  });
};

// 获取近7天任务执行统计
export const getWeekStatistics = () => {
  return request({
    url: '/statis/week',
    method: 'get'
  });
};

// 获取系统信息
export const getSystemInfo = (uuid = '') => {
  return request({
    url: '/statis/system',
    method: 'get',
    params: { uuid }
  });
};


export const getRecentJobLogs = (params = {}) => {
  return request({
    url: '/job/log',
    method: 'post',
    data: {
      page: params.page || 1,
      pageSize: params.pageSize || 10,
      // 默认查询所有状态，按执行时间倒序
      order: 'start_time desc'
    }
  });
};