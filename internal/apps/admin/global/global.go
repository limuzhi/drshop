package global

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var (
	GVA_Validate *validator.Validate
	GVA_Trans    ut.Translator
)
