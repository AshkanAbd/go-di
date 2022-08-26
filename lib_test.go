package dependency_injection

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

type TestType struct {
	counter int
}

type TestInterface interface {
	GetCounter() int
}

func (t *TestType) GetCounter() int {
	return t.counter
}

func TestGetSingletonService(t *testing.T) {
	collection := InitServiceCollection()

	counter := 0

	err := AddSingleton[TestType](collection, func(s *Scope) any {
		counter++
		return TestType{
			counter: counter,
		}
	})
	assert.Nil(t, err)

	collection.Lock()

	{
		// Create first scope
		scope1, err := collection.CreateScope()
		assert.Nil(t, err)

		// Retrieve services for first time in scope-1
		testType1, err := GetService[TestType](scope1)
		assert.Nil(t, err)
		assert.Equal(t, 1, testType1.counter)

		// Retrieve services again in scope-1
		testType2, err := GetService[TestType](scope1)
		assert.Nil(t, err)
		assert.Equal(t, testType1, testType2)
		assert.Equal(t, 1, testType1.counter)
		assert.Equal(t, 1, testType2.counter)
	}
	{
		// Create second scope
		scope2, err := collection.CreateScope()
		assert.Nil(t, err)

		// Retrieve services for first time in scope-2
		testType1, err := GetService[TestType](scope2)
		assert.Nil(t, err)
		assert.Equal(t, 1, testType1.counter)

		// Retrieve services again in scope-2
		testType2, err := GetService[TestType](scope2)
		assert.Nil(t, err)
		assert.Equal(t, testType1, testType2)
		assert.Equal(t, 1, testType1.counter)
		assert.Equal(t, 1, testType2.counter)
	}
}

func TestGetScopedService(t *testing.T) {
	collection := InitServiceCollection()
	counter := 0

	err := AddScoped[TestType](collection, func(s *Scope) any {
		counter++
		return TestType{
			counter: counter,
		}
	})
	assert.Nil(t, err)

	collection.Lock()
	var firstScopeTest1 TestType
	var firstScopeTest2 TestType

	{
		// Create first scope:
		scope1, err := collection.CreateScope()
		assert.Nil(t, err)

		// Retrieve services for first time in scope-1
		firstScopeTest1, err = GetService[TestType](scope1)
		assert.Nil(t, err)
		assert.Equal(t, 1, firstScopeTest1.counter)

		// Retrieve services again in scope-1
		firstScopeTest2, err = GetService[TestType](scope1)
		assert.Nil(t, err)
		assert.Equal(t, firstScopeTest1, firstScopeTest2)
		assert.Equal(t, 1, firstScopeTest1.counter)
		assert.Equal(t, 1, firstScopeTest2.counter)
	}
	{
		// Create second scope
		scope2, err := collection.CreateScope()
		assert.Nil(t, err)

		// Retrieve services for first time in scope-2
		testType1, err := GetService[TestType](scope2)
		assert.Nil(t, err)
		assert.Equal(t, 2, testType1.counter)

		// Retrieve services again in scope-2
		testType2, err := GetService[TestType](scope2)
		assert.Nil(t, err)
		assert.Equal(t, testType1, testType2)
		assert.NotEqual(t, testType1, firstScopeTest1)
		assert.NotEqual(t, testType2, firstScopeTest2)

		assert.Equal(t, testType1.counter, firstScopeTest1.counter+1)
		assert.Equal(t, testType2.counter, firstScopeTest2.counter+1)

		assert.Equal(t, 2, testType1.counter)
		assert.Equal(t, 2, testType2.counter)
	}
}

