package models

import (
  "fmt"
  "sort"
  "strings"
  "crypto/md5"
  "encoding/hex"
  "encoding/xml"
  "github.com/astaxie/beego"
)

/* 订单存储结构 */
type Order struct {
  ID               string `json:"id"`
  Appid            string `json:"appid"`
  UserId           string `json:"uid"`
  Quantity         int    `json:"quantity"`
  PrePurchaseDate  string `json:"original_purchase_date"`
  PreTransactionId string `json:"original_transaction_id"`
  ProductId        string `json:"product_id"`
  ProductName      string `json:"product_name"`
  ProductCount     int    `json:"product_count"`
  ProductAmount    int    `json:"product_amount"`
  PurchaseDate     string `json:"purchase_date"`
  TransactionId    string `json:"transaction_id"`
  PurchaseType     string `json:"purchase_type"`
  Platform         string `json:"platform"`
  IP               string `json:"ip"`
}

/* 统一下单请求 */
type UnifyOrderReq struct {
  Appid          string `xml:"appid"`
  MchId          string `xml:"mch_id"`
  DeviceInfo     string `xml:"device_info"`
  NonceStr       string `xml:"nonce_str"`
  Sign           string `xml:"sign"`
  SignType       string `xml:"sign_type"`
  Body           string `xml:"body"`
  Detail         string `xml:"detail"`
  Attach         string `xml:"attach"`
  OutTradeNo     string `xml:"out_trade_no"`
  FeeType        string `xml:"fee_type"`
  TotalFee       int    `xml:"total_fee"`
  SpbillCreateIp string `xml:"spbill_create_ip"`
  TimeStart      string `xml:"time_start"`
  TimeExpire     string `xml:"time_expire"`
  GoodTag        string `xml:"good_tag"`
  NotifyUrl      string `xml:"notify_url"`
  TradeType      string `xml:"trade_type"`
  ProductId      string `xml:"product_id"`
  LimitPay       string `xml:"limit_pay"`
  OpenId         string `xml:"openid"`
  SceneInfo      string `xml:"scene_info"`
}

/* 统一下单反馈 */
type UnifyOrderResp struct {
  ReturnCode string `xml:"return_code"`
  ReturnMsg  string `xml:"return_msg"`
  AppId      string `xml:"appid"`
  MchId      string `xml:"mch_id"`
  DeviceInfo string `xml:"device_info"`
  NonceStr   string `xml:"nonce_str"`
  Sign       string `xml:"sign"`
  ResultCode string `xml:"result_code"`
  ErrCode    string `xml:"err_code"`
  ErrCodeDes string `xml:"err_code_des"`
  TradeType  string `xml:"trade_type"`
  PrepayId   string `xml:"prepay_id"`
  MwebUrl    string `xml:"mweb_url"`
}

/* 支付结果通知 */
type OrderResult struct {
  ReturnCode         string `xml:"return_code"`
  ReturnMsg          string `xml:"return_msg"`
  Appid              string `xml:"appid"`
  MchId              string `xml:"mch_id"`
  DeviceInfo         string `xml:"device_info"`
  NonceStr           string `xml:"nonce_str"`
  Sign               string `xml:"sign"`
  SignType           string `xml:"sign_type"`
  ResultCode         string `xml:"result_code"`
  ErrCode            string `xml:"err_code"`
  ErrCodeDes         string `xml:"err_code_des"`
  OpenId             string `xml:"openid"`
  IsSubscribe        string `xml:"is_subscribe"`
  TradeType          string `xml:"trade_type"`
  BankType           string `xml:"bank_type"`
  TotalFee           int    `xml:"total_fee"`
  SettlementTotalFee int    `xml:"settlement_total_fee"`
  FeeType            string `xml:"fee_type"`
  CashFee            int    `xml:"cash_fee"`
  CashFeeType        string `xml:"cash_fee_type"`
  CouponFee          int    `xml:"coupon_fee"`
  CouponCount        int    `xml:"coupon_count"`
  CouponTypeN        string `xml:"coupon_type_$n"`
  CouponIdN          string `xml:"coupon_id_$n"`
  CouponFeeN         int    `xml:"coupon_fee_$n"`
  TransactionId      string `xml:"transaction_id"`
  OutTradeNo         string `xml:"out_trade_no"`
  Attach             string `xml:"attach"`
  TimeEnd            string `xml:"time_end"`
}

