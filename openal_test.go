package main

import "testing"

func TestEnvelope(t *testing.T) {
	thing := defEnv.toFloat(500)
	_ = thing
}
