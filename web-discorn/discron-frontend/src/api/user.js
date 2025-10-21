import request from './request';

export const getUserList = (data) => {
  return request({
    url: '/user/search',
    method: 'post',
    data
  });
};


export const addUser = (data) => {
  return request({
    url: '/register',
    method: 'post',
    data
  });
};


export const updateUser = (data) => {
  return request({
    url: '/user',
    method: 'put',
    data
  });
};


export const deleteUser = (ids) => {
  return request({
    url: '/user/del',
    method: 'post',
    data: { ids }
  });
};


export const changeUserStatus = (id, status) => {
  return request({
    url: '/user/status',
    method: 'patch',
    data: { id, status }
  });
};


export const getUserDetail = (id) => {
  return request({
    url: `/user/${id}`,
    method: 'get'
  });
};