// src/utils/batch-operation.ts
import { ElMessageBox, ElMessage } from 'element-plus'
import request from '@/utils/request'

/**
 * 批量更新用户状态（封禁/解封）
 */
export const batchUpdateUserStatus = async (userIds: number[], status: 'normal' | 'disabled'): Promise<void> => {
    if (userIds.length === 0) {
        ElMessage.warning('请至少选择一项')
        return
    }

    const actionText = status === 'normal' ? '解封' : '封禁'

    try {
        await ElMessageBox.confirm(
            `确定要批量${actionText}选中的 ${userIds.length} 位用户吗？`,
            '批量操作警告',
            { type: 'warning', confirmButtonText: '确定执行', cancelButtonText: '取消' }
        )

        await request.put('/v1/admin/users/batch/status', { ids: userIds, status })
        ElMessage.success(`批量${actionText}成功`)
    } catch (error) {
        if (error !== 'cancel') console.error(error)
        throw error // 抛出错误以便组件刷新列表
    }
}

/**
 * 批量审核互助内容
 */
export const batchAuditPosts = async (postIds: number[], status: 'approved' | 'rejected'): Promise<void> => {
    if (postIds.length === 0) {
        ElMessage.warning('请至少选择一项')
        return
    }

    const actionText = status === 'approved' ? '通过' : '拒绝'

    try {
        await ElMessageBox.confirm(
            `确定要批量${actionText}选中的 ${postIds.length} 条内容吗？`,
            '批量审核',
            { type: 'info', confirmButtonText: '确定', cancelButtonText: '取消' }
        )

        await request.put('/v1/admin/contents/batch/audit', { ids: postIds, status })
        ElMessage.success(`批量${actionText}成功`)
    } catch (error) {
        if (error !== 'cancel') console.error(error)
        throw error
    }
}