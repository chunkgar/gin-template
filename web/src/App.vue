<script setup lang="ts">
import { RouterView, useRouter, useRoute } from 'vue-router'
import { ref, watch, onMounted } from 'vue'
import { logout, getToken } from '@/utils/auth'
import { ElMessage, ElMessageBox } from 'element-plus'
import axios from 'axios'
import { useAdminStore } from '@/stores/admin'

const router = useRouter()
const route = useRoute()
const activeMenu = ref(route.path)
const adminStore = useAdminStore()
const baseURL = import.meta.env.VITE_BASE_URL || ''

watch(() => route.path, (newPath) => {
  activeMenu.value = newPath
})

const handleLogout = async () => {
  try {
    await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    adminStore.clearUser()
    logout()
    ElMessage.success('退出成功')
    router.push('/login')
  } catch (error: unknown) {
    // 用户取消操作
    if (error instanceof Error) {
      ElMessage.error(error.message || '退出失败')
    } else {
      ElMessage.error('退出失败')
    }
  }
}

const initAdminUser = async () => {
  try {
    const token = getToken()
    if (!token) return

    const response = await axios.get(`${baseURL}/api/admin/auth/profile`, {
      headers: { Authorization: `Bearer ${token}` }
    })

    const data = response.data
    if (data.code === 100001 && data.data) {
      adminStore.setUser(data.data)
    }
  } catch (error: unknown) {
    // ignore
    if (axios.isAxiosError(error)) {
      if (error.response) {
        ElMessage.error(error.response.data.error || '获取用户信息失败')
      } else {
        ElMessage.error('获取用户信息失败')
      }
    } else {
      ElMessage.error('获取用户信息失败')
    }
  }
}

onMounted(() => {
  initAdminUser()
})
</script>

<template>
  <el-container class="layout-container">
    <!-- 侧边栏 -->
    <el-aside width="200px" class="sidebar">
      <div class="logo">
        <h2>管理后台</h2>
      </div>
      <el-menu
        :default-active="activeMenu"
        router
        class="sidebar-menu"
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409eff"
      >
        <el-menu-item index="/dashboard">
          <el-icon><House /></el-icon>
          <span>仪表板</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <!-- 主内容区域 -->
    <el-container>
      <el-header class="header">
        <div class="header-content">
          <div class="breadcrumb">
            <el-breadcrumb separator="/">
              <el-breadcrumb-item :to="{ path: '/dashboard' }">首页</el-breadcrumb-item>
            </el-breadcrumb>
          </div>
          <div class="user-info" v-if="adminStore.user">
            <el-button type="danger" size="small" @click="handleLogout">退出</el-button>
          </div>
        </div>
      </el-header>
      <el-main class="main-content">
        <RouterView />
      </el-main>
    </el-container>
  </el-container>
</template>

<style scoped>
.layout-container {
  height: 100vh;
}

.sidebar {
  background-color: #304156;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  border-bottom: 1px solid #2b3848;
}

.logo h2 {
  margin: 0;
  font-size: 18px;
}

.sidebar-menu {
  border: none;
}

.header {
  background-color: #fff;
  border-bottom: 1px solid #e6e6e6;
  padding: 0 20px;
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 100%;
}

.breadcrumb {
  flex: 1;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.main-content {
  background-color: #f5f5f5;
  padding: 20px;
}
</style>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  background-color: #f5f5f5;
}

html, body, #app {
  height: 100%;
}
</style>
