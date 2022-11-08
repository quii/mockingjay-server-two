someContentType: "application/xml"
anotherContentType: "application/json"

fixtures: [
	{
		endpoint: {
			description: "matching of request headers"
			request: {
				method: "GET"
				path: "/"
				headers: {
					"Accept": [someContentType]
				}
			}
			response: {
				status: 200
				body: "whatever"
			}
		}
		matchingRequests: [
			{
				description: "order of headers and extra headers dont matter",
				request: {
					method: "GET"
					path: "/"
					headers: {
						"Accept": [anotherContentType, someContentType]
					}
				}
			}
		]
		nonMatchingRequests: [
			{
				description: "wont match if header isn't present",
				request: {
					method: "GET"
					path: "/"
					headers: {
						"Accept": [anotherContentType]
					}
				}
			}
		]
	}
]