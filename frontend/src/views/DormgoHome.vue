<!-- src/views/DormgoHome.vue -->
<template>
  <div class="dormgo-home">
    <!-- 搜索框 -->
    <div class="search-section">
      <div class="search-container">
        <div class="search-box">
          <span class="search-icon">🔍</span>
          <input
            v-model="searchKeyword"
            type="text"
            placeholder="搜索帖子标题或内容..."
            class="search-input"
            @input="handleSearch"
          />
          <button v-if="searchKeyword" class="clear-btn" @click="clearSearch">×</button>
        </div>
        <button class="search-btn" @click="handleSearch">搜索</button>
      </div>
    </div>

    <!-- 上边栏 - 分类筛选 -->
    <div class="top-nav">
      <div class="nav-content">
        <div class="category-filters">
          <span class="filter-label">分类：</span>
          <button
            v-for="category in categories"
            :key="category.id"
            class="filter-btn"
            :class="{ active: selectedCategory === category.id }"
            @click="selectCategory(category.id)"
          >
            {{ category.name }}
            <span class="category-count">{{ getCategoryCount(category.id) }}</span>
          </button>
        </div>
        
        <div class="nav-actions">
          <button class="publish-btn" @click="handlePublish">
            <span class="btn-icon">📝</span>
            发布帖子
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
            :key="dorm.id"
            class="dorm-item"
            :class="{ active: selectedDorm === dorm.id }"
            @click="selectDorm(dorm.id)"
          >
            <span class="dorm-name">{{ dorm.name }}</span>
            <span class="post-count">{{ getDormCount(dorm.name) }}</span>
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

        <div class="posts-container">
          <Card 
            v-for="post in filteredPosts"
            :key="post.id"
            :post="post"
            class="post-card"
          />
        </div>
        
        <!-- 加载更多 -->
        <div class="load-more" v-if="hasMorePosts && filteredPosts.length > 0">
          <button class="load-more-btn" @click="loadMorePosts">
            加载更多帖子
          </button>
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
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import Card from '@/components/Card.vue'

// 响应式数据
const selectedDorm = ref('all') // 默认显示全部宿舍楼
const selectedCategory = ref('all')
const searchKeyword = ref('')
const posts = ref([])

// 宿舍楼数据
const dormList = ref([
  { id: 'all', name: '全部宿舍' },
  { id: 'rong9', name: '榕园9号' },
  { id: 'rong8', name: '榕园8号' },
  { id: 'rong7', name: '榕园7号' },
  { id: 'rong6', name: '榕园6号' },
  { id: 'rong5', name: '榕园5号' },
  { id: 'xin1', name: '馨园1号' },
  { id: 'xin2', name: '馨园2号' },
  { id: 'xin3', name: '馨园3号' },
  { id: 'zhi1', name: '芷园1号' },
  { id: 'zhi2', name: '芷园2号' }
])

// 分类数据
const categories = ref([
  { id: 'all', name: '全部' },
  { id: 'food', name: '约饭' },
  { id: 'sports', name: '约球' },
  { id: 'help', name: '求助' },
  { id: 'trade', name: '交易' },
  { id: 'study', name: '学习' }
])

// 模拟帖子数据 - 每个帖子都有分类和宿舍楼信息
const mockPosts = [
  {
    id: '1',
    userName: '小明同学',
    userAvatar: '/avatars/1.jpg',
    time: '2小时前',
    category: 'food',
    dormBuilding: '榕园9号',
    title: '今晚有人一起去食堂吃饭吗？',
    content: '一个人吃饭太无聊了，想找几个同学一起去食堂，可以聊聊天，交流一下学习心得。',
    images: ['/posts/food1.jpg'],
    commentCount: 5,
    viewCount: 32,
    likeCount: 8
  },
  {
    id: '2',
    userName: '篮球少年',
    userAvatar: '/avatars/2.jpg',
    time: '5小时前',
    category: 'sports',
    dormBuilding: '榕园8号',
    title: '明天下午篮球场约球，3V3缺两人',
    content: '明天下午4点在东区篮球场，现有4人，还缺2个，有兴趣的同学欢迎加入！篮球运动对身体很好。',
    images: ['/posts/sports1.jpg', '/posts/sports2.jpg'],
    commentCount: 12,
    viewCount: 45,
    likeCount: 15
  },
  {
    id: '3',
    userName: '学习委员',
    userAvatar: '/avatars/3.jpg',
    time: '1天前',
    category: 'study',
    dormBuilding: '榕园9号',
    title: '高数复习小组招人',
    content: '准备期末高数考试，组建复习小组，每周三、五晚上在图书馆讨论区一起学习数学。',
    images: [],
    commentCount: 8,
    viewCount: 28,
    likeCount: 12
  },
  {
    id: '4',
    userName: '急需帮助',
    userAvatar: '/avatars/4.jpg',
    time: '1天前',
    category: 'help',
    dormBuilding: '馨园1号',
    title: '电脑突然蓝屏开不了机，求大神帮忙',
    content: '今天早上电脑突然蓝屏，重启后还是不行，有没有懂电脑的同学能帮忙看看？非常感谢！',
    images: ['/posts/help1.jpg'],
    commentCount: 15,
    viewCount: 67,
    likeCount: 6
  },
  {
    id: '5',
    userName: '二手交易',
    userAvatar: '/avatars/5.jpg',
    time: '2天前',
    category: 'trade',
    dormBuilding: '芷园1号',
    title: '转让几乎全新的机械键盘',
    content: 'Cherry MX红轴，买来用了不到一个月，因换笔记本用不上了，原价450，现300出。',
    images: ['/posts/trade1.jpg', '/posts/trade2.jpg', '/posts/trade3.jpg'],
    commentCount: 3,
    viewCount: 23,
    likeCount: 4
  },
  {
    id: '6',
    userName: '美食家',
    userAvatar: '/avatars/6.jpg',
    time: '3小时前',
    category: 'food',
    dormBuilding: '榕园7号',
    title: '周末想约火锅，有人一起吗？',
    content: '发现学校附近新开了一家重庆火锅，想周末去尝尝，找几个爱吃火锅的同学一起！',
    images: ['/posts/food2.jpg'],
    commentCount: 7,
    viewCount: 29,
    likeCount: 10
  },
  {
    id: '7',
    userName: '羽毛球爱好者',
    userAvatar: '/avatars/7.jpg',
    time: '4小时前',
    category: 'sports',
    dormBuilding: '馨园2号',
    title: '晚上羽毛球馆约球',
    content: '晚上7点体育馆羽毛球场地，现有3人，还缺1人，欢迎喜欢羽毛球的朋友加入！',
    images: [],
    commentCount: 4,
    viewCount: 18,
    likeCount: 6
  },
  {
    id: '8',
    userName: '考研党',
    userAvatar: '/avatars/8.jpg',
    time: '6小时前',
    category: 'study',
    dormBuilding: '榕园6号',
    title: '寻找考研自习伙伴',
    content: '准备2024年考研，想在图书馆长期自习，寻找志同道合的小伙伴互相监督。',
    images: [],
    commentCount: 9,
    viewCount: 34,
    likeCount: 11
  }
]

