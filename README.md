# API Rate Limiter - Vincent Gonzalez

## Runtime Environment
The minimum version of Go that is required to run this solution is `1.18`. This version is also defined in the `go.mod` file of this solution. However, I have compiled and run this solution using version `1.19.2` of Go.

## How to run
A makefile has been included in the base directory. Use `make buildRun` in a terminal in the base directory to build and execute the program. Use `make build` to only compile the program. Use `make run` to run the program after it has already been compiled.

### Makefile
Contains the CLI commands for the project. It also contains some environment variables that are used within the CLI commands.

The file includes the following commands:
- `build` - changes to the source directory, compiles the program, and places the executable into a build directory.
- `buildRun` - combines the `build` and `run` commands into one.
- `run` - executes the compiled program. The following flags allow the user to customize the rate limit and request limit of the rate limiter when using the test harness code:
    - `-rateLimit` sets the rate interval in milliseconds
    - `-requestLimit` sets the maximum number of requests per rate interval
- `loadTest` - a command that will run a Vegeta load test and is provided for convenience

## Dependencies
This program uses the following libraries and tools:
### Vegeta
[Vegeta](https://github.com/tsenart/vegeta) is used to load test the rate limiter. It is not required to run the program. It is included as a convenience only.
```
brew install vegeta
```
### golang.org/x/time/rate
```
go get golang.org/x/time/rate
```

## Files
### main.go
Contains the test harness code and testing request middleware.
### rate_limiter.go
Contains the definition of the rate limiter `RateLimiter`. A factory function `NewRateLimiter` creates new instances of `RateLimiter` and is included in this file.
### connection.go
Contains a helper struct `Connection` that defines a request to be handled by the rate limiter.
