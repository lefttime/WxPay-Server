package controllers

import (
  "os"
  "fmt"
  "strings"
  "github.com/astaxie/beego"
  "github.com/lefttime/PaymentSystem/models"

  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/mysql"
)

type Order struct {
  beego.Controller
}

func CreateOrder( order models.Order ) bool {
  msg, err := checkValid( order )
  if err==false {
    fmt.Println( msg )
    return false
  }

  if models.DB.Create( &order )==nil {
    return false
  }
  return true
}

func FindOrderBy( transactionId string ) models.Order {
  var order models.Order
  models.DB.Where( "transaction_id = ?", transactionId ).First( &order )
  return order
}

func FindOrderByUserId( uid string ) models.Order {
  var order models.Order
  models.DB.Where( "uid = ?", uid ).First( &order )
  return order
}

func UpdateOrder( order models.Order ) bool {
  var orderRef models.Order
  models.DB.Where( "id = ?", order.ID ).First( &orderRef )
  if len( orderRef.ID ) > 0 {
    models.DB.Model( &order ).Updates( map[string]interface{}{ "purchase_date": order.PurchaseDate, "transaction_id": order.TransactionId, "purchase_type": order.PurchaseType, "platform": order.Platform, "ip": order.IP } )
    return true
  }

  return false
}

func checkValid( order models.Order ) (string, bool) {
  if len( order.UserId )==0 {
    return "用户标识异常", false
  }

  if len( order.Appid )==0 {
    return "应用标识异常", false
  }

  if order.Quantity <= 0 {
    return "数量异常", false
  }

  if len( order.PrePurchaseDate )==0 {
    return "预付款日期异常", false
  }

  if len( order.PreTransactionId )==0 {
    return "预付款订单号异常", false
  }

  if len( order.ProductId )==0 {
    return "预购产品标识异常", false
  }

  return "OK", true
}

func init() {
  url := "{user}:{password}@tcp({host}:{port})/{database}?charset={charset}&parseTime=True&loc=Local"
  url  = strings.Replace( url, "{database}", beego.AppConfig.String( "Pay::Database" ), -1 )
  url  = strings.Replace( url, "{user}",     beego.AppConfig.String( "Pay::User" ),     -1 )
  url  = strings.Replace( url, "{password}", beego.AppConfig.String( "Pay::Password" ), -1 )
  url  = strings.Replace( url, "{host}",     beego.AppConfig.String( "Pay::Host" ),     -1 )
  url  = strings.Replace( url, "{port}",     beego.AppConfig.String( "Pay::Port" ),     -1 )
  url  = strings.Replace( url, "{charset}",  beego.AppConfig.String( "Pay::Charset" ),  -1 )

  var err error
  models.DB, err = gorm.Open( beego.AppConfig.String( "Pay::Dialect" ), url )
  if err != nil {
    fmt.Println( err.Error() )
    os.Exit( -1 )
  } else {
    fmt.Println( "Connect Database success." )
  }
}