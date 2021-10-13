package clux

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

var _defaultJsClient nats.JetStreamContext
var _defaultNatsClient *nats.Conn

func NatsDefaultConfig() []nats.Option {
	return []nats.Option{
		nats.PingInterval(5 * time.Second),
		nats.MaxPingsOutstanding(3),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Println("Got disconnected! Reason: ", err)
		}),
		nats.MaxReconnects(60),
		// nats.ReconnectWait()
		nats.ReconnectHandler(func(nc *nats.Conn) {
			if nc.Reconnects > 58 {
				panic("too much times")
			}
			log.Println("Got reconnected to", "[", nc.Reconnects, "]", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			log.Println("Connection closed. Reason: ", nc.LastError())
		}),
		nats.DiscoveredServersHandler(func(nc *nats.Conn) {
			log.Println("Discover closed. Reason: ", nc.ConnectedAddr(), nc.ConnectedUrl())
		}),
		nats.CustomReconnectDelay(func(attempts int) time.Duration {
			return time.Second * 2
		}),
	}
}

func NatsInit(url string, options ...nats.Option) error {
	log.Println("connect to ", url)
	arr := make([]nats.Option, 0)
	arr = append(arr, NatsDefaultConfig()...)
	arr = append(arr, options...)
	var err error
	_defaultNatsClient, err = nats.Connect(url, options...)
	if err != nil {
		return err
	}

	_defaultJsClient, err = _defaultNatsClient.JetStream()
	if err != nil {
		return err
	}
	return nil
}

func NatsJs() nats.JetStreamContext {
	return _defaultJsClient
}

func NatsClient() *nats.Conn {
	return _defaultNatsClient
}

func NatsCreateStream(cfg *nats.StreamConfig) error {
	// Check if the ORDERS stream already exists; if not, create it.
	if _defaultJsClient == nil {
		return errors.New("please InitNats")
	}
	stream, err := _defaultJsClient.StreamInfo(cfg.Name)
	if err != nil && !strings.Contains(err.Error(), "stream not found") {
		return err
	}
	//stream not found
	if stream == nil {
		log.Printf("creating stream %q and subjects %q", cfg.Name, cfg.Subjects)
		_, err = _defaultJsClient.AddStream(cfg)
		if err != nil {
			return err
		}
	} else {
		_, err = _defaultJsClient.UpdateStream(cfg)
		if err != nil {
			return err
		}
	}
	return nil
}

//CloseNats CloseNats
func NatsClose() {
	if _defaultNatsClient == nil {
		return
	}
	_defaultNatsClient.Drain()
	_defaultNatsClient.Close()
}
