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
	"fmt"
	"go-test/proto"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
)

var jwtKey = []byte("your-secret-key")

func main() {
	app := fiber.New()

	// Login route - gets JWT from Auth Service
	app.Post("/login", func(c *fiber.Ctx) error {
		type LoginInput struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		var input LoginInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
		}

		// Connect to user auth service
		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Printf("gRPC dial error: %v", err)
			return c.Status(500).JSON(fiber.Map{"error": "Auth service unreachable"})
		}
		defer conn.Close()

		client := proto.NewUserAuthServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		resp, err := client.Login(ctx, &proto.LoginRequest{
			Username: input.Username,
			Password: input.Password,
		})
		if err != nil || !resp.Success {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
		}

		return c.JSON(fiber.Map{
			"message": "Login successful",
			"token":   resp.Token,
		})
	})

	// Inventory route - JWT protected
	app.Get("/inventory", func(c *fiber.Ctx) error {
		tokenStr := c.Get("Authorization")

		if tokenStr == "" {
			return c.Status(401).SendString("Missing token")
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return jwtKey, nil
		})

		if err != nil {
			fmt.Println("Token parse error:", err)
			c.Status(401).SendString("Unauthorized")
			return c.Status(401).SendString("Unauthorized")
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			fmt.Println("Authenticated User:", claims["username"])
			// Proceed to the next middleware or handler
		} else {
			c.Status(401).SendString("Unauthorized")
			return c.Status(401).SendString("Unauthorized")
		}

		return c.JSON(fiber.Map{
			"inventory": []string{"Item1", "Item2", "Item3"},
		})
	})

	log.Fatal(app.Listen(":3000"))
}
