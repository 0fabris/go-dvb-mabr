module github.com/0fabris/go-dvb-mabr

go 1.21

replace github.com/0fabris/go-dvb-route => ../go-dvb-route

require (
	github.com/0fabris/go-dvb-route v0.0.0-00010101000000-000000000000
	github.com/seancfoley/ipaddress-go v1.5.4
)

require github.com/seancfoley/bintree v1.2.1 // indirect
