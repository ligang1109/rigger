package rconf

import (
	"github.com/goinbox/golog"
)

type RiggerConf struct {
	VarConfItem    *VarConf
	TplConfItem    *TplConf
	ActionConfItem *ActionConf

	logger golog.ILogger
}

func NewRiggerConf(rconfDir string, extArgs map[string]string, logger golog.ILogger) (*RiggerConf, error) {
	vc, err := NewVarConf(rconfDir+"/var.json", extArgs, logger)
	if err != nil {
		return nil, err
	}

	tc, err := NewTplConf(rconfDir+"/tpl.json", logger)
	if err != nil {
		return nil, err
	}

	ac, err := NewActionConf(rconfDir+"/action.json", logger)
	if err != nil {
		return nil, err
	}

	rf := &RiggerConf{
		VarConfItem:    vc,
		TplConfItem:    tc,
		ActionConfItem: ac,

		logger: logger,
	}

	return rf, nil
}

func (r *RiggerConf) Parse() error {
	r.logger.Debug([]byte("parse var conf"))
	err := r.VarConfItem.Parse()
	if err != nil {
		return err
	}

	r.logger.Debug([]byte("parse tpl conf"))
	err = r.TplConfItem.Parse(r.VarConfItem)
	if err != nil {
		return err
	}

	r.logger.Debug([]byte("parse action conf"))
	err = r.ActionConfItem.Parse(r.VarConfItem)
	if err != nil {
		return err
	}

	return nil
}
