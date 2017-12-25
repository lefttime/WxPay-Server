package models

type SignGenerator interface {
  GenerateMD5() string
}

type XmlFormater interface {
  ToXml() string
}