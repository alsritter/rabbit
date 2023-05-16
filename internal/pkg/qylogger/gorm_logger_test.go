package qylogger

import (
	"testing"

	"go.elastic.co/apm/module/apmsql"
)

func TestLogger_Trace(t *testing.T) {
	t.Log(apmsql.QuerySignature("UPDATE `goods` SET `is_on_sale`=0 WHERE store_id IN (NULL)"))
}
