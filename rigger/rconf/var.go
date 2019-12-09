package rconf

import (
	"github.com/goinbox/golog"
	"github.com/goinbox/gomisc"
	"github.com/goinbox/shell"

	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	FUNC_VALUE_ARG  = "ARG"
	FUNC_VALUE_MATH = "MATH"
)

var definedValueRegex = regexp.MustCompile("\\${([^}]+)}")
var funcValueRegex = regexp.MustCompile("__([A-Z]+)__\\(([^)]+)\\)")
var mathFuncRegex = regexp.MustCompile("([0-9]+)([+\\-*/])([0-9]+)")

type VarConf struct {
	Vars map[string]string

	confPath       string
	user           string
	extArgs        map[string]string
	unparsedValues map[string]string

	logger golog.ILogger
}

func NewVarConf(confPath string, extArgs map[string]string, logger golog.ILogger) (*VarConf, error) {
	if !gomisc.FileExist(confPath) {
		return nil, errors.New("varConf not exist: " + confPath)
	}

	return &VarConf{
		Vars: make(map[string]string),

		confPath:       confPath,
		user:           os.Getenv("USER"),
		extArgs:        extArgs,
		unparsedValues: make(map[string]string),

		logger: logger,
	}, nil
}

func (v *VarConf) Parse() error {
	var varConfJson map[string]interface{}

	err := gomisc.ParseJsonFile(v.confPath, &varConfJson)
	if err != nil {
		return err
	}

	for key, item := range varConfJson {
		vs, err := v.parseVarJsonItemtoString(item)
		if err != nil {
			return err
		}

		v.unparsedValues[key] = vs
	}

	for len(v.unparsedValues) > 0 {
		for key, value := range v.unparsedValues {
			value, delay, err := v.ParseValueByDefined(value)
			if err != nil {
				return err
			}

			if !delay {
				vs, err := v.parseValueByFunc(value)
				if err != nil {
					return err
				}

				v.Vars[key] = vs
				delete(v.unparsedValues, key)
			}
		}
	}

	return nil
}

func (v *VarConf) parseVarJsonItemtoString(item interface{}) (string, error) {
	var r string

	switch item.(type) {
	case string:
		r = item.(string)
	case map[string]interface{}:
		mv := item.(map[string]interface{})
		v, ok := mv[v.user]
		if !ok {
			v = mv["default"]
		}
		r = v.(string)
	default:
		return "", errors.New("item's type not support")
	}

	return strings.TrimSpace(r), nil
}

/**
* return parsed value and whether delay parsed, if delay parsed, bool is true
 */
func (v *VarConf) ParseValueByDefined(value string) (string, bool, error) {
	matches := definedValueRegex.FindAllStringSubmatch(value, -1)

	if len(matches) == 0 {
		return value, false, nil
	}

	var rs []string
	for _, item := range matches {
		k := item[1]
		vs, ok := v.Vars[k]
		if ok {
			v.logger.Debug([]byte("find var k:" + k + " in vars"))
		} else {
			find := false
			vs, find = v.parseValueByEnv(k)
			if find {
				v.logger.Debug([]byte("find var k:" + k + " in env"))
			} else {
				_, ok := v.unparsedValues[k]
				if ok {
					return "", true, nil
				}
				return "", false, errors.New("Undefined field: " + k)
			}
		}

		rs = append(rs, item[0])
		rs = append(rs, vs)
	}

	return strings.NewReplacer(rs...).Replace(value), false, nil
}

func (v *VarConf) parseValueByEnv(name string) (string, bool) {
	vs, find := os.LookupEnv(name)
	if find {
		return vs, true
	}

	vs = string(shell.RunCmd("echo $" + name).Output)
	vs = strings.TrimSpace(vs)
	if vs == "" {
		return "", false
	}

	return vs, true
}

func (v *VarConf) parseValueByFunc(value string) (string, error) {
	match := funcValueRegex.FindStringSubmatch(value)
	var err error

	if len(match) != 0 {
		switch match[1] {
		case FUNC_VALUE_ARG:
			v.logger.Debug([]byte("parse value: " + value + " by arg func"))
			value, err = v.parseByArgFunc(match[2])
		case FUNC_VALUE_MATH:
			v.logger.Debug([]byte("parse value: " + value + " by math func"))
			value, err = v.parseByMathFunc(match[2])
		default:
			err = errors.New("Not support func " + match[1])
		}
	}

	return value, err
}

func (v *VarConf) parseByArgFunc(argName string) (string, error) {
	vs, ok := v.extArgs[argName]
	if !ok {
		return "", errors.New("Not has arg " + argName)
	}

	return vs, nil
}

func (v *VarConf) parseByMathFunc(express string) (string, error) {
	match := mathFuncRegex.FindStringSubmatch(express)

	if len(match) == 0 {
		return "", errors.New("Invalid match express " + express)
	}

	lv, _ := strconv.ParseInt(match[1], 10, 64)
	rv, _ := strconv.ParseInt(match[3], 10, 64)
	var value int64

	switch match[2] {
	case "+":
		value = lv + rv
	case "-":
		value = lv - rv
	case "*":
		value = lv * rv
	case "/":
		value = lv / rv
	}

	return strconv.FormatInt(value, 10), nil
}
