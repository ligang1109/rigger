package rconf

import (
	"errors"
	"github.com/goinbox/golog"
	"github.com/goinbox/gomisc"
	"strconv"
)

type MkdirItem struct {
	Dir  string
	Mask string
	Sudo bool
}

type ActionConf struct {
	Mkdir []*MkdirItem
	Exec  []string

	confPath string
	logger   golog.ILogger
}

func NewActionConf(confPath string, logger golog.ILogger) (*ActionConf, error) {
	if !gomisc.FileExist(confPath) {
		return nil, errors.New("actionConf not exist: " + confPath)
	}

	return &ActionConf{
		confPath: confPath,
		logger:   logger,
	}, nil
}

func (a *ActionConf) Parse(vc *VarConf) error {
	err := gomisc.ParseJsonFile(a.confPath, a)
	if err != nil {
		return err
	}

	var delay bool
	for i, item := range a.Mkdir {
		pkeyPrefix := "parse action mkdir " + strconv.Itoa(i) + " item "

		a.logger.Debug([]byte(pkeyPrefix + "dir: " + item.Dir))
		item.Dir, delay, err = vc.ParseValueByDefined(item.Dir)
		if delay {
			err = errors.New("must not delay")
		}
		if err != nil {
			return errors.New(pkeyPrefix + "dir: " + item.Dir + " error: " + err.Error())
		}

		a.logger.Debug([]byte(pkeyPrefix + "mask: " + item.Mask))
		item.Mask, delay, err = vc.ParseValueByDefined(item.Mask)
		if delay {
			err = errors.New("must not delay")
		}
		if err != nil {
			return errors.New(pkeyPrefix + "mask: " + item.Mask + " error: " + err.Error())
		}
	}
	for i, cmd := range a.Exec {
		pkeyPrefix := "parse action exec " + strconv.Itoa(i) + " cmd "

		a.logger.Debug([]byte(pkeyPrefix + "cmd: " + cmd))
		a.Exec[i], delay, err = vc.ParseValueByDefined(cmd)
		if delay {
			err = errors.New("must not delay")
		}
		if err != nil {
			return errors.New(pkeyPrefix + "cmd: " + cmd + " error: " + err.Error())
		}
	}

	return nil
}
