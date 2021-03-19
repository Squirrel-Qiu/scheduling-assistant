package conf

import (
	"bytes"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Listen string `toml:"Listen"`
	Mysql  mysql  `toml:"mysql"`
	Debug  bool   `toml:"debug"`
}

type mysql struct {
	Username string `toml:"Username"`
	Password string `toml:"Password"`
	Url      string `toml:"Url"`
}

func ReadConf() (url, dbUser, dbPassword, listenAddr string, debug bool) {
	file, err := os.Open("./config1.toml")
	defer file.Close()
	if err != nil {
		panic(err)
	}

	var conf Config

	buf := bytes.NewBufferString("")
	_, err = buf.ReadFrom(file)
	if err != nil {
		panic(err)
	}

	_, err = toml.Decode(buf.String(), &conf)
	if err != nil {
		panic(err)
	}

	return conf.Mysql.Url, conf.Mysql.Username, conf.Mysql.Password, conf.Listen, conf.Debug
}