/* 支付结果反馈 */
type OrderResultResp struct {
  ReturnCode string `xml:"return_code"`
  ReturnMsg  string `xml:"return_msg"`
}

/* 订单查询请求 */
type QueryOrderReq struct {
  Appid         string `xml:"appid"`
  MchId         string `xml:"mch_id"`
  TransactionId string `xml:"transaction_id"`
  OutTradeNo    string `xml:"out_trade_no"`
  NonceStr      string `xml:"nonce_str"`
  Sign          string `xml:"sign"`
  SignType      string `xml:"sign_type"`
}

/* 订单查询反馈 */
type QueryOrderResp struct {
  ReturnCode         string `xml:"return_code"`
  ReturnMsg          string `xml:"return_msg"`
  Appid              string `xml:"appid"`
  MchId              string `xml:"mch_id"`
  NonceStr           string `xml:"nonce_str"`
  Sign               string `xml:"sign"`
  ResultCode         string `xml:"result_code"`
  ErrCode            string `xml:"err_code"`
  ErrCodeDes         string `xml:"err_code_des"`
  DeviceInfo         string `xml:"device_info"`
  OpenId             string `xml:"open_id"`
  IsSubScribe        string `xml:"is_subscribe"`
  TradeType          string `xml:"trade_type"`
  TradeState         string `xml:"trade_state"`
  BankType           string `xml:"bank_type"`
  TotalFee           int    `xml:"total_fee"`
  SettlementTotalFee int    `xml:"settlement_total_fee"`
  FeeType            string `xml:"fee_type"`
  CashFee            int    `xml:"cash_fee"`
  CashFeeType        string `xml:"cash_fee_type"`
  CouponFee          int    `xml:"coupon_fee"`
  CouponCount        int    `xml:"coupon_count"`
  CouponTypeN        string `xml:"coupon_type_$n"`
  CouponIdN          string `xml:"coupon_id_$n"`
  CouponFeeN         int    `xml:"coupon_fee_$n"`
  TransactionId      string `xml:"transaction_id"`
  OutTradeNo         string `xml:"out_trade_no"`
  Attach             string `xml:"attach"`
  TimeEnd            string `xml:"time_end"`
  TradeStateDesc     string `xml:"trade_state_desc"`
}

/*---------- 转化XML ----------*/
func (data *UnifyOrderReq) ToXml() string {
  bytes_req, err := xml.Marshal( data )
  // bytes_req, err := xml.MarshalIndent( data, "", "  " )
  if err != nil {
    fmt.Println( "转换为xml错误:", err )
  }
  str_req := strings.Replace( string(bytes_req), "UnifyOrderReq", "xml", -1 )
  return str_req
}

func (data *UnifyOrderResp) ToXml() string {
  // bytes_req, err := xml.Marshal( data )
  bytes_req, err := xml.MarshalIndent( data, "", "  " )
  if err != nil {
    fmt.Println( "转换为xml错误", err )
  }
  str_req := strings.Replace( string(bytes_req), "UnifyOrderResp", "xml", -1 )
  return str_req
}

func (data *OrderResult) ToXml() string {
  // bytes_req, err := xml.Marshal( data )
  bytes_req, err := xml.MarshalIndent( data, "", "  " )
  if err != nil {
    fmt.Println( "转换为xml错误", err )
  }
  str_req := strings.Replace( string(bytes_req), "UnifyOrderResp", "xml", -1 )
  return str_req
}

func (data *OrderResultResp) ToXml() string {
  // bytes_req, err := xml.Marshal( data )
  bytes_req, err := xml.MarshalIndent( data, "", "  " )
  if err != nil {
    fmt.Println( "转换为xml错误", err )
  }
  str_req := strings.Replace( string(bytes_req), "UnifyOrderResp", "xml", -1 )
  return str_req
}

