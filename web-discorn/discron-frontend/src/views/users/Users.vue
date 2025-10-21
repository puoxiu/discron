<template>
  <div class="users-page">
    <h2 class="page-title">用户管理</h2>
    <el-card>
      <template #header>
        <div class="card-header">
          <span>用户列表</span>
          <el-button type="primary" @click="handleAddUser">
            <el-icon>
              <Plus />
            </el-icon>
            添加用户
          </el-button>
        </div>
      </template>

      <!-- 搜索区域 -->
      <div class="search-container">
        <el-row :gutter="20" class="mb-4">
          <el-col :span="6">
            <el-input v-model="searchForm.userName" placeholder="请输入用户名" clearable />
          </el-col>
          <el-col :span="6">
            <el-input v-model="searchForm.email" placeholder="请输入邮箱" clearable />
          </el-col>
          <el-col :span="4">
            <el-select v-model="searchForm.role" placeholder="请选择角色" clearable>
              <el-option label="普通用户" value="1" />
              <el-option label="管理员" value="2" />
            </el-select>
          </el-col>
          <el-col :span="2">
            <el-button type="primary" @click="fetchUsers">搜索</el-button>
          </el-col>
        </el-row>
      </div>

      <!-- 表格区域 -->
      <div class="table-container">
        <el-table :data="users" stripe style="width: 100%">
          <el-table-column prop="id" label="用户ID" width="80" />
          <el-table-column prop="username" label="用户名" />
          <el-table-column prop="email" label="邮箱" />
          <el-table-column prop="role" label="角色">
            <template #default="scope">
              <span>{{ scope.row.role === 1 ? '普通用户' : '管理员' }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="created" label="创建时间">
            <template #default="scope">
              <span>{{ formatTime(scope.row.created) }}</span>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="150" fixed="right">
            <template #default="scope">
              <el-button type="primary" link @click="handleEditUser(scope.row)">
                编辑
              </el-button>
              <el-button type="danger" link @click="handleDeleteUser(scope.row.id)">
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>

        <!-- 添加用户弹窗 -->
        <el-dialog title="添加用户" :model-value="addDialogVisible" width="500px" @close="resetAddForm">
          <el-form :model="addForm" :rules="addFormRules" ref="addFormRef" label-width="100px">
            <el-form-item label="用户名" prop="username">
              <el-input v-model="addForm.username" placeholder="请输入用户名" />
            </el-form-item>
            <el-form-item label="密码" prop="password">
              <el-input type="password" v-model="addForm.password" placeholder="请输入密码" />
            </el-form-item>
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="addForm.email" placeholder="请输入邮箱" />
            </el-form-item>
            <el-form-item label="角色" prop="role">
              <el-select v-model="addForm.role" placeholder="请选择角色">
                <el-option label="普通用户" :value="1" />
                <el-option label="管理员" :value="2" />
              </el-select>
            </el-form-item>
          </el-form>
          <template #footer>
            <el-button @click="addDialogVisible = false">取消</el-button>
            <el-button type="primary" @click="submitAddUser">提交</el-button>
          </template>
        </el-dialog>
      </div>

      <!-- 分页控件 -->
      <div class="pagination-container">
        <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50]" :total="pagination.total" layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange" @current-change="handleCurrentChange" />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { Plus } from '@element-plus/icons-vue';
import { getUserList, deleteUser, addUser } from '../../api/user';

// 状态定义
const users = ref([]);
const total = ref(0);

// 添加用户弹窗状态
const addDialogVisible = ref(false);

// 搜索表单
const searchForm = ref({
  userName: '',
  email: '',
  role: '',
  id: 0
});

// 添加用户表单
const addForm = ref({
  username: '',
  password: '',
  email: '',
  role: 1
});

const addFormRef = ref(null);

// 表单校验规则
const addFormRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  email: [{ required: true, message: '请输入邮箱', trigger: 'blur' }],
  role: [{ required: true, message: '请选择角色', trigger: 'change' }]
};
// 打开添加用户弹窗
const handleAddUser = () => {
  addDialogVisible.value = true;
};
// 重置表单
const resetAddForm = () => {
  addForm.value = {
    username: '',
    password: '',
    email: '',
    role: 1
  };
  addFormRef.value?.clearValidate();
};

// 分页参数
const pagination = ref({
  page: 1,
  pageSize: 10,
  total: 0
});

// 格式化时间戳为日期
const formatTime = (timestamp) => {
  if (!timestamp) return '';
  const date = new Date(timestamp * 1000); // 后端返回的是秒级时间戳
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  });
};

// 获取用户列表
const fetchUsers = async () => {
  try {
    // 构造请求参数（与后端ReqUserSearch结构体对应）
    const params = {
      page: pagination.value.page,
      page_size: pagination.value.pageSize,
      username: searchForm.value.userName,
      email: searchForm.value.email,
      role: searchForm.value.role ? Number(searchForm.value.role) : 0,
      id: searchForm.value.id
    };

    const res = await getUserList(params);
    if (res.code === 200) {
      console.log(res.data.list);
      users.value = res.data.list || [];
      pagination.value.total = res.data.total || 0;
    } else {
      ElMessage.error(res.msg || '获取用户列表失败');
    }
  } catch (error) {
    ElMessage.error('获取用户列表失败');
    console.error('获取用户列表失败:', error);
  }
};

// 分页尺寸变更
const handleSizeChange = (size) => {
  pagination.value.pageSize = size;
  pagination.value.page = 1; // 重置为第一页
  fetchUsers();
};

// 当前页变更
const handleCurrentChange = (page) => {
  pagination.value.page = page;
  fetchUsers();
};

// 处理添加用户
const submitAddUser = () => {
  addFormRef.value.validate(async (valid) => {
    if (!valid) return;
    try {
      const res = await addUser(addForm.value);
      if (res.code === 200) {
        ElMessage.success('添加成功');
        addDialogVisible.value = false;
        resetAddForm();
        fetchUsers(); // 刷新用户列表
      } else {
        ElMessage.error(res.msg || '添加失败');
      }
    } catch (error) {
      ElMessage.error('添加失败');
      console.error(error);
    }
  });
};


// 处理编辑用户
const handleEditUser = (user) => {
  // 实际开发中这里会打开编辑用户的弹窗并回显数据
  ElMessage.info(`编辑用户: ${user.userName}`);
};

// 处理删除用户
const handleDeleteUser = async (userId) => {
  try {
    await ElMessageBox.confirm(
      '确定要删除该用户吗？',
      '删除确认',
      {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning'
      }
    );

    // 调用删除接口
    console.log("删除用户ID:", userId);
    const res = await deleteUser([userId]);
    if (res.code === 200) {
      ElMessage.success('删除成功');
      fetchUsers(); // 刷新用户列表
    } else {
      ElMessage.error(res.msg || '删除失败');
    }

  } catch (error) {
    if (error === 'cancel') return; // 用户取消操作
    ElMessage.error('删除失败');
    console.error('删除用户失败:', error);
  }
};

// 页面加载时获取用户列表
onMounted(() => {
  fetchUsers();
});
</script>

<style scoped>
.users-page {
  width: 100%;
  padding: 20px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-container {
  margin-bottom: 16px;
}

.table-container {
  overflow-x: auto;
  margin-bottom: 16px;
}

.pagination-container {
  text-align: right;
  margin-top: 16px;
}
</style>