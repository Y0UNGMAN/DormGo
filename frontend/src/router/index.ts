import { createRouter, createWebHistory } from 'vue-router'
// 组件引入
import AdminLogin from '@/components/AdminLogin.vue'
import AdminLayout from '@/components/AdminLayout.vue'
import AdminProfile from '@/views/AdminProfile.vue'
import AdminSideMenu from '@/views/AdminSideMenu.vue'
import DormgoHome from '@/views/DormgoHome.vue'
import PostDetail from '@/views/PostDetail.vue'
import PostPublish from '@/views/PostPublish.vue'

// 创建路由实例
const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'AdminLogin',
      component: AdminLogin//登录页面无需布局
    },
    {
      path: '/layout',
      name: 'AdminLayout',
      component: AdminLayout,
      //meta: { requiresAuth: true }, // 需登录访问
      children: [
        {
          path: 'admin/profile',
          name: 'AdminProfile',
          component: AdminProfile,
          meta: { title: '个人中心' }
        },
      ]
    },
    {
      path: '/dormgo',
      name: 'DormgoHome',
      component: DormgoHome,
    },
    {
      path: '/post/:id',
      name: 'PostDetail',
      component: PostDetail,
      props: true
    },
    {
      path: '/publish',
      name: 'PublishPost',
      component: PostPublish
    }
  ]
})

// 路由守卫
router.beforeEach((to, from, next) => {
  // 设置页面标题
  if (to.meta.title) {
    document.title = to.meta.title + ' - 寝友Go管理员后台'
  } else {
    document.title = '寝友Go管理员后台'
  }

  // 权限校验
  if (to.meta.requiresAuth) {
    const token = localStorage.getItem('adminToken')
    if (token) {
      next()
    } else {
      next({ path: '/login' })
    }
  } else {
    next()
  }
})

export default router