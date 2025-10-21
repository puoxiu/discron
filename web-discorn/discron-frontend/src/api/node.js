import request from './request';

// 搜索节点列表
export const getNodeList = (data) => {
  return request({
    url: '/node/search',
    method: 'post',
    data
  });
};

// 删除节点
export const deleteNode = (uuid) => {
  return request({
    url: '/node/del',
    method: 'post',
    data: { uuid }
  });
};

// 获取系统信息
export const getSystemInfo = (uuid) => {
  return request({
    url: '/statis/system',
    method: 'get',
    params: { uuid }
  });
};