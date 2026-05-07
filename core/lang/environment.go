package lang

import (
	"fmt"
	"mimic/core/lang/env"
	"mimic/core/util"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type Function func(args []string) string

type Environment struct {
	Vars    map[string]string
	Prompts map[string]string
	Funcs   map[string]Function
}

func buildDefaultVars() map[string]string {
	return map[string]string{}
}

func buildDefaultPrompts() map[string]string {
	return map[string]string{}
}

func buildDefaultFuncs() map[string]Function {
	return map[string]Function{
		"upper": func(args []string) string {
			if len(args) != 1 {
				return ""
			}

			return env.Upper(args[0])
		},
		"lower": func(args []string) string {
			if len(args) != 1 {
				return ""
			}

			return env.Lower(args[0])
		},

		"proper": func(args []string) string {
			if len(args) != 1 {
				return ""
			}

			return env.Proper(args[0])
		},
		"title": func(args []string) string {
			if len(args) != 1 {
				return ""
			}

			return env.Title(args[0])
		},
		"capitalize": func(args []string) string {
			if len(args) != 1 {
				return ""
			}

			return env.Capitalize(args[0])
		},

		"pascal": func(args []string) string {
			if len(args) != 1 {
				return ""
			}

			return env.Pascal(args[0])
		},
		"camel": func(args []string) string {
			if len(args) != 1 {
				return ""
			}

			return env.Camel(args[0])
		},
		"flat": func(args []string) string {
			if len(args) != 1 {
				return ""
			}

			return env.Flat(args[0])
		},

		"kebab": func(args []string) string {
			if len(args) != 1 {
				return ""
			}

			return env.Kebab(args[0])
		},
		"snake": func(args []string) string {
			if len(args) != 1 {
				return ""
			}

			return env.Snake(args[0])
		},

		"replace": func(args []string) string {
			if len(args) != 3 {
				return ""
			}

			return env.Replace(args[0], args[1], args[2])
		},
		"normalize": func(args []string) string {
			if len(args) != 1 {
				return ""
			}

			return env.Normalize(args[0])
		},
		"join": func(args []string) string {
			if len(args) != 2 {
				return ""
			}

			return env.Join(args[0], args[1])
		},
	}
}

func NewEnvironment() *Environment {
	return &Environment{
		Vars:    buildDefaultVars(),
		Prompts: buildDefaultPrompts(),
		Funcs:   buildDefaultFuncs(),
	}
}

func (e *Environment) DefineGlobalVars() {
	now := time.Now()

	e.Vars["__TIMESTAMP__"] = fmt.Sprintf("%d", now.Unix())

	e.Vars["__DATE__"] = now.Format("2006-01-02")
	e.Vars["__TIME__"] = now.Format("15:04:05")
	e.Vars["__DATETIME__"] = now.Format("2006-01-02T15:04:05Z")

	e.Vars["__YEAR__"] = now.Format("2006")
	e.Vars["__MONTH__"] = now.Format("01")
	e.Vars["__DAY__"] = now.Format("02")

	e.Vars["__HOUR__"] = now.Format("15")
	e.Vars["__MINUTE__"] = now.Format("04")
	e.Vars["__SECOND__"] = now.Format("05")

	ns := now.Nanosecond()

	e.Vars["__MILLISECOND__"] = fmt.Sprintf("%03d", ns/1_000_000)
	e.Vars["__MICROSECOND__"] = fmt.Sprintf("%06d", ns/1_000)
	e.Vars["__NANOSECOND__"] = fmt.Sprintf("%09d", ns)
}

func (e *Environment) DefineLocalVars(pathName string, fileInfo os.FileInfo, data []byte) {
	e.Vars["__UUID__"] = uuid.NewString()

	e.Vars["__16_DIGIT__"] = util.RandDigit(16)
	e.Vars["__8_DIGIT__"] = util.RandDigit(8)
	e.Vars["__4_DIGIT__"] = util.RandDigit(4)
	e.Vars["__2_DIGIT__"] = util.RandDigit(2)

	e.Vars["__BASEPATH__"] = filepath.Dir(pathName)
	e.Vars["__BASENAME__"] = filepath.Base(pathName)

	delete(e.Vars, "__DIRNAME__")
	delete(e.Vars, "__FILENAME__")
	delete(e.Vars, "__FILEDATA__")

	if fileInfo.IsDir() {
		e.Vars["__DIRNAME__"] = pathName
	} else {
		e.Vars["__FILENAME__"] = pathName
		e.Vars["__FILEDATA__"] = string(data)
	}
}
