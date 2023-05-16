package helpers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"alsritter.icu/rabbit-template/internal/pkg/strtools"
	"alsritter.icu/rabbit-template/internal/pkg/timetools"
)

type CardType int32

const (
	// 证件类型: 身份证
	CardTypeIDCard CardType = 0
	// 证件类型: 护照
	CardTypePassport CardType = 1
	// 证件类型: 港澳居民来往内地通行证
	CardTypeGAToMainlandPass CardType = 2
	// 证件类型: 台胞证
	CardTypeTWToMainlandPass CardType = 3
	// 证件类型: 其他
	CardTypeOther CardType = 255
)

var (
	regionMap = map[int]string{
		11: "北京",
		12: "天津",
		13: "河北",
		14: "山西",
		15: "内蒙古",
		21: "辽宁",
		22: "吉林",
		23: "黑龙江",
		31: "上海",
		32: "江苏",
		33: "浙江",
		34: "安徽",
		35: "福建",
		36: "江西",
		37: "山东",
		41: "河南",
		42: "湖北",
		43: "湖南",
		44: "广东",
		45: "广西",
		46: "海南",
		50: "重庆",
		51: "四川",
		52: "贵州",
		53: "云南",
		54: "西藏",
		61: "陕西",
		62: "甘肃",
		63: "青海",
		64: "宁夏",
		65: "新疆",
		71: "台湾",
		81: "香港",
		82: "澳门",
		83: "台湾",
		91: "国外",
	}
	idCardVerifyFactory = []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	idCardVerifyNumber  = []uint8{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}
)

// CardValidator 证件校验器接口
type CardValidator interface {
	IsValid() bool
}

func NewCardValidator(cardType CardType, cardNo string) CardValidator {
	switch cardType {
	case CardTypeIDCard:
		return idCardValidator{cardNo: cardNo}
	case CardTypePassport:
		return passportValidator{}
	case CardTypeGAToMainlandPass:
		return gAMenToMainlandPassportValidator{cardNo: cardNo}
	case CardTypeTWToMainlandPass:
		return tWToMainlandPassportValidator{cardNo: cardNo}
	case CardTypeOther:
		return otherCardValidator{}
	default:
		return invalidCardTypeValidator{}
	}
}

// 身份证验证器
type idCardValidator struct {
	cardNo string
}

func (i idCardValidator) isRegionValid() bool {
	regionCode, err := strconv.Atoi(i.cardNo[:2])
	if err != nil {
		return false
	}
	_, ok := regionMap[regionCode]
	return ok
}

func (i idCardValidator) isDateValid() bool {
	var err error
	if len(i.cardNo) == 15 {
		_, err = timetools.TimeParseString("19"+i.cardNo[6:12], timetools.TimeFormatLayout_YmD)
	} else {
		_, err = timetools.TimeParseString(i.cardNo[6:14], timetools.TimeFormatLayout_YmD)
	}
	return err == nil
}

func (i idCardValidator) isVerified() bool {
	return i.cardNo[17] == getVerifyBit(i.cardNo[:17])
}

func (i idCardValidator) IsValid() bool {
	// 为空
	if i.cardNo == "" {
		return false
	}

	// 包含中文
	if strtools.IsIncludeChinese(i.cardNo) {
		return true
	}

	// 非15位或者18位不校验
	if len(i.cardNo) != 15 && len(i.cardNo) != 18 {
		return true
	}

	// 长度以及格式校验
	isMatch18, _ := regexp.MatchString("^\\d{17}(\\d|X|x)$", i.cardNo)
	isMatch15, _ := regexp.MatchString("^\\d{15}$", i.cardNo)
	if !isMatch18 && !isMatch15 {
		return false
	}

	// 地区校验
	if !i.isRegionValid() {
		return false
	}

	// 日期校验
	if !i.isDateValid() {
		return false
	}

	// 15位转18位
	if len(i.cardNo) == 15 {
		i.cardNo = TransTo18CardNo(i.cardNo)
	}

	// 全转大写
	i.cardNo = strings.ToUpper(i.cardNo)

	// 验证18位校验码
	return i.isVerified()

}

