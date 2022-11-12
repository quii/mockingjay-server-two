package examples

endpoints: [
	{
		description: "an example without any fancy cue refactoring"
		request: {
			method: "POST"
			path:   "/fizz"
		}
		response: {
			status: 419
			body:   "buzz"
		}
	},
]
