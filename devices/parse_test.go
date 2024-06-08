package devices

import (
	"testing"

	"github.com/dermaddis/z2m_cli/sliceutil"
)

func TestAliases(t *testing.T) {
    // Check if all alias-name mappings actually map to an existing name
	for alias, realName := range deviceNameAliases {
        if !sliceutil.Contains(deviceNames, realName) {
            t.Errorf("realName %q does not exist (%s -> %s)", realName, alias, realName)
        }
	}
}
