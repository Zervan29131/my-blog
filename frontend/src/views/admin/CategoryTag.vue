<template>
  <div class="manage-container">
    <el-row :gutter="20">
      <!-- 左侧：分类管理 -->
      <el-col :span="12" :xs="24">
        <el-card class="box-card">
          <template #header>
            <div class="card-header">
              <span>分类管理</span>
              <el-button type="primary" size="small" @click="openCategoryDialog">
                <el-icon><Plus /></el-icon> 新增分类
              </el-button>
            </div>
          </template>
          
          <el-table :data="categories" v-loading="loadingCategory" style="width: 100%">
            <el-table-column prop="ID" label="ID" width="60" />
            <el-table-column prop="name" label="名称" />
            <el-table-column label="操作" width="100" align="right">
              <template #default="scope">
                <el-popconfirm 
                  title="确定删除该分类吗？" 
                  @confirm="handleDeleteCategory(scope.row.ID)"
                >
                  <template #reference>
                    <el-button type="danger" link size="small">删除</el-button>
                  </template>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>

      <!-- 右侧：标签管理 -->
      <el-col :span="12" :xs="24" class="tag-col">
        <el-card class="box-card">
          <template #header>
            <div class="card-header">
              <span>标签管理</span>
              <el-button type="success" size="small" @click="openTagDialog">
                <el-icon><Plus /></el-icon> 新增标签
              </el-button>
            </div>
          </template>

          <el-table :data="tags" v-loading="loadingTag" style="width: 100%">
            <el-table-column prop="ID" label="ID" width="60" />
            <el-table-column prop="name" label="名称" />
            <el-table-column label="操作" width="100" align="right">
              <template #default="scope">
                <el-popconfirm 
                  title="确定删除该标签吗？" 
                  @confirm="handleDeleteTag(scope.row.ID)"
                >
                  <template #reference>
                    <el-button type="danger" link size="small">删除</el-button>
                  </template>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <!-- 分类弹窗 -->
    <el-dialog v-model="categoryDialogVisible" title="新增分类" width="30%">
      <el-form :model="categoryForm" label-width="80px">
        <el-form-item label="分类名称">
          <el-input v-model="categoryForm.name" placeholder="请输入分类名称" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="categoryDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitCategory">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 标签弹窗 -->
    <el-dialog v-model="tagDialogVisible" title="新增标签" width="30%">
      <el-form :model="tagForm" label-width="80px">
        <el-form-item label="标签名称">
          <el-input v-model="tagForm.name" placeholder="请输入标签名称" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="tagDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitTag">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { getCategoryList, createCategory, deleteCategory, type Category } from '../../api/category'
import { getTagList, createTag, deleteTag, type Tag } from '../../api/tag'

// --- State Definitions ---
const categories = ref<Category[]>([])
const tags = ref<Tag[]>([])
const loadingCategory = ref(false)
const loadingTag = ref(false)

const categoryDialogVisible = ref(false)
const tagDialogVisible = ref(false)

const categoryForm = reactive({ name: '' })
const tagForm = reactive({ name: '' })

// --- Category Logic ---
const fetchCategories = async () => {
  loadingCategory.value = true
  try {
    const res: any = await getCategoryList()
    categories.value = res.data || []
  } catch (error) {
    console.error(error)
  } finally {
    loadingCategory.value = false
  }
}

const openCategoryDialog = () => {
  categoryForm.name = ''
  categoryDialogVisible.value = true
}

const submitCategory = async () => {
  if (!categoryForm.name) return ElMessage.warning('请输入分类名称')
  try {
    await createCategory({ name: categoryForm.name })
    ElMessage.success('添加成功')
    categoryDialogVisible.value = false
    fetchCategories()
  } catch (error) {
    console.error(error)
  }
}

const handleDeleteCategory = async (id: number) => {
  try {
    await deleteCategory(id)
    ElMessage.success('删除成功')
    fetchCategories()
  } catch (error) {
    console.error(error)
  }
}

// --- Tag Logic ---
const fetchTags = async () => {
  loadingTag.value = true
  try {
    const res: any = await getTagList()
    tags.value = res.data || []
  } catch (error) {
    console.error(error)
  } finally {
    loadingTag.value = false
  }
}

const openTagDialog = () => {
  tagForm.name = ''
  tagDialogVisible.value = true
}

const submitTag = async () => {
  if (!tagForm.name) return ElMessage.warning('请输入标签名称')
  try {
    await createTag({ name: tagForm.name })
    ElMessage.success('添加成功')
    tagDialogVisible.value = false
    fetchTags()
  } catch (error) {
    console.error(error)
  }
}

const handleDeleteTag = async (id: number) => {
  try {
    await deleteTag(id)
    ElMessage.success('删除成功')
    fetchTags()
  } catch (error) {
    console.error(error)
  }
}

// --- Initialization ---
onMounted(() => {
  fetchCategories()
  fetchTags()
})
</script>

<style scoped>
.manage-container {
  max-width: 1200px;
  margin: 20px auto;
  padding: 0 20px;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.tag-col {
  margin-top: 0;
}
@media (max-width: 768px) {
  .tag-col {
    margin-top: 20px;
  }
}
</style>