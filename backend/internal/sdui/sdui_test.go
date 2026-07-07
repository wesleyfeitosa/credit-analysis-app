package sdui

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// filterComponent returns the "filter" component from a screen, failing the
// test if none is present.
func filterComponent(t *testing.T, s Screen) Component {
	t.Helper()
	for _, c := range s.Components {
		if c.Type == "filter" {
			return c
		}
	}
	t.Fatal("no filter component in screen")
	return Component{}
}

func TestCreditAnalysesScreen_NoSavedFilters(t *testing.T) {
	s := CreditAnalysesScreen(nil)

	filter := filterComponent(t, s)
	assert.Nil(t, filter.Values, "filter should have no values when nothing is saved")

	// Nil RawMessage must be omitted from the JSON, not serialized as null.
	raw, err := json.Marshal(s)
	require.NoError(t, err)
	assert.NotContains(t, string(raw), `"values"`)
}

func TestCreditAnalysesScreen_WithSavedFilters(t *testing.T) {
	saved := json.RawMessage(`{"status":"APROVADO","clientName":"Maria"}`)

	s := CreditAnalysesScreen(saved)

	filter := filterComponent(t, s)
	assert.JSONEq(t, string(saved), string(filter.Values))

	// The values ride on the filter component only, never the table.
	for _, c := range s.Components {
		if c.Type == "table" {
			assert.Nil(t, c.Values, "table component must not carry values")
		}
	}
}
