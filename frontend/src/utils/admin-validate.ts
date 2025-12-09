// src/utils/admin-validate.ts
/**
 * 管理员账号格式验证（学号/工号规则）
 * 规则：必须是数字，长度在 5-20 位之间
 */
export const validateStudentId = (studentId: string): boolean => {
    const regex = /^\d{5,20}$/
    return regex.test(studentId)
}

/**
 * 密码强度校验
 * 返回值：'weak' | 'medium' | 'strong'
 * 规则：
 * Weak: < 6字符
 * Medium: > 6字符，包含字母和数字
 * Strong: > 8字符，包含大小写字母、数字和特殊符号
 */
export const checkPasswordStrength = (password: string): 'weak' | 'medium' | 'strong' => {
    if (password.length < 6) return 'weak'

    const hasDigit = /\d/.test(password)
    const hasLower = /[a-z]/.test(password)
    const hasUpper = /[A-Z]/.test(password)
    const hasSpecial = /[!@#$%^&*(),.?":{}|<>]/.test(password)

    if (password.length >= 8 && hasDigit && hasLower && hasUpper && hasSpecial) {
        return 'strong'
    }

    if (hasDigit && (hasLower || hasUpper)) {
        return 'medium'
    }

    return 'weak'
}