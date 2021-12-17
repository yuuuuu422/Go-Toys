package models

// for crack
type Service struct {
	Ip string
	Port int
	Protocol string
	Username string
	Password string

}

type ScanResult struct {
	Service Service
	Iscrack bool
}

// before cracking, determine whether the port is alive
type IpAddr struct {
	Ip string
	Port int
	Protocol string
}
