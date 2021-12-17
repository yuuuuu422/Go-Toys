package plugins

import "CrackerBoy/models"

type scanFunc func(service models.Service)(result models.ScanResult)

var FuncMap map[string]scanFunc

func init(){
	FuncMap =make(map[string]scanFunc)
	FuncMap["ssh"]=ScanSsh
	FuncMap["mysql"]=ScanMysql
}