package smoke

import (
	"encoding/json"
	"testing"

	"github.com/breathbath/goalert/expflag"
	"github.com/breathbath/goalert/test/smoke/harness"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestExperimentalFlag_None tests the GraphQL API when no experimental flags are set.
func TestExperimentalFlag_None(t *testing.T) {
	t.Parallel()

	h := harness.NewHarness(t, "", "")
	defer h.Close()

	var resp struct {
		ExperimentalFlags expflag.FlagSet
	}

	err := json.Unmarshal(h.GraphQLQuery2(`{experimentalFlags}`).Data, &resp)
	require.NoError(t, err)
	assert.Len(t, resp.ExperimentalFlags, 0)
}

// TestExperimentalFlag_Example tests the GraphQL API when the example experimental flag is set.
func TestExperimentalFlag_Example(t *testing.T) {
	t.Parallel()

	h := harness.NewHarnessWithFlags(t, "", "", expflag.FlagSet{expflag.Example})
	defer h.Close()

	var resp struct {
		ExperimentalFlags expflag.FlagSet
	}

	err := json.Unmarshal(h.GraphQLQuery2(`{experimentalFlags}`).Data, &resp)
	require.NoError(t, err)
	require.Len(t, resp.ExperimentalFlags, 1)
	assert.True(t, resp.ExperimentalFlags.Has(expflag.Example))
}