func TestGetTransientService(t *testing.T) {
	collection := InitServiceCollection()

	counter := 0

	err := AddTransient[TestType](collection, func(s *Scope) any {
		counter++
		return TestType{
			counter: counter,
		}
	})

	assert.Nil(t, err)

	collection.Lock()
	var firstScopeTest1 TestType
	var firstScopeTest2 TestType

	{
		// Create first scope:
		scope1, err := collection.CreateScope()
		assert.Nil(t, err)

		// Retrieve services for first time in scope-1
		firstScopeTest1, err = GetService[TestType](scope1)
		assert.Nil(t, err)
		assert.Equal(t, 1, firstScopeTest1.counter)

		// Retrieve services again in scope-1
		firstScopeTest2, err = GetService[TestType](scope1)
		assert.Nil(t, err)
		assert.NotEqual(t, firstScopeTest1, firstScopeTest2)
		assert.Equal(t, 1, firstScopeTest1.counter)
		assert.Equal(t, 2, firstScopeTest2.counter)
	}
	{
		// Create second scope
		scope2, err := collection.CreateScope()
		assert.Nil(t, err)

		// Retrieve services for first time in scope-2
		testType1, err := GetService[TestType](scope2)
		assert.Nil(t, err)
		assert.Equal(t, 3, testType1.counter)

		// Retrieve services again in scope-2
		testType2, err := GetService[TestType](scope2)
		assert.Nil(t, err)
		assert.NotEqual(t, testType1, testType2)
		assert.NotEqual(t, testType1, firstScopeTest1)
		assert.NotEqual(t, testType1, firstScopeTest2)
		assert.NotEqual(t, testType2, firstScopeTest1)
		assert.NotEqual(t, testType2, firstScopeTest2)

		assert.Equal(t, testType1.counter, firstScopeTest1.counter+2)
		assert.Equal(t, testType2.counter, firstScopeTest2.counter+2)

		assert.Equal(t, 3, testType1.counter)
		assert.Equal(t, 4, testType2.counter)
	}
}

func TestAddServiceToLockedCollection(t *testing.T) {
	collection := InitServiceCollection()

	err := AddSingleton[TestType](collection, func(s *Scope) any {
		return TestType{
			counter: 0,
		}
	})
	assert.Nil(t, err)

	collection.Lock()

	err = AddSingleton[TestType](collection, func(s *Scope) any {
		return TestType{
			counter: -1,
		}
	})
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Errorf("service collection is locked, you can't register an other service"), err)
}

func TestCreateScopeOnUnlockServiceCollection(t *testing.T) {
	collection := InitServiceCollection()

	err := AddSingleton[TestType](collection, func(s *Scope) any {
		return TestType{
			counter: 0,
		}
	})

	assert.Nil(t, err)

	scope, err := collection.CreateScope()
	assert.Nil(t, scope)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Errorf("you have to lock service collection to create a scope"), err)
}

func TestGetNotRegisteredService(t *testing.T) {
	collection := InitServiceCollection()

	collection.Lock()

	scope, err := collection.CreateScope()
	assert.Nil(t, err)

	_, err = GetService[TestType](scope)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Errorf("service dependency_injection.TestType is not registered in service collection"), err)
}

func TestRegisterInterface(t *testing.T) {
	collection := InitServiceCollection()

	counter := 0

	err := AddScoped[TestInterface](collection, func(s *Scope) any {
		counter++
		return &TestType{
			counter: counter,
		}
	})
	assert.Nil(t, err)

	collection.Lock()
	var firstScopeTest1 TestInterface
	var firstScopeTest2 TestInterface

	{
		// Create first scope:
		scope1, err := collection.CreateScope()
		assert.Nil(t, err)

		// Retrieve services for first time in scope-1
		firstScopeTest1, err = GetService[TestInterface](scope1)
		assert.Nil(t, err)
		assert.Equal(t, 1, firstScopeTest1.GetCounter())

		// Retrieve services again in scope-1
		firstScopeTest2, err = GetService[TestInterface](scope1)
		assert.Nil(t, err)
		assert.Equal(t, firstScopeTest1, firstScopeTest2)
		assert.Equal(t, 1, firstScopeTest1.GetCounter())
		assert.Equal(t, 1, firstScopeTest2.GetCounter())
	}
	{
		// Create second scope
		scope2, err := collection.CreateScope()
		assert.Nil(t, err)

		// Retrieve services for first time in scope-2
		testType1, err := GetService[TestInterface](scope2)
		assert.Nil(t, err)
		assert.Equal(t, 2, testType1.GetCounter())

		// Retrieve services again in scope-2
		testType2, err := GetService[TestInterface](scope2)
		assert.Nil(t, err)
		assert.Equal(t, testType1, testType2)
		assert.NotEqual(t, testType1, firstScopeTest1)
		assert.NotEqual(t, testType2, firstScopeTest2)

		assert.Equal(t, testType1.GetCounter(), firstScopeTest1.GetCounter()+1)
		assert.Equal(t, testType2.GetCounter(), firstScopeTest2.GetCounter()+1)

		assert.Equal(t, 2, testType1.GetCounter())
		assert.Equal(t, 2, testType2.GetCounter())
	}
}

