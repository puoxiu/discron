import request from './request';


// 获取日志列表
export const getLogs = () => {
  return request({
    url: '/job/log',
    method: 'get'
  });
};
