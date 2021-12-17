package plugins

import (
	"CrackerBoy/models"
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"time"
)

func ScanSsh(service models.Service)(result models.ScanResult)  {
	result.Service=service
	result.Iscrack=false
	config:=&ssh.ClientConfig{
		User: service.Username,
		Auth:[]ssh.AuthMethod{
			ssh.Password(service.Password),
		},
		Timeout: time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error{return nil},
	}
	client,err:=ssh.Dial("tcp",fmt.Sprintf("%v:%v",service.Ip,service.Port),config)
	if err == nil {
		defer client.Close()
		session, err := client.NewSession()
		errRet := session.Run("echo hello")
		if err == nil && errRet == nil {
			result.Iscrack = true
			defer session.Close()
		}
	}
	return result
}