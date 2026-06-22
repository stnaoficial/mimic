package tests

import (
	"maps"
	"mimic/internal/lang"
	"mimic/tests/util"
	"testing"
)

func Test_MustCompileKnownStringVariables(t *testing.T) {
	env := lang.NewEnvironment()

	maps.Insert(env.Vars, maps.All(map[string]string{
		"var0": "value0",
		"var1": "value1",
		"var2": "value2",
	}))

	comp := lang.NewCompilerConfigurable(env, lang.NewExpression(), true)

	result, err := comp.Compile(lang.NewBuffer("<test>", "{{ var0 }} {{ var1 }} {{ var2 }}"))

	if err != nil {
		t.Error(err)
	}

	util.AssertEquals(t, result, "value0 value1 value2")
}

func Test_MustCompileKnownNumericVariables(t *testing.T) {
	env := lang.NewEnvironment()

	maps.Insert(env.Vars, maps.All(map[string]string{
		"var0": "0",
		"var1": "1",
		"var2": "2",
	}))

	comp := lang.NewCompilerConfigurable(env, lang.NewExpression(), true)

	result, err := comp.Compile(lang.NewBuffer("<test>", "{{ var0 }} {{ var1 }} {{ var2 }}"))

	if err != nil {
		t.Error(err)
	}

	util.AssertEquals(t, result, "0 1 2")
}

func Test_MustCompileKnownCompositeVariables(t *testing.T) {
	env := lang.NewEnvironment()

	maps.Insert(env.Vars, maps.All(map[string]string{
		"var0": "1 is number one",
	}))

	comp := lang.NewCompilerConfigurable(env, lang.NewExpression(), true)

	result, err := comp.Compile(lang.NewBuffer("<test>", "{{ var0 }}"))

	if err != nil {
		t.Error(err)
	}

	util.AssertEquals(t, result, "1 is number one")
}

func Test_MustNotCompileUnknownVariables(t *testing.T) {
	env := lang.NewEnvironment()

	maps.Insert(env.Vars, maps.All(map[string]string{
		"var0": "value0",
	}))

	comp := lang.NewCompilerConfigurable(env, lang.NewExpression(), true)

	_, err := comp.Compile(lang.NewBuffer("<test>", "{{ var1 }}"))

	util.AssertErrorReturn(t, err)
}
