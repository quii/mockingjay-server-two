baseRequest: {
	method: "GET"
}

fixtures: [
	{
		endpoint: {
			description: "matching of paths"
			request:     baseRequest & {
				path: "/hello/world"
			}
			response: {
				status: 200
				body:   "hello, world"
			}
		}
		matchingRequests: [
			{
				description: "path matches"
				request:     baseRequest & {
					path: "/hello/world"
				}
			},
		]
		nonMatchingRequests: [
			{
				description: "path is different"
				request:     baseRequest & {
					path: "/hello/mars"
				}
			},
		]
	},
]
