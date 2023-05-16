package strtools

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func Limit(str, joint string, length int) string {
	if str == "" {
		return ""
	}

	strByte := []rune(str)
	if utf8.RuneCountInString(str) > length {
		return RemoveInvisibleStr(string(append(strByte[:length], []rune(joint)...)))
	}
	return str
}

// 掩码手机
func MaskMobile(mobile string) string {
	pos := 4
	if len(mobile) < 4 {
		pos = len(mobile)
	}
	suffix := mobile[len(mobile)-pos:]
	prefix := strings.Repeat("*", len(mobile)-pos)

	return string(prefix) + suffix
}

// 掩码证件号码
func MaskIdCardNo(idCardNo string) string {
	// 可能有中文，要转换为字符切片
	r := []rune(idCardNo)
	l := len(idCardNo)
	rl := len(r)
	// 如果有中文特殊处理
	if rl > 0 && l != rl {
		// 只显示第一位
		return string(r[0]) + strings.Repeat("*", len(r)-1)
	}

	// 如果没有中文，则进入正常的流程
	// 15位身份证
	if l == 15 {
		return idCardNo[:8] + "******" + idCardNo[14:]
	}

	// 18位身份证
	if l == 18 {
		return idCardNo[:10] + "******" + idCardNo[16:17] + "*"
	}

	// 其他 港澳等
	if l > 8 {
		return idCardNo[:5] + strings.Repeat("*", l-5)
	}
	// 不明数据全加密
	return strings.Repeat("*", l)
}

// 删除字符串中开头和结尾显示为“空格”的字符，以及字符串中所有的不可见字符（无法打印出来的）
func TrimInvisibleStr(str string) string {
	return strings.TrimSpace(strings.TrimFunc(str, func(r rune) bool {
		return !unicode.IsGraphic(r)
	}))
}

// RemoveInvisibleStr 移除不可见字符
func RemoveInvisibleStr(str string) string {
	rs := make([]rune, 0, len(str))
	for _, r := range str {
		if unicode.IsGraphic(r) && !unicode.IsSpace(r) && r != '\ufffd' {
			rs = append(rs, r)
		}
	}
	return string(rs)
}

// 字符串是否包含中文字符
func IsIncludeChinese(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Han, r) {
			return true
		}
	}
	return false
}

// 新的掩码证件号码
func NewMaskIdCardNo(idCardNo string) string {
	// 可能有中文，要转换为字符切片
	r := []rune(idCardNo)
	l := len(idCardNo)
	rl := len(r)
	// 如果有中文特殊处理
	if rl > 0 && l != rl {
		// 只显示第一位
		return string(r[0]) + strings.Repeat("*", len(r)-1)
	}

	// 如果没有中文，则进入正常的流程
	// 4位及以下证件号码：显示全部
	if l <= 4 {
		return idCardNo
	}
	// 5-12显示后4位
	if l >= 5 && l <= 12 {
		return strings.Repeat("*", len(r)-4) + string(r[l-4:])
	}
	// 13位及以上显示前6位和后3位
	if l >= 13 {
		return string(r[:6]) + strings.Repeat("*", len(r)-9) + string(r[len(r)-3:])
	}

	// 不明数据全加密
	return strings.Repeat("*", l)
}
