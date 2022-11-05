// set some defaults
basePath: "/hello/"
endpoints: [... { request: { method: *"GET" | _}}]
endpoints: [... { response: { status: *200 | _}}]

endpoints: [
	{
			request: {
				path: basePath + "world"
			}
			response: {
				body: """
hello
world!"""
			}
	},
		{
			request: {
				path: basePath + "chris"
			}
			response: {
				body: "hello chris!"
			}
	},
	{
		description: "joke"
		request: {
			path: "tellmeajoke"
			method: "POST"
		}
		response: {
			status: 201
			body: "lmao"
		}
	}
]