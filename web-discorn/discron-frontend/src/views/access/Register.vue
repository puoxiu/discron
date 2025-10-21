<template>
  <div class="register-container">
    <el-card header="用户注册" class="register-card">
      <el-form 
        ref="registerFormRef" 
        :model="registerForm" 
        :rules="registerRules" 
        label-width="80px"
      >
        <el-form-item label="用户名" prop="userName">
          <el-input v-model="registerForm.userName" placeholder="请输入用户名"></el-input>
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input 
            v-model="registerForm.password" 
            type="password" 
            placeholder="请输入密码"
          ></el-input>
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="registerForm.email" placeholder="请输入邮箱（可选）"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleRegister">注册</el-button>
          <el-button type="default" @click="$router.push('/login')">返回登录</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../../stores/auth';
import { ElMessage } from 'element-plus';

const router = useRouter();
const authStore = useAuthStore();
const registerFormRef = ref(null);

// 表单数据（对应后端ReqUserRegister结构体）
const registerForm = reactive({
  userName: '', // 用户名
  password: '', // 密码（后端会MD5加密）
  email: ''     // 邮箱（可选）
});

// 表单验证规则
const registerRules = {
  userName: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度3-20字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码至少6位', trigger: 'blur' }
  ],
  email: [
    { type: 'email', message: '请输入正确邮箱格式', trigger: 'blur', required: false }
  ]
};

// 注册处理
const handleRegister = async () => {
  const valid = await registerFormRef.value.validate();
  if (!valid) return;

  try {
    await authStore.register(registerForm);
    ElMessage.success('注册成功，请登录');
    router.push('/login'); // 注册成功跳转到登录页
  } catch (error) {
    console.error('注册失败:', error);
  }
};
</script>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f5f7fa;
}
.register-card {
  width: 400px;
  padding: 20px;
}
</style>