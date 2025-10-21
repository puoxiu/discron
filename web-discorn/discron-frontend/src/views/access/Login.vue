<template>
    <div class="login-container">
        <el-card header="用户登录" class="login-card">
            <el-form ref="loginFormRef" :model="loginForm" :rules="loginRules" label-width="80px">
                <el-form-item label="用户名" prop="userName">
                    <el-input v-model="loginForm.userName" placeholder="请输入用户名"></el-input>
                </el-form-item>
                <el-form-item label="密码" prop="password">
                    <el-input v-model="loginForm.password" type="password" placeholder="请输入密码"></el-input>
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" @click="handleLogin">登录</el-button>
                    <el-button type="default" @click="$router.push('/register')">注册</el-button>
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
const loginFormRef = ref(null);

// 表单数据
const loginForm = reactive({
    userName: '',
    password: ''
});

// 表单验证规则
const loginRules = {
    userName: [
        { required: true, message: '请输入用户名', trigger: 'blur' }
    ],
    password: [
        { required: true, message: '请输入密码', trigger: 'blur' }
    ]
};

// 登录处理
const handleLogin = async () => {
  // 1. 表单验证：用 try/catch 捕获验证失败的错误
  try {
    await loginFormRef.value.validate(); 
  } catch (error) {
    console.error('表单验证失败:', error);
    return; // 验证失败，直接退出
  }

  // 2. 验证通过后，执行登录逻辑
  try {
    await authStore.login(loginForm);
    ElMessage.success('登录成功');
    console.log('登录成功2:', authStore.userInfo); // 此时会正常打印
    router.push('/');
  } catch (error) {
    console.error('登录失败:', error);
    ElMessage.error('登录失败，请检查账号密码'); // 建议增加失败提示
  }
};
</script>

<style scoped>
.login-container {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
    background-color: #f5f7fa;
}

.login-card {
    width: 400px;
    padding: 20px;
}
</style>