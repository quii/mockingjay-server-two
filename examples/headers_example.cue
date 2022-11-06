package examples

// content negotiation examples, note the method and path for all these requests are the same, only variance are headers and response

endpoints: [... { request: { method: *"GET" | _, path: "/hello/Ruth"}}]
endpoints: [... { response: { status: *201 | _}}]

endpoints: [
	{
		description: "match on request header"
		request: {
			headers: {
				"Accept": ["application/xml"]
			}
		}
		response: {
			body: """
<hello>Ruth</hello>
"""
			headers: {
				"Content-Type": ["application/xml"]
			}
		}
	},
		{
		description: "headers are not case-sensitive"
		request: {
			headers: {
				"aCCepT": ["application/json"]
			}
		}
		response: {
			body: """
{"hello":"Ruth"}
"""
			headers: {
				"Content-Type": ["application/json"]
			}
		}
	}
]