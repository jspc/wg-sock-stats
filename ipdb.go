package main

import (
	"net"

	"github.com/oschwald/geoip2-golang"
)

// IPDB holds the two geoip2 databases we care about;
// City and ASN.
//
// City provides some kind of locational information for an
// IP address, whereas ASN provides ownership of that IP in
// the hopes that we can use it to provide some kind of ISP
// information
type IPDB struct {
	city *geoip2.Reader
	asn  *geoip2.Reader
}

func NewIPDB(cityPath, asnPath string) (i IPDB, err error) {
	i.city, err = geoip2.Open(cityPath)
	if err != nil {
		return
	}

	i.asn, err = geoip2.Open(asnPath)

	return
}

func (i IPDB) City(n net.IP) (*geoip2.City, error) {
	return i.city.City(n)
}

func (i IPDB) ASN(n net.IP) (*geoip2.ASN, error) {
	return i.asn.ASN(n)
}
