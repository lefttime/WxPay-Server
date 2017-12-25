package models

type Product struct {
  ID     string `json:"id"`
  Name   string `json:"name"`
  Amount int    `json:"amount"`
  Count  int    `json:"cnt"`
}