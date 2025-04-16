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
	"log"
	"time"

	"go-test/proto" // Update this import path to match your project

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
)

func main() {
	// Initialize Fiber app
	app := fiber.New()

	// Route for fetching inventory
	app.Get("/inventory", func(c *fiber.Ctx) error {
		username := c.Query("username")
		password := c.Query("password")

		if !authenticateUser(username, password) {
			return c.Status(401).JSON(fiber.Map{
				"message": "Authentication failed",
			})
		}

		return c.JSON(fiber.Map{
			"inventory": []string{"Item1", "Item2", "Item3"},
		})
	})

	app.Listen(":3000")
}

func authenticateUser(username, password string) bool {
	// Connect to the User Authentication Service via gRPC
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := proto.NewUserAuthServiceClient(conn)

	// Make the gRPC call to the UserAuthService.Login method
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &proto.LoginRequest{
		Username: username,
		Password: password,
	}

	resp, err := client.Login(ctx, req)
	if err != nil {
		log.Fatalf("could not login: %v", err)
	}

	return resp.GetSuccess()
}
