import "strings"

basePath: "/hello/"
defaultStatus: 200

// this will set a default method if it is not set
endpoints: [... { request: { method: *"GET" | _}}]

endpoints: [
	{
			description: "hello, world"
			request: {
				path: basePath + "world"
			}
			response: {
				status: defaultStatus
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
				status: defaultStatus
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