# go-di

#### Scope based dependency injection for golang

## Content

- [How to install](#how-to-install)
- [Features](#features)
- [How to use](#how-to-use)
- [Examples](#examples)

### How to install

```shell
go get -u github.com/ashkanabd/go-di
```

### Features:

- Easy to set up and use.
- Scope based dependency injection like `NestJS` and `.Net` frameworks.
- Able to integrate with golang popular frameworks like [gin](https://github.com/gin-gonic/gin), [fiber](https://github.com/gofiber/fiber), etc.
- Define your services once, use as much as you want.
- Supports interface registration and can provide structs that implement registered interface.
- Supports pointer registration. 
- Supports nested service resolving, regardless of lifetime (You have access to scope in providers).
- Supports goroutines, no data race issues. Implemented with `RWMutex` to optimize performance.

### How to use

#### 1. Import go-di package:

```go
import di "github.com/ashkanabd/go-di"
```

#### 2. Initialize service collection:

At fist step you need to initialize service collection:

```go
collection := di.InitServiceCollection()
```

#### 3. Register your services in service collection:

You can register your service with 3 different lifetimes:

1. `Singleton`: Services that will initialize once when requested and will retrieve whenever requested again
2. `Scoped`: Services which will initialize once in the scope that requested and will be retrieved if the same scope requests.
   But if another scope requests same service, will be initialized again
3. `Transient`: Services which will be initialized whenever requested using their providers.

If you want to register an interface, your provider function must return a pointer to the instance of a struct that implements the interface.

You can register your services with following functions:

```go
// For singleton services
error := di.AddSingleton[ServiceType/ServiceInterface](collection, func (s *di.Scope) any {
   // provider of your service
   return ServiceType{}
})
// For scoped services
error := di.AddScoped[ServiceType/ServiceInterface](collection, func (s *di.Scope) any {
   // provider of your service
   return ServiceType{}
})
// For transient services
error := di.AddTransient[ServiceType/ServiceInterface](collection, func (s *di.Scope) any {
   // provider of your service
   return ServiceType{}
})
```

#### 4. Lock your service collection:

In order to create scope from your service collection, you have to lock it to prevent adding more services while you are requesting for services.

```go
collection.Lock()
```

#### 5. Creating scope:

To get services which you registered before, you need to create a scope:

```go
scope, error := collection.CreateScope()
```

#### 6. Getting services:

Now you can get your services from dependency injection:

```go
service, error := di.GetService[ServiceType](scope)
```

### Examples

Here is implemented examples in different frameworks:

1. Integration with golang built-in http web server: [link](./examples/http)
2. Integration with [gin](https://github.com/gin-gonic/gin) framework: [link](./examples/gin)
3. Integration with [fiber](https://github.com/gofiber/fiber) framework: [link](./examples/fiber)
