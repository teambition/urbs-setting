package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/teambition/urbs-setting/src/api"
	"github.com/teambition/urbs-setting/src/service"
	"github.com/teambition/urbs-setting/src/util"
)

var help = flag.Bool("help", false, "show help info")
var version = flag.Bool("version", false, "show version info")
var file = flag.String("file", "", "SQL file to execute")

func main() {
	flag.Parse()
	if *help || *version || *file == "" {
		data, _ := json.Marshal(api.GetVersion())
		fmt.Println(string(data))
		os.Exit(0)
	}

	data, err := ioutil.ReadFile(*file)
	if err == nil {
		err = util.DigInvoke(func(sql *service.SQL) error {
			return sql.DB.Exec(string(data)).Error
		})
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
