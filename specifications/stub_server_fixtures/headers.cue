baseRequest: {
	method: "GET"
	path:   "/hello/ruth"
}

someContentType:    "application/xml"
anotherContentType: "application/json"

fixtures: [
	{
		endpoint: {
			description: "matching of request headers"
			request:     baseRequest & {
				headers: {
					"Accept": [someContentType]
				}
			}
			response: {
				status: 201
				body:   "<hello>Ruth</hello>"
			}
		}
		matchingRequests: [
			{
				description: "order of headers and extra headers dont matter"
				request:     baseRequest & {
					headers: {
						"Accept": [anotherContentType, someContentType]
					}
				}
			},
		]
		nonMatchingRequests: [
			{
				description: "wont match if header isn't present"
				request:     baseRequest & {
					headers: {
						"Accept": [anotherContentType]
					}
				}
			},
		]
	},
]
