package examples

endpoints: [
	{
		  description: "an example fixture"
			request: {
				method: "POST"
				path: "/fizz"
			}
			response: {
				status: 419
				body: "buzz"
			}
	}
]