// 护照验证器
type passportValidator struct {
}

func (p passportValidator) IsValid() bool {
	return true
}

type gAMenToMainlandPassportValidator struct {
	cardNo string
}

func (g gAMenToMainlandPassportValidator) IsValid() bool {
	if len(g.cardNo) != 11 {
		return true
	}
	match, _ := regexp.MatchString("^[H|h|M|m]\\d{10}$", g.cardNo)
	return match
}

type tWToMainlandPassportValidator struct {
	cardNo string
}

func (t tWToMainlandPassportValidator) IsValid() bool {
	if len(t.cardNo) != 8 {
		return true
	}
	match, _ := regexp.MatchString("^\\d{8}$", t.cardNo)
	return match
}

type otherCardValidator struct {
}

func (o otherCardValidator) IsValid() bool {
	return true
}

type invalidCardTypeValidator struct {
}

func (i invalidCardTypeValidator) IsValid() bool {
	return false
}

func IsValidCardType(cardType CardType) bool {
	return cardType == CardTypeIDCard || cardType == CardTypePassport || cardType == CardTypeGAToMainlandPass || cardType == CardTypeTWToMainlandPass || cardType == CardTypeOther
}

func TransTo18CardNo(cardNo string) string {
	cardNo18 := cardNo[:6] + "19" + cardNo[6:]
	cardNo18 = string(append([]uint8(cardNo18), getVerifyBit(cardNo18)))
	return cardNo18
}

func getVerifyBit(idCardBase string) uint8 {
	checkSum := 0
	for index, c := range idCardBase {
		checkSum += int(c-'0') * idCardVerifyFactory[index]
	}
	return idCardVerifyNumber[checkSum%11]
}

const (
	GenderUnknown = 0 // 未知性别
	GenderMan     = 1 // 男人
	GenderWoman   = 2 // 女人
)

const (
	AgeUnknown = -1 // 未知年龄
)

// IdCard 证件
type IdCard struct {
	CardType CardType // 证件类型 身份证:0 护照:1 港澳居民来往内地通行证:2 台胞证:3 其他:255
	IdCardNo string   // 证件号码
}

func NewIdCardCalc(cardType CardType, idCardNo string) *IdCard {
	return &IdCard{
		CardType: cardType,
		IdCardNo: idCardNo,
	}
}

// GetAge 获取年龄
func (s *IdCard) GetAge() int32 {
	if s.CardType == CardTypeIDCard {
		idCardNo := s.get18IdCard()
		if idCardNo == "" {
			return AgeUnknown
		}
		birthday, err := timetools.TimeParseString(idCardNo[6:14], timetools.TimeFormatLayout_YmD)
		if err != nil {
			return AgeUnknown
		}
		curTime := timetools.TimeCurrentTimePointer()
		age := curTime.Year() - birthday.Year()
		if curTime.Month() < birthday.Month() || (curTime.Month() == birthday.Month() && curTime.Day() > birthday.Day()) {
			age--
		}
		if age < 0 {
			return AgeUnknown
		}
		return int32(age)
	}
	return AgeUnknown
}

// GetGender 获取性别
func (s *IdCard) GetGender() int32 {
	if s.CardType == CardTypeIDCard {
		idCardNo := s.get18IdCard()
		if idCardNo == "" {
			return GenderUnknown
		}
		sexBit := idCardNo[16] - '0'
		if sexBit&1 == 1 {
			return GenderMan
		}
		return GenderWoman
	}
	return GenderUnknown
}

// GetBirthday 获取生日日期 返回 2022-01-01
func (s *IdCard) GetBirthday() string {
	if s.CardType == CardTypeIDCard {
		idCardNo := s.get18IdCard()
		if idCardNo == "" {
			return ""
		}
		return fmt.Sprintf("%s-%s-%s", idCardNo[6:10], idCardNo[10:12], idCardNo[12:14])
	}
	return ""
}

func (s *IdCard) get18IdCard() string {
	switch len(s.IdCardNo) {
	case 15:
		return TransTo18CardNo(s.IdCardNo)
	case 18:
		return s.IdCardNo
	default:
		return ""
	}
}
