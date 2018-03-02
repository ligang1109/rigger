package rconf

import (
	"github.com/goinbox/golog"

	"fmt"
	"os"
	"testing"
)

func TestRconf(t *testing.T) {
	prjHome := os.Getenv("GOPATH") + "/demo"
	rconfDir := prjHome + "/conf/rigger"
	logger, _ := golog.NewSimpleLogger(golog.NewStdoutWriter(), golog.LEVEL_DEBUG, golog.NewConsoleFormater())
	extArgs := map[string]string{
		"prj_home": prjHome,
	}

	riggerConf, err := NewRiggerConf(rconfDir, extArgs, logger)
	if err != nil {
		t.Error(err)
	}

	err = riggerConf.Parse()
	if err != nil {
		t.Error(err)
	}

	for key, value := range riggerConf.VarConfItem.Vars {
		fmt.Println("var conf", key, value)
	}

	for key, item := range riggerConf.TplConfItem.Tpls {
		fmt.Println("tpl conf", key, item)
	}

	for i, item := range riggerConf.ActionConfItem.Mkdir {
		fmt.Println("action conf mkdir", i, item)
	}

	for i, cmd := range riggerConf.ActionConfItem.Exec {
		fmt.Println("action conf exec", i, cmd)
	}
}
