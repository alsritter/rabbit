package helpers

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
)

var (
	priceLingJiaoLingFenRegex = regexp.MustCompile(`零角零分$`)
	priceLingJiaoRegex        = regexp.MustCompile(`零角`)
	priceLingFenRegex         = regexp.MustCompile(`零分$`)
	priceLingQianBaiShiRegex  = regexp.MustCompile(`零[仟佰拾]`)
	priceLing2Regex           = regexp.MustCompile(`零{2,}`)
	priceLingYiRegex          = regexp.MustCompile(`零亿`)
	priceLingWanRegex         = regexp.MustCompile(`零万`)
	priceLingMultiYuanRegex   = regexp.MustCompile(`零*元`)
	priceYiLingWanRegex       = regexp.MustCompile(`亿零{0, 3}万`)
	priceLingYuanRegex        = regexp.MustCompile(`零元`)
	priceUnitList             = [14]string{"仟", "佰", "拾", "亿", "仟", "佰", "拾", "万", "仟", "佰", "拾", "元", "角", "分"}
	priceUpperUnitMap         = map[string]string{"0": "零", "1": "壹", "2": "贰", "3": "叁", "4": "肆", "5": "伍", "6": "陆", "7": "柒", "8": "捌", "9": "玖"}
)

// PriceParseMoneyCapital 将float64 转换成中文金额字符串
func PriceParseMoneyCapital(num float64) string {
	// 是否为负数
	isNegative := num < 0
	strNum := ""
	str := ""
	if isNegative {
		strNum = strconv.FormatFloat(num*-100, 'f', 0, 64)
		str = "负"
	} else {
		strNum = strconv.FormatFloat(num*100, 'f', 0, 64)
	}
	s := priceUnitList[len(priceUnitList)-len(strNum):]

	for k, v := range strNum {
		str = str + priceUpperUnitMap[string(v)] + s[k]
	}

	str = priceLingJiaoLingFenRegex.ReplaceAllString(str, "整")
	str = priceLingJiaoRegex.ReplaceAllString(str, "零")
	str = priceLingFenRegex.ReplaceAllString(str, "整")
	str = priceLingQianBaiShiRegex.ReplaceAllString(str, "零")
	str = priceLing2Regex.ReplaceAllString(str, "零")
	str = priceLingYiRegex.ReplaceAllString(str, "亿")
	str = priceLingWanRegex.ReplaceAllString(str, "万")
	str = priceLingMultiYuanRegex.ReplaceAllString(str, "元")
	str = priceYiLingWanRegex.ReplaceAllString(str, "^元")
	str = priceLingYuanRegex.ReplaceAllString(str, "零")

	if str == "整" {
		str = "零元整"
	}
	return str
}

// PriceCheckPriceStr 价格检查
// 检查方案：
// 	1.不能为空
// 	2.是float类型
// 	3.大于0
// 	4. 不能超过2位小数
func PriceCheckPriceStr(priceStr string) (ok bool, errMsg string) {
	if priceStr == "" {
		return false, "空字符串"
	}

	priceFloat, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return false, "非数值"
	}

	if priceFloat <= 0 {
		return false, "不能小于0"
	}

	splits := strings.Split(priceStr, ".")
	if len(splits) == 2 && len(splits[1]) > 2 {
		return false, "超过了2位小数"
	}

	return true, ""
}

// PriceCommaf 千分位显示
// e.g. Commaf(834142.32) -> 834,142.32
func PriceCommaf(v float64) string {
	buf := &bytes.Buffer{}
	if v < 0 {
		buf.Write([]byte{'-'})
		v = 0 - v
	}

	comma := []byte{','}

	parts := strings.Split(strconv.FormatFloat(v, 'f', -1, 64), ".")
	pos := 0
	if len(parts[0])%3 != 0 {
		pos += len(parts[0]) % 3
		buf.WriteString(parts[0][:pos])
		buf.Write(comma)
	}
	for ; pos < len(parts[0]); pos += 3 {
		buf.WriteString(parts[0][pos : pos+3])
		buf.Write(comma)
	}
	buf.Truncate(buf.Len() - 1)

	if len(parts) > 1 {
		buf.Write([]byte{'.'})
		buf.WriteString(parts[1])
	}
	return buf.String()
}
