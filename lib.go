package dependency_injection

import (
	"reflect"
)

// Lifetimes
const (
	// SINGLETON lifetime represent to the services which will initialize once when requested
	// and same instance will be retrieved whenever requested again in the entire application.
	SINGLETON = iota
	// SCOPED lifetime represent to the services that will initialize in every new scope when requested
	// and same value will be retrieved in if the same scope requests it.
	// Note that requesting same service in different scopes will cause initializing new instance.
	SCOPED
	// TRANSIENT lifetime represent to the services that will initialize whenever requested.
	TRANSIENT
)

// ServiceType used to store the configuration of the services in ServiceCollection
type ServiceType struct {
	lifetime int
	provider func(s *Scope) any
}

// getReflectType returns reflect type of given generic type.
func getReflectType[T any]() reflect.Type {
	var t T

	reflectType := reflect.TypeOf(t)

	if reflectType == nil {
		var t *T
		return reflect.TypeOf(t)
	}
	return reflectType
}
