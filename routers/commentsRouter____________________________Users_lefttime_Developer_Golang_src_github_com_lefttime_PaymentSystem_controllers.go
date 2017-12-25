package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/lefttime/PaymentSystem/controllers:WxPay"] = append(beego.GlobalControllerRouter["github.com/lefttime/PaymentSystem/controllers:WxPay"],
		beego.ControllerComments{
			Method: "RequestUnifiedOrder",
			Router: `/h5/create`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

}
