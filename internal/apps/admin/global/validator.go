package global

import "drpshop/internal/pkg/utilvalidator"

func InitValidate() {
	GVA_Validate, GVA_Trans = utilvalidator.InitValidate()
}

