package models

type Player struct {
  ID       string `json:"id"`
  Name     string `json:"name"`
  Diamonds int    `json:"diamonds"`
  Avatar   string `json:"avatar"`
  OpenId   string `json:"openid"`
}