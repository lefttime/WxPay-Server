package controllers

import (
  "fmt"
  "time"
  "bytes"
  "strings"
  "strconv"
  "net/http"
  "io/ioutil"
  "math/rand"
  "encoding/xml"
  "encoding/json"
  "github.com/astaxie/beego"
  "github.com/lefttime/PaymentSystem/models"
)

type WxPay struct {
  beego.Controller
}

func (c *WxPay) RequestUnifiedOrder() {
  type DataForm struct {
    PlayerId  string `json:"playerId"`
    Area      string `json:"area"`
    ProductId string `json:"cid"`
    Title     string `json:"title"`
    IP        string `json:"ipaddr"`
  }

  type DataResponse struct {
    Status int    `json:"return_code"`
    Data   string `json:"data"`
  }

  dataRes := DataResponse{ Status: -1 }

  var dataForm DataForm
  c.Ctx.Request.ParseForm()
  formValue := c.Ctx.Request.FormValue( "orderInfo" )
  json.Unmarshal( []byte(formValue), &dataForm )

  payUrl := requestUnifiedOrder( dataForm.ProductId, dataForm.PlayerId, dataForm.IP )
  if payUrl != "" {
    dataRes.Status = 0
    dataRes.Data   = payUrl
  }

  c.Ctx.Output.JSON( dataRes, false, false )
}

func (c *WxPay) ResponseOrderResult() {
  handleOrderResult( c )
}

func requestUnifiedOrder( pid string, uid string, ipAddr string ) string {
  product := GetProductByInnerId( pid )
  if product.ID==""  {
    return ""
  }
  player := GetPlayerById( uid )

  var dataReq models.UnifyOrderReq
  dataReq.Appid          = beego.AppConfig.String( "wxAppId"   )
  dataReq.MchId          = beego.AppConfig.String( "partnerId" )
  dataReq.DeviceInfo     = "WEB"
  dataReq.NonceStr       = getRandomString( 32 )
  dataReq.SignType       = "MD5"
  dataReq.Body           = product.Name
  dataReq.TotalFee       = product.Amount * 100
  dataReq.SpbillCreateIp = ipAddr
  dataReq.NotifyUrl      = beego.AppConfig.String( "baseUrl" ) + "/pay/h5/notify"
  dataReq.TradeType      = "MWEB"
  dataReq.OutTradeNo     = strconv.FormatInt(time.Now().Unix(), 10) + getRandomString( 3 )
  dataReq.OpenId         = player.OpenId
  dataReq.SceneInfo      = `{"h5_info":{"type":"IOS","app_name":"摩语卡五星","bundle_id":"com.xc.xfkwx"}}`
  dataReq.Sign           = dataReq.GenerateMD5()

  xmlStr := dataReq.ToXml()
  bytes_req := []byte( xmlStr )
  req, err := http.NewRequest( "POST", "https://api.mch.weixin.qq.com/pay/unifiedorder", bytes.NewReader( bytes_req ) )
  if err != nil {
    fmt.Println( "New Http Request发生错误，原因:", err )
    return ""
  }

  req.Header.Set( "Accept", "application/xml" )
  req.Header.Set( "Content-Type", "application/xml;charset=utf-8" )

  client := http.Client{}
  resp, _err := client.Do( req )
  if _err != nil {
    fmt.Println( "请求微信支付统一下单接口发送错误, 原因:", _err )
    return ""
  }

  respBytes, err := ioutil.ReadAll( resp.Body )
  if err != nil {
    fmt.Println( "解析返回body错误", err )
    return ""
  }

  xmlResp := models.UnifyOrderResp{}
  _err = xml.Unmarshal( respBytes, &xmlResp )
  if xmlResp.ReturnCode=="FAIL" {
    fmt.Println( "微信支付统一下单不成功，原因:", xmlResp.ReturnMsg, " str_req-->", xmlStr )
    return ""
  }

  var order models.Order
  order.Appid            = beego.AppConfig.String( "wxAppId" )
  order.UserId           = player.ID
  order.Quantity         = 1
  order.PrePurchaseDate  = strconv.FormatInt( time.Now().Unix() * 1000, 10 )
  order.PreTransactionId = xmlResp.PrepayId
  order.ProductId        = product.ID
  order.ProductCount     = product.Count
  order.ProductAmount    = product.Amount * 100
  order.TransactionId    = dataReq.OutTradeNo
  order.PurchaseType     = "H5"
  order.IP               = ipAddr
  if CreateOrder( order )==false {
    return ""
  }

  return xmlResp.MwebUrl
}

func requestOrderQuery() {
  var dataReq models.QueryOrderReq
  dataReq.Appid            = beego.AppConfig.String( "wxAppId"   )
  dataReq.MchId            = beego.AppConfig.String( "partnerId" )
  dataReq.TransactionId    = ""
  dataReq.OutTradeNo       = strconv.FormatInt(time.Now().Unix(), 10) + getRandomString( 3 ) // 订单号
  dataReq.NonceStr         = getRandomString( 32 )
  dataReq.Sign             = dataReq.GenerateMD5()

  params := make( map[string]interface{}, 0 )
  params[ "appid" ]            = dataReq.Appid
  params[ "mch_id" ]           = dataReq.MchId
  params[ "transaction_id"]    = dataReq.TransactionId
  params[ "out_trade_no" ]     = dataReq.OutTradeNo
  params[ "nonce_str" ]        = dataReq.NonceStr

  bytes_req, err := xml.Marshal( dataReq )
  if err != nil {
    fmt.Println( "转换为xml错误:", err )
  }

  str_req := strings.Replace( string(bytes_req), "UnifyOrderReq", "xml", -1 )
  bytes_req = []byte(str_req)

  http.NewRequest( "POST", "https://api.mch.weixin.qq.com/pay/orderquery", bytes.NewReader( bytes_req ) )
}

