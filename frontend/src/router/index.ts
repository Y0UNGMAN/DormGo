import { createRouter, createWebHistory } from 'vue-router'

// 组件引入
// Vue 项目里常用 @ 作为 src 根目录的别名
import AdminLogin from '@/components/AdminLogin.vue'
import AdminLayout from '@/components/AdminLayout.vue'
import AdminProfile from '@/views/AdminProfile.vue'
import DormgoHome from '@/views/DormgoHome.vue'
import PostDetail from '@/views/PostDetail.vue'
import PostPublish from '@/views/PostPublish.vue'
import Login from '@/views/Login.vue'
import PersonalHome from '@/views/PersonalHome.vue'


// 创建路由实例
const router = createRouter({
  history: createWebHistory(),  // Histoty模式的URL更干净
  routes: [
    {
      path: '/',  // URL地址
      name: 'AdminLogin', // 给该路由起的名字
      component: AdminLogin
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
      meta: { title: 'Dormgo首页' },
      component: DormgoHome
    },
    {
      // 动态路由
      path: '/post/:id',  // :id → 路径参数，可以动态传入不同帖子 ID
      name: 'PostDetail',
      meta: { title: '帖子详情' },
      component: PostDetail,
      props: true // 自动把路径参数传给组件 props
    },
    {
      path: '/publish',
      name: 'PublishPost',
      meta: { title: '发布帖子' },
      component: PostPublish
    },
    {
      path: '/login',
      name: 'Login',
      meta: { title: '登录' },
      component: Login,
    },
    {
      path: '/home',
      name: 'Home',
      meta: { title: '个人中心' },
      component: PersonalHome,
    }
  ]
})

// 路由守卫
router.beforeEach((to, _from, next) => {
  // 设置页面标题
  if (to.meta.title) {
    document.title = to.meta.title as string
  } else {
    document.title = ''
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