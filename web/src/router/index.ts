import { createRouter, createWebHistory } from 'vue-router'
import Login from '@/views/Login.vue'
import Dashboard from '@/views/Dashboard.vue'
import { checkLoginStatus } from '@/utils/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: Login,
      meta: { requiresAuth: false },
    },
    {
      path: '/dashboard',
      name: 'Dashboard',
      component: Dashboard,
      meta: { requiresAuth: true },
    },
    {
      path: '/',
      redirect: '/dashboard',
    },
  ],
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const isLoggedIn = checkLoginStatus()

  if (to.meta.requiresAuth && !isLoggedIn) {
    // 重定向到登录页面，并保存当前路径
    next({
      path: '/login',
      query: { redirect: to.fullPath },
    })
  } else if (to.path === '/login' && isLoggedIn) {
    // 如果已登录且访问登录页面，重定向到dashboard
    next('/dashboard')
  } else {
    next()
  }
})

export default router