func (data *QueryOrderReq) ToXml() string {
  // bytes_req, err := xml.Marshal( data )
  bytes_req, err := xml.MarshalIndent( data, "", "  " )
  if err != nil {
    fmt.Println( "转换为xml错误", err )
  }
  str_req := strings.Replace( string(bytes_req), "UnifyOrderResp", "xml", -1 )
  return str_req
}

func (data *QueryOrderResp) ToXml() string {
  // bytes_req, err := xml.Marshal( data )
  bytes_req, err := xml.MarshalIndent( data, "", "  " )
  if err != nil {
    fmt.Println( "转换为xml错误", err )
  }
  str_req := strings.Replace( string(bytes_req), "UnifyOrderResp", "xml", -1 )
  return str_req
}

/*---------- 生成签名 ----------*/
func (data *UnifyOrderReq) GenerateMD5() string {
  params := make( map[string]interface{}, 0 )
  params[ "appid" ]            = data.Appid
  params[ "mch_id" ]           = data.MchId
  params[ "device_info"]       = data.DeviceInfo
  params[ "nonce_str" ]        = data.NonceStr
  params[ "sign_type" ]        = data.SignType
  params[ "body" ]             = data.Body
  params[ "detail" ]           = data.Detail
  params[ "attach" ]           = data.Attach
  params[ "out_trade_no" ]     = data.OutTradeNo
  params[ "free_type" ]        =  data.FeeType
  params[ "total_fee" ]        = data.TotalFee
  params[ "spbill_create_ip" ] = data.SpbillCreateIp
  params[ "time_start"  ]      = data.TimeStart
  params[ "time_expire" ]      = data.TimeExpire
  params[ "good_tag" ]         = data.GoodTag
  params[ "notify_url" ]       = data.NotifyUrl
  params[ "trade_type" ]       = data.TradeType
  params[ "product_id" ]       = data.ProductId
  params[ "limit_pay" ]        = data.LimitPay
  params[ "openid" ]           = data.OpenId
  params[ "scene_info" ]       = data.SceneInfo
  return generateSign( params, beego.AppConfig.String( "wxPayKey" ) )
}

func (data *UnifyOrderResp) GenerateMD5() string {
  params := make( map[string]interface{}, 0 )
  params[ "return_code" ]  = data.ReturnCode
  params[ "return_msg" ]   = data.ReturnMsg
  params[ "appid"]         = data.AppId
  params[ "mch_id" ]       = data.MchId
  params[ "device_info" ]  = data.DeviceInfo
  params[ "nonce_str" ]    = data.NonceStr
  params[ "sign" ]         = data.Sign
  params[ "result_code" ]  = data.ResultCode
  params[ "err_code" ]     = data.ErrCode
  params[ "err_code_des" ] = data.ErrCodeDes
  params[ "trade_type" ]   = data.TradeType
  params[ "prepay_id" ]    = data.PrepayId
  params[ "mweb_url" ]     = data.MwebUrl
  return generateSign( params, beego.AppConfig.String( "wxPayKey" ) )
}

