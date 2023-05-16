package wechat

import (
	"io"
	"os"
	"testing"

	"github.com/ArtisanCloud/PowerWeChat/v2/src/kernel/power"
)

func TestNewMiniProgram(t *testing.T) {
	mp, _ := NewMiniProgram("wx37a8e11656f222c3", "130aae9ec5d2d5a9fc3382c086b7c13a")

	resp, err := mp.WXACode.GetUnlimited(
		"goods_id=2",
		"pages/goods/goods-detail/index",
		430,
		false,
		&power.HashMap{"r": 0, "g": 0, "b": 0},
		false)
	if err != nil {
		panic(err)
	}

	bd, _ := io.ReadAll(resp.Body)
	os.WriteFile("./output2.png", bd, 0666)
}
