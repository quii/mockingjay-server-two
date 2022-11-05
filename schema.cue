package main

endpoints: [...#Endpoint]

#Endpoint: {
	description: string | *"\(request.method) \(request.path)"
	request: {
			method: *"GET" | "POST" | "PATCH" | "PUT" | "DELETE"
			path: string
	}
	response: {
			status: >=200 & <=599
			body: string
	}
}