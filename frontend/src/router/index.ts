import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
// === 布局组件 ===
const AdminLayout = () => import('@/components/AdminLayout.vue')
// === 公共视图 (改为引用 Components 中的 Login) ===
const Login = () => import('@/components/Login.vue')
// === 用户侧视图 ===
const PersonalHome = () => import('@/views/PersonalHome.vue')
// === 帖子相关视图 ===
const PostPublish = () => import('@/views/PostPublish.vue')
const PostDetail = () => import('@/views/PostDetail.vue')
// === 管理员侧视图 ===
const Statistics = () => import('@/views/Statistics.vue')
const AdminProfile = () => import('@/views/AdminProfile.vue')
const UserList = () => import('@/views/UserList.vue')
const ContentAudit = () => import('@/views/ContentAudit.vue')
const NoticeManage = () => import('@/views/NoticeManage.vue')
const ViolationHandle = () => import('@/views/ViolationHandle.vue')
const DormManage = () => import('@/views/DormManage.vue')
const SystemConfig = () => import('@/views/SystemConfig.vue')
const SensitiveWordManage = () => import('@/views/SensitiveWordManage.vue')
const PostTypeConfig = () => import('@/views/PostTypeConfig.vue')
const ResetPwd = () => import('@/views/ResetPwd.vue')
const DormGoHome = () => import('@/views/DormGoHome.vue')
const LoginIn = () => import('@/views/LoginIn.vue')
const ChatPage = () => import('@/views/ChatPage.vue')
const Messages = () => import('@/views/MessageList.vue')
const OtherUserProfile = () => import('@/views/OtherUserProfile.vue')

const routes: RouteRecordRaw[] = [
  {
    path: '/loginin',
    name: 'LoginIn',
    meta: { title: '登录' },
    component: LoginIn
  },


  // 1. 认证路由 (Login Route)
  {
    path: '/',
    name: 'Login',
    component: Login,
    meta: { title: '欢迎登录 - 寝友Go' }
  },

  // 2. 根路径重定向
  // {
  //   path: '/',
  //   // 修复：将未使用的参数 to 改为 _to 以避免 TS 6133 警告
  //   redirect: (_to) => {
  //     const adminToken = localStorage.getItem('adminToken')
  //     const userToken = localStorage.getItem('userToken')
  //     if (adminToken) return '/admin/statistics'
  //     if (userToken) return '/dormgo'
  //     return '/LoginIn'
  //   }
  // },

  // 5. 管理员侧路由 (Admin Routes)
  {
    path: '/admin',
    component: AdminLayout,
    meta: { requiresAuth: true, role: 'admin' },
    children: [
      {
        path: '',
        redirect: '/admin/statistics'
      },
      {
        path: 'statistics',
        name: 'AdminStatistics',
        component: Statistics,
        meta: { title: '数据统计' }
      },
      {
        path: 'profile',
        name: 'AdminProfile',
        component: AdminProfile,
        meta: { title: '个人中心' }
      },
      {
        path: 'users',
        name: 'UserManage',
        component: UserList,
        meta: { title: '用户管理' }
      },
      {
        path: 'dorms',
        name: 'DormManage',
        component: DormManage,
        meta: { title: '楼栋管理' }
      },
      {
        path: 'audit',
        name: 'ContentAudit',
        component: ContentAudit,
        meta: { title: '内容审核' }
      },
      {
        path: 'notices',
        name: 'NoticeManage',
        component: NoticeManage,
        meta: { title: '通知管理' }
      },
      {
        path: 'violations',
        name: 'ViolationHandle',
        component: ViolationHandle,
        meta: { title: '违规处理' }
      },
      {
        path: 'config',
        name: 'SystemConfig',
        component: SystemConfig,
        meta: { title: '系统配置' }
      },
      {
        path: 'sensitive-words',
        name: 'SensitiveWordManage',
        component: SensitiveWordManage,
        meta: { title: '敏感词管理' }
      },
      {
        path: 'post-types',
        name: 'PostTypeConfig',
        component: PostTypeConfig,
        meta: { title: '互助类型配置' }
      },
      {
        path: 'reset-pwd',
        name: 'AdminResetPwd',
        component: ResetPwd,
        meta: { title: '重置密码' }
      }
    ]
  },

  // 6. 404路由
  {
    path: '/:pathMatch(.*)*',
    redirect: '/login'
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
    path: '/chat/:id', 
    name: 'Chat',
    component: ChatPage,
    meta: { title : "私信页" }
  },
  {
    path: '/messages', 
    name: 'Messages',
    component: Messages,
    meta: { title : "消息列表" }
  },
  {
    path: '/user/:id', 
    name: 'OtherUserProfile',
    component: OtherUserProfile,
    meta: { title : "个人主页" }
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

// === 全局路由守卫 ===
// 修复：将未使用的 from 参数改为 _from 以避免 TS 6133 警告
router.beforeEach((to, _from, next) => {
  // 设置标题
  document.title = to.meta.title ? `${to.meta.title}` : '寝友Go'

  const userToken = localStorage.getItem('userToken')
  const adminToken = localStorage.getItem('adminToken')
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)

  if (requiresAuth) {
    // 检查角色权限
    const role = to.meta.role

    if (role === 'admin') {
      if (!adminToken) {
        // 无管理员权限，跳转登录
        next({ path: '/login', query: { redirect: to.fullPath, role: 'admin' } })
      } else {
        next()
      }
    } else if (role === 'user') {
      if (!userToken) {
        // 无用户权限，跳转登录
        next({ path: '/login', query: { redirect: to.fullPath } })
      } else {
        next()
      }
    } else {
      next() // 不需要特定角色的认证路由
    }
  } else {
    // 已登录用户访问登录页自动跳转
    if (to.path === '/login') {
      if (adminToken) next('/admin/statistics')
      else if (userToken) next('/dormgo')
      else next()
    } else {
      next()
    }
  }
})

export default router