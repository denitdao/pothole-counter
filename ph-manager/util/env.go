package util

import "github.com/magiconair/properties"

var p = properties.MustLoadFile("application.properties", properties.UTF8)

func GetProperty(key string) string {
	return p.MustGetString(key)
}
