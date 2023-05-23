package resource

//go:generate fyne bundle -o bundled.go ./assets/alipay.jpg
//go:generate fyne bundle -o bundled.go -append ./assets/wxpay.jpg

var (
	ResourceWxPay  = resourceWxpayJpg
	ResourceAliPay = resourceAlipayJpg
)
