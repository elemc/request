package request

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

// Value - тип строки, с возможностью его преобразования в иные значения
type Value string

// MustInt - функция для преобразования строки значения в целое число
func (rv Value) MustInt() int {
	res, err := strconv.Atoi(string(rv))
	if err != nil {
		return 0
	}
	return res
}

// MustUint64 - функция для преобразования строки значения в целое число
func (rv Value) MustUint64() uint64 {
	res, err := strconv.ParseUint(string(rv), 10, 64)
	if err != nil {
		res = 0
	}
	return res
}

// MustInt64 - функция для преобразования строки значения в целое число
func (rv Value) MustInt64() int64 {
	res, err := strconv.ParseInt(string(rv), 10, 64)
	if err != nil {
		res = 0
	}
	return res
}

// MustBool - функция для преобразования строки значения в булево
func (rv Value) MustBool() bool {
	res, err := strconv.ParseBool(string(rv))
	if err != nil {
		return false
	}
	return res
}

// String - функция возвращает строковое значение
func (rv Value) String() string {
	return string(rv)
}

// MustMAC - функция проверит и вернет MAC адрес в правильном виде или возвратит пустую строку
func (rv Value) MustMAC() string {
	hw, err := net.ParseMAC(rv.String())
	if err != nil {
		return ""
	}
	itsNullMAC := true
	for _, b := range hw {
		if b != 0 {
			itsNullMAC = false
		}
	}
	if itsNullMAC {
		return ""
	}
	return fmt.Sprintf("%02x-%02x-%02x-%02x-%02x-%02x", hw[0], hw[1], hw[2], hw[3], hw[4], hw[5])
}

// MustTime - функция преобразует значение строки к значению time.Time
func (rv Value) MustTime() time.Time {
	var (
		result time.Time
		err    error
	)

	// пробуем RFC3389Nano
	if result, err = time.Parse(time.RFC3339Nano, string(rv)); err == nil {
		return result
	}

	// пробуем RFC3389
	if result, err = time.Parse(time.RFC3339, string(rv)); err == nil {
		return result
	}

	// пробуем обратный
	if result, err = time.Parse("2006-01-02", string(rv)); err == nil {
		return result
	}

	// пробуем преобразовать в int и интепретируем, как timestamp
	if value, err := strconv.ParseInt(string(rv), 10, 64); err == nil && value != 0 {
		return time.Unix(value, 0)
	}

	return time.Time{}
}
