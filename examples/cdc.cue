package examples

endpoints: [
	{
		description: "consumer driven contract verification"
		request: {
			method: "GET"
			path:   "/todos/1"
		}
		response: {
			status: 419
			body: """
																												{
																												  "userId": 1,
																												  "id": 1,
																												  "title": "delectus aut autem",
																												  "completed": false
																												}"""
		}
//		cdcs: [{
//			baseURL: "https://jsonplaceholder.typicode.com"
//		}]
	},
]
