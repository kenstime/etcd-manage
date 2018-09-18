package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/shiguanghuxian/etcd-manage/conf"
	"github.com/shiguanghuxian/etcd-manage/routers"

	"github.com/shiguanghuxian/etcd-manage/e3ch"

	"github.com/gin-gonic/gin"
	"go.etcd.io/etcd/version"
)

const (
	PROGRAM_NAME    = "e3w"
	PROGRAM_VERSION = "0.0.2"
)

var configFilepath string

func init() {
	flag.StringVar(&configFilepath, "conf", "conf/config.default.ini", "config file path")
	rev := flag.Bool("rev", false, "print rev")
	flag.Parse()

	if *rev {
		fmt.Printf("[%s v%s]\n[etcd %s]\n",
			PROGRAM_NAME, PROGRAM_VERSION,
			version.Version,
		)
		os.Exit(0)
	}
}

func main() {
	config, err := conf.Init(configFilepath)
	if err != nil {
		panic(err)
	}

	client, err := e3ch.NewE3chClient(config)
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	router.Use(MiddleWare())
	router.UseRawPath = true
	routers.InitRouters(router, config, client)
	router.Run(":" + config.Port)
}

func MiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Origin,Content-Type, Authorization, Access-Control-Allow-Origin, X-Etcd-Username, X-Etcd-Password")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if c.Request.Method == "OPTIONS" {
			c.Status(http.StatusOK)
		}

	}
}
