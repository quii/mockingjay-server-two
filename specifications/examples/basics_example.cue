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
	},
]
