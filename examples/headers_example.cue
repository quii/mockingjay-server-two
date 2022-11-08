package examples

// content negotiation examples, note the method and path for all these requests are the same, only variance are headers and response

endpoints: [... { request: { method: *"GET" | _, path: "/hello/Ruth"}}]
endpoints: [... { response: { status: *201 | _}}]
someContentType: "application/xml"
anotherContentType: "application/json"

endpoints: [
	{
		description: "match on request header"
		request: {
			headers: {
				"Accept": [someContentType]
			}
		}
		response: {
			body: """
<hello>Ruth</hello>
"""
			headers: {
				"Content-Type": [someContentType]
			}
		}
	},
		{
		description: "headers are not case-sensitive"
		request: {
			headers: {
				"aCCepT": [anotherContentType]
			}
		}
		response: {
			body: """
{"hello":"Ruth"}
"""
			headers: {
				"Content-Type": [anotherContentType]
			}
		}
	}
]