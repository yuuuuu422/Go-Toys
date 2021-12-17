package util

import (
	"CrackerBoy/models"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ReadIpList(fileName string) (ipList []models.IpAddr) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	//fmt.Println(scanner.Text())
	for scanner.Scan() {
		line := strings.Replace(scanner.Text(), " ", "", -1)
		tmp1 := strings.Split(line, "|")
		ipPort := tmp1[0]
		tmp2 := strings.Split(ipPort, ":")
		ip := tmp2[0]
		port, _ := strconv.Atoi(tmp2[1])
		if len(tmp1) == 2 {
			protocol := tmp1[1]
			if models.SupportProtocols[protocol] {
				ipAddr := models.IpAddr{Ip: ip, Port: port, Protocol: protocol}
				ipList = append(ipList, ipAddr)
			}
		}else {
			protocol,exit:=models.PortNames[port]
			if !exit{
				fmt.Println(port,"not exit")
			}else {
				if models.SupportProtocols[protocol] {
					ipAddr := models.IpAddr{Ip: ip, Port: port, Protocol: protocol}
					ipList = append(ipList, ipAddr)
				}
			}
		}


	}
	return ipList
}

func ReadUsernameDict(fileName string)(usernames []string){
	file,err:=os.Open(fileName)
	if err!=nil{
		panic(err)
	}
	defer  file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		username := strings.TrimSpace(scanner.Text())
		usernames=append(usernames,username )
	}
	return usernames
}

func ReadPasswordDict(fileName string)(passwords []string){
	file,err:=os.Open(fileName)
	if err!=nil{
		panic(err)
	}
	defer  file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		password := strings.TrimSpace(scanner.Text())
		passwords=append(passwords,password )
	}
	return passwords
}

func FindDict(addr models.IpAddr)(dict[]string){
	dict=append(dict,fmt.Sprintf("./dict/%s_user.txt",addr.Protocol),fmt.Sprintf("./dict/%s_pass.txt",addr.Protocol))
	return dict
}