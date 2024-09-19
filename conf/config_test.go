package conf_test

import (
	"os"
	"testing"

	"github.com/fubieliangpu/WorkOrderDeployment/conf"
)

func TestToYAML(t *testing.T) {
	t.Log(conf.Default().ToYAML())
}

func TestToLoadFromYAML(t *testing.T) {
	err := conf.LoadConfigFromYaml("../apps/rcdevice/impl/user.yml")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(conf.C().ToYAML())
}

func TestToLoadFromEnv(t *testing.T) {
	os.Setenv("DATASOURCE_USERNAME", "env test")
	err := conf.LoadConfigFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(conf.C().ToYAML())
}
