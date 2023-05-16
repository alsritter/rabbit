package strtools

import (
	"encoding/json"
	"fmt"
)

type StringOrNumber string

func (son *StringOrNumber) UnmarshalJSON(data []byte) error {
	var num float64
	err := json.Unmarshal(data, &num)
	if err == nil {
		*son = StringOrNumber(fmt.Sprintf("%.0f", num))
		return nil
	}
	var str string
	err = json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	*son = StringOrNumber(str)
	return nil
}
