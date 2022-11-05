# 2. Use cue and define MJ config format

Date: 2022-11-05

## Status

Accepted

## Context

A big DX complaint of using MJ1 was using YAML for configuration. The whitespaces, confusing around strings, etc caused a lot of pain beyond simple configurations. 

For MJ2 we need something simpler, ideally with something that has tools to cut down the verbosity of configuration

## Decision

[cue](https://cuelang.org) will be the configuration format for MJ2. It allows:

- MJ to define a schema which it, and consumers can use to easily validate their configuration. The schema language is simple and powerful.
- A much friendlier, concise syntax for defining configuration which isn't whitespace sensitive
- An extremely powerful templating language too, so consumers of MJ2 can use cue creatively in their own context to template their configurations, DRY, simplify, etc, without any need for MJ2 to be involved. All MJ2 has to do is provide a schema.
- Easy to marshal into Go structs

### Example schema

```cue
package main

endpoints: [...#Endpoint]

#Endpoint: {
	description: string | *"\(request.method) \(request.path)"
	request: {
			method: *"GET" | "POST" | "PATCH" | "PUT" | "DELETE" | "OPTIONS" | "HEAD"
			path: string
	}
	response: {
			status: >=200 & <=599
			body: string
	}
}
```

Note how cue can:
- Constrain and validate not only the structure of the config, but the values too
- Generate a description from the rest of the configuration, if it's not provided

### Example configuration a user of MJ2 could provide

```cue
basePath: "/hello/"
endpoints: [... { request: { method: *"GET" | _}}]
endpoints: [... { response: { status: *200 | _}}]

endpoints: [
	{
			request: {
				path: basePath + "world"
			}
			response: {
				body: """
hello
world!"""
			}
	},
		{
			request: {
				path: basePath + "chris"
			}
			response: {
				body: "hello chris!"
			}
	},
	{
		description: "joke"
		request: {
			path: "tellmeajoke"
			method: "POST"
		}
		response: {
			status: 201
			body: "lmao"
		}
	}
]
```

I have only dived in to cue for a few hours but found it easy to DRY up my config, adding some default values, etc.

### Program that loads a config from a consumer and validates it using the schema

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

## Consequences

- More fun
- Making my own format obviously has its disadvantages. Interop with openapi could be useful, but I don't care that much honestly.
- If i really need to go to openapi, it shouldn't be too bad, so long as i don't tightly couple configuration from the actual workings of MJ2
