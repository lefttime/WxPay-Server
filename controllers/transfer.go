package controllers

import (
  "github.com/lefttime/PaymentSystem/models"
)

func CreateTransfer( transfer models.Transfer ) bool {
  if models.DB.Create( &transfer )==nil {
    return false
  }
  return true
}

func DeleteTransferBy( partnerTradeNo string ) {
  transfer := FindTransferByTradeNo( partnerTradeNo )
  models.DB.Delete( &transfer )
}

func FindTransferBy( paymentNo string ) models.Transfer {
  var transfer models.Transfer
  models.DB.Where( "payment_no = ?", paymentNo ).First( &transfer )
  return transfer
}

func FindTransferByTradeNo( partnerTradeNo string ) models.Transfer {
  var transfer models.Transfer
  models.DB.Where( "partner_trade_no = ?", partnerTradeNo ).First( &transfer )
  return transfer
}

func UpdateTransfer( transfer models.Transfer ) bool {
  var transferRef models.Transfer
  models.DB.Where( "id = ?", transfer.ID ).First( &transferRef )
  if transferRef.ID != 0 {
    models.DB.Model( &transfer ).Updates( map[string]interface{}{ "payment_no": transfer.PaymentNo, "payment_time": transfer.PaymentTime } )
    return true
  }

  return false
}
