package mj

import _ "embed"

var (
	//go:embed schema.cue
	Schema []byte

	//go:embed fixture_schema.cue
	FixtureSchema []byte
)
