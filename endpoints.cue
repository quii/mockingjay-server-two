endpoints: [
	{
			description: "hello, world"
			request: {
				method: "GET"
				path: "/hello-world"
			}
			response: {
				status: 200
				body: """
hello
world!"""
			}
	}
]