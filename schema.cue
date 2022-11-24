package mockingjay

// As a user of MJ, you'll only need to provide endpoints
endpoints?: [...#Endpoint]

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

// The other schema are for tests, which you might want to do if you wish to file an issue
stubFixtures?: [...#StubFixture]
cdcFixtures?: [...#CDCFixture]

#CDCFixture: {
	description:        string
	shouldBeCompatible: bool
	got:                #Response
	want:               #Response
}

#StubFixture: {
	endpoint: #Endpoint
	matchingRequests: [...#RequestDescription]
	nonMatchingRequests: [...#RequestDescription]
}

#RequestDescription: {
	request:     #Request
	description: string
}
