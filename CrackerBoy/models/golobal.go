package models

import (
	"time"
)

var (
	IpList     = "input.txt"
	ResultList = "output.txt"
	TimeOut    = 3 * time.Second
	Thread     = 500

	PortNames = map[int]string{
		22:   "ssh",
		3306: "mysql",
	}
	Dict = map[string][]string{
		"mysql": {"./dict/mysql_user.txt","./dict/mysql_pass.txt"},
		"ssh": {"./dict/ssh_user.txt","./dict/ssh_pass.txt"},
	}
	SupportProtocols map[string]bool
)

func init() {
	SupportProtocols=make(map[string]bool)
	for _,index:=range PortNames{
		SupportProtocols[index]=true
	}
}
