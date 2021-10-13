package clux

import (
	"errors"

	"github.com/nats-io/nats.go"
	"github.com/urfave/cli/v2"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var _verbose = false

func Verbose() bool {
	return _verbose
}
func NameSpace() string {
	return ns
}

func Action(opts ...ActionOpt) *cli.App {
	cfg := &actionOpts{}
	cfg.app = cli.NewApp()
	for i := 0; i < len(opts); i++ {
		opts[i].config(cfg)
	}
	// cfg.app.Name = cfg.name
	cfg.app.Metadata = make(map[string]interface{})

	cfg.app.Metadata["ns"] = ns
	cfg.app.Flags = make([]cli.Flag, 0)
	cfg.app.Flags = append(cfg.app.Flags, &cli.StringFlag{
		Name:  "ns",
		Usage: "NameSpace",
		Value: ns,
	}, &cli.BoolFlag{
		Name:    "verbose",
		Usage:   "verbose",
		Aliases: []string{"V", "v"},
		Value:   false,
	})
	// if cfg.useNats {
	cfg.app.Flags = append(cfg.app.Flags, &cli.StringFlag{
		Name:  "nats_addr",
		Usage: "nats server addr",
		Value: "nats-headless:4222",
		// Required: true,
	},
		&cli.StringFlag{
			Name:  "nats_user",
			Usage: "nats server user",
			Value: "",
			// Required: true,
		},
		&cli.StringFlag{
			Name:  "nats_pwd",
			Usage: "nats server pwd",
			Value: "",
		},
	)
	// }
	// if cfg.useEtcd {
	cfg.app.Flags = append(cfg.app.Flags, &cli.StringSliceFlag{
		Name:  "etcd_addr",
		Usage: "etcd server addr",
		Value: cli.NewStringSlice("etcd-headless:2379"),
	},
		&cli.StringFlag{
			Name:  "etcd_pwd",
			Usage: "etcd root Password",
			Value: "",
		},
	)
	// }
	// if cfg.useDb {
	cfg.app.Flags = append(cfg.app.Flags, &cli.StringFlag{
		Name:  "db_driver",
		Usage: "Database driver",
		Value: "mysql",
	},
		&cli.StringFlag{
			Name:  "db_dsn",
			Usage: "db connect str",
			Value: "",
		},
	)
	// }
	cfg.app.Before = func(c *cli.Context) error {
		_verbose = c.Bool("verbose")
		if cfg.useNats {
			ops := make([]nats.Option, 0)
			natsUser := c.String("nats_user")
			natsPwd := c.String("nats_pwd")
			if len(natsUser) > 0 && len(natsPwd) > 0 {
				ops = append(ops, nats.UserInfo(natsUser, natsPwd))
			}
			if err := NatsInit(c.String("nats_addr"), ops...); err != nil {
				return err
			}
		}
		if cfg.useEtcd {
			opts := make([]EtcdOpt, 0)
			opts = append(opts, WithEtcdConfigEndpoints(c.StringSlice("etcd_addr")))
			if len(c.String("etcd_pwd")) > 0 {
				opts = append(opts, WithEtcdConfigUserPwd("root", c.String("etcd_pwd")))
			}
			if _cli, err := GetEtcdClient(opts...); err != nil {
				return err
			} else {
				_defaultEtcdClient = _cli
			}
		}

		if cfg.useDb {
			var err error
			switch c.String("db_driver") {
			case "mysql":
				db, err = gorm.Open(mysql.Open(c.String("db_dsn")), &gorm.Config{})
				if err != nil {
					return err
				}
				break
			case "sqlite":
				db, err = gorm.Open(sqlite.Open(c.String("db_dsn")), &gorm.Config{})
				if err != nil {
					return err
				}
				break
			default:
				return errors.New("un support driver")
			}

		}
		return nil
	}
	return cfg.app
}

type actionOpts struct {
	useNats bool
	useEtcd bool
	useDb   bool
	app     *cli.App
}

// actionOptFn is a function
type actionOptFn func(opts *actionOpts) error

func (opt actionOptFn) config(opts *actionOpts) error {
	return opt(opts)
}

// ActionOpt configures options
type ActionOpt interface {
	config(opts *actionOpts) error
}

// UseNats
func UseNats() ActionOpt {
	return actionOptFn(func(opts *actionOpts) error {
		opts.useNats = true
		return nil
	})
}

func UseEtcd() ActionOpt {
	return actionOptFn(func(opts *actionOpts) error {
		opts.useEtcd = true
		return nil
	})
}

func UseDb() ActionOpt {
	return actionOptFn(func(opts *actionOpts) error {
		opts.useDb = true
		return nil
	})
}

// Name string
func Name(name string) ActionOpt {
	return actionOptFn(func(opts *actionOpts) error {
		opts.app.Name = name
		return nil
	})
}
