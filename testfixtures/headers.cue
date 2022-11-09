theRequest: {
	method: "GET"
	path: "/hello/ruth"
}

fixtures: [... {endpoint:{ response: { status: 201 , body: "<hello>Ruth</hello>"}}}]
someContentType: "application/xml"
anotherContentType: "application/json"

fixtures: [
	{
		endpoint: {
			description: "matching of request headers"
			request: theRequest & {
				headers: {
					"Accept": [someContentType]
				}
			}
		}
		matchingRequests: [
			{
				description: "order of headers and extra headers dont matter",
				request: theRequest & {
					headers: {
						"Accept": [anotherContentType, someContentType]
					}
				}
			}
		]
		nonMatchingRequests: [
			{
				description: "wont match if header isn't present",
				request: theRequest & {
					headers: {
						"Accept": [anotherContentType]
					}
				}
			}
		]
	}
]