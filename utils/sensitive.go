package utils

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// 敏感字段列表
var sensitiveFields = []string{
	"password", "pwd", "passwd",
	"mobile", "phone", "tel",
	"idcard", "id_card",
	"email",
	"card_no", "cardno", "bankcard",
}

// MaskSensitiveData 对敏感数据进行脱敏处理
func MaskSensitiveData(data interface{}) interface{} {
	if data == nil {
		return nil
	}

	switch v := data.(type) {
	case string:
		return v
	case map[string]interface{}:
		return maskMap(v)
	case []interface{}:
		return maskSlice(v)
	default:
		// 尝试将其他类型转换为map处理
		if jsonBytes, err := json.Marshal(data); err == nil {
			var mapData map[string]interface{}
			if err := json.Unmarshal(jsonBytes, &mapData); err == nil {
				return maskMap(mapData)
			}
		}
		return data
	}
}

// maskMap 处理map类型的数据
func maskMap(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range data {
		// 检查键名是否为敏感字段
		if isSensitiveField(k) {
			result[k] = maskValue(v, k)
		} else {
			// 递归处理值
			result[k] = MaskSensitiveData(v)
		}
	}
	return result
}

// maskSlice 处理切片类型的数据
func maskSlice(data []interface{}) []interface{} {
	result := make([]interface{}, len(data))
	for i, v := range data {
		result[i] = MaskSensitiveData(v)
	}
	return result
}

// isSensitiveField 判断字段名是否为敏感字段
func isSensitiveField(field string) bool {
	field = strings.ToLower(field)
	for _, sensitive := range sensitiveFields {
		if strings.Contains(field, sensitive) {
			return true
		}
	}
	return false
}

// maskValue 对敏感值进行脱敏处理
func maskValue(value interface{}, field string) string {
	if value == nil {
		return ""
	}

	str := fmt.Sprintf("%v", value)
	if str == "" {
		return ""
	}

	// 密码字段完全脱敏
	field = strings.ToLower(field)
	if strings.Contains(field, "password") || strings.Contains(field, "pwd") || strings.Contains(field, "passwd") {
		return "******"
	}

	// 邮箱脱敏
	if strings.Contains(str, "@") {
		return maskEmail(str)
	}

	// 手机号脱敏
	if regexp.MustCompile(`^1[3-9]\d{9}$`).MatchString(str) {
		return maskMobile(str)
	}

	// 身份证号脱敏
	if regexp.MustCompile(`^\d{17}[\dXx]$`).MatchString(str) {
		return maskIDCard(str)
	}

	// 其他敏感信息，保留前后各1/4的字符
	return maskString(str)
}

// maskEmail 邮箱脱敏
func maskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return maskString(email)
	}
	name := parts[0]
	if len(name) <= 2 {
		return "*" + name[len(name)-1:] + "@" + parts[1]
	}
	return name[:2] + "***" + "@" + parts[1]
}

// maskMobile 手机号脱敏
func maskMobile(mobile string) string {
	if len(mobile) != 11 {
		return maskString(mobile)
	}
	return mobile[:3] + "****" + mobile[7:]
}

// maskIDCard 身份证号脱敏
func maskIDCard(idCard string) string {
	if len(idCard) != 18 {
		return maskString(idCard)
	}
	return idCard[:6] + "********" + idCard[14:]
}

// maskString 通用字符串脱敏
func maskString(s string) string {
	length := len(s)
	if length == 0 {
		return ""
	}
	if length <= 2 {
		return "**"
	}

	// 保留前后各1/4的字符，中间用*代替
	showLen := length / 4
	if showLen < 1 {
		showLen = 1
	}
	stars := strings.Repeat("*", length-showLen*2)
	return s[:showLen] + stars + s[length-showLen:]
}
