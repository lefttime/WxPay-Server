package routers

import (
  "github.com/lefttime/PaymentSystem/controllers"
	"github.com/astaxie/beego"
)

func init() {
  beego.Router( "/pay/h5/create",        &controllers.WxPay{},      "post:RequestUnifiedOrder" )
  beego.Router( "/pay/h5/notify",        &controllers.WxPay{},      "post:ResponseOrderResult" )
  beego.Router( "/pay/transfer/request", &controllers.WxTransfer{}, "post:RequestTransfer"     )
  beego.Router( "/pay/transfer/resend",  &controllers.WxTransfer{}, "post:ResendTransfer"      )
  beego.Router( "/pay/transfer/query",   &controllers.WxTransfer{}, "post:QueryTransfer"       )

  beego.Router( "/pay/playerInfo",  &controllers.Player{}         )
  beego.Router( "/pay/productInfo", &controllers.Product{}        )
  beego.Router( "/*",               &controllers.MainController{} )
}
