package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/hkloudou/clux"
)

func GetIntranetIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				// fmt.Println("ip:", ipnet.IP.String())
				return ipnet.IP.String()
			}
		}
	}
	panic("GetIntranetIp err")
}
func main() {
	if tmp, err := clux.NewServiceRegister(5); err != nil {
		panic(err)
	} else {
		ip := GetIntranetIp()
		tmp.PutService("/system/proxy/nodes/"+ip, fmt.Sprintf("http://%s:8080", ip))
		tmp.ListenLeaseRespChan()
		time.Sleep(10 * time.Second) //关闭通道后过10秒钟停掉
	}
}
