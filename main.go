package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	uflog "ufile-pack/gosdk/log"
	"ufile-pack/logic"
	"ufile-pack/model"
)

func loadConfig(path string) (err error) {
	err = model.LoadConf(path)
	if err != nil {
		return err
	}
	dir := model.G_Config.Log.LogDir
	prefix := model.G_Config.Log.LogPrefix
	suffix := model.G_Config.Log.LogSuffix
	logSize := model.G_Config.Log.LogSize
	logLevel := model.G_Config.Log.LogLevel
	uflog.InitLogger(dir, prefix, suffix, logSize, logLevel)
	if model.G_Config.US3Config.PrivateKey == "" || model.G_Config.US3Config.PublicKey == "" {
		return errors.New("PublicKey or PrivateKey is Null！")
	}
	return nil
}

func main() {
	err := loadConfig("server_conf.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	http.HandleFunc("/", logic.HttpRouter)
	fmt.Println("service us3-pack-server ip:", model.G_Config.Http.Ip, ", port:", model.G_Config.Http.Port)
	address := net.JoinHostPort(model.G_Config.Http.Ip, strconv.Itoa(model.G_Config.Http.Port))
	server := &http.Server{Addr: address}
	server.SetKeepAlivesEnabled(true)
	server.ListenAndServe()

	c := make(chan os.Signal)
	//监听所有信号
	signal.Notify(c)
	fmt.Println("start!")
	s := <-c
	fmt.Println("stop,signal : ", s)

}
