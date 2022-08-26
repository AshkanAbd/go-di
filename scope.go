package dependency_injection

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"reflect"
	"sync"
)

// Scope is a struct to request services from service collection.
// Each scope has an isolated object pool for scoped services but will use a shared object pool singleton services.
// To initialize Scope struct please CreateScope method of ServiceCollection.
type Scope struct {
	// The ServiceCollection that Scope was created for.
	collection *ServiceCollection
	// The object pool that will be used to retrieve scoped services if the service was initialized before.
	// All new provided scoped services will store in this object pool for future services requests.
	scopeServicePool map[reflect.Type]any
	// A mutex to handle data race while providing or initializing scoped services.
	mutex sync.RWMutex
}

// Initialize or retrieve singleton services from ServiceCollection singleton object pool.
func provideSingletonService[T any](s *Scope, reflectType reflect.Type, serviceType *ServiceType) T {
	log.Debugf("Injecting signleton service <%v>\n", reflectType.String())

	s.collection.mutex.RLock()
	value, available := s.collection.singletonServicePool[reflectType]
	s.collection.mutex.RUnlock()

	if available {
		log.Debugf("Value retrived from singleton pool for service <%v>\n", reflectType.String())
		return value.(T)
	}

	value = serviceType.provider(s)

	s.collection.mutex.Lock()
	s.collection.singletonServicePool[reflectType] = value
	s.collection.mutex.Unlock()

	log.Debugf("Providing signleton value for service <%v>\n", reflectType.String())

	return value.(T)
}

// Initialize or retrieve scoped service from Scope object pool.
func provideScopedService[T any](s *Scope, reflectType reflect.Type, serviceType *ServiceType) T {
	log.Debugf("Injecting scoped service <%v>\n", reflectType.String())

	s.mutex.RLock()
	value, available := s.scopeServicePool[reflectType]
	s.mutex.RUnlock()

	if available {
		log.Debugf("Value retrived from scope pool for service <%v>\n", reflectType.String())
		return value.(T)
	}

	value = serviceType.provider(s)

	s.mutex.Lock()
	s.scopeServicePool[reflectType] = value
	s.mutex.Unlock()

	log.Debugf("Providing scoped value for service <%v>\n", reflectType.String())

	return value.(T)
}

// Initialize transient services.
func provideTransientService[T any](s *Scope, reflectType reflect.Type, serviceType *ServiceType) T {
	log.Debugf("Providing value to transient service <%v>\n", reflectType.String())
	return serviceType.provider(s).(T)
}

// GetService is responsible to retrieve or initialize requested service based on it's lifetime.
func GetService[T any](s *Scope) (T, error) {
	reflectType := getReflectType[T]()

	serviceType, exists := s.collection.registeredServicePool[reflectType]

	if !exists {
		var t T
		return t, fmt.Errorf("service %v is not registered in service collection", reflectType.String())
	}

	switch serviceType.lifetime {
	case SINGLETON:
		return provideSingletonService[T](s, reflectType, serviceType), nil
	case SCOPED:
		return provideScopedService[T](s, reflectType, serviceType), nil
	case TRANSIENT:
		return provideTransientService[T](s, reflectType, serviceType), nil
	}

	var t T
	return t, fmt.Errorf("invalid lifetime for service %v", reflectType.String())
}
