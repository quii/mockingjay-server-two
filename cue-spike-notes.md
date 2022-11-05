# cue spike notes

## goal

see how good cue is for defining MJ configuration

- cue is a superset of JSON
- load cue into structs and validate against a schema ✅
- the config file itself is way more powerful than json, it has templating

## resources

- [Getting started](https://cuelang.org/docs/install/)

## notes

### working thing
```cue
package mj
#Request: {
	description: string
	method: *"GET" | "POST" | "PATCH" | "PUT" | "DELETE"
	path: string
}

#Response: {
	status: int
	body: string
}

#Endpoint: {
	request: #Request
	response: #Response
}

#Server: [...#Endpoint]

#Server & [ #Endpoint & {
	request: #Request & {
		description: "hello, world"
		method: "GET"
		path: "/hello-world"
	}
	response: #Response & {
		status: 200
		body: "hello world!"
	}
}]
```

Run `cue vet` to validate data with schema (try changing field names). `cue export` spits it to JSON

### next working thing

```cue
// schema.cue

[...#Endpoint]

#Endpoint: {
	request: {
			description: string
			method: *"GET" | "POST" | "PATCH" | "PUT" | "DELETE"
			path: string
	}
	response: {
			status: int
			body: string
	}
}
```

validate a potential mj config file

```cue
[
	{
			request: {
				description: "hello, world"
				method: "GET"
				path: "/hello-world"
			}
			response: {
				status: 200
				body: "hello world!"
			}
	}
]
```

Validate with `cue vet example_config.cue schema.cue -c`

## In go

```go
package main

import (
	_ "embed"
	"log"

	"github.com/cue-exp/cueconfig"
)

type Response struct {
	Status int
	Body   string
}

type Request struct {
	Method string
	Path   string
}

type Endpoint struct {
	Description string
	Request     Request
	Response    Response
}

type Endpoints struct {
	Endpoints []Endpoint
}

var (
	//go:embed schema.cue
	schema []byte
)

func main() {
	var endpoints Endpoints

	if err := cueconfig.Load("endpoints.cue", schema, nil, nil, &endpoints); err != nil {
		log.Fatal(err)
	}

	log.Println(endpoints)
}
```

This gets validated by the schema, which is quite swish. 