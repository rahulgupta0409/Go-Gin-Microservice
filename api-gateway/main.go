package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// APIKey - Replace this with your actual API key
const APIKey = "your-api-key"

var UserMicroservice = "http://localhost:9000"
var ProductMicroservice = "http://localhost:8082"

func main() {
	apiGatewayURL := "http://localhost:8086"
	proxy1 := createReverseProxy(UserMicroservice)
	proxy2 := createReverseProxy(ProductMicroservice)

	r := gin.Default()

	// Add middleware to authenticate and rate limit requests
	//r.Use(authenticate())
	// Add middleware to rate limit requests
	r.Use(rateLimit())
	r.Use(CORSMiddleware())

	// *path is a wildcard parameter in Gin. It's a path parameter
	//that matches any number of URL segments in the request path.
	//It's used to capture the rest of the request path
	//after /microservice1/, including any additional URL segments
	//or trailing slashes.
	r.Any("/microservice1/*path", proxy1)
	r.Any("/microservice2/*path", proxy2)

	log.Printf("API Gateway listening on %s", apiGatewayURL)
	log.Fatal(r.Run(":8086"))
}

func createReverseProxy(targetURL string) func(*gin.Context) {
	target, err := url.Parse(targetURL)
	if err != nil {
		log.Fatalf("Error parsing target URL: %v", err)
	}

	// Create a reverse proxy for the target URL. It creates a new
	//reverse proxy that forwards requests to a single target host.
	//The function takes a *url.URL as an argument, which represents
	//the target host to which the requests should be forwarded.
	proxy := httputil.NewSingleHostReverseProxy(target)

	return func(c *gin.Context) {
		log.Printf("Request received for %s", c.Request.URL.Path)

		// Modify the request path to exclude the microservice prefix
		c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/microservice1")
		c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/microservice2")

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

// authenticate - Middleware to authenticate requests
func authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey != APIKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"error": "Unauthorized"})
			return
		}
		c.Next()
	}
}

// rateLimit - Middleware to rate limit requests
func rateLimit() gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Every(1*time.Minute), 5)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests,
				gin.H{"error": "Too many requests"})
			return
		}
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
