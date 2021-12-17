package plugins

import (
	"CrackerBoy/models"
	"database/sql"
	_"github.com/go-sql-driver/mysql"
	"fmt"
)
func ScanMysql(service models.Service)(result models.ScanResult){
	result.Service=service
	result.Iscrack=false
	db,_:=sql.Open("mysql",fmt.Sprintf("%v:%v@tcp(%v:%v)/mysql",service.Username,service.Password,service.Ip,service.Port))
	defer db.Close()
	err:=db.Ping()
	if err==nil{
		result.Iscrack=true
	}
	return  result
}
