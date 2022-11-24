someContentType:    "application/xml"
anotherContentType: "application/json"
cdcFixtures: [... {want: {body: *"" | _, status: *200 | _}}]
cdcFixtures: [... {got: {body: *"" | _, status: *200 | _}}]

cdcFixtures: [
	{
		description:        "bodies match exactly"
		shouldBeCompatible: true
		want: {
			body: "buzz"
		}
		got: {
			body: "buzz"
		}
	},
	{
		description:        "bodies don't match (and aren't json either)"
		shouldBeCompatible: false
		want: {
			body: "fizz"
		}
		got: {
			body: "buzz"
		}
	},
	{
		description:        "status match"
		shouldBeCompatible: true
		want: {
			status: 200
		}
		got: {
			status: 200
		}
	},
	{
		description:        "status do not match"
		shouldBeCompatible: false
		want: {
			status: 200
		}
		got: {
			status: 201
		}
	},
	{
		description:        "exact headers"
		shouldBeCompatible: true
		want: {
			headers: {
				"Content-Type": [someContentType]
			}
		}
		got: {
			headers: {
				"Content-Type": [someContentType]
			}
		}
	},
	{
		description:        "superflous headers"
		shouldBeCompatible: true
		want: {
			headers: {
				"Content-Type": [someContentType]
			}
		}
		got: {
			headers: {
				"Content-Type": [anotherContentType, someContentType]
			}
		}
	},
	{
		description:        "if header exists, but with wrong value"
		shouldBeCompatible: false
		want: {
			headers: {
				"Content-Type": [someContentType]
			}
		}
		got: {
			headers: {
				"Content-Type": [anotherContentType]
			}
		}
	},
	{
		description:        "json structure is the same, but value is different"
		shouldBeCompatible: true
		want: {
			body: """
				{"message": "hello, Kat"}
				"""
		}
		got: {
			body: """
				{"message": "hello, Marc"}
				"""
		}
	},
	{
		description:        "json, but structure is different"
		shouldBeCompatible: false
		want: {
			body: """
				{"message": "hello, Kat"}
				"""
		}
		got: {
			body: """
				{"greeting": "hello, Kat"}
				"""
		}
	},
	{
		description:        "json arrays with item in response"
		shouldBeCompatible: true
		want: {
			body: """
				[{"message": "hello, Kat"}]
				"""
		}
		got: {
			body: """
				[{"message": "hello, Kat"}]
				"""
		}
	},
	{
		description:        "json arrays, but response from real system has empty array"
		shouldBeCompatible: false
		want: {
			body: """
				[{"message": "hello, Kat"}]
				"""
		}
		got: {
			body: """
				[]
				"""
		}
	},
	{
		description:        "missing keys in json response"
		shouldBeCompatible: false
		want: {
			body: """
				{"total":0,"max_score":null}
				"""
		}
		got: {
			body: """
				{"total":0}
				"""
		}
	},
	{
		description:        "extra keys in json response"
		shouldBeCompatible: true
		want: {
			body: """
				{"total":0}
				"""
		}
		got: {
			body: """
				{"total":0, "other_thing":true}
				"""
		}
	},
	{
		description:        "wrong value type in json response"
		shouldBeCompatible: false
		want: {
			body: """
				{"total":0}
				"""
		}
		got: {
			body: """
				{"total":"zero"}
				"""
		}
	},
]
