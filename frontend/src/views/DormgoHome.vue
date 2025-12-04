<!-- src/views/DormgoHome.vue -->
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
          :src="currentUser.avatar" 
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

        <!-- 帖子列表 - 直接显示所有帖子 -->
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
          <div class="empty-icon">📝</div>
          <p class="empty-text" v-if="searchKeyword">没有找到包含 "{{ searchKeyword }}" 的帖子</p>
          <p class="empty-text" v-else-if="selectedCategory !== 'all'">在{{ getDormName(selectedDorm) }}没有找到{{ getCategoryName(selectedCategory) }}相关的帖子</p>
          <p class="empty-text" v-else>暂无帖子，快去发布第一个吧！</p>
          <button class="empty-btn" @click="handlePublish">发布帖子</button>
          <button v-if="hasActiveFilters" class="empty-btn secondary" @click="clearAllFilters">查看全部帖子</button>
        </div>

        <!-- 已显示所有帖子的提示 -->
        <div class="all-posts-loaded" v-if="filteredPosts.length > 0">
          <div class="loaded-text">已显示所有帖子</div>
          <div class="loaded-count">共 {{ filteredPosts.length }} 个帖子</div>
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
import axios from 'axios'
const router = useRouter()

// 响应式数据
const selectedDorm = ref('all') // 默认显示全部宿舍楼
const selectedCategory = ref('all')
const searchKeyword = ref('')

// 当前用户信息
const currentUser = ref({
  id: '1',
  name: '当前用户',
  avatar: 'https://dorm-go.oss-cn-guangzhou.aliyuncs.com/avator/midnight.jpg'
})

const dormList = ref([]);   // 宿舍列表
const fetchDormList = async() => {    
  try {
    const response = await axios.get('http://127.0.0.1:8080/api/v1/post/dorms');
    if (response.data && response.data.data) {
      // 成功获取数据，并赋值给响应式变量 dormList
      dormList.value = response.data.data;
      console.log('宿舍楼列表:', dormList.value);
    } else {
      // 如果数据结构不符合预期
      throw new Error('接口返回数据结构异常');
    }
  } catch (error) {
    console.error('获取宿舍楼列表失败:', error);
  }
};

const postTypes = ref([]);    // 宿舍分类
const fetchPostTypes = async () => {     
  try {
    const response = await axios.get('http://127.0.0.1:8080/api/v1/post/post_type');
    if (response.data && response.data.data) {
      // 成功获取数据，并赋值给响应式变量 postTypes
      postTypes.value = response.data.data;
      console.log('帖子分类列表:', postTypes.value);
    } else {
      // 如果数据结构不符合预期
      throw new Error('接口返回数据结构异常');
    }
  } catch (error) {
    console.error('获取帖子分类失败:', error);
  }
};

// 计算属性：筛选帖子（包含搜索、分类、宿舍楼筛选）

const posts = ref([]);    // 帖子
const fetchPostList = async () => {
  try {
    const response = await axios.get('http://127.0.0.1:8080/api/v1/post/posts');
    if (response.data && response.data.data) {
      // 成功获取数据，并赋值给响应式变量 posts
      posts.value = response.data.data;
      console.log('帖子列表:', posts.value);
    } else {
      // 如果数据结构不符合预期
      throw new Error('接口返回数据结构异常');
    }
  } catch (error) {
    console.error('获取帖子列表失败:', error);
  }
};

const filteredPosts = computed(() => {  // 筛选帖子
  let filtered = posts.value
  
  // 搜索筛选：在标题和内容中搜索
  if (searchKeyword.value.trim()) {
    const keyword = searchKeyword.value.trim().toLowerCase()
    filtered = filtered.filter(post => 
      post.title.toLowerCase().includes(keyword) ||
      post.content.toLowerCase().includes(keyword)
    )
  }
  
  // 按分类筛选
  if (selectedCategory.value !== 'all') {
    filtered = filtered.filter(post => post.typeid === selectedCategory.value)
  }
  
  // 按宿舍楼筛选
  if (selectedDorm.value !== 'all') {
    const dormName = getDormName(selectedDorm.value)
    filtered = filtered.filter(post => post.Dorm.dormname === dormName)
  }
  
  return filtered
})