func TestGetServiceWithInvalidLifetime(t *testing.T) {
	serviceReflectType := getReflectType[TestType]()
	collection := InitServiceCollection()

	provider := func(s *Scope) any {
		return TestType{
			counter: 0,
		}
	}

	collection.registeredServicePool[serviceReflectType] = &ServiceType{
		lifetime: -1,
		provider: provider,
	}

	collection.Lock()

	scope, err := collection.CreateScope()
	assert.Nil(t, err)

	value, err := GetService[TestType](scope)
	assert.NotNil(t, value)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Errorf("invalid lifetime for service dependency_injection.TestType"), err)
}

func TestDataRaceInSingleton(t *testing.T) {
	wg := sync.WaitGroup{}
	collection := InitServiceCollection()

	err := AddSingleton[TestType](collection, func(s *Scope) any {
		return TestType{
			counter: 0,
		}
	})
	assert.Nil(t, err)

	collection.Lock()

	scopeCount := 10
	scopeSlice := make([]*Scope, scopeCount)

	for i := 0; i < scopeCount; i++ {
		scope, err := collection.CreateScope()
		assert.Nil(t, err)
		assert.NotNil(t, scope)
		scopeSlice[i] = scope
	}

	scope, err := collection.CreateScope()
	assert.Nil(t, err)
	assert.NotNil(t, scope)

	wg.Add(scopeCount * 5)

	f := func(s *Scope) {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			value, err := GetService[TestType](s)
			assert.Nil(t, err)
			assert.NotNil(t, value)
		}
	}

	for i := 0; i < scopeCount; i++ {
		go f(scopeSlice[i])
		go f(scopeSlice[i])
		go f(scopeSlice[i])
		go f(scopeSlice[i])
		go f(scopeSlice[i])
	}

	wg.Wait()
}

func TestDataRaceInScoped(t *testing.T) {
	wg := sync.WaitGroup{}
	collection := InitServiceCollection()

	err := AddScoped[TestType](collection, func(s *Scope) any {
		return TestType{
			counter: 0,
		}
	})
	assert.Nil(t, err)

	collection.Lock()

	scopeCount := 10
	scopeSlice := make([]*Scope, scopeCount)

	for i := 0; i < scopeCount; i++ {
		scope, err := collection.CreateScope()
		assert.Nil(t, err)
		assert.NotNil(t, scope)
		scopeSlice[i] = scope
	}

	wg.Add(scopeCount * 5)

	f := func(s *Scope) {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			value, err := GetService[TestType](s)
			assert.Nil(t, err)
			assert.NotNil(t, value)
		}
	}

	for i := 0; i < scopeCount; i++ {
		go f(scopeSlice[i])
		go f(scopeSlice[i])
		go f(scopeSlice[i])
		go f(scopeSlice[i])
		go f(scopeSlice[i])
	}

	wg.Wait()
}

func TestDataRaceInTransient(t *testing.T) {
	wg := sync.WaitGroup{}
	collection := InitServiceCollection()

	err := AddTransient[TestType](collection, func(s *Scope) any {
		return TestType{
			counter: 0,
		}
	})
	assert.Nil(t, err)

	collection.Lock()

	scopeCount := 10
	scopeSlice := make([]*Scope, scopeCount)

	for i := 0; i < scopeCount; i++ {
		scope, err := collection.CreateScope()
		assert.Nil(t, err)
		assert.NotNil(t, scope)
		scopeSlice[i] = scope
	}

	wg.Add(scopeCount * 5)

	f := func(s *Scope) {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			value, err := GetService[TestType](s)
			assert.Nil(t, err)
			assert.NotNil(t, value)
		}
	}

	for i := 0; i < scopeCount; i++ {
		go f(scopeSlice[i])
		go f(scopeSlice[i])
		go f(scopeSlice[i])
		go f(scopeSlice[i])
		go f(scopeSlice[i])
	}

	wg.Wait()
}
