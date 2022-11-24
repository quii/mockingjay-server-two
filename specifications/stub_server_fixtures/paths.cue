baseRequest: {
	method: "GET"
}

baseResponse: {
	status: 200
	body:   "hello, world"
}

fixtures: [
	{
		endpoint: {
			description: "matching of paths"
			request:     baseRequest & {
				path: "/hello/world"
			}
			response: baseResponse
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
	{
		endpoint: {
			description: "matching of regex"
			request:     baseRequest & {
				path:      "/happy-birthday/elodie"
				regexPath: "\/happy-birthday\/[a-z]+"
			}
			response: baseResponse
		}
		matchingRequests: [
			{
				description: "matches regex"
				request:     baseRequest & {
					path: "/happy-birthday/milo"
				}
			},
			{
				description: "matches another one"
				request:     baseRequest & {
					path: "/happy-birthday/sarah"
				}
			},
		]
		nonMatchingRequests: [
			{
				description: "path doesn't match regex"
				request:     baseRequest & {
					path: "/bonjour/milo"
				}
			},
		]
	},
]
