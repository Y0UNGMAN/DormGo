import {createRouter , createWebHashHistory} from 'vue-router'

//定义路由
const routers = [
    {
        path:"/",
        name:"login",
        meta:{title:'登录'},
        component:()=>import('../views/Login.vue')
    },
    {
        path:"/home",
        name:"Home",
        meta:{title:'首页'},
        component:()=>import('../views/Home.vue')
    },
    {
        path:"/management",
        name:"Management",
        meta:{title:'用户管理'},
        component:()=>import('../views/Management.vue')
    },

]

//创建路由实例并传递router配置
const router = createRouter({
    history: createWebHashHistory(),
    routes:routers
})

export default router