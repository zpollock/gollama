package main

import (
	"flag"
	"gollama/cmd/routes"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	host = flag.String("host", "0.0.0.0", "Set the IP address to listen")
	port = flag.Int("port", 8081, "Set the port to listen")
)

func main() {
	flag.Parse()
	r := gin.Default()

	routes.SetupRoutes(r)

	addr := *host + ":" + strconv.Itoa(*port)
	r.Run(addr)
}