func handleOrderResult( c *WxPay ) {
  body, err := ioutil.ReadAll( c.Ctx.Request.Body )
  if err != nil {
    return
  }
  defer c.Ctx.Request.Body.Close()

  var dataReq models.OrderResult
  err = xml.Unmarshal( body, &dataReq )
  if err != nil {
    http.Error( c.Ctx.ResponseWriter.ResponseWriter.( http.ResponseWriter ), http.StatusText( http.StatusBadRequest ), http.StatusBadRequest )
    return
  }

  var dataResp models.OrderResultResp
  sign := dataReq.GenerateMD5()
  if sign==dataReq.Sign {
    dataResp.ReturnCode = "SUCCESS"
    dataResp.ReturnMsg  = "OK"

    if dataReq.ReturnCode=="SUCCESS" {
      order := FindOrderBy( dataReq.OutTradeNo )
      if len(order.ID) > 0 {
        order.TransactionId = dataReq.TransactionId
        order.PurchaseDate  = strconv.FormatInt( time.Now().Unix() * 1000, 10 )
        UpdateOrder( order )
        notifyServers( order, dataReq )
      }
    }
  } else {
    dataResp.ReturnCode = "FAIL"
    dataResp.ReturnMsg  = "Failed to verify sign, please retry!"
  }

  bytes, _err := xml.Marshal( dataResp )
  str_resp    := strings.Replace( string( bytes ), "OrderResultResp", "xml", -1 )

  if _err != nil {
    http.Error( c.Ctx.ResponseWriter.ResponseWriter.( http.ResponseWriter ), http.StatusText( http.StatusBadRequest ), http.StatusBadRequest )
    return
  }

  c.Ctx.ResponseWriter.ResponseWriter.(http.ResponseWriter).WriteHeader( http.StatusOK )
  fmt.Fprint( c.Ctx.ResponseWriter.ResponseWriter.(http.ResponseWriter), str_resp )
  c.ServeXML()
}

func notifyServers( order models.Order, result models.OrderResult ) {
  notifyGameServer( order, result )
  notifyAgencyServer( order )
}

func notifyAgencyServer( order models.Order ) {
  postData := "orderId=" + order.TransactionId
  postData += "&userId=" + order.UserId
  postData += "&money="  + strconv.Itoa( order.ProductAmount )

  fmt.Println( "http://opxy.moy2017.com/gameRecharge?" + postData )
  req, err := http.Get( "http://opxy.moy2017.com/gameRecharge?" + postData )
  if err != nil {
    fmt.Println( "发生错误，原因: ", err )
    return
  }

  defer req.Body.Close()
  body, err := ioutil.ReadAll(req.Body)
  fmt.Println( string(body) )
}

func notifyGameServer( order models.Order, result models.OrderResult ) {
  host     := beego.AppConfig.String( "Game::Host" )
  port     := beego.AppConfig.String( "Game::Port" )
  serverId := "10003"
  admin    := "h5Pay"
  openId   := result.OpenId
  cardType := "2"
  cardNum  := strconv.Itoa( order.ProductCount )
  operType := "1"

  url := "http://{host}:{port}/cgi-bin/admin?serverid={serverId}&msg=charge&admin={admin}&openId={openId}&cardType={cardType}&cardNum={cardNum}&operType={operType}"
  url = strings.Replace( url, "{host}",      host,     -1 )
  url = strings.Replace( url, "{port}",      port,     -1 )
  url = strings.Replace( url, "{serverId}",  serverId, -1 )
  url = strings.Replace( url, "{admin}",     admin,    -1 )
  url = strings.Replace( url, "{openId}",    openId,   -1 )
  url = strings.Replace( url, "{cardType}",  cardType, -1 )
  url = strings.Replace( url, "{cardNum}",   cardNum,  -1 )
  url = strings.Replace( url, "{operType}",  operType, -1 )

  req, err := http.NewRequest( "GET", url, nil )
  if err != nil {
    fmt.Println( "New Http Request发生错误，原因:", err )
    return
  }

  req.Header.Set( "Accept", "application/json" )
  req.Header.Set( "Content-Type", "application/json;charset=utf-8" )

  client := http.Client{}
  resp, _err := client.Do( req )
  if _err != nil {
    fmt.Println( "请求发送错误, 原因:", _err )
    return
  }

  respBytes, err := ioutil.ReadAll( resp.Body )
  if err != nil {
    fmt.Println( "解析返回body错误", err )
    return
  }

  type DataResp struct {
    ErrCode int    `json:"errorCode"`
    ErrMsg  string `json:"errorMsg"`
  }

  var feedback DataResp
  if json.Unmarshal( respBytes, &feedback )==nil && feedback.ErrCode==0 {
    fmt.Println( "充值成功" )
  }
}

func getRandomString( total int) string {
  str    := "02468acehjloqsvxz12123sdf"
  bytes  := []byte( str )
  result := []byte{}
  rdn    := rand.New( rand.NewSource( time.Now().UnixNano() ) )
  for idx := 0; idx < total; idx++ {
    result = append( result, bytes[rdn.Intn( len( bytes ) )] )
  }
  return string( result )
}
