package main

import (
	"log"
	"os"
	"time"

	"github.com/hkloudou/clux"
)

type GetAddrRequest struct {
	Name         string `binding:"required"`
	Chain        string `binding:"required"`
	ContractAddr string `binding:"required"`
	Uid          uint64 `binding:"required"`
}
type GetAddrResponse struct {
	Name  string `binding:"required"`
	Chain string `binding:"required"`
	Addr  string `binding:"required"`
}

func main() {
	app := clux.Action(clux.UseNats())

	if err := app.Run(os.Args); err != nil {
		log.Println(err)
		time.Sleep(10 * time.Second)
		panic(err)
	}
	// clux.NatsClient().QueueSubscribe(clux.NameSpace()+".wallet.getaddr.trc20", "worker", func(m *nats.Msg) {
	// 	// m.RespondMsg(handle(m))
	// 	c := clux.NewContext(m)
	// 	c.Writer.Header().Set("addr", "yy")
	// 	var req GetAddrRequest
	// 	// log.Println()
	// 	if err := c.ShouldBindHeaderRaw(&req); err != nil {
	// 		c.AbortWithError(err)
	// 		return
	// 	}
	// 	c.JSON(GetAddrResponse{
	// 		Addr:  "123",
	// 		Name:  "ut",
	// 		Chain: "tr",
	// 	}, nil)
	// })
	<-make(chan bool)
}
