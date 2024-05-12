/*
* Created on 09 May 2024
* @author Sai Sumanth
 */
package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/sai7xp/xbank/utils"
)

// Custom Validator
var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		// check currency is supported
		return utils.IsSupportedCurrency(currency)
	}
	return false
}
