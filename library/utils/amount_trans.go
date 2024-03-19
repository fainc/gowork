package utils

import (
	"github.com/shopspring/decimal"
)

type amountTrans struct{}

var amountTransVar = amountTrans{}

func AmountTrans() *amountTrans {
	return &amountTransVar
}

// Int642Yuan int64数据转元，p参数控制小数截断，例：1 / 3 = 0.33333... ，p截断2位，返回0.33
func (rec *amountTrans) Int642Yuan(x int64, u float64, p int32) (ret float64, str string) {
	ret, _ = decimal.NewFromInt(x).Div(decimal.NewFromFloat(u)).Truncate(p).Float64()
	str = decimal.NewFromFloat(ret).StringFixed(p)
	return
}

// Yuan2Int64 元转int64，只保留转换存储单位后整数部分，小数部分舍弃，例：2.9999 * 100 以分单位存储，返回299
func (rec *amountTrans) Yuan2Int64(x float64, u float64) (ret int64) {
	ret = decimal.NewFromFloat(x).Mul(decimal.NewFromFloat(u)).IntPart()
	return
}

// Yuan2Int 元转int，只保留转换存储单位后整数部分，小数部分舍弃，例：2.9999 * 100 以分单位存储，返回299
func (rec *amountTrans) Yuan2Int(x float64, u float64) (ret int) {
	ret = int(rec.Yuan2Int64(x, u))
	return
}
