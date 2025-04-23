// package main

// import (
// 	//"fmt"
// 	//"net/http"
// 	"log"

// 	_ "go-test/docs"
// 	"go-test/routes"

// 	"github.com/gofiber/fiber/v2"
// 	fiberSwagger "github.com/swaggo/fiber-swagger"
// )

// // func handler(w http.ResponseWriter, r *http.Request) {
// // 	fmt.Fprintf(w, "Hello, World!")
// // }

// // @title POS System API
// // @version 1.0
// // @description Swagger docs for POS system
// // @host localhost:8080
// // @BasePath /
// func main() {
// 	// http.HandleFunc("/", handler)
// 	// fmt.Println("Starting server on :8080")
// 	// err := http.ListenAndServe(":8080", nil)
// 	// if err != nil {
// 	// 	fmt.Println("Error starting server:", err)
// 	// }

// 	// Create a new Fiber app
// 	app := fiber.New()

// 	// // Define an API route
// 	// app.Get("/api", func(c *fiber.Ctx) error {
// 	// 	return c.JSON(fiber.Map{
// 	// 		"message": "Hello from Fiber API!",
// 	// 	})
// 	// })

// 	// // Define a dynamic route with parameters
// 	// app.Get("/api/user/:id", func(c *fiber.Ctx) error {
// 	// 	id := c.Params("id")
// 	// 	return c.JSON(fiber.Map{
// 	// 		"user_id": id,
// 	// 	})
// 	// })

// 	app.Get("/swagger/*", fiberSwagger.WrapHandler) // Swagger route

// 	routes.Register(app)

// 	// Start the server on port 3000
// 	log.Fatal(app.Listen(":8080"))
// }

package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"go-test/proto"
	"io/ioutil"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var jwtKey = []byte("your-secret-key")

func main() {
	app := fiber.New()

	app.Post("/login", func(c *fiber.Ctx) error {
		type LoginInput struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		var input LoginInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
		}

		cert, err := tls.LoadX509KeyPair("client.crt", "client.key")
		if err != nil {
			log.Fatalf("Failed to load client cert/key: %v", err)
		}

		//Load CA cert
		caCert, err := ioutil.ReadFile("ca.crt")
		if err != nil {
			log.Fatalf("Failed to load CA cert: %v", err)
		}

		caPool := x509.NewCertPool()
		if !caPool.AppendCertsFromPEM(caCert) {
			log.Fatal("Failed to add CA cert to pool")
		}

		//TLS config for mTLS
		tlsConfig := &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caPool,
			InsecureSkipVerify: false,
			ServerName:         "127.0.0.1",
		}
		creds := credentials.NewTLS(tlsConfig)

		fmt.Println("Dialing gRPC server with mTLS...")
		conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(creds))
		if err != nil {
			log.Printf("gRPC dial error: %v", err)
			return c.Status(500).JSON(fiber.Map{"error": "Auth service unreachable"})
		}
		defer conn.Close()

		client := proto.NewUserAuthServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		resp, err := client.Login(ctx, &proto.LoginRequest{
			Username: input.Username,
			Password: input.Password,
		})
		if err != nil {
			log.Printf("gRPC Login error: %v", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Login failed. Please check your credentials or try again later.",
			})
		}

		return c.JSON(fiber.Map{
			"message": "Login successful",
			"token":   resp.Token,
		})
	})

	app.Get("/inventory", func(c *fiber.Ctx) error {
		tokenStr := c.Get("Authorization")
		if tokenStr == "" {
			return c.Status(401).SendString("Missing token")
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method")
			}
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			return c.Status(401).SendString("Unauthorized")
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			fmt.Println("Authenticated User:", claims["username"])
		}

		return c.JSON(fiber.Map{
			"inventory": []string{"Item1", "Item2", "Item3"},
		})
	})

	log.Fatal(app.Listen(":3000"))
}
