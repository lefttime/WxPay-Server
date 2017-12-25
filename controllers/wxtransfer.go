package controllers

import (
  "fmt"
  "time"
  "bytes"
  "errors"
  "strconv"
  "net/http"
  "io/ioutil"
  "crypto/tls"
  "crypto/x509"
  "encoding/xml"
  "encoding/json"
  "path/filepath"
  "github.com/astaxie/beego"
  "github.com/lefttime/PaymentSystem/utils"
  "github.com/lefttime/PaymentSystem/models"
)

type WxTransfer struct {
  beego.Controller
}

func (c *WxTransfer) RequestTransfer() {
  type DataForm struct {
    PlayerId string `json:"playerid"`
    Realname string `json:"realname"`
    Amount   int    `json:"amount"`
  }

  var dataForm DataForm
  c.Ctx.Request.ParseForm()
  formValue := c.Ctx.Request.FormValue( "transInfo" )
  json.Unmarshal( []byte(formValue), &dataForm )

  type DataResponse struct {
    Status    int    `json:"status"`
    PaymentNo string `json:"payment_no"`
    ErrMsg    string `json:"err_msg"`
  }

  dataRes  := DataResponse{ Status: -1 }
  userInfo := GetPlayerById( dataForm.PlayerId )

  if len(userInfo.ID) < 1 {
    dataRes.ErrMsg = "不存在该用户"
  } else if len(dataForm.Realname) < 1 {
    dataRes.ErrMsg = "未填写用户名"
  } else {
    err := checkRequestValid( dataForm.PlayerId, dataForm.Realname, dataForm.Amount )
    if err != nil {
      dataRes.ErrMsg = err.Error()
    } else {
      transfer := generateTransfer( userInfo.OpenId, dataForm.Realname, dataForm.Amount )
      if len(transfer.Appid)==0 {
        dataRes.ErrMsg = "创建数据模型失败"
      } else {
        paymentNo, err := requestTransfer( transfer )
        if err != nil {
          dataRes.ErrMsg = err.Error()
        } else {
          dataRes.Status    = 0
          dataRes.PaymentNo = paymentNo
        }
      }
    }
  }

  c.Ctx.Output.JSON( dataRes, false, false )
}

func (c *WxTransfer) QueryTransfer() {
  type DataForm struct {
    PartnerTradeNo string `json:"partner_trade_no"`
  }

  var dataForm DataForm
  c.Ctx.Request.ParseForm()
  formValue := c.Ctx.Request.FormValue( "transInfo" )
  json.Unmarshal( []byte(formValue), &dataForm )

  type DataResponse struct {
    Status    int    `json:"status"`
    ErrMsg    string `json:"err_msg"`
  }

  dataRes  := DataResponse{ Status: -1 }

  transfer := FindTransferByTradeNo( dataForm.PartnerTradeNo )
  if transfer.ID==0 {
    dataRes.ErrMsg = "不存在该订单"
  } else {
    status, err := queryTransfer( transfer )
    dataRes.Status = status
    if err != nil {
      dataRes.ErrMsg = err.Error()
    }
  }

  fmt.Println( "QueryTransfer -------> ", dataRes )

  c.Ctx.Output.JSON( dataRes, false, false )
}

func (c *WxTransfer) ResendTransfer() {
  type DataForm struct {
    PartnerTradeNo string `json:"partner_trade_no"`
  }

  var dataForm DataForm
  c.Ctx.Request.ParseForm()
  formValue := c.Ctx.Request.FormValue( "transInfo" )
  json.Unmarshal( []byte(formValue), &dataForm )

  type DataResponse struct {
    Status        int    `json:"status"`
    ErrMsg        string `json:"err_msg"`
    TransactionId string `json:"transaction_id"`
  }

  dataRes  := DataResponse{ Status: -1 }

  transfer := FindTransferByTradeNo( dataForm.PartnerTradeNo )
  if transfer.ID==0 {
    dataRes.ErrMsg = "不存在该订单"
  } else if transfer.Status==1 {
    dataRes.Status = 0
    dataRes.ErrMsg = "SUCCESS"
  } else {
    transactionId, err := requestTransfer( transfer )
    if err != nil {
      dataRes.ErrMsg = err.Error()
    } else {
      dataRes.Status = 0
      dataRes.ErrMsg = "SUCCESS"
    }
    dataRes.TransactionId = transactionId
  }

  fmt.Println( "ResendTransfer -------> ", dataRes )

  c.Ctx.Output.JSON( dataRes, false, false )
}

/* transactionId, err */
func requestTransfer( transfer models.Transfer ) (string, error){
  requestURL := beego.AppConfig.String( "Transfer::RequestURL" )

  var dataReq models.TransferReq
  dataReq.MchAppid       = transfer.Appid
  dataReq.Mchid          = beego.AppConfig.String( "partnerId" )
  dataReq.DeviceInfo     = "WEB"
  dataReq.NonceStr       = utils.GetRandomString( 32 )
  dataReq.PartnerTradeNo = transfer.PartnerTradeNo
  dataReq.Openid         = transfer.Openid
  dataReq.CheckName      = "FORCE_CHECK"
  dataReq.ReUserName     = transfer.ReUserName
  dataReq.Amount         = transfer.Amount
  dataReq.Desc           = "提现"
  dataReq.SpbillCreateIp = "127.0.0.1"
  dataReq.Sign = dataReq.GenerateMD5()

  resp, err := securePost( requestURL, []byte(dataReq.ToXml()) )
  if err != nil {
    return "", err
  }

  respBytes, err := ioutil.ReadAll( resp.Body )
  if err != nil {
    return "", errors.New( "解析返回数据错误" )
  }

  xmlResp := models.TransferRes{}
  err = xml.Unmarshal( respBytes, &xmlResp )
  if xmlResp.ReturnCode=="FAIL" || xmlResp.ResultCode=="FAIL" {
    transactionId := ""
    if xmlResp.ErrCode != "SYSTEMERROR" {
      DeleteTransferBy( xmlResp.PartnerTradeNo )
    } else {
      transactionId = dataReq.PartnerTradeNo
    }

    if xmlResp.ReturnCode=="FAIL" {
      return transactionId, errors.New( xmlResp.ReturnMsg )
    } else {
    return transactionId, errors.New( xmlResp.ErrCodeDes )
    }
  }

  if xmlResp.ReturnCode=="SUCCESS" && xmlResp.ResultCode=="SUCCESS" {
    transferRef := FindTransferByTradeNo( xmlResp.PartnerTradeNo )
    transferRef.Status      = 1
    transferRef.PaymentNo   = xmlResp.PaymentNo
    transferRef.PaymentTime = xmlResp.PaymentTime
    UpdateTransfer( transferRef )
    return xmlResp.PaymentNo, nil
  }

  return "", errors.New( "请求异常" )
}

