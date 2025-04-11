package main

import (
	"flag"
	"fmt"

	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/server"
	"github.com/piupuer/go-helper/pkg/log"
)

var configPath = flag.String("f", "../../conf", "the config path")
var env = flag.String("env", "dev", "the env config")

func main() {

	flag.Parse()
	configFile := fmt.Sprintf("%s/%s_%s.yaml", *configPath, "xiaozhi-server", *env)

	err := server.NewServer(configFile, *env)
	if err != nil {
		log.Fatal(err)
	}

}
