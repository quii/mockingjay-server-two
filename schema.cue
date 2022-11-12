package mockingjay

endpoints: [...#Endpoint]

#Endpoint: {
	description: string | *"\(request.method) \(request.path)"
	request:     #Request
	response:    #Response
	cdcs?: [#CDC]
}

#Request: {
	method:     *"GET" | "POST" | "PATCH" | "PUT" | "DELETE" | "OPTIONS" | "HEAD"
	path:       string
	regexPath?: string | ""
	headers?: [string]: [...string]
	body: string | *""
}

#Response: {
	status: >=200 & <=599
	body:   string
	headers?: [string]: [...string]
}

#CDC: {
	baseURL:   string
	retries:   int | *0
	timeoutMS: int | *5000
}
