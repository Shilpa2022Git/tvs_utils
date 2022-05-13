package utils

import (
	"fmt"
	"sync"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Db       database             `toml:"database"`
	Srv      server               `toml:"server"`
	Comm     servicecommunication `toml:"servicecommunication"`
	AuthInfo authkey              `toml:"authkey"`
}

type database struct {
	Server   string
	Port     string
	Name     string
	User     string
	Password string
}

type server struct {
	Port string
}

type servicecommunication struct {
	Port string
}

type authkey struct {
	Secrekey string
}

var conf *Config
var lock = &sync.Mutex{}

func NewConfig() *Config {
	if conf == nil {
		lock.Lock()
		defer lock.Unlock()

		if conf == nil {

			if _, err := toml.DecodeFile("./infrastructure/config.toml", &conf); err != nil {
				fmt.Println(err)
			}

			fmt.Printf("%#v\n", conf)
			fmt.Println("db value \n", conf.Db.Port)
			return conf
		} else {
			return conf
		}
	}

	return conf
}
