someContentType:    "application/xml"
anotherContentType: "application/json"

fixtures: [
	{
		description: "when statuses and bodies match exactly, they are compatible"
		want: {
			status: 419
			body:   "buzz"
		}
		got: {
			status: 419
			body:   "buzz"
		}
		shouldBeCompatible: true
	},
	{
		description: "exact headers mean compatible"
		want : {
			status: 200
			headers: {
					"Content-Type": [someContentType]
			}
		}
		got : {
			status: 200
			headers: {
					"Content-Type": [someContentType]
			}
		}
		shouldBeCompatible: true
	},
		{
		description: "if header exists, but with wrong value, its not compatible"
		want : {
			status: 200
			headers: {
					"Content-Type": [someContentType]
			}
		}
		got : {
			status: 200
			headers: {
					"Content-Type": [anotherContentType]
			}
		}
		shouldBeCompatible: false
	},
	{
		description: "json structure is the same, but value is different"
		want : {
			status: 200
			body: """
			{"message": "hello, Kat"}
			"""
		}
		got : {
			status: 200
			body: """
			{"message": "hello, Marc"}
			"""
		}
		shouldBeCompatible: true
	}
]