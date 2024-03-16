package main

import (
	"log"
	"net"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

type Service struct {
	address string
	network string
	timeout time.Duration
}

var services = []Service{
	{
		address: "localhost:8081",
		network: "tcp",
		timeout: 5 * time.Second,
	},
	{
		address: "localhost:8082",
		network: "tcp",
		timeout: 5 * time.Second,
	},
	{
		address: "localhost:8083",
		network: "tcp",
		timeout: 5 * time.Second,
	},
}

var currentServiceIndex = -1

func main() {
	router := gin.Default()
	getNextServiceURL := func() *url.URL {
		return findNextAvailableService()
	}

	router.GET("/usuarios", func(c *gin.Context) {
		proxyRequest(c, getNextServiceURL())
	})
	router.GET("/usuarios/:id", func(c *gin.Context) {
		proxyRequest(c, getNextServiceURL())
	})
	router.POST("/usuarios", func(c *gin.Context) {
		proxyRequest(c, getNextServiceURL())
	})
	router.GET("/usuarios/login", func(c *gin.Context) {
		proxyRequest(c, getNextServiceURL())
	})
	router.DELETE("/usuarios/:id", func(c *gin.Context) {
		proxyRequest(c, getNextServiceURL())
	})
	router.PUT("/usuarios/:id", func(c *gin.Context) {
		proxyRequest(c, getNextServiceURL())
	})

	log.Println("Starting gateway on :8080")
	router.Run(":8080")
}

func findNextAvailableService() *url.URL {
	currentServiceIndex = (currentServiceIndex + 1) % len(services)
	if services[currentServiceIndex].check() {
		parsedURL, _ := url.Parse("http://" + services[currentServiceIndex].address)
		return parsedURL
	}
	return findNextAvailableService()
}

func (service Service) check() bool {
	conn, err := net.DialTimeout(service.network, service.address, service.timeout)
	if err != nil {
		// fmt.Printf("Failed to connect to %s: %v\n", service.address, err)
		return false
	}
	defer conn.Close()

	return true
}

func proxyRequest(c *gin.Context, targetURL *url.URL) {
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.ServeHTTP(c.Writer, c.Request)
}
