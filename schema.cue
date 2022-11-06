package main

endpoints: [...#Endpoint]

#Endpoint: {
	description: string | *"\(request.method) \(request.path)"
	request: #Request
	response: #Response
}

#Request: {
			method: *"GET" | "POST" | "PATCH" | "PUT" | "DELETE" | "OPTIONS" | "HEAD"
			path: string
			headers?: [string] : [...string]
			body: string | *""
}

#Response: {
		status: >=200 & <=599
		body: string
		headers?: [string] : [...string]
}