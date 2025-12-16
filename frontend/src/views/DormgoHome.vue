<template>
  <div class="dormgo-home">
    <!-- 顶部栏：搜索框和用户头像 -->
    <div class="top-bar">
      <div class="message-entry" @click="goToMessageList">
        <div class="icon-wrapper">
          <span class="msg-icon">💬</span> <div v-if="totalUnread > 0" class="dot-badge"></div> 
        </div>
        <span class="msg-text">私信</span>
      </div>
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
        <div class="avatar-wrapper" @click="goToProfile">
            <img 
              :src="currentUser.avatarurl" 
              alt="用户头像" 
              class="user-avatar"
            />
            <div v-if="unreadCount > 0" class="global-badge">
                {{ unreadCount > 99 ? '99+' : unreadCount }}
            </div>
          </div>
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
          <button v-if="hasActiveFilters" class="empty-btn secondary" @click="clearAllFilters">查看全部帖子</button>
        </div>
        <div v-if="isLoading" class="loading-state">
           ⏳ 正在加载更多帖子...
        </div>
        <!-- 已显示所有帖子的提示 -->
        <div class="all-posts-loaded" v-if="filteredPosts.length > 0 && !hasMore ">
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
import { ref, computed,onUnmounted ,onMounted ,onActivated, onDeactivated} from 'vue'
import { useRouter } from 'vue-router'
import Card from '@/components/Card.vue'
import axios from 'axios'
import api from '@/api/index.js';
import { useUserStore } from '@/stores/user';

defineOptions({
  name: 'DormgoHome'
})

const userStore = useUserStore();
const router = useRouter()
const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'

// 响应式数据
const selectedDorm = ref('all')
const selectedCategory = ref('all')
const searchKeyword = ref('')

// 当前用户信息
const currentUser = computed(() => userStore.currentUser)
const currentUserId = computed(() => userStore.currentUserId)

const dormList = ref([]);   // 宿舍列表
const page = ref(1);
const hasMore = ref(true);
const isLoading = ref(false);
const totalUnread = ref(0) // 总未读数

const unreadCount = computed(() => userStore.unreadCount);
let pollingTimer = null;

const goToMessageList = () => {
  router.push('/messages') // 记得在 router/index.js 里配这个路由
}
// 获取总未读数
const fetchTotalUnread = async () => {
    if (!userStore.isLoggedIn) return
    try {
        const res = await api.get('/api/v1/message/unread_count')
        if (res.data.code === 200) {
            totalUnread.value = res.data.data
        }
    } catch (e) {
        console.error(e)
    }
}

const fetchDormList = async() => {    
  try {
    const response = await api.get('/api/v1/post/dorms');
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
    const response = await api.get('/api/v1/post/post_type');
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
const fetchPostList = async (isRefresh = false) => {
  if (isLoading.value || (!hasMore.value && !isRefresh)) return;
  isLoading.value = true; // 上锁
  if (isRefresh) {
    page.value = 1;
    hasMore.value = true; // 重置 hasMore
  }
  try {
    const response = await api.get(`/api/v1/post/posts?page=${page.value}&size=10`);
    if (response.data && response.data.code === 200) {
      const newPosts = response.data.data;
      if (isRefresh) {
         posts.value = newPosts;
      } else {
         posts.value = [...posts.value, ...newPosts];
      }
      hasMore.value = response.data.has_more;
      if (hasMore.value) {
        page.value++; 
      }
      console.log(`加载第 ${page.value-1} 页成功, 当前总数: ${posts.value.length}`);
      console.log('帖子列表:', posts.value);
    } 
  } catch (error) {
    console.error('获取帖子列表失败:', error);
  }finally {
    isLoading.value = false; 
  }
};

const handleScroll = () => {

  const scrollTop = document.documentElement.scrollTop || document.body.scrollTop; 
  const clientHeight = document.documentElement.clientHeight; 
  const scrollHeight = document.documentElement.scrollHeight; 

  // 距离底部还有 100px 时，提前加载
  if (scrollTop + clientHeight >= scrollHeight - 100) {
    // 如果还有更多数据，且当前不在加载中
    if (hasMore.value && !isLoading.value) {
        console.log("触底加载下一页...");
        fetchPostList();
    }
  }
}


const filteredPosts = computed(() => {  // 筛选帖子
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
  console.log('跳转到个人资料页面')
  router.push('/personalhome')
}
onActivated(() => {
  if (userStore.isLoggedIn && !pollingTimer) {
      fetchTotalUnread()
      userStore.fetchUnreadCount();
      // 重新开启轮询
      pollingTimer = setInterval(() => {
          fetchTotalUnread()
          userStore.fetchUnreadCount();
      }, 3000); 
  }
})

// 4. 将清除逻辑移到 Deactivated (页面隐藏时)
onDeactivated(() => {
  if (pollingTimer) {
    clearInterval(pollingTimer);
    pollingTimer = null; // 清空变量
    console.log("暂停首页轮询");
  }
})
// 生命周期
onMounted(() => {
  fetchPostTypes();
  fetchDormList();
  fetchPostList(true);
  window.addEventListener('scroll', handleScroll);
  
});

onUnmounted(() => {
  if (pollingTimer) clearInterval(pollingTimer);
  window.addEventListener('scroll', handleScroll);
});

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
  justify-content: flex-start;
  gap: 20px;
  margin-bottom: 20px;
}

/* 搜索框样式 - 居中 */
.search-section {
  flex: 1;
  display: flex;
  justify-content: center; /* 新增：水平居中 */
  padding-right: 40px;
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

/* 新增样式 */
.avatar-wrapper {
  position: relative; /* 为了让红点绝对定位 */
  display: inline-block;
  cursor: pointer;
}

.global-badge {
  position: absolute;
  top: -2px;
  right: -2px;
  background-color: #ff4d4f;
  color: white;
  border-radius: 10px;
  padding: 0 6px;
  font-size: 12px;
  height: 18px;
  line-height: 18px;
  min-width: 18px;
  text-align: center;
  border: 2px solid white; /* 增加一点白边，更好看 */
  font-weight: bold;
  z-index: 10;
}
/* 在 style scoped 底部添加 */
.loading-state {
  text-align: center;
  padding: 20px;
  color: #1890ff;
  font-size: 14px;
}

.message-entry {
  display: flex;
  flex-direction: column;
  align-items: center;
  cursor: pointer;
  margin-right: 10px; /* 与搜索框拉开距离 */
}

.icon-wrapper {
  position: relative;
  font-size: 24px;
}

.msg-icon {
    /* 也可以换成 svg 或 img */
    font-size: 22px; 
}

.dot-badge {
  position: absolute;
  top: -2px;
  right: -2px;
  width: 8px;
  height: 8px;
  background: #ff4d4f;
  border-radius: 50%;
}

.msg-text {
  font-size: 10px;
  color: #666;
  margin-top: -4px;
}

</style>