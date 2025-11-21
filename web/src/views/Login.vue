<template>
  <div class="login-container">
    <div class="login-form">
      <h2>管理员登录</h2>
      <el-form :model="loginForm" :rules="loginRules" ref="loginFormRef" label-width="80px" class="login-form-content">
        <el-form-item label="账号" prop="username">
          <el-input v-model="loginForm.username" placeholder="请输入账号" prefix-icon="User" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="loginForm.password" type="password" placeholder="请输入密码" prefix-icon="Lock" show-password />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleLogin" class="login-btn">
            登录
          </el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { setToken, isLoggedIn } from '@/utils/auth'
import axios from 'axios'
import { defineOptions } from 'vue'

defineOptions({
  name: 'AdminLogin'
})

const router = useRouter()
const loginFormRef = ref()
const loading = ref(false)

const loginForm = reactive({
  username: '',
  password: ''
})

const baseURL = import.meta.env.VITE_BASE_URL || ''

const loginRules = {
  username: [
    { required: true, message: '请输入账号', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少6位', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  if (!loginFormRef.value) return

  const valid = await loginFormRef.value.validate()
  if (!valid) return

  loading.value = true

  try {
    const response = await axios.post(`${baseURL}/api/admin/login`, loginForm)
    const data = response.data
    if (data.code === 100001) {
      console.log("登录成功", data)
      setToken(data.data.token)
      isLoggedIn.value = true
      ElMessage.success('登录成功')
      const redirect = router.currentRoute.value.query.redirect as string
      router.push(redirect || '/dashboard')
    } else {
      ElMessage.error(data.error || '登录失败')
    }
  } catch (error: unknown) {
    if (axios.isAxiosError(error)) {
      if (error.response) {
        ElMessage.error(error.response.data.error || '登录失败')
      } else {
        ElMessage.error('网络错误，请稍后重试')
      }
    }
    console.error('登录错误:', error)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #f5f5f5;
}

.login-form {
  background: white;
  padding: 40px;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  width: 400px;
}

.login-form h2 {
  text-align: center;
  margin-bottom: 30px;
  color: #333;
}

.login-form-content {
  width: 100%;
}

.login-btn {
  width: 100%;
}
</style>