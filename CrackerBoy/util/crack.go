package util

import (
	"CrackerBoy/models"
	"CrackerBoy/plugins"
	"fmt"
	"github.com/urfave/cli"
	"sync"
)

func CreatTask(addr models.IpAddr,usernames []string,passwords []string)(tasks []models.Service){
	for _,username :=range usernames{
		for _,password:=range passwords {
			service:=models.Service{Ip: addr.Ip,Port: addr.Port,Protocol: addr.Protocol,Username: username,Password: password}
			tasks=append(tasks,service)
		}
	}
	return tasks
}

func ExecTask(tasks []models.Service){
	taskChan:=make(chan models.Service,models.Thread)

	//Create multi-threaded
	for i:=0;i<models.Thread;i++{
		go Crack(taskChan,&Wg)
	}

	for _,task:=range tasks{
		Wg.Add(1)
		taskChan<-task
	}
	Wg.Wait()
	close(taskChan)
}

func Crack(taskChan chan models.Service,wg *sync.WaitGroup){
	for task:=range taskChan{
		scan:=plugins.FuncMap[task.Protocol]
		if scan(task).Iscrack==true{
			fmt.Println(task)
		}
		wg.Done()
	}
}

func Begin(cli *cli.Context){
	if cli.IsSet("thread"){
		models.Thread=cli.Int("thread")
	}

	ipList:=ReadIpList("input.txt")
	AliveAddr:=CheckAddr(ipList)
	tasks:=make([]models.Service,0)
	//这个地方写的很垃圾... 后续改进吧
	for _,addr :=range(AliveAddr){
		userdict:=models.Dict[addr.Protocol][0]
		passdict:=models.Dict[addr.Protocol][1]
		usernames:=ReadUsernameDict(userdict)
		passwords:=ReadPasswordDict(passdict)
		tasks= append(tasks, CreatTask(addr,usernames,passwords)...)
	}
	ExecTask(tasks)
}


