package dependency_injection

import (
	"fmt"
	"reflect"
	"sync"
)

// ServiceCollection is collection of services to use them in your application.
// This struct contains two pools, first one for registered services and second one for provided singleton services.
// You may need to initialize ServiceCollection only one time in your application
type ServiceCollection struct {
	// Pool of registered services
	registeredServicePool map[reflect.Type]*ServiceType
	// A shared object pool between different scopes to collect provided singleton services.
	// This pool will be used to retrieve singleton objects if they are initialized before.
	// All new initialized singleton services will be stored in this object pool.
	singletonServicePool map[reflect.Type]any
	// Lock of service collection
	locked bool
	// A mutex to handle data race while providing or initializing singleton services.
	mutex sync.RWMutex
}

// InitServiceCollection initialize a service collection
func InitServiceCollection() *ServiceCollection {
	return &ServiceCollection{
		registeredServicePool: make(map[reflect.Type]*ServiceType),
		singletonServicePool:  make(map[reflect.Type]any),
		locked:                false,
		mutex:                 sync.RWMutex{},
	}
}

// AddSingleton registers a service as singleton
func AddSingleton[T any](collection *ServiceCollection, provider func(s *Scope) any) error {
	if err := collection.checkLock(); err != nil {
		return err
	}

	collection.add(getReflectType[T](), SINGLETON, provider)

	return nil
}

// AddScoped registers a service as scoped
func AddScoped[T any](collection *ServiceCollection, provider func(s *Scope) any) error {
	if err := collection.checkLock(); err != nil {
		return err
	}

	collection.add(getReflectType[T](), SCOPED, provider)

	return nil
}

// AddTransient registers a service as transient
func AddTransient[T any](collection *ServiceCollection, provider func(scope *Scope) any) error {
	if err := collection.checkLock(); err != nil {
		return err
	}

	collection.add(getReflectType[T](), TRANSIENT, provider)

	return nil
}

// Add a service to service collection with given lifetime and provider
func (collection *ServiceCollection) add(t reflect.Type, lifetime int, provider func(scope *Scope) any) {
	collection.registeredServicePool[t] = &ServiceType{
		lifetime: lifetime,
		provider: provider,
	}
}

// checks lock of the service collection to avoid adding more services when application starts to work
func (collection *ServiceCollection) checkLock() error {
	if collection.locked {
		return fmt.Errorf("service collection is locked, you can't register an other service")
	}

	return nil
}

// Lock preparing service collection to create scope
func (collection *ServiceCollection) Lock() {
	collection.locked = true
}

// CreateScope creates a new scope in application to retrieve services
func (collection *ServiceCollection) CreateScope() (*Scope, error) {
	if !collection.locked {
		return nil, fmt.Errorf("you have to lock service collection to create a scope")
	}

	return &Scope{
		collection:       collection,
		scopeServicePool: make(map[reflect.Type]any),
		mutex:            sync.RWMutex{},
	}, nil
}
