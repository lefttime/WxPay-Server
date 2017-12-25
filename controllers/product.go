package controllers

import (
  "encoding/json"
  "github.com/astaxie/beego"
  "github.com/lefttime/PaymentSystem/models"
)

var products []models.Product

type Product struct {
  beego.Controller
}

func (c *Product) Get() {
  area := c.Ctx.Input.Query( "area" )

  type DataResponse struct {
    Status int    `json:"status"`
    Data   string `json:"data"`
  }

  dataResp := DataResponse{ Status: -1 }
  if area=="xiangyang" {
    jsonData, err := json.Marshal( products )
    if err==nil {
      dataResp.Status = 0
      dataResp.Data   = string( jsonData )
    }
  }
  c.Ctx.Output.JSON( dataResp, false, false )
}

func GetProductById( id int ) models.Product {
  if id < 0 || id >= len( products ) {
    return models.Product{ ID: "", Name: "", Amount: 0, Count: 0 }
  }

  return products[id]
}

func GetProductByInnerId( idx string ) models.Product {
  for _, product := range products {
    if product.ID==idx {
      return product
    }
  }

  return models.Product{ ID:"" }
}

func init() {
  products = append( products, models.Product{ ID: "com.xc.6",   Name: "6颗钻石",   Amount: 6,   Count: 6   } )
  products = append( products, models.Product{ ID: "com.xc.18",  Name: "18颗钻石",  Amount: 18,  Count: 18  } )
  products = append( products, models.Product{ ID: "com.xc.30",  Name: "30颗钻石",  Amount: 30,  Count: 30  } )
  products = append( products, models.Product{ ID: "com.xc.68",  Name: "68颗钻石",  Amount: 68,  Count: 68  } )
  products = append( products, models.Product{ ID: "com.xc.128", Name: "128颗钻石", Amount: 128, Count: 128 } )
  products = append( products, models.Product{ ID: "com.xc.328", Name: "328颗钻石", Amount: 328, Count: 328 } )
}
