<!-- eslint-disable vue/multi-word-component-names -->
<template>
  <div class="dashboard-container">
    <div class="welcome-card">
      <el-card>
        <template #header>
          <div class="card-header">
            <span>仪表板</span>
          </div>
        </template>
        <p>欢迎使用管理后台系统！</p>
        <p>当前时间：{{ currentTime }}</p>
        <p>用户角色：{{ userInfo?.role }}</p>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { logout, getToken } from '@/utils/auth'
import axios from 'axios'
import { useAdminStore } from '@/stores/admin'

const router = useRouter()
const currentTime = ref('')
const adminStore = useAdminStore()
const userInfo = computed(() => adminStore.user)
const baseURL = import.meta.env.VITE_BASE_URL || ''

const updateTime = () => {
  currentTime.value = new Date().toLocaleString('zh-CN')
}

const fetchUserInfo = async () => {
  try {
    const token = getToken()
    if (!token) {
      router.push('/login')
      return
    }

    const response = await axios.get(`${baseURL}/api/admin/auth/profile`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })

    const data = response.data
    if (data.code === 100001) {
      adminStore.setUser(data.data)
    } else {
      ElMessage.error('获取用户信息失败')
      logout()
    }
  } catch (error: unknown) {
    // 处理未知错误
    if (error instanceof Error) {
      ElMessage.error(error.message || '获取用户信息失败')
    } else {
      ElMessage.error('获取用户信息失败')
    }
    logout()
  }
}

onMounted(() => {
  updateTime()
  setInterval(updateTime, 1000)
  fetchUserInfo()
})
</script>

<style scoped>
.dashboard-container {
  padding: 20px;
}

.welcome-card {
  max-width: 600px;
  margin: 20px auto;
}

.card-header {
  font-weight: bold;
}
</style>
