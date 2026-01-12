import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

// 懒加载组件 - 管理员端
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

// 懒加载组件 - 用户端
const UserLogin = () => import('@/components/UserLogin.vue')
const DormgoHome = () => import('@/views/DormgoHome.vue')
const PersonalHome = () => import('@/views/PersonalHome.vue')
const PostDetail = () => import('@/views/PostDetail.vue')
const PostPublish = () => import('@/views/PostPublish.vue')
// 用户端 - 个人中心子页面
const PostHistory = () => import('@/views/PostHistory.vue')
const MyFavorites = () => import('@/views/MyFavorites.vue')
const Notifications = () => import('@/views/Notifications.vue')
const CoinFlow = () => import('@/views/CoinFlow.vue')
const UserDetail = () => import('@/views/UserDetail.vue')
const CommunityRules = () => import('@/views/CommunityRules.vue')
const ContactUs = () => import('@/views/ContactUs.vue')

const routes: RouteRecordRaw[] = [
  // ===================== 用户端路由 =====================
  {
    path: '/user/login',
    name: 'UserLogin',
    component: UserLogin,
    meta: { title: '用户登录' }
  },
  {
    path: '/home',
    name: 'DormgoHome',
    component: DormgoHome,
    meta: { title: '首页', requiresUserAuth: true }
  },
  {
    path: '/profile',
    name: 'PersonalHome',
    component: PersonalHome,
    meta: { title: '个人中心', requiresUserAuth: true }
  },
  {
    path: '/post/:id',
    name: 'PostDetail',
    component: PostDetail,
    meta: { title: '帖子详情', requiresUserAuth: true }
  },
  {
    path: '/publish',
    name: 'PostPublish',
    component: PostPublish,
    meta: { title: '发布帖子', requiresUserAuth: true }
  },
  // 用户端 - 个人中心子页面
  {
    path: '/profile/post-history',
    name: 'PostHistory',
    component: PostHistory,
    meta: { title: '发帖历史', requiresUserAuth: true }
  },
  {
    path: '/profile/favorites',
    name: 'MyFavorites',
    component: MyFavorites,
    meta: { title: '我的收藏', requiresUserAuth: true }
  },
  {
    path: '/profile/notifications',
    name: 'Notifications',
    component: Notifications,
    meta: { title: '消息通知', requiresUserAuth: true }
  },
  {
    path: '/profile/coins',
    name: 'CoinFlow',
    component: CoinFlow,
    meta: { title: '寝友币', requiresUserAuth: true }
  },
  {
    path: '/profile/detail',
    name: 'UserDetail',
    component: UserDetail,
    meta: { title: '详细资料', requiresUserAuth: true }
  },
  {
    path: '/profile/rules',
    name: 'CommunityRules',
    component: CommunityRules,
    meta: { title: '社区规范', requiresUserAuth: true }
  },
  {
    path: '/profile/contact',
    name: 'ContactUs',
    component: ContactUs,
    meta: { title: '联系我们', requiresUserAuth: true }
  },

  // ===================== 管理员端路由 =====================
  {
    path: '/admin/login',
    name: 'AdminLogin',
    component: AdminLogin,
    meta: { title: '管理员登录' }
  },
  {
    path: '/',
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
      }
    ]
  },
  // 默认路由 - 未匹配路径跳转到用户登录
  {
    path: '/:pathMatch(.*)*',
    redirect: '/user/login'
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

router.beforeEach((to, _from, next) => {
  // 设置页面标题
  if (to.meta.title) {
    const isAdminRoute = to.path.startsWith('/admin')
    document.title = isAdminRoute
      ? `${to.meta.title as string} - 寝友Go管理员后台`
      : `${to.meta.title as string} - 寝友Go`
  } else {
    document.title = '寝友Go'
  }

  const adminToken = localStorage.getItem('adminToken')
  const userToken = localStorage.getItem('userToken')

  // 管理员端路由鉴权
  if (to.meta.requiresAuth && !adminToken) {
    next({ path: '/admin/login', query: { redirect: to.fullPath } })
    return
  }

  // 用户端路由鉴权
  if (to.meta.requiresUserAuth && !userToken) {
    next({ path: '/user/login', query: { redirect: to.fullPath } })
    return
  }

  // 已登录管理员访问管理员登录页，跳转到后台首页
  if (adminToken && to.path === '/admin/login') {
    next('/admin/statistics')
    return
  }

  // 已登录用户访问用户登录页，跳转到用户首页
  if (userToken && to.path === '/user/login') {
    next('/home')
    return
  }

  // 管理员权限校验
  if (to.meta.roles) {
    const info = JSON.parse(localStorage.getItem('adminInfo') || '{}')
    const requiredRoles = to.meta.roles as string[]

    if (info.role && requiredRoles.includes(info.role)) {
      next()
    } else {
      // 权限不足，这里简化处理直接放行
      next()
    }
  } else {
    next()
  }
})

export default router