<template>
  <div class="dormgo-home">
    <!-- 顶部栏：搜索框和用户头像 -->
    <div class="top-bar">
      <!-- 搜索框 - 居中 -->
      <div class="search-section">
        <div class="search-container">
          <div class="search-box">
            <span class="search-icon">🔍</span>
            <input
              v-model="searchKeyword"
              type="text"
              placeholder="搜索帖子标题/内容"
              class="search-input"
              @input="handleSearch"
            />
            <button v-if="searchKeyword" class="clear-btn" @click="clearSearch">×</button>
          </div>
          <button class="search-btn" @click="handleSearch">搜索</button>
        </div>
      </div>

      <!-- 用户头像 -->
      <div class="user-avatar-section">
        <img 
          :src="currentUser.avatar || defaultAvatar" 
          alt="用户头像" 
          class="user-avatar"
          @click="goToProfile"
        />
      </div>
    </div>

    <!-- 上边栏 - 分类筛选（居中） -->
    <div class="top-nav">
      <div class="nav-content">
        <div class="category-filters">
          <span class="filter-label">分类：</span>
          <button
            class="filter-btn"
            :class="{ active: selectedCategory === 'all' }"
            @click="selectCategory('all')"
          >
            全部
             <span class="category-count">{{ posts.length }}</span>
          </button>
          <button
            v-for="posttype in postTypes"
            :key="posttype.typeid"
            class="filter-btn"
            :class="{ active: selectedCategory === posttype.typeid }"
            @click="selectCategory(posttype.typeid)"
          >
            {{ posttype.typename }}
            <span class="category-count">{{ getCategoryCount(posttype.typeid) }}</span>
          </button>
        </div>
      </div>
    </div>
    
    <div class="main-content">
      <!-- 左边栏 - 宿舍楼切换 -->
      <div class="sidebar">
        <h3 class="sidebar-title">宿舍楼号</h3>
        <div class="dorm-list">
          <div
            class="dorm-item"
            :class="{ active: selectedDorm === 'all' }"
            @click="selectDorm('all')"
          >
             <span class="dorm-name">全部宿舍</span>
          </div>
          <div
            v-for="dorm in dormList"
            :key="dorm.dormid"
            class="dorm-item"
            :class="{ active: selectedDorm === dorm.dormid }"
            @click="selectDorm(dorm.dormid)"
          >
            <span class="dorm-name">{{ dorm.dormname }}</span>
            <span class="post-count">{{ getDormCount(dorm.dormname) }}</span>
          </div>
        </div>
      </div>
      
      <!-- 帖子列表区域 -->
      <div class="posts-area">
        <!-- 筛选状态提示 -->
        <div class="filter-status" v-if="hasActiveFilters">
          <span class="status-text">
            <span v-if="searchKeyword">搜索: "{{ searchKeyword }}" </span>
            <span v-if="selectedCategory !== 'all'">分类: {{ getCategoryName(selectedCategory) }} </span>
            <span v-if="selectedDorm !== 'all'">宿舍楼: {{ getDormName(selectedDorm) }}</span>
            <span class="post-count-text">({{ filteredPosts.length }}个结果)</span>
          </span>
          <button class="clear-filters" @click="clearAllFilters">清除筛选</button>
        </div>

        <!-- 帖子列表 -->
        <div class="posts-container">
          <Card 
            v-for="post in filteredPosts"
            :key="post.id"
            :post="post"
            class="post-card"
          />
        </div>
        
        <!-- 空状态 -->
        <div class="empty-state" v-if="filteredPosts.length === 0">
          <div class="empty-icon">🍃</div>
          <p class="empty-text">这里什么都没有，快去发布第一条内容吧</p>
          <button class="empty-btn" @click="handlePublish">发布帖子</button>
        </div>
      </div>
    </div>

    <!-- 发布帖子按钮 - 固定在右下角 -->
    <button class="publish-fab" @click="handlePublish">
      <span class="fab-icon">📝</span>
      发布帖子
    </button>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Card from '@/components/Card.vue'
import request from '@/utils/request' // 【修改点】引入封装好的 request，而不是 axios

const router = useRouter()
const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'

// 响应式数据
const selectedDorm = ref('all')
const selectedCategory = ref('all')
const searchKeyword = ref('')
const currentUser = ref({})

