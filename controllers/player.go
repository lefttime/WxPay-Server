package controllers

import (
  "fmt"
  "strings"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "github.com/astaxie/beego"
  "github.com/lefttime/PaymentSystem/models"
)

type Player struct {
  beego.Controller
}

func (c *Player) Get() {
  area := c.Ctx.Input.Query( "area" )
  uid  := c.Ctx.Input.Query( "uid"  )

  type DataResponse struct {
    Status int    `json:"status"`
    Data   string `json:"data"`
  }

  dataResp := DataResponse{ Status: -1 }
  if area=="xiangyang" && len( uid ) > 0 {
    info := GetPlayerById( uid )
    if len( info.OpenId ) > 0 {
      jsonData, err := json.Marshal( info )
      if err==nil {
        dataResp.Status = 0
        dataResp.Data   = string( jsonData )
      }
    }
  }
  c.Ctx.Output.JSON( dataResp, false, false )
}

func GetPlayerById( uid string ) models.Player {
  result := models.Player{}

  host := beego.AppConfig.String( "Game::Host" )
  port := beego.AppConfig.String( "Game::Port" )
  url := "http://{host}:{port}/cgi-bin/admin?serverid=3&msg=userinfo&userid={userid}"
  url = strings.Replace( url, "{host}",   host, -1 )
  url = strings.Replace( url, "{port}",   port, -1 )
  url = strings.Replace( url, "{userid}", uid,  -1 )

  req, err := http.NewRequest( "GET", url, nil )
  if err != nil {
    fmt.Println( "New Http Request发生错误，原因:", err )
    return result
  }

  req.Header.Set( "Accept", "application/json" )
  req.Header.Set( "Content-Type", "application/json;charset=utf-8" )

  client := http.Client{}
  resp, _err := client.Do( req )
  if _err != nil {
    fmt.Println( "请求发送错误, 原因:", _err )
    return result
  }

  respBytes, err := ioutil.ReadAll( resp.Body )
  if err != nil {
    fmt.Println( "解析返回body错误", err )
    return result
  }

  if json.Unmarshal( respBytes, &result ) != nil {
    fmt.Println( result )
  }

  return result
}
