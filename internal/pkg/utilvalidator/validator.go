package utilvalidator

import (
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	ch_translations "github.com/go-playground/validator/v10/translations/zh"
	"regexp"
)

// 初始化Validator数据校验
func InitValidate() (*validator.Validate, ut.Translator) {
	chinese := zh.New()
	uni := ut.New(chinese, chinese)
	trans, _ := uni.GetTranslator("zh")
	validate := validator.New()
	_ = ch_translations.RegisterDefaultTranslations(validate, trans)
	_ = validate.RegisterValidation("checkMobile", checkMobile)
	return validate, trans
}

func checkMobile(fl validator.FieldLevel) bool {
	reg := `^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(fl.Field().String())
}