// 数据源
const dormList = ref([])
const postTypes = ref([])
const posts = ref([])

onMounted(() => {
  // 从本地存储恢复用户信息
  const userStr = localStorage.getItem('userInfo')
  if (userStr) {
     currentUser.value = JSON.parse(userStr)
  }
  
  fetchPostTypes()
  fetchDormList()
  fetchPostList()
})

// 数据获取 (修改为使用 request)
const fetchDormList = async() => {    
  try {
    // 【修改点】直接写路径，request 会自动拼接到 http://127.0.0.1:8080
    const response = await request.get('/api/v1/post/dorms')
    // 后端返回结构为 { code:200, msg:"success", data: [...] }
    dormList.value = response.data || []
  } catch (error) { console.error('获取宿舍失败', error) }
}

const fetchPostTypes = async () => {     
  try {
    const response = await request.get('/api/v1/post/post_type')
    postTypes.value = response.data || []
  } catch (error) { console.error('获取分类失败', error) }
}

const fetchPostList = async () => {
  try {
    const response = await request.get('/api/v1/post/posts')
    posts.value = response.data || []
  } catch (error) { console.error('获取帖子失败', error) }
}

// 筛选逻辑
const filteredPosts = computed(() => {
  let filtered = posts.value
  
  if (searchKeyword.value.trim()) {
    const keyword = searchKeyword.value.trim().toLowerCase()
    filtered = filtered.filter(post => 
      post.title?.toLowerCase().includes(keyword) ||
      post.content?.toLowerCase().includes(keyword)
    )
  }
  
  if (selectedCategory.value !== 'all') {
    filtered = filtered.filter(post => post.typeid === selectedCategory.value)
  }
  
  if (selectedDorm.value !== 'all') {
    const dormName = getDormName(selectedDorm.value)
    filtered = filtered.filter(post => post.Dorm?.dormname === dormName)
  }
  
  return filtered
})

const hasActiveFilters = computed(() => {
  return searchKeyword.value.trim() !== '' || selectedCategory.value !== 'all' || selectedDorm.value !== 'all'
})

// 辅助函数
const getCategoryName = (id) => {
  const category = postTypes.value.find(cat => cat.typeid === id)
  return category ? category.typename : '全部'
}

const getDormName = (dormId) => {
  const dorm = dormList.value.find(d => d.dormid === dormId)
  return dorm ? dorm.dormname : '全部宿舍'
}

const getCategoryCount = (categoryId) => {
  return posts.value.filter(post => post.typeid === categoryId).length
}

const getDormCount = (dormName) => {
  return posts.value.filter(post => post.Dorm?.dormname === dormName).length
}

// 事件处理
const selectDorm = (dormId) => selectedDorm.value = dormId
const selectCategory = (typeid) => selectedCategory.value = typeid
const clearSearch = () => searchKeyword.value = ''
const clearAllFilters = () => {
  searchKeyword.value = ''
  selectedCategory.value = 'all'
  selectedDorm.value = 'all'
}

// 路由跳转
const handlePublish = () => {
  router.push('/post/publish')
}

const goToProfile = () => {
  router.push('/user/profile')
}
const handleSearch = () => {
  // 搜索逻辑主要依赖 computed
}
</script>

<style scoped>
.dormgo-home {
  min-height: 100vh;
  background: #f5f5f5;
  padding: 20px;
  position: relative;
}

/* 顶部栏布局 */
.top-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  margin-bottom: 20px;
}

/* 搜索框样式 - 居中 */
.search-section {
  flex: 1;
  display: flex;
  justify-content: center;
}

.search-container {
  display: flex;
  gap: 12px;
  max-width: 600px;
  width: 100%;
}

.search-box {
  flex: 1;
  position: relative;
  display: flex;
  align-items: center;
}

.search-icon {
  position: absolute;
  left: 12px;
  font-size: 16px;
  color: #999;
}

.search-input {
  width: 100%;
  padding: 12px 40px 12px 40px;
  border: 1px solid #e8e8e8;
  border-radius: 8px;
  font-size: 14px;
  transition: all 0.3s;
}

.search-input:focus {
  outline: none;
  border-color: #1890ff;
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.2);
}

