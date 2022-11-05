package examples

// content negotiation examples, note the method and path for all these requests are the same, only variance are headers and response

endpoints: [... { request: { method: *"GET" | _, path: "/hello/Ruth"}}]
endpoints: [... { response: { status: *201 | _}}]

endpoints: [
	{
		description: "match on xml"
		request: {
			headers: {
				"Content-Type": ["application/xml"]
			}
		}
		response: {
			body: """
<hello>Ruth</hello>
"""
		}
	},
		{
		description: "match on json (headers not case-sensitive)"
		request: {
			headers: {
				"content-type": ["application/json"]
			}
		}
		response: {
			body: """
{"hello":"Ruth"}
"""
		}
	}
]