// 计算属性：是否有活跃的筛选条件
const hasActiveFilters = computed(() => {
  return searchKeyword.value.trim() !== '' || selectedCategory.value !== 'all' || selectedDorm.value !== 'all'
})

// 获取分类名称
const getCategoryName = (id) => {
  const category = postTypes.value.find(cat => cat.typeid === id)
  return category ? category.typename : '全部'
}

// 获取宿舍楼名称
const getDormName = (dormId) => {
  const dorm = dormList.value.find(d => d.dormid === dormId)
  return dorm ? dorm.dormname : '全部宿舍'
}

// 获取分类的帖子数量
const getCategoryCount = (categoryId) => {
  if (categoryId === 'all') {
    return posts.value.length
  }
  return posts.value.filter(post => post.typeid === categoryId).length
}

// 获取宿舍楼的帖子数量
const getDormCount = (dormName) => {
  if (dormName === '全部宿舍') {
    return posts.value.length
  }
  return posts.value.filter(post => post.dormname === dormName).length
}

// 事件处理函数
const selectDorm = (dormId) => {
  selectedDorm.value = dormId
  console.log('切换宿舍楼:', getDormName(dormId))
}

const selectCategory = (typeid) => {
  selectedCategory.value = typeid
  console.log('切换分类:', getCategoryName(typeid))
}

const handleSearch = () => {
  console.log('搜索关键词:', searchKeyword.value)
}

const clearSearch = () => {
  searchKeyword.value = ''
}

const clearAllFilters = () => {
  searchKeyword.value = ''
  selectedCategory.value = 'all'
  selectedDorm.value = 'all'
}

const handlePublish = () => {
  console.log('跳转到发布帖子页面')
  router.push('/publish')
}

const goToProfile = () => {
  console.log('跳转到个人资料页面')
  // 这里跳转到个人资料页面
  // router.push('/profile')
}

// 生命周期
onMounted(() => {
  fetchPostTypes();
  // loadPosts();
  fetchDormList();
  fetchPostList();
})
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
  justify-content: center; /* 新增：水平居中 */
}

.search-container {
  display: flex;
  gap: 12px;
  max-width: 600px; /* 限制最大宽度 */
  width: 100%; /* 确保容器宽度 */
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
  padding: 0;
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.clear-btn:hover {
  color: #666;
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
  white-space: nowrap;
}

.search-btn:hover {
  background: #40a9ff;
}

/* 用户头像样式 */
.user-avatar-section {
  flex-shrink: 0;
}

.user-avatar {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  object-fit: cover;
  cursor: pointer;
  border: 2px solid #e8e8e8;
  transition: all 0.3s ease;
}

.user-avatar:hover {
  border-color: #1890ff;
  transform: scale(1.05);
}

/* 上边栏样式 - 分类筛选居中 */
.top-nav {
  background: white;
  border-radius: 12px;
  padding: 16px 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  margin-bottom: 20px;
}

.nav-content {
  display: flex;
  justify-content: center; /* 新增：水平居中 */
  align-items: center;
}

.category-filters {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  justify-content: center; /* 新增：分类按钮居中 */
}

.filter-label {
  font-size: 14px;
  color: #666;
  font-weight: 500;
  white-space: nowrap;
}

.filter-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  background: #f5f5f5;
  border: 1px solid #e8e8e8;
  border-radius: 20px;
  font-size: 13px;
  color: #666;
  cursor: pointer;
  transition: all 0.3s ease;
  white-space: nowrap;
}

.filter-btn:hover {
  background: #e6f7ff;
  border-color: #1890ff;
  color: #1890ff;
}

.filter-btn.active {
  background: #1890ff;
  border-color: #1890ff;
  color: white;
}

