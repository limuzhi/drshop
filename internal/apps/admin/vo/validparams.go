package vo

import (
	"drpshop/internal/apps/admin/global"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

func DefalutGetValidParams(c *gin.Context, params interface{}) error {
	if err := c.ShouldBind(params); err != nil {
		return err
	}
	log.Println("params==", params)
	if err := global.GVA_Validate.Struct(params); err != nil {
		errs := err.(validator.ValidationErrors)
		sliceErrs := []string{}
		for _, e := range errs {
			sliceErrs = append(sliceErrs, e.Translate(global.GVA_Trans))
		}
		return errors.New(strings.Join(sliceErrs, ","))
	}
	return nil
}
