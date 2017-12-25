package models

import (
  "fmt"
  "strings"
  "encoding/xml"
  "github.com/astaxie/beego"
  "github.com/lefttime/PaymentSystem/utils"
)

type Transfer struct {
  ID             int    `json:"id"`
  Appid          string `json:"appid"`
  Openid         string `json:"openid"`
  PartnerTradeNo string `json:"partner_trade_no"`
  ReUserName     string `json:"re_user_name"`
  Amount         int    `json:"amount"`
  Desc           string `json:"desc"`
  PaymentNo      string `json:"payment_no"`
  PaymentTime    string `json:"payment_time"`
  Status         int    `json:"status"`
}

/* 企业付款请求结构 */
type TransferReq struct {
  MchAppid       string `xml:"mch_appid"`        // 微信分配的账号ID（企业号corpid即为此appId）
  Mchid          string `xml:"mchid"`            // 微信支付分配的商户号
  DeviceInfo     string `xml:"device_info"`      // 微信支付分配的终端设备号
  NonceStr       string `xml:"nonce_str"`        // 随机字符串，不长于32位
  Sign           string `xml:"sign"`
  PartnerTradeNo string `xml:"partner_trade_no"` // 商户订单号，需保持唯一性(只能是字母或者数字，不能包含有符号)
  Openid         string `xml:"openid"`           // 商户appid下，某用户的openid
  CheckName      string `xml:"check_name"`       // NO_CHECK：不校验真实姓名 FORCE_CHECK：强校验真实姓名
  ReUserName     string `xml:"re_user_name"`     // 收款用户真实姓名。 如果check_name设置为FORCE_CHECK，则必填用户真实姓名
  Amount         int    `xml:"amount"`           // 企业付款金额，单位为分
  Desc           string `xml:"desc"`             // 企业付款操作说明信息。必填
  SpbillCreateIp string `xml:"spbill_create_ip"` // 调用接口的机器Ip地址
}

/* 企业付款反馈结构 */
type TransferRes struct {
  ReturnCode     string `xml:"return_code"`
  ReturnMsg      string `xml:"return_msg"`
  MchAppid       string `xml:"mch_appid"`
  Mchid          string `xml:"mchid"`
  DeviceInfo     string `xml:"device_info"`
  NonceStr       string `xml:"nonce_str"`
  ResultCode     string `xml:"result_code"`
  ErrCode        string `xml:"err_code"`
  ErrCodeDes     string `xml:"err_code_des"`
  PartnerTradeNo string `xml:"partner_trade_no"`
  PaymentNo      string `xml:"payment_no"`       // 企业付款成功，返回的微信订单号
  PaymentTime    string `xml:"payment_time"`     // 企业付款成功时间
}

/* 企业付款查询结构 */
type TransferQuery struct {
  NonceStr       string `xml:"nonce_str"`        // 随机字符串，不长于32位
  Sign           string `xml:"sign"`
  PartnerTradeNo string `xml:"partner_trade_no"` // 商户订单号，需保持唯一性(只能是字母或者数字，不能包含有符号)
  Mchid          string `xml:"mch_id"`           // 微信支付分配的商户号
  Appid          string `xml:"appid"`
}

/* 企业付款结果结构 */
type TransferRet struct {
  ReturnCode     string `xml:"return_code"`
  ReturnMsg      string `xml:"return_msg"`
  ResultCode     string `xml:"result_code"`
  ErrCode        string `xml:"err_code"`
  ErrCodeDes     string `xml:"err_code_des"`
  PartnerTradeNo string `xml:"partner_trade_no"`
  MchId          string `xml:"mch_id"`
  DetailId       string `xml:"detail_id"`
  Status         string `xml:"status"`
  Reason         string `xml:"reason"`
  Openid         string `xml:"openid"`
  TransferName   string `xml:"transfer_name"`
  PaymentAmount  int    `xml:"payment_amount"`
  TransferTime   string `xml:"transfer_time"`
  Des            string `xml:"desc"`
}

/*---------- 转化XML ----------*/
func (data *TransferReq) ToXml() string {
  bytes_req, err := xml.Marshal( data )
  // bytes_req, err := xml.MarshalIndent( data, "", "  " )
  if err != nil {
    fmt.Println( "转换为xml错误:", err )
  }
  str_req := strings.Replace( string(bytes_req), "TransferReq", "xml", -1 )
  return str_req
}

func (data *TransferReq) FormatXml() string {
  bytes_req, err := xml.MarshalIndent( data, "", "  " )
  if err != nil {
    fmt.Println( "转换为xml错误:", err )
  }
  str_req := strings.Replace( string(bytes_req), "TransferReq", "xml", -1 )
  return str_req
}

func (data *TransferQuery) ToXml() string {
  bytes_req, err := xml.Marshal( data )
  // bytes_req, err := xml.MarshalIndent( data, "", "  " )
  if err != nil {
    fmt.Println( "转换为xml错误:", err )
  }
  str_req := strings.Replace( string(bytes_req), "TransferQuery", "xml", -1 )
  return str_req
}

func (data *TransferRet) ToXml() string {
  bytes_req, err := xml.Marshal( data )
  // bytes_req, err := xml.MarshalIndent( data, "", "  " )
  if err != nil {
    fmt.Println( "转换为xml错误:", err )
  }
  str_req := strings.Replace( string(bytes_req), "TransferRet", "xml", -1 )
  return str_req
}

/*---------- 生成签名 ----------*/
func (data *TransferReq) GenerateMD5() string {
  params := make( map[string]interface{}, 0 )
  params["mch_appid"]        = data.MchAppid
  params["mchid"]            = data.Mchid
  params["device_info"]      = data.DeviceInfo
  params["nonce_str"]        = data.NonceStr
  params["partner_trade_no"] = data.PartnerTradeNo
  params["openid"]           = data.Openid
  params["check_name"]       = data.CheckName
  params["re_user_name"]     = data.ReUserName
  params["amount"]           = data.Amount
  params["desc"]             = data.Desc
  params["spbill_create_ip"] = data.SpbillCreateIp

  return utils.GenerateSign( params, beego.AppConfig.String( "wxPayKey" ) )
}

func (data *TransferQuery) GenerateMD5() string {
  params := make( map[string]interface{}, 0 )
  params["nonce_str"]        = data.NonceStr
  params["partner_trade_no"] = data.PartnerTradeNo
  params["mch_id"]           = data.Mchid
  params["nonce_str"]        = data.NonceStr
  params["appid"]            = data.Appid

  return utils.GenerateSign( params, beego.AppConfig.String( "wxPayKey" ) )
}

func (data *TransferRet) GenerateMD5() string {
  params := make( map[string]interface{}, 0 )
  params["return_code"]      = data.ReturnCode
  params["return_msg"]       = data.ReturnMsg
  params["result_code"]      = data.ResultCode
  params["err_code"]         = data.ErrCode
  params["err_code_des"]     = data.ErrCodeDes
  params["partner_trade_no"] = data.PartnerTradeNo
  params["mch_id"]           = data.MchId
  params["detail_id"]        = data.DetailId
  params["status"]           = data.Status
  params["reason"]           = data.Reason
  params["openid"]           = data.Openid
  params["transfer_name"]    = data.TransferName
  params["payment_amount"]   = data.PaymentAmount
  params["transfer_time"]    = data.TransferTime
  params["desc"]             = data.Des

  return utils.GenerateSign( params, beego.AppConfig.String( "wxPayKey" ) )
}
