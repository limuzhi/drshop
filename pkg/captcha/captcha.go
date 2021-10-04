package captcha

import (
	"github.com/google/uuid"
	"github.com/mojocn/base64Captcha"
)

//configJsonBody json request body.
type configJsonBody struct {
	Id            string
	CaptchaType   string
	VerifyValue   string
	DriverAudio   *base64Captcha.DriverAudio
	DriverString  *base64Captcha.DriverString
	DriverChinese *base64Captcha.DriverChinese
	DriverMath    *base64Captcha.DriverMath
	DriverDigit   *base64Captcha.DriverDigit
}


func DriverStringFunc(store base64Captcha.Store, height, width, length int) (id, b64s string, err error) {
	e := configJsonBody{}
	e.Id = uuid.New().String()
	e.DriverString = &base64Captcha.DriverString{
		Height:          height,
		Width:           width,
		NoiseCount:      50,
		ShowLineOptions: 20,
		Length:          length,
		Source:          "abcdefghjkmnpqrstuvwxyz23456789",
		Fonts:           []string{"chromohv.ttf"},
	}
	driver := e.DriverString.ConvertFonts()
	cap := base64Captcha.NewCaptcha(driver, store)
	return cap.Generate()
}

func DriverDigitFunc(store base64Captcha.Store, height, width, length int) (id, b64s string, err error) {
	e := configJsonBody{}
	e.Id = uuid.New().String()
	//h 80 w 240 l 6
	e.DriverDigit = base64Captcha.NewDriverDigit(height, width, length, 0.7, 80)
	driver := e.DriverDigit
	cap := base64Captcha.NewCaptcha(driver, store)
	return cap.Generate()
}
