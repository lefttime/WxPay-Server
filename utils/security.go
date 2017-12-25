package utils

import (
  "fmt"
  "sort"
  "strings"
  "crypto/md5"
  "encoding/hex"
)

func GenerateSign( params map[string]interface{}, key string ) string {
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