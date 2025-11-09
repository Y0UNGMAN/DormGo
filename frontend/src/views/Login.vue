<template>
  <div>
    <h1>寝友Go</h1>
    <h1>登录界面</h1>
  </div>
  <div class="btn-zone">
    <router-link to="/Home">首页</router-link>
    <router-link to="/Management">用户管理</router-link>
    <div>后端返回：{{ msg }}</div>
  </div>
</template>

<script setup lang="ts">
import {onMounted , ref} from 'vue'
import api from '../api/index'

const msg = ref('')

onMounted(async()=>{
    try{
        const res = await api.get('/ping')
        msg.value = res.data.message
    }catch{
        msg.value = '请求失败'
    }
})
</script>

<style scoped>
.btn-zone{
    display:flex;
    flex-direction: column;
    gap: 10px; 
}
</style>
