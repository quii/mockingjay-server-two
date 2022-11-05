package examples

// request body examples

endpoints: [... { request: { method: "POST", path: "/posts"}}]
endpoints: [... { response: { body: "whatever"}}]

endpoints: [
	{
		description: "match on body"
		request: {
			body: "whatever"
		}
		response: {
			status: 200
		}
	},
		{
		description: "match on a different body"
		request: {
			body: "lol"
		}
		response: {
			status: 400
		}
	},
]