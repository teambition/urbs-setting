package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		var err error

		if conf.Config.CertFile != "" && conf.Config.KeyFile != "" {
			logging.Info(fmt.Sprintf("server start on https://%v", conf.Config.SrvAddr))
			err = app.ListenTLS(conf.Config.SrvAddr, conf.Config.CertFile, conf.Config.KeyFile)
		} else {
			logging.Info(fmt.Sprintf("server start on http://%v", conf.Config.SrvAddr))
			err = app.Listen(conf.Config.SrvAddr)
		}

		if err != http.ErrServerClosed {
			logging.Panic(err)
		}

		logging.Warning("Server gracefully stopped")
	}()

	logging.Warning(fmt.Sprintf("HTTP service quited, received signal: %v", <-signals))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.Close(ctx); err != nil {
		logging.Err(err)
	}
}
