package util

import (
	"CrackerBoy/models"
	"fmt"
	"net"
	"sync"
)

var (
	AliveAddr []models.IpAddr
	Mutex     sync.Mutex
	Wg        sync.WaitGroup
)

func init() {
	AliveAddr = make([]models.IpAddr, 0)
}

func IsAlive(IpAddr models.IpAddr) bool {
	_, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", IpAddr.Ip, IpAddr.Port), models.TimeOut)
	if err == nil {
		return true
	}
	return false
}

func CheckAddr(IpAddr []models.IpAddr) (AliveAddr []models.IpAddr) {
	Wg.Add(len(IpAddr))
	for _, addr := range IpAddr {
		go func(addr models.IpAddr) {
			defer Wg.Done()
			if IsAlive(addr) {
				Mutex.Lock()
				AliveAddr = append(AliveAddr, addr)
				Mutex.Unlock()
			}
		}(addr)
	}
	Wg.Wait()
	return AliveAddr
}
