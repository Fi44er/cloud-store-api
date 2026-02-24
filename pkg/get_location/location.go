package getlocation

import (
	_ "embed"
	"net"

	"github.com/oschwald/geoip2-golang"
)

//go:embed GeoLite2-City.mmdb
var geoIPData []byte

func GetIPLocation(ipString string) (*geoip2.City, error) {
	db, err := geoip2.FromBytes(geoIPData)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	ip := net.ParseIP(ipString)
	record, err := db.City(ip)
	if err != nil {
		return nil, err
	}

	return record, nil
}
