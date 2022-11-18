package examples

// request body examples

endpoints: [... {request: {method: "POST", path: "/posts"}}]
endpoints: [... {response: {body: "whatever"}}]

endpoints: [
	{
		description: "match on request body"
		request: {
			body: "whatever"
		}
		response: {
			status: 200
		}
	},
	{
		description: "match on a different request body"
		request: {
			body: "lol"
		}
		response: {
			status: 400
		}
	},
]
