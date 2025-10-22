import request from './request';


// 获取日志列表
export const getLogs = (data) => {
  return request({
    url: '/job/log',
    method: 'post',
    data
  });
};
