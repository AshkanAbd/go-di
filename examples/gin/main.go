package main

import (
	"fmt"
	di "github.com/ashkanabd/go-di"
	"github.com/gin-gonic/gin"
	"net/http"
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

	// Initialize gin
	r := gin.Default()

	// Register create scope middleware to create scope on incoming requests
	r.Use(CreateScopeMiddleware)

	// Register routes on gin
	r.GET("/testSingleton", handleSingleton)
	r.GET("/testScoped", handleScoped)
	r.GET("/testTransient", handleTransient)

	r.Run(":3000")
}

func CreateScopeMiddleware(c *gin.Context) {
	// Create scope before request
	scope, err := collection.CreateScope()

	if err != nil {
		c.Writer.WriteString(fmt.Sprintf("Can't create scope: %v\n", err.Error()))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.Set("scope", scope)

	c.Next()
}

func handleSingleton(ctx *gin.Context) {
	// Get scope from gin context
	scope := ctx.MustGet("scope").(*di.Scope)

	// Request a singleton service
	singletonService, _ := di.GetService[*SingletonService](scope)

	ctx.Writer.WriteString(fmt.Sprintf("Pre value was: %v\n", singletonService.counter))

	valueStr := ctx.Query("add")

	// ignoring error
	value, _ := strconv.Atoi(valueStr)
	ctx.Writer.WriteString(fmt.Sprintf("Adding %v from query string <add>\n", value))

	// Change singleton service value
	singletonService.counter += value

	// Request singleton service again
	singletonService, _ = di.GetService[*SingletonService](scope)

	ctx.Writer.WriteString(fmt.Sprintf("New value: %v\n", singletonService.counter))
}

func handleScoped(ctx *gin.Context) {
	// Get scope from gin context
	scope := ctx.MustGet("scope").(*di.Scope)

	// Request a scoped service
	scopedService, _ := di.GetService[*ScopedService](scope)

	ctx.Writer.WriteString(fmt.Sprintf("Pre value was: %v\n", scopedService.counter))

	valueStr := ctx.Query("add")

	// ignoring error
	value, _ := strconv.Atoi(valueStr)
	ctx.Writer.WriteString(fmt.Sprintf("Adding %v from query string <add>\n", value))

	// Change scope service value
	scopedService.counter += value

	// Request scoped service again
	scopedService, _ = di.GetService[*ScopedService](scope)

	ctx.Writer.WriteString(fmt.Sprintf("New value: %v\n", scopedService.counter))
}

func handleTransient(ctx *gin.Context) {
	// Get scope from gin context
	scope := ctx.MustGet("scope").(*di.Scope)

	// Request a transient service
	transientService, _ := di.GetService[*TransientService](scope)

	ctx.Writer.WriteString(fmt.Sprintf("Pre value was: %v\n", transientService.counter))

	valueStr := ctx.Query("add")

	// ignoring error
	value, _ := strconv.Atoi(valueStr)
	ctx.Writer.WriteString(fmt.Sprintf("Adding %v from query string <add>\n", value))

	// Change transient service value
	transientService.counter += value

	// Request transient service again
	transientService, _ = di.GetService[*TransientService](scope)

	ctx.Writer.WriteString(fmt.Sprintf("New value: %v\n", transientService.counter))
}
