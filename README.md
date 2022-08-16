# varlog-parser

This project is a basic utility exposing a RESTful HTTP service to list entries from a log file. Given a preconfigured directory on start, the service will accept a file name and a few parameters to list the latest entries of that file.

## Usage

### Starting the Server

Once you compile the binary, you may pass the following flags:
```text
Usage of varlogd:
  -httpPort int
        The port on which the http server will listen. Default is 8080. (default 8080)
  -logPath /var/log
        Tells the service where to look for requested files. Default is /var/log. (default "/var/log")
```

### API Documentation

API Documentation is written in OpenAPI3, and is located in the `/api` directory of this project. You can copy / import this file into a live editor, such as [Swagger's Online Editor](https://editor.swagger.io/), and see more information about the endpoints, parameters and response types. 

## Build and Developer Instructions

### Precursor

Before beginning, the following needs to be installed:
- [Golang](https://go.dev/doc/install)
- [GolangCI-Lint](https://golangci-lint.run/usage/install/)
- `make`

### Build and Run

Once the Go and golangci-lint are installed, you should be able to run `make compile` which will attempt to run a few additional targets. This will generate a `varlogd` binary in the root of the project. This can be called directly by referring to the options in the Usage section above.

### Additional Targets

For your convenience, all necessary operations are centralized into a `Makefile` located in the root of the project. See that file's comments for specific information.

### API Development

This project uses [OpenAPI v3](https://spec.openapis.org/oas/v3.1.0), [Deepmap's OpenAPI Code generator](https://github.com/deepmap/oapi-codegen), and [Chi](https://go-chi.io/#/) to manage API and REST Endpoint development. 

Workflow:

1. All changes to any part of the REST API must start from the OpenAPI specs in the projects `/api` directory. 
2. Run `version=v1 make generate-api` which will call the code generator and overwrite the `api.gen.go` file(s). 
3. Make required changes to the code that is implementing the server endpoints.  

## TODOs, FIXMEs, and Wishes

1. Add additional query parameters to support paging.
2. Add fields to the response to indicate total records, filtered records, and current page.
3. Benchmark testing for data structure supporting file parsing. Currently uses a very dumb / naive implementation of a RingBuffer based on a string slice. Strings and slices can be burdensome on the garbage collector, and might be alleviated using some combination of doubly-linked list and sync.Pool.
4. Provide a Docker file for build to remove the need to setup Go and other tools.
5. Provide a Docker file for deploy so it can be deployed as a container more easily.
6. Automated API testing using the OpenAPI spec to generate a Client.
