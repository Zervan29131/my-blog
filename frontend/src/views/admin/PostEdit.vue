<template>
  <div class="edit-container">
    <el-card class="edit-card">
      <template #header>
        <div class="card-header">
          <h2>{{ isEdit ? '编辑文章' : '写新文章' }}</h2>
          <el-button @click="$router.back()">返回</el-button>
        </div>
      </template>

      <el-form :model="form" label-width="80px">
        <el-form-item label="标题">
          <el-input v-model="form.title" placeholder="请输入文章标题" />
        </el-form-item>

        <el-form-item label="分类ID">
          <!-- 简单起见暂时手动输入ID，后续可改为下拉选择 -->
          <el-input-number v-model="form.category_id" :min="1" />
          <span class="tips"> (请确保数据库有对应的分类ID)</span>
        </el-form-item>

        <el-form-item label="摘要">
          <el-input v-model="form.summary" type="textarea" :rows="2" placeholder="简短的摘要..." />
        </el-form-item>

        <el-form-item label="内容">
          <el-input 
            v-model="form.content" 
            type="textarea" 
            :rows="15" 
            placeholder="支持 Markdown 格式..." 
            class="content-editor"
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSubmit" :loading="loading">
            {{ isEdit ? '保存修改' : '立即发布' }}
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { createPost, updatePost, getPost } from '../../api/post'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const form = reactive({
  title: '',
  content: '',
  summary: '',
  category_id: 1 // 默认分类
})

// 判断是否是编辑模式（URL 中是否有 id 参数）
const isEdit = computed(() => !!route.params.id)

onMounted(async () => {
  if (isEdit.value) {
    const id = Number(route.params.id)
    try {
      const res: any = await getPost(id)
      // 回显数据
      form.title = res.data.title
      form.content = res.data.content
      form.summary = res.data.summary
      form.category_id = res.data.category_id
    } catch (error) {
      console.error(error)
      ElMessage.error('获取文章详情失败')
    }
  }
})

const handleSubmit = async () => {
  if (!form.title || !form.content) {
    ElMessage.warning('标题和内容不能为空')
    return
  }

  loading.value = true
  try {
    if (isEdit.value) {
      // 更新
      await updatePost(Number(route.params.id), form)
      ElMessage.success('更新成功')
    } else {
      // 创建
      await createPost(form)
      ElMessage.success('发布成功')
    }
    router.push('/') // 回到列表页
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.edit-container {
  max-width: 800px;
  margin: 20px auto;
  padding: 0 20px;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.tips {
  margin-left: 10px;
  color: #999;
  font-size: 12px;
}
</style>