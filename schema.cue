package main

endpoints: [...#Endpoint]

#Endpoint: {
	description: string | *"\(request.method) \(request.path)"
	request: {
			method: *"GET" | "POST" | "PATCH" | "PUT" | "DELETE" | "OPTIONS" | "HEAD"
			path: string
			headers?: [string] : [...string]
			body: string | *""
	}
	response: {
			status: >=200 & <=599
			body: string
	}
}