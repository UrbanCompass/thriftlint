package thriftlint

import (
	"github.com/stretchr/testify/require"

	"testing"
)

func TestChecks(t *testing.T) {
	checks := Checks{
		MakeCheck("alpha", nil),
		MakeCheck("alpha.beta", nil),
		MakeCheck("beta.gamma", nil),
		MakeCheck("beta.zeta", nil),
	}

	actual := checks.CloneAndDisable("alpha")
	expected := Checks{checks[2], checks[3]}
	require.Equal(t, expected, actual)

	actual = checks.CloneAndDisable("alpha.beta")
	expected = Checks{checks[0], checks[2], checks[3]}
	require.Equal(t, expected, actual)
}
