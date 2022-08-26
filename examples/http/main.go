package main

import (
	"fmt"
	di "github.com/ashkanabd/go-di"
	"io"
	"net/http"
	"strconv"
)

func main() {
	// Initializing service collection
	collection := di.InitServiceCollection()

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

	// Registering routes and handle functions
	http.HandleFunc("/testSingleton", func(w http.ResponseWriter, r *http.Request) {
		scope, err := collection.CreateScope()
		if err != nil {
			writeInternalError(w, fmt.Sprintf("Can't create scope: %v", err.Error()))
			return
		}

		handleTestSingleton(w, r, scope)
	})
	http.HandleFunc("/testScoped", func(w http.ResponseWriter, r *http.Request) {
		scope, err := collection.CreateScope()
		if err != nil {
			writeInternalError(w, fmt.Sprintf("Can't create scope: %v", err.Error()))
			return
		}

		handleTestScoped(w, r, scope)
	})
	http.HandleFunc("/testTransient", func(w http.ResponseWriter, r *http.Request) {
		scope, err := collection.CreateScope()
		if err != nil {
			writeInternalError(w, fmt.Sprintf("Can't create scope: %v", err.Error()))
			return
		}

		handleTestTransient(w, r, scope)
	})

	fmt.Println("Staring to listen at port 3000")
	fmt.Println("Registered routes: ")
	fmt.Println("/testSingleton")
	fmt.Println("/testScoped")
	fmt.Println("/testTransient")
	// Staring to listen
	http.ListenAndServe(":3000", nil)
}

func handleTestSingleton(w http.ResponseWriter, r *http.Request, s *di.Scope) {
	// Request a singleton service
	singletonService, _ := di.GetService[*SingletonService](s)

	io.WriteString(w, fmt.Sprintf("Pre value was: %v\n", singletonService.counter))

	valueStr := getFromQuery(r, "add", "1")

	// ignoring error
	value, _ := strconv.Atoi(valueStr)
	io.WriteString(w, fmt.Sprintf("Adding %v from query string <add>\n", value))

	// Change singleton service value
	singletonService.counter += value

	// Request singleton service again
	singletonService, _ = di.GetService[*SingletonService](s)

	io.WriteString(w, fmt.Sprintf("New value: %v\n", singletonService.counter))
}

func handleTestScoped(w http.ResponseWriter, r *http.Request, s *di.Scope) {
	// Request a scoped service
	scopedService, _ := di.GetService[*ScopedService](s)

	io.WriteString(w, fmt.Sprintf("Pre value was: %v\n", scopedService.counter))

	valueStr := getFromQuery(r, "add", "1")

	// ignoring error
	value, _ := strconv.Atoi(valueStr)
	io.WriteString(w, fmt.Sprintf("Adding %v from query string <add>\n", value))

	// Change scoped service value
	scopedService.counter += value

	// Request scoped service again
	scopedService, _ = di.GetService[*ScopedService](s)

	io.WriteString(w, fmt.Sprintf("New value: %v\n", scopedService.counter))
}

func handleTestTransient(w http.ResponseWriter, r *http.Request, s *di.Scope) {
	// Request a transient service
	transientService, _ := di.GetService[*TransientService](s)

	io.WriteString(w, fmt.Sprintf("Pre value was: %v\n", transientService.counter))

	valueStr := getFromQuery(r, "add", "1")

	// ignoring error
	value, _ := strconv.Atoi(valueStr)
	io.WriteString(w, fmt.Sprintf("Adding %v from query string <add>\n", value))

	// Change transient service value
	transientService.counter += value

	// Request transient service again
	transientService, _ = di.GetService[*TransientService](s)

	io.WriteString(w, fmt.Sprintf("New value: %v\n", transientService.counter))
}
