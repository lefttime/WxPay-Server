package utils

import (
  "time"
  "math/rand"
)

func GetRandomString( total int) string {
  str    := "456789abssldgijlnprtvxz"
  bytes  := []byte( str )
  result := []byte{}
  rdn    := rand.New( rand.NewSource( time.Now().UnixNano() ) )
  for idx := 0; idx < total; idx++ {
    result = append( result, bytes[rdn.Intn( len( bytes ) )] )
  }
  return string( result )

