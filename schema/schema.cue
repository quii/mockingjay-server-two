[...#Endpoint]

#Endpoint: {
	description: string
	request: {
			method: *"GET" | "POST" | "PATCH" | "PUT" | "DELETE"
			path: string
	}
	response: {
			status: >=200 & <=599
			body: string
	}
}