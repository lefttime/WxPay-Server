package controllers

import (
  "fmt"
  "strings"
  "net/http"
  "path/filepath"
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
  orpath := c.Ctx.Request.URL.Path
  fmt.Println( orpath )
  if strings.Index( orpath, "MP_verify_BGGJT9p7dpvr3s6Y.txt" ) >= 0 {
    path := filepath.Join( `static`, "MP_verify_BGGJT9p7dpvr3s6Y.txt" )
    http.ServeFile( c.Ctx.ResponseWriter, c.Ctx.Request, path )
  } else if strings.Index( orpath, "result.html" ) >= 0 {
    c.TplName = "result.tpl"
  } else {
    c.TplName = "pay.tpl"
  }
}
