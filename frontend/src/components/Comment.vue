<template>
  <div class="comment-section">
    <div class="section-title">
      <h3>评论 ({{ total }})</h3>
    </div>

    <!-- 1. 顶部主要发表框 -->
    <div class="comment-form-card">
      <el-form :model="form" :rules="rules" ref="formRef">
        <el-row :gutter="15">
          <el-col :span="12" :xs="24">
            <el-form-item prop="nickname">
              <el-input v-model="form.nickname" placeholder="昵称 (必填)" prefix-icon="User" />
            </el-form-item>
          </el-col>
          <el-col :span="12" :xs="24">
            <el-form-item prop="email">
              <el-input v-model="form.email" placeholder="邮箱 (选填, 用于头像)" prefix-icon="Message" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item prop="content">
          <el-input 
            v-model="form.content" 
            type="textarea" 
            :rows="3" 
            placeholder="写下你的想法..." 
            resize="none"
          />
        </el-form-item>
        <div class="form-footer">
          <el-button type="primary" @click="submitComment(null)" :loading="submitting">发表评论</el-button>
        </div>
      </el-form>
    </div>

    <!-- 2. 评论列表 -->
    <div class="comment-list">
      <el-empty v-if="comments.length === 0" description="暂无评论，快来抢沙发！" />
      
      <div v-else v-for="item in commentTree" :key="item.ID" class="comment-item">
        <!-- 父评论 -->
        <div class="comment-main">
          <el-avatar :size="40" :src="getAvatar(item.nickname)" />
          <div class="comment-body">
            <div class="comment-header">
              <span class="nickname">{{ item.nickname }}</span>
              <span class="date">{{ formatDate(item.created_at) }}</span>
            </div>
            <div class="comment-content">{{ item.content }}</div>
            <div class="comment-actions">
              <el-button link type="primary" size="small" @click="toggleReply(item.ID)">
                {{ replyTargetId === item.ID ? '取消回复' : '回复' }}
              </el-button>
            </div>
            
            <!-- 回复框 (当点击回复时显示) -->
            <div v-if="replyTargetId === item.ID" class="reply-form">
              <el-input v-model="replyContent" placeholder="回复内容..." size="small">
                <template #append>
                  <el-button @click="submitComment(item.ID)" :loading="submitting">发送</el-button>
                </template>
              </el-input>
            </div>
          </div>
        </div>

        <!-- 子评论 (递归展示或简单的二级展示) -->
        <div v-if="item.children && item.children.length > 0" class="sub-comments">
          <div v-for="child in item.children" :key="child.ID" class="comment-item sub-item">
             <el-avatar :size="30" :src="getAvatar(child.nickname)" />
             <div class="comment-body">
                <div class="comment-header">
                  <span class="nickname">{{ child.nickname }}</span>
                  <span class="date">{{ formatDate(child.created_at) }}</span>
                </div>
                <div class="comment-content">
                  <span class="reply-tag">@{{ item.nickname }}</span> {{ child.content }}
                </div>
             </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { User, Message } from '@element-plus/icons-vue'
import { getComments, createComment, type Comment } from '../api/comment'

const props = defineProps<{
  postId: number
}>()

const comments = ref<Comment[]>([])
const submitting = ref(false)
const replyTargetId = ref<number | null>(null)
const replyContent = ref('')
const formRef = ref()

const form = reactive({
  nickname: '',
  email: '',
  content: '',
  website: ''
})

const rules = {
  nickname: [{ required: true, message: '请输入昵称', trigger: 'blur' }],
  content: [{ required: true, message: '请输入评论内容', trigger: 'blur' }]
}

const total = computed(() => comments.value.length)

// 将扁平数组转换为树形结构
const commentTree = computed(() => {
  const map: Record<number, Comment> = {}
  const roots: Comment[] = []
  
  // 深拷贝并建立索引
  const raw = JSON.parse(JSON.stringify(comments.value)) as Comment[]
  raw.forEach(c => {
    c.children = []
    map[c.ID] = c
  })

  // 组装树
  raw.forEach(c => {
    if (c.parent_id && map[c.parent_id]) {
      map[c.parent_id].children?.push(c)
    } else {
      roots.push(c)
    }
  })
  
  return roots
})

// 获取头像 (使用 UI Avatars 服务)
const getAvatar = (name: string) => {
  return `https://ui-avatars.com/api/?name=${encodeURIComponent(name)}&background=random&color=fff`
}

const formatDate = (str: string) => {
  return new Date(str).toLocaleString()
}

// 加载评论
const loadComments = async () => {
  if (!props.postId) return
  try {
    const res: any = await getComments(props.postId)
    comments.value = res.data || []
  } catch (error) {
    console.error(error)
  }
}

// 切换回复框
const toggleReply = (id: number) => {
  if (replyTargetId.value === id) {
    replyTargetId.value = null
  } else {
    replyTargetId.value = id
    replyContent.value = ''
  }
}

// 提交评论 (parentId 为 null 则是主评论)
const submitComment = async (parentId: number | null) => {
  // 如果是主评论，校验表单
  if (!parentId) {
    await formRef.value.validate(async (valid: boolean) => {
      if (!valid) return
      await doSubmit({
        post_id: props.postId,
        parent_id: null,
        ...form
      })
      // 清空主表单
      form.content = ''
    })
  } else {
    // 如果是回复
    if (!replyContent.value.trim()) return ElMessage.warning('请输入回复内容')
    // 回复时如果没有填昵称，先弹窗提示或使用匿名 (这里简单处理，假设用户已在上方填过昵称，或者你需要弹窗让用户填昵称)
    // 为了简化，这里直接复用上方表单的昵称，如果没填则提示
    if (!form.nickname) return ElMessage.warning('请先在上方填写昵称')
    
    await doSubmit({
      post_id: props.postId,
      parent_id: parentId,
      nickname: form.nickname,
      email: form.email,
      content: replyContent.value
    })
    replyTargetId.value = null
  }
}

const doSubmit = async (data: any) => {
  submitting.value = true
  try {
    await createComment(data)
    ElMessage.success('评论成功')
    loadComments()
  } catch (error) {
    console.error(error)
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadComments()
})
</script>

<style scoped>
.comment-section {
  margin-top: 60px;
  padding-top: 40px;
  border-top: 1px solid #f0f0f0;
}

.section-title h3 {
  margin-bottom: 20px;
  font-size: 20px;
  color: #303133;
}

.comment-form-card {
  background: #f9fafb;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 40px;
}

.form-footer {
  text-align: right;
}

.comment-list {
  display: flex;
  flex-direction: column;
  gap: 30px;
}

.comment-item {
  display: flex;
  flex-direction: column;
}

.comment-main {
  display: flex;
  gap: 15px;
}

.comment-body {
  flex: 1;
}

.comment-header {
  margin-bottom: 5px;
}

.nickname {
  font-weight: 600;
  color: #333;
  margin-right: 10px;
}

.date {
  font-size: 12px;
  color: #999;
}

.comment-content {
  color: #606266;
  line-height: 1.6;
  font-size: 15px;
}

.comment-actions {
  margin-top: 5px;
}

.reply-form {
  margin-top: 10px;
}

/* 子评论样式 */
.sub-comments {
  margin-left: 55px;
  margin-top: 15px;
  background: #f9fafb;
  padding: 15px;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.sub-item {
  flex-direction: row;
  gap: 10px;
}

.reply-tag {
  color: #409eff;
  font-weight: 500;
  margin-right: 5px;
}
</style>