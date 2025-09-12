package config

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Setenv("ENVIRONMENT", "development")
	code := m.Run()
	os.Exit(code)
}