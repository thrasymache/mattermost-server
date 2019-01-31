package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDefaultsFromStruct(t *testing.T) {
	s := struct {
		TestSettings struct {
			IntValue    int
			BoolValue   bool
			StringValue string
		}
		PointerToTestSettings *struct {
			Value int
		}
	}{}

	defaults := getDefaultsFromStruct(s)

	assert.Equal(t, defaults["TestSettings.IntValue"], 0)
	assert.Equal(t, defaults["TestSettings.BoolValue"], false)
	assert.Equal(t, defaults["TestSettings.StringValue"], "")
	assert.Equal(t, defaults["PointerToTestSettings.Value"], 0)
	assert.NotContains(t, defaults, "PointerToTestSettings")
	assert.Len(t, defaults, 4)
}