.category-count {
  background: rgba(255, 255, 255, 0.3);
  padding: 1px 6px;
  border-radius: 10px;
  font-size: 11px;
  min-width: 16px;
  text-align: center;
}

.filter-btn:not(.active) .category-count {
  background: #e8e8e8;
  color: #666;
}

/* 发布按钮 - 固定在右下角 */
.publish-fab {
  position: fixed;
  bottom: 30px;
  right: 30px;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 16px 24px;
  background: #1890ff;
  color: white;
  border: none;
  border-radius: 50px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  box-shadow: 0 4px 16px rgba(24, 144, 255, 0.3);
  transition: all 0.3s ease;
  z-index: 1000;
}

.publish-fab:hover {
  background: #40a9ff;
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(24, 144, 255, 0.4);
}

.fab-icon {
  font-size: 18px;
}

/* 主内容区样式 */
.main-content {
  display: flex;
  gap: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

/* 左边栏样式 */
.sidebar {
  width: 200px;
  background: white;
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  height: fit-content;
  position: sticky;
  top: 20px;
  flex-shrink: 0;
}

.sidebar-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin: 0 0 16px 0;
  padding-bottom: 12px;
  border-bottom: 1px solid #f0f0f0;
}

.dorm-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.dorm-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 1px solid transparent;
}

.dorm-item:hover {
  background: #f5f5f5;
}

.dorm-item.active {
  background: #e6f7ff;
  border-color: #1890ff;
  color: #1890ff;
}

.dorm-name {
  font-size: 14px;
  font-weight: 500;
}

.post-count {
  background: #f0f0f0;
  color: #666;
  padding: 2px 6px;
  border-radius: 10px;
  font-size: 11px;
  min-width: 20px;
  text-align: center;
}

.dorm-item.active .post-count {
  background: #1890ff;
  color: white;
}

/* 帖子区域样式 */
.posts-area {
  flex: 1;
  min-width: 0;
}

.filter-status {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #e6f7ff;
  border: 1px solid #91d5ff;
  border-radius: 8px;
  padding: 12px 16px;
  margin-bottom: 16px;
}

.status-text {
  font-size: 14px;
  color: #1890ff;
}

.post-count-text {
  color: #666;
}

.clear-filters {
  background: none;
  border: 1px solid #1890ff;
  color: #1890ff;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.3s;
}

.clear-filters:hover {
  background: #1890ff;
  color: white;
}

.posts-container {
  display: flex;
  flex-direction: column;
}

.post-card {
  margin-bottom: 16px;
}

/* 已显示所有帖子的提示 */
.all-posts-loaded {
  text-align: center;
  padding: 40px 20px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  margin-top: 20px;
}

.loaded-text {
  font-size: 16px;
  color: #666;
  margin-bottom: 8px;
}

.loaded-count {
  font-size: 14px;
  color: #999;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.empty-text {
  font-size: 16px;
  color: #666;
  margin-bottom: 20px;
}

.empty-btn {
  padding: 10px 20px;
  background: #1890ff;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  cursor: pointer;
  transition: background 0.3s;
  margin: 0 4px;
}

.empty-btn:hover {
  background: #40a9ff;
}

.empty-btn.secondary {
  background: white;
  color: #666;
  border: 1px solid #e8e8e8;
}

.empty-btn.secondary:hover {
  background: #f5f5f5;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .dormgo-home {
    padding: 12px;
  }
  
  .top-bar {
    flex-direction: column;
    gap: 12px;
  }
  
  .search-section {
    max-width: 100%;
  }
  
  .search-container {
    flex-direction: column;
  }
  
  .main-content {
    flex-direction: column;
  }
  
  .sidebar {
    width: 100%;
    position: static;
  }
  
  .nav-content {
    flex-direction: column;
    gap: 12px;
    align-items: stretch;
  }
  
  .category-filters {
    justify-content: center;
  }
  
  .publish-fab {
    bottom: 20px;
    right: 20px;
    padding: 14px 20px;
    font-size: 14px;
  }
}
</style>