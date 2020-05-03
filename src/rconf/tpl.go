package rconf

import (
	"github.com/goinbox/golog"
	"github.com/goinbox/gomisc"

	"errors"
)

type TplItem struct {
	Tpl  string
	Dst  string
	Ln   string
	Sudo bool
}

type TplConf struct {
	Tpls map[string]*TplItem

	confPath string
	logger   golog.ILogger
}

func NewTplConf(confPath string, logger golog.ILogger) (*TplConf, error) {
	if !gomisc.FileExist(confPath) {
		return nil, errors.New("tplConf not exist: " + confPath)
	}

	return &TplConf{
		Tpls: make(map[string]*TplItem),

		confPath: confPath,
		logger:   logger,
	}, nil
}

func (t *TplConf) Parse(vc *VarConf) error {
	err := gomisc.ParseJsonFile(t.confPath, &t.Tpls)
	if err != nil {
		return err
	}

	var delay bool
	for key, item := range t.Tpls {
		pkeyPrefix := "parse tpl: " + key + " item "

		t.logger.Debug([]byte(pkeyPrefix + "tpl: " + item.Tpl))
		item.Tpl, delay, err = vc.ParseValueByDefined(item.Tpl)
		if delay {
			err = errors.New("must not delay")
		}
		if err != nil {
			return errors.New(pkeyPrefix + "tpl: " + item.Tpl + " error: " + err.Error())
		}

		t.logger.Debug([]byte(pkeyPrefix + "dst: " + item.Dst))
		item.Dst, delay, err = vc.ParseValueByDefined(item.Dst)
		if delay {
			err = errors.New("must not delay")
		}
		if err != nil {
			return errors.New(pkeyPrefix + "dst: " + item.Dst + " error: " + err.Error())
		}

		t.logger.Debug([]byte(pkeyPrefix + "ln: " + item.Ln))
		item.Ln, delay, err = vc.ParseValueByDefined(item.Ln)
		if delay {
			err = errors.New("must not delay")
		}
		if err != nil {
			return errors.New(pkeyPrefix + "ln: " + item.Ln + " error: " + err.Error())
		}
	}

	return nil
}
