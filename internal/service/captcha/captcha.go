package captcha

import (
	"github.com/mojocn/base64Captcha"
	"image/color"
)

var driverString base64Captcha.DriverMath
var store = base64Captcha.DefaultMemStore

// 配置验证码信息
var captchaConfig = base64Captcha.DriverMath{
	Height:          60,
	Width:           200,
	NoiseCount:      0,
	ShowLineOptions: 2 | 4,
	//Length:          4,
	//Source:          "1234567890qwertyuioplkjhgfdsazxcvbnm",
	BgColor: &color.RGBA{
		R: 0,
		G: 85,
		B: 131,
		A: 1,
	},
	Fonts: []string{"wqy-microhei.ttc"},
	//Fonts: []string{"3Dumb.ttf"},
}

// Genarate
//
//	@Description: 生成图形验证码
//	@return id		图形验证码id
//	@return base64Image		base64格式图形
//	@return answer	答案
//	@return err				错误信息
func CaptchaGenarate() (id string, base64Image string, answer string, err error) {
	driverString = captchaConfig
	c := base64Captcha.NewCaptcha(driverString.ConvertFonts(), store)
	id, base64Image, answer, err = c.Generate()
	if err != nil {
		return id, base64Image, answer, err
	}
	return id, base64Image, answer, nil
}

// Verify
//
//	@Description: 	验证图形验证码
//	@param id		验证码id
//	@param capt		验证码
//	@return bool	验证结果，true：验证通过，false:验证失败
func CaptchaVerify(id string, capt string) bool {
	return store.Verify(id, capt, true)
}
