package api

import (
	"log"

	"github.com/ShenPingYuan/go-webdemo/util"
	"github.com/go-playground/validator/v10"
)

//var validate *validator.Validate

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	//打印下参数
	log.Printf("fieldLevel.Field().Interface():%v", fieldLevel.Field().Interface())
	//打印下参数类型
	log.Printf("fieldLevel.Field().Type():%v", fieldLevel.Field().Type())
	currency, ok := fieldLevel.Field().Interface().(string)
	if ok {
		return util.IsValidCurrency(currency)
	}
	return false
}
