package config

import "testing"

func TestInit(t *testing.T) {
	Init()
	t.Log(GetGlobalConfig())
}
