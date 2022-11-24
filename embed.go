package mj

import _ "embed"

var (
	//go:embed schema.cue
	Schema []byte
)