.clear-btn {
  position: absolute;
  right: 12px;
  background: none;
  border: none;
  font-size: 18px;
  color: #999;
  cursor: pointer;
}

.search-btn {
  padding: 12px 24px;
  background: #1890ff;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.3s;
}

.search-btn:hover {
  background: #40a9ff;
}

/* 用户头像样式 */
.user-avatar-section { flex-shrink: 0; }
.user-avatar {
  width: 44px; height: 44px; border-radius: 50%; object-fit: cover;
  cursor: pointer; border: 2px solid #e8e8e8; transition: all 0.3s ease;
}
.user-avatar:hover { border-color: #1890ff; transform: scale(1.05); }

/* 上边栏样式 */
.top-nav {
  background: white; border-radius: 12px; padding: 16px 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1); margin-bottom: 20px;
}
.nav-content { display: flex; justify-content: center; align-items: center; }
.category-filters { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; justify-content: center; }
.filter-label { font-size: 14px; color: #666; font-weight: 500; }
.filter-btn {
  display: flex; align-items: center; gap: 6px; padding: 8px 12px;
  background: #f5f5f5; border: 1px solid #e8e8e8; border-radius: 20px;
  font-size: 13px; color: #666; cursor: pointer; transition: all 0.3s ease;
}
.filter-btn:hover { background: #e6f7ff; border-color: #1890ff; color: #1890ff; }
.filter-btn.active { background: #1890ff; border-color: #1890ff; color: white; }
.category-count { background: rgba(255, 255, 255, 0.3); padding: 1px 6px; border-radius: 10px; font-size: 11px; }
.filter-btn:not(.active) .category-count { background: #e8e8e8; color: #666; }

/* 发布按钮 */
.publish-fab {
  position: fixed; bottom: 30px; right: 30px; display: flex; align-items: center; gap: 8px;
  padding: 16px 24px; background: #1890ff; color: white; border: none; border-radius: 50px;
  font-size: 16px; font-weight: 600; cursor: pointer;
  box-shadow: 0 4px 16px rgba(24, 144, 255, 0.3); transition: all 0.3s ease; z-index: 1000;
}
.publish-fab:hover { background: #40a9ff; transform: translateY(-2px); }

/* 主内容区 */
.main-content { display: flex; gap: 20px; max-width: 1200px; margin: 0 auto; }
.sidebar {
  width: 200px; background: white; border-radius: 12px; padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1); height: fit-content; position: sticky; top: 20px;
}
.sidebar-title { font-size: 16px; font-weight: 600; color: #333; margin: 0 0 16px 0; border-bottom: 1px solid #f0f0f0; padding-bottom: 12px; }
.dorm-list { display: flex; flex-direction: column; gap: 8px; }
.dorm-item {
  display: flex; align-items: center; justify-content: space-between; padding: 10px 12px;
  border-radius: 8px; cursor: pointer; transition: all 0.3s ease;
}
.dorm-item:hover { background: #f5f5f5; }
.dorm-item.active { background: #e6f7ff; color: #1890ff; }
.post-count { background: #f0f0f0; color: #666; padding: 2px 6px; border-radius: 10px; font-size: 11px; }
.dorm-item.active .post-count { background: #1890ff; color: white; }

/* 帖子区域 */
.posts-area { flex: 1; min-width: 0; }
.filter-status {
  display: flex; justify-content: space-between; align-items: center;
  background: #e6f7ff; border: 1px solid #91d5ff; border-radius: 8px;
  padding: 12px 16px; margin-bottom: 16px;
}
.status-text { color: #1890ff; font-size: 14px; }
.clear-filters { background: none; border: 1px solid #1890ff; color: #1890ff; padding: 4px 8px; border-radius: 4px; font-size: 12px; cursor: pointer; }
.posts-container { display: flex; flex-direction: column; gap: 16px; }

/* 空状态 */
.empty-state { text-align: center; padding: 60px 20px; background: white; border-radius: 12px; }
.empty-icon { font-size: 48px; margin-bottom: 16px; }
.empty-text { color: #666; margin-bottom: 20px; }
.empty-btn { padding: 10px 20px; background: #1890ff; color: white; border: none; border-radius: 8px; cursor: pointer; }
</style>