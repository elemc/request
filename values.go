package request

import (
	"fmt"
	"net"
	"strconv"
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

// MustMAC функция проверит и вернет MAC адрес в правильном виде или возвратит пустую строку
func (rv Value) MustMAC() string {
	hw, err := net.ParseMAC(rv.String())
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%02x-%02x-%02x-%02x-%02x-%02x", hw[0], hw[1], hw[2], hw[3], hw[4], hw[5])
}
