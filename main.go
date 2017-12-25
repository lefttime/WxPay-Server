package main

import (
	_ "github.com/lefttime/PaymentSystem/routers"
	"github.com/astaxie/beego"
  // "github.com/lefttime/PaymentSystem/controllers"
)

func main() {
	beego.Run()
  // controllers.WxPayRequest()
}