// 计算属性：筛选帖子（包含搜索、分类、宿舍楼筛选）
const filteredPosts = computed(() => {
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
    filtered = filtered.filter(post => post.category === selectedCategory.value)
  }
  
  // 按宿舍楼筛选
  if (selectedDorm.value !== 'all') {
    const dormName = getDormName(selectedDorm.value)
    filtered = filtered.filter(post => post.dormBuilding === dormName)
  }
  
  return filtered
})

// 计算属性：是否还有更多帖子
const hasMorePosts = computed(() => {
  return filteredPosts.value.length < 50
})

// 计算属性：是否有活跃的筛选条件
const hasActiveFilters = computed(() => {
  return searchKeyword.value.trim() !== '' || selectedCategory.value !== 'all' || selectedDorm.value !== 'all'
})

// 获取分类名称
const getCategoryName = (categoryId) => {
  const category = categories.value.find(cat => cat.id === categoryId)
  return category ? category.name : '全部'
}

// 获取宿舍楼名称
const getDormName = (dormId) => {
  const dorm = dormList.value.find(d => d.id === dormId)
  return dorm ? dorm.name : '全部宿舍'
}

// 获取分类的帖子数量
const getCategoryCount = (categoryId) => {
  if (categoryId === 'all') {
    return posts.value.length
  }
  return posts.value.filter(post => post.category === categoryId).length
}

// 获取宿舍楼的帖子数量
const getDormCount = (dormName) => {
  if (dormName === '全部宿舍') {
    return posts.value.length
  }
  return posts.value.filter(post => post.dormBuilding === dormName).length
}

// 事件处理函数
const selectDorm = (dormId) => {
  selectedDorm.value = dormId
  console.log('切换宿舍楼:', getDormName(dormId))
}

const selectCategory = (categoryId) => {
  selectedCategory.value = categoryId
  console.log('切换分类:', getCategoryName(categoryId))
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
  console.log('发布帖子')
  // 这里可以跳转到发布页面或打开发布弹窗
}

const loadMorePosts = () => {
  console.log('加载更多帖子')
  // 这里可以调用API加载更多帖子
}

// 生命周期
onMounted(() => {
  // 模拟获取帖子数据
  posts.value = mockPosts
})
</script>

<style scoped>
/* 样式部分与之前相同，保持不变 */
.dormgo-home {
  min-height: 100vh;
  background: #f5f5f5;
  padding: 20px;
}

.search-section {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  margin-bottom: 20px;
}

.search-container {
  display: flex;
  gap: 12px;
  max-width: 600px;
  margin: 0 auto;
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

.top-nav {
  background: white;
  border-radius: 12px;
  padding: 16px 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  margin-bottom: 20px;
}

.nav-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.category-filters {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
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

.nav-actions {
  display: flex;
  gap: 12px;
}

.publish-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: #1890ff;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.3s;
}

.publish-btn:hover {
  background: #40a9ff;
}

.btn-icon {
  font-size: 16px;
}

.main-content {
  display: flex;
  gap: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

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

.load-more {
  text-align: center;
  padding: 20px;
}

.load-more-btn {
  padding: 10px 20px;
  background: white;
  border: 1px solid #e8e8e8;
  border-radius: 8px;
  color: #666;
  cursor: pointer;
  transition: all 0.3s;
}

.load-more-btn:hover {
  background: #f0f0f0;
  border-color: #d9d9d9;
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

@media (max-width: 768px) {
  .dormgo-home {
    padding: 12px;
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
  
  .search-container {
    flex-direction: column;
  }
}
</style>