func queryTransfer( transfer models.Transfer ) (int, error) {
  var dataReq models.TransferQuery
  dataReq.NonceStr       = utils.GetRandomString( 32 )
  dataReq.PartnerTradeNo = transfer.PartnerTradeNo
  dataReq.Mchid          = beego.AppConfig.String( "partnerId" )
  dataReq.Appid          = beego.AppConfig.String( "wxAppId"   )
  dataReq.Sign = dataReq.GenerateMD5()

  queryURL := beego.AppConfig.String( "Transfer::QueryURL" )

  resp, err := securePost( queryURL, []byte(dataReq.ToXml()) )
  if err != nil {
    return -1, err
  }

  respBytes, err := ioutil.ReadAll( resp.Body )
  if err != nil {
    return -1, errors.New( "解析返回数据错误" )
  }

  xmlResp := models.TransferRet{}
  err = xml.Unmarshal( respBytes, &xmlResp )
  if xmlResp.ReturnCode=="FAIL" || xmlResp.ResultCode=="FAIL" {
    status := -1
    if xmlResp.ErrCode != "SYSTEMERROR" {
      DeleteTransferBy( xmlResp.PartnerTradeNo )
    } else {
      status = 1
    }

    if xmlResp.ReturnCode=="FAIL" {
      return status, errors.New( xmlResp.ReturnMsg  )
    } else {
      return status, errors.New( xmlResp.ErrCodeDes )
    }
  }

  if xmlResp.ReturnCode=="SUCCESS" && xmlResp.ResultCode=="SUCCESS" {
    return 0, nil
  }

  return -1, errors.New( "请求异常" )
}

func generateTransfer( openid string, realname string, amount int ) models.Transfer {
  var transfer models.Transfer
  transfer.Appid          = beego.AppConfig.String( "wxAppId" )
  transfer.Openid         = openid
  transfer.PartnerTradeNo = strconv.FormatInt(time.Now().Unix(), 10) + getRandomString( 3 )
  transfer.ReUserName     = realname
  transfer.Amount         = amount
  transfer.Desc           = "提现"
  transfer.Status         = 0
  if CreateTransfer( transfer ) != true {
    transfer = models.Transfer{}
  }

  return transfer
}

func checkRequestValid( uid string, realname string, amount int ) error {
  requestUrl := "http://127.0.0.1/gameMoneyCount?gameUserId=" + uid
  resp, err := http.Get( requestUrl )
  if err != nil {
    return err
  }
  defer resp.Body.Close()

  body, err := ioutil.ReadAll( resp.Body )
  var info map[string]interface{}
  err = json.Unmarshal( body, &info )
  if err != nil {
    return err
  }

  if info["code"]==-1 {
    return errors.New( info["msg"].(string) )
  }

  result := info["result"].(map[string]interface{})
  if strconv.Itoa( int(result["gameUserId"].(float64))) != uid {
   return errors.New( "不存在该用户" )
  } else if result["name"].(string) != realname {
    return errors.New( "玩家真实姓名有误" )
  } else if int(result["moneyCount"].(float64) * 100) < amount {
    return errors.New( "可提现金额不足" )
  }
  return nil
}

/* 企业付款双向认证操作 */
var _tlsConfig *tls.Config
func getTLSConfig() (*tls.Config, error) {
  if _tlsConfig != nil {
    return _tlsConfig, nil
  }

  certPath := filepath.Join( `static/cert`, beego.AppConfig.String( "Transfer::CertFile" ) )
  keyPath  := filepath.Join( `static/cert`, beego.AppConfig.String( "Transfer::KeyFile"  ) )
  caPath   := filepath.Join( `static/cert`, beego.AppConfig.String( "Transfer::CAFile"   ) )

  cert , err := tls.LoadX509KeyPair( certPath, keyPath )
  if err != nil {
    fmt.Println( "load wechat key fail", err )
    return nil, err
  }

  caData, err := ioutil.ReadFile( caPath )
  if err != nil {
    fmt.Println( "read wechat ca fail", err )
    return nil, err
  }

  pool := x509.NewCertPool()
  pool.AppendCertsFromPEM( caData )

  _tlsConfig = &tls.Config {
    Certificates: []tls.Certificate{ cert },
    RootCAs:      pool,
  }

  return _tlsConfig, nil
}

func securePost( requestURL string, xmlContent []byte ) (*http.Response, error) {
  tlsConfig, err := getTLSConfig()
  if err != nil {
    return nil, err
  }

  tr         := &http.Transport{ TLSClientConfig: tlsConfig }
  client     := &http.Client{ Transport: tr }
  return client.Post(
    requestURL,
    "text/xml",
    bytes.NewBuffer( xmlContent ),
  )
}
