package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/teambition/urbs-setting/src/api"
	"github.com/teambition/urbs-setting/src/conf"
	"github.com/teambition/urbs-setting/src/logging"
)

var help = flag.Bool("help", false, "show help info")
var version = flag.Bool("version", false, "show version info")

func main() {
	flag.Parse()
	if *help || *version {
		data, _ := json.Marshal(api.GetVersion())
		fmt.Println(string(data))
		os.Exit(0)
	}

	if len(conf.Config.SrvAddr) == 0 {
		conf.Config.SrvAddr = ":8081"
	}

	app := api.NewApp()
	ctx := conf.Config.GlobalCtx
	host := "http://" + conf.Config.SrvAddr
	if conf.Config.CertFile != "" && conf.Config.KeyFile != "" {
		host = "https://" + conf.Config.SrvAddr
	}
	logging.Infof("Urbs-Setting start on %s", host)
	logging.Errf("Urbs-Setting closed %v", app.ListenWithContext(
		ctx, conf.Config.SrvAddr, conf.Config.CertFile, conf.Config.KeyFile))
}
