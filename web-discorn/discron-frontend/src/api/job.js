import request from './request';


// 获取任务列表
export const getJobs = (data) => {
  return request({
    url: '/job/search',
    method: 'post',
    data
  });
};

// 获取日志列表
export const getLogs = (data) => {
  return request({
    url: '/job/log',
    method: 'post',
    data
  });
};
