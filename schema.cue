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