import "strings"

basePath: "/hello/"

// set some defaults
endpoints: [... { request: { method: *"GET" | _}}]
endpoints: [... { response: { status: *200 | _}}]

endpoints: [
	{
			description: "hello, world"
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
			description: "hello, chris"
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