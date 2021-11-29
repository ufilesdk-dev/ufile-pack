package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var (
	G_Config Config
)

type Config struct {
	Http      ConfigHTTP `json:"http"` //对外http 服务
	Log       ConfigLog  `json:"log"`  //默认日志配置
	US3Config ConfigSDK  `json:"us3_config"`
}

type ConfigHTTP struct {
	Ip      string `json:"Ip"`
	Port    int    `json:"Port"`
	timeout uint32 `json:"timeout"`
}

type ConfigLog struct {
	LogDir    string `json:"LogDir"`
	LogPrefix string `json:"LogPrefix"`
	LogSuffix string `json:"LogSuffix"`
	LogSize   int64  `json:"LogSize"`
	LogLevel  string `json:"LogLevel"`
}

type ConfigSDK struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

func LoadConf(conf_file string) error {

	b, e := ioutil.ReadFile(conf_file)
	if e != nil {
		fmt.Println("read file error")
		return e
	}

	err := json.Unmarshal([]byte(b), &G_Config)
	if err != nil {
		fmt.Println("config file error: ", err)
		return err
	}
	var out bytes.Buffer
	json.Indent(&out, b, "", "\t")
	fmt.Printf("Config: \n%s\n", string(out.Bytes()))

	return nil
}