func (data *OrderResult) GenerateMD5() string {
  params := make( map[string]interface{}, 0 )
  params[ "return_code" ]          = data.ReturnCode
  params[ "return_msg" ]           = data.ReturnMsg
  params[ "appid" ]                = data.Appid
  params[ "mch_id" ]               = data.MchId
  params[ "device_info" ]          = data.DeviceInfo
  params[ "nonce_str" ]            = data.NonceStr
  params[ "sign" ]                 = data.Sign
  params[ "sign_type" ]            = data.SignType
  params[ "result_code" ]          = data.ResultCode
  params[ "err_code" ]             = data.ErrCode
  params[ "err_code_des" ]         = data.ErrCodeDes
  params[ "openid" ]               = data.OpenId
  params[ "is_subscribe" ]         = data.IsSubscribe
  params[ "trade_type" ]           = data.TradeType
  params[ "bank_type" ]            = data.BankType
  params[ "total_fee" ]            = data.TotalFee
  params[ "fee_type" ]             = data.FeeType
  params[ "cash_fee_type" ]        = data.CashFeeType
  params[ "coupon_type_$n" ]       = data.CouponTypeN
  params[ "coupon_id_$n" ]         = data.CouponIdN
  params[ "transaction_id" ]       = data.TransactionId
  params[ "out_trade_no" ]         = data.OutTradeNo
  params[ "attach" ]               = data.Attach
  params[ "time_end" ]             = data.TimeEnd

  if data.SettlementTotalFee > 0 {
    params[ "settlement_total_fee" ] = data.SettlementTotalFee
  }
  if data.CashFee > 0 {
    params[ "cash_fee" ] = data.CashFee
  }
  if data.CouponFee > 0 {
    params[ "coupon_fee" ] = data.CouponFee
  }
  if data.CouponCount > 0 {
    params[ "coupon_count" ] = data.CouponCount
  }
  if data.CouponFeeN > 0 {
    params[ "coupon_fee_$n" ] = data.CouponFeeN
  }

  return generateSign( params, beego.AppConfig.String( "wxPayKey" ) ) 
}

func (data *QueryOrderReq) GenerateMD5() string {
  params := make( map[string]interface{}, 0 )
  params["appid"]          = data.Appid
  params["mch_id"]         = data.MchId
  params["transaction_id"] = data.TransactionId
  params["out_trade_no"]   = data.OutTradeNo
  params["nonce_str"]      = data.NonceStr
  params["sign"]           = data.Sign
  params["sign_type"]      = data.SignType
  return generateSign( params, beego.AppConfig.String( "wxPaKey" ) )
}

func (data *QueryOrderResp) GenerateMD5() string {
  params := make( map[string]interface{}, 0 )
  params["return_code"]          = data.ReturnCode
  params["return_msg"]           = data.ReturnMsg
  params["appid"]                = data.Appid
  params["mch_id"]               = data.MchId
  params["nonce_str"]            = data.NonceStr
  params["sign"]                 = data.Sign
  params["result_code"]          = data.ResultCode
  params["err_code"]             = data.ErrCode
  params["err_code_des"]         = data.ErrCodeDes
  params["device_info"]          = data.DeviceInfo
  params["open_id"]              = data.OpenId
  params["is_subscribe"]         = data.IsSubScribe
  params["trade_type"]           = data.TradeType
  params["trade_state"]          = data.TradeState
  params["bank_type"]            = data.BankType
  params["total_fee"]            = data.TotalFee
  params["settlement_total_fee"] = data.SettlementTotalFee
  params["fee_type"]             = data.FeeType
  params["cash_fee"]             = data.CashFee
  params["cash_fee_type"]        = data.CashFeeType
  params["coupon_fee"]           = data.CouponFee
  params["coupon_count"]         = data.CouponCount
  params["coupon_type_$n"]       = data.CouponTypeN
  params["coupon_id_$n"]         = data.CouponIdN
  params["coupon_fee_$n"]        = data.CouponFeeN
  params["transaction_id"]       = data.TransactionId
  params["out_trade_no"]         = data.OutTradeNo
  params["attach"]               = data.Attach
  params["time_end"]             = data.TimeEnd
  params["trade_state_desc"]     = data.TradeStateDesc
  return generateSign( params, beego.AppConfig.String( "wxPaKey" ) )
}

func generateSign( params map[string]interface{}, key string ) string {
  sortedKeys := make( []string, 0 )
  for sKey, _ := range params {
    if sKey != "sign" {
      sortedKeys = append( sortedKeys, sKey )
    }
  }
  sort.Strings( sortedKeys )

  var signStrings string
  for _, sKey := range sortedKeys {
    param := fmt.Sprintf( "%v", params[sKey] )
    if param != "" {
      signStrings = signStrings + sKey + "=" + param + "&"
    }
  }
  signStrings = signStrings + "key=" + key

  md5Ctx := md5.New()
  md5Ctx.Write( []byte(signStrings) )
  cipherStr := md5Ctx.Sum( nil )
  result := strings.ToUpper( hex.EncodeToString( cipherStr ) )
  return result
}
