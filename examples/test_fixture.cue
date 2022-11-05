package examples

basePath: "/hello/"
endpoints: [... { request: { method: *"GET" | _}}]
endpoints: [... { response: { status: *200 | _}}]

endpoints: [
	{
			request: {
				path: basePath + "world"
			}
			response: {
				body: "hello, world!"
			}
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