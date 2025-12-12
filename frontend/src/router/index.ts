import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

// 懒加载组件
const AdminLogin = () => import('@/components/AdminLogin.vue')
const AdminLayout = () => import('@/components/AdminLayout.vue')
const AdminProfile = () => import('@/views/AdminProfile.vue')
const UserList = () => import('@/views/UserList.vue')
const ContentAudit = () => import('@/views/ContentAudit.vue')
const NoticeManage = () => import('@/views/NoticeManage.vue')
const ViolationHandle = () => import('@/views/ViolationHandle.vue')
const DormManage = () => import('@/views/DormManage.vue')
const ResetPwd = () => import('@/views/ResetPwd.vue')
const Statistics = () => import('@/views/Statistics.vue')
const SystemConfig = () => import('@/views/SystemConfig.vue')
const SensitiveWordManage = () => import('@/views/SensitiveWordManage.vue')
const PostTypeConfig = () => import('@/views/PostTypeConfig.vue')
const DormGoHome = () => import('@/views/DormGoHome.vue')
const PostDetail = () => import('@/views/PostDetail.vue')
const PostPublish = () => import('@/views/PostPublish.vue')
const LoginIn = () => import('@/views/Login.vue')
const PersonalHome = () => import('@/views/PersonalHome.vue')
const ChatPage = () => import('@/views/ChatPage.vue')
const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'AdminLogin',
    component: AdminLogin,
    meta: { title: '管理员登录' }
  },
  {
    path: '/admin',
    name: 'AdminLayout',
    component: AdminLayout,
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        redirect: '/admin/statistics'
      },
      {
        path: '/admin/profile',
        name: 'AdminProfile',
        component: AdminProfile,
        meta: { title: '个人中心' }
      },
      {
        path: '/admin/user-list',
        name: 'UserList',
        component: UserList,
        meta: { title: '用户管理' }
      },
      {
        path: '/admin/content-audit',
        name: 'ContentAudit',
        component: ContentAudit,
        meta: { title: '内容审核' }
      },
      {
        path: '/admin/notice',
        name: 'NoticeManage',
        component: NoticeManage,
        meta: { title: '通知管理' }
      },
      {
        path: '/admin/violation',
        name: 'ViolationHandle',
        component: ViolationHandle,
        meta: { title: '违规处理' }
      },
      {
        path: '/admin/dorm-manage',
        name: 'DormManage',
        component: DormManage,
        meta: { title: '楼栋管理' }
      },
      {
        path: '/admin/statistics',
        name: 'Statistics',
        component: Statistics,
        meta: { title: '数据统计' }
      },
      {
        path: '/admin/system-config',
        name: 'SystemConfig',
        component: SystemConfig,
        meta: { title: '系统配置' }
      },
      // 新增路由
      {
        path: '/admin/sensitive-word',
        name: 'SensitiveWordManage',
        component: SensitiveWordManage,
        meta: { title: '敏感词管理', roles: ['super_admin'] } // 权限示例
      },
      {
        path: '/admin/post-type-config',
        name: 'PostTypeConfig',
        component: PostTypeConfig,
        meta: { title: '互助类型配置', roles: ['super_admin'] }
      },
      {
        path: '/admin/reset-pwd',
        name: 'ResetPwd',
        component: ResetPwd,
        meta: { title: '重置密码' }
      },
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/login'
  },
  {
      path: '/dormgo',
      name: 'DormgoHome',
      meta: { title: 'Dormgo首页' },
      component: DormGoHome
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
    path: '/',
    name: 'LoginIn',
    meta: { title: '登录' },
    component: LoginIn
  },
  {
    path: '/personalhome',
    name: 'PersonalHome',
    meta: { title: '个人中心' },
    component: PersonalHome
  },
  {
    path: '/chat/:id', // :id 是对方的用户ID
    name: 'Chat',
    component: ChatPage,
    meta: { title : "私信页" }
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

router.beforeEach((to, _from, next) => {
  document.title = to.meta.title
    ? `${to.meta.title as string} - 寝友Go管理员后台`
    : '寝友Go管理员后台'

  const token = localStorage.getItem('adminToken')

  if (to.meta.requiresAuth && !token) {
    next({ path: '/login', query: { redirect: to.fullPath } })
  } else if (token && to.path === '/login') {
    next('/')
  } else {
    // 简单的权限校验示例
    if (to.meta.roles) {
      const info = JSON.parse(localStorage.getItem('adminInfo') || '{}')
      // 修复：添加类型断言 (as string[]) 解决 TS 报错
      const requiredRoles = to.meta.roles as string[]

      if (info.role && requiredRoles.includes(info.role)) {
        next()
      } else {
        // 权限不足，这里简化处理直接放行或可跳转403
        // 实际项目中通常会 next('/403')
        next()
      }
    } else {
      next()
    }
  }
})

export default router