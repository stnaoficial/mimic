package tests

import (
	"maps"
	"mimic/internal/lang"
	"mimic/tests/util"
	"testing"
)

func Test_MustCompileKnownVariables(t *testing.T) {
	env := lang.NewEnvironment()

	maps.Insert(env.Vars, maps.All(map[string]string{
		"name0": "value0",
		"name1": "value1",
		"name2": "value2",
	}))

	comp := lang.NewCompilerConfigurable(env, lang.NewExpression(), true)

	result, err := comp.Compile(lang.NewBuffer("<test>", "{{ name0 }} {{ name1 }} {{ name2 }}"))

	if err != nil {
		t.Error(err)
	}

	util.AssertEquals(t, result, "value0 value1 value2")
}

func Test_MustNotCompileUnknownVariables(t *testing.T) {
	env := lang.NewEnvironment()

	maps.Insert(env.Vars, maps.All(map[string]string{
		"name0": "value0",
	}))

	comp := lang.NewCompilerConfigurable(env, lang.NewExpression(), true)

	_, err := comp.Compile(lang.NewBuffer("<test>", "{{ name1 }}"))

	util.AssertErrorReturn(t, err)
}
