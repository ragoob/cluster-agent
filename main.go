package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	routers "github.com/kube-carbonara/cluster-agent/routers"
	"github.com/labstack/echo/v4"
	"github.com/rancher/remotedialer"
	"github.com/sirupsen/logrus"
)

func init() {
}

var (
	addr  string
	id    string
	debug bool
)

func handleRouting(e *echo.Echo) {
	namespacesRouter := routers.NameSpacesRouter{}
	podsRouter := routers.PodsRouter{}
	namespacesRouter.Handle(e)
	podsRouter.Handle(e)
}

func main() {
	// set by default for dev env
	if os.Getenv("SERVER_ADDRESS") == "" {
		os.Setenv("SERVER_ADDRESS", "127.0.0.1:8099")
	}

	clusterGuid := os.Getenv("CLIENT_ID")
	flag.StringVar(&addr, "connect", fmt.Sprintf("ws://%s/connect", os.Getenv("SERVER_ADDRESS")), "Address to connect to")
	flag.StringVar(&id, "id", clusterGuid, "Client ID")
	flag.BoolVar(&debug, "debug", true, "Debug logging")
	flag.Parse()

	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	headers := http.Header{
		"X-Tunnel-ID": []string{id},
	}

	e := echo.New()
	e.GET("/", func(context echo.Context) error {
		return context.String(http.StatusOK, "Hello, World!")
	})

	time.AfterFunc(5*time.Second, func() {
		remotedialer.ClientConnect(context.Background(), addr, headers, nil, func(string, string) bool { return true }, nil)
	})

	handleRouting(e)
	e.Logger.Fatal(e.Start(":1323"))
}