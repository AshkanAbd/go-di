package main

import (
	"fmt"
	di "github.com/ashkanabd/go-di"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"strconv"
)

var (
	collection *di.ServiceCollection
)

func main() {
	// Initializing service collection
	collection = di.InitServiceCollection()

	// Registering services
	_ = di.AddSingleton[*SingletonService](collection, func(s *di.Scope) any {
		return &SingletonService{
			counter: 0,
		}
	})
	_ = di.AddScoped[*ScopedService](collection, func(s *di.Scope) any {
		return &ScopedService{
			counter: 0,
		}
	})
	_ = di.AddTransient[*TransientService](collection, func(s *di.Scope) any {
		return &TransientService{
			counter: 0,
		}
	})

	// Locking service collection
	collection.Lock()

	// Initialize fiber app
	fiberApp := fiber.New()
	fiberApp.Use(logger.New())

	// Register create scope middleware in fiber
	fiberApp.Use(CreateScopeMiddleware)

	// Register routes
	fiberApp.Get("/testSingleton", handleSingleton)
	fiberApp.Get("/testScoped", handleScoped)
	fiberApp.Get("/testTransient", handleTransient)

	fiberApp.Listen(":3000")
}

func CreateScopeMiddleware(c *fiber.Ctx) error {
	// Create scope before request
	scope, err := collection.CreateScope()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			SendString(fmt.Sprintf("Can't create scope: %v\n", err.Error()))
	}

	c.Locals("scope", scope)
	return c.Next()
}

func handleSingleton(ctx *fiber.Ctx) error {
	// Get scope from fiber locals
	scope := ctx.Locals("scope").(*di.Scope)

	// Request a singleton service
	singletonService, _ := di.GetService[*SingletonService](scope)

	ctx.WriteString(fmt.Sprintf("Pre value was: %v\n", singletonService.counter))

	valueStr := ctx.Query("add")

	// ignoring error
	value, _ := strconv.Atoi(valueStr)
	ctx.WriteString(fmt.Sprintf("Adding %v from query string <add>\n", value))

	// Change singleton service value
	singletonService.counter += value

	// Request singleton service again
	singletonService, _ = di.GetService[*SingletonService](scope)

	ctx.WriteString(fmt.Sprintf("New value: %v\n", singletonService.counter))

	return ctx.SendStatus(fiber.StatusOK)
}

func handleScoped(ctx *fiber.Ctx) error {
	// Get scope from fiber locals
	scope := ctx.Locals("scope").(*di.Scope)

	// Request a scoped service
	scopedService, _ := di.GetService[*ScopedService](scope)

	ctx.WriteString(fmt.Sprintf("Pre value was: %v\n", scopedService.counter))

	valueStr := ctx.Query("add")

	// ignoring error
	value, _ := strconv.Atoi(valueStr)
	ctx.WriteString(fmt.Sprintf("Adding %v from query string <add>\n", value))

	// Change scope service value
	scopedService.counter += value

	// Request scoped service again
	scopedService, _ = di.GetService[*ScopedService](scope)

	ctx.WriteString(fmt.Sprintf("New value: %v\n", scopedService.counter))

	return ctx.SendStatus(fiber.StatusOK)
}

func handleTransient(ctx *fiber.Ctx) error {
	// Get scope from fiber locals
	scope := ctx.Locals("scope").(*di.Scope)

	// Request a transient service
	transientService, _ := di.GetService[*TransientService](scope)

	ctx.WriteString(fmt.Sprintf("Pre value was: %v\n", transientService.counter))

	valueStr := ctx.Query("add")

	// ignoring error
	value, _ := strconv.Atoi(valueStr)
	ctx.WriteString(fmt.Sprintf("Adding %v from query string <add>\n", value))

	// Change transient service value
	transientService.counter += value

	// Request transient service again
	transientService, _ = di.GetService[*TransientService](scope)

	ctx.WriteString(fmt.Sprintf("New value: %v\n", transientService.counter))

	return ctx.SendStatus(fiber.StatusOK)
}
