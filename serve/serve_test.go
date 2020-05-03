package serve

import (
	"github.com/BurntSushi/toml"
	"testing"
)

func TestSetupConfig(t *testing.T) {
	var c SConfig
	if _, err := toml.DecodeFile("../testdata/serve-config-sample.toml", &c); err != nil {
		t.Error(err)
		return
	}

	t.Logf("%#v", c)
}
