package examples

basePath: "/hello/"
endpoints: [... {request: {method: *"GET" | _}}]
endpoints: [... {response: {status: *200 | _}}]

endpoints: [
	{
		request: {
			path: basePath + "world"
		}
		response: {
			body: """
				{
				  "message": "Deliberately wrong value to show the CDC still passes due to structure being correct downstream"
				}
				
				"""
		}
		cdcs: [{
			baseURL: "https://sandbox.api.service.nhs.uk/hello-world/"
		}]
	},
	{
		description: "hello pepper endpoint"
		request: {
			path: basePath + "pepper"
		}
		response: {
			body: "hello pepper!"
		}
		cdcs: [{
			baseURL: "https://sandbox.api.service.nhs.uk/hello-world/"
			ignore:  true // useful if you know an endpoint is not working but still want to run MJ successfully (it'll non 0 exit otherwise)
		}]
	},
]
