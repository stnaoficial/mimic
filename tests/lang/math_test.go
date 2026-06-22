package tests

import (
	"maps"
	"mimic/internal/lang"
	"mimic/tests/util"
	"testing"
)

func Test_MustSumTwoValues(t *testing.T) {
	comp := lang.NewCompiler(true)

	result, err := comp.Compile(lang.NewBuffer("<test>", "{{ 1 + 1 }}"))

	if err != nil {
		t.Fatal(err)
	}

	util.AssertEquals(t, result, "2")
}

func Test_MustSubtractTwoValues(t *testing.T) {
	comp := lang.NewCompiler(true)

	result, err := comp.Compile(lang.NewBuffer("<test>", "{{ 10 - 5 }}"))

	if err != nil {
		t.Fatal(err)
	}

	util.AssertEquals(t, result, "5")
}

func Test_MustMultiplyTwoValues(t *testing.T) {
	comp := lang.NewCompiler(true)

	result, err := comp.Compile(lang.NewBuffer("<test>", "{{ 5 * 5 }}"))

	if err != nil {
		t.Fatal(err)
	}

	util.AssertEquals(t, result, "25")
}

func Test_MustDivideTwoValues(t *testing.T) {
	comp := lang.NewCompiler(true)

	result, err := comp.Compile(lang.NewBuffer("<test>", "{{ 10 / 2 }}"))

	if err != nil {
		t.Fatal(err)
	}

	util.AssertEquals(t, result, "5")
}

func Test_MustRespectOperatorPrecedence(t *testing.T) {
	comp := lang.NewCompiler(true)

	result, err := comp.Compile(lang.NewBuffer("<test>", "{{ 1 + 2 * 3 }}"))

	if err != nil {
		t.Fatal(err)
	}

	util.AssertEquals(t, result, "7")
}

func Test_MustRespectParentheses(t *testing.T) {
	comp := lang.NewCompiler(true)

	result, err := comp.Compile(lang.NewBuffer("<test>", "{{ (1 + 2) * 3 }}"))

	if err != nil {
		t.Fatal(err)
	}

	util.AssertEquals(t, result, "9")
}

func Test_MustChainOperations(t *testing.T) {
	comp := lang.NewCompiler(true)

	result, err := comp.Compile(lang.NewBuffer("<test>", "{{ 1 + 2 + 3 + 4 }}"))

	if err != nil {
		t.Fatal(err)
	}

	util.AssertEquals(t, result, "10")
}

func Test_MustHandleUnaryMinus(t *testing.T) {
	comp := lang.NewCompiler(true)

	result, err := comp.Compile(lang.NewBuffer("<test>", "{{ -5 + 10 }}"))

	if err != nil {
		t.Fatal(err)
	}

	util.AssertEquals(t, result, "5")
}

func Test_MustSumVariableAndNumber(t *testing.T) {
	env := lang.NewEnvironment()

	maps.Insert(env.Vars, maps.All(map[string]string{
		"var0": "10",
	}))

	comp := lang.NewCompilerConfigurable(env, lang.NewExpression(), true)

	result, err := comp.Compile(
		lang.NewBuffer("<test>", "{{ var0 + 5 }}"),
	)

	if err != nil {
		t.Fatal(err)
	}

	util.AssertEquals(t, result, "15")
}

func Test_MustSubtractVariableAndNumber(t *testing.T) {
	env := lang.NewEnvironment()

	maps.Insert(env.Vars, maps.All(map[string]string{
		"var0": "10",
	}))

	comp := lang.NewCompilerConfigurable(env, lang.NewExpression(), true)

	result, err := comp.Compile(
		lang.NewBuffer("<test>", "{{ var0 - 3 }}"),
	)

	if err != nil {
		t.Fatal(err)
	}

	util.AssertEquals(t, result, "7")
}

func Test_MustMultiplyVariableAndNumber(t *testing.T) {
	env := lang.NewEnvironment()

	maps.Insert(env.Vars, maps.All(map[string]string{
		"var0": "10",
	}))

	comp := lang.NewCompilerConfigurable(env, lang.NewExpression(), true)

	result, err := comp.Compile(lang.NewBuffer("<test>", "{{ var0 * 2 }}"))

	if err != nil {
		t.Fatal(err)
	}

	util.AssertEquals(t, result, "20")
}

func Test_MustDivideVariableAndNumber(t *testing.T) {
	env := lang.NewEnvironment()

	maps.Insert(env.Vars, maps.All(map[string]string{
		"var0": "10",
	}))

	comp := lang.NewCompilerConfigurable(env, lang.NewExpression(), true)

	result, err := comp.Compile(lang.NewBuffer("<test>", "{{ var0 / 2 }}"))

	if err != nil {
		t.Fatal(err)
	}

	util.AssertEquals(t, result, "5")
}
