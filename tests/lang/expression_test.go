package tests

import (
	"mimic/internal/lang"
	"mimic/tests/util"
	"testing"
)

func Test_MustCompileEmptyExpression(t *testing.T) {
	comp := lang.NewCompiler(true)

	result, err := comp.Compile(lang.NewBuffer("<test>", "{{ }}"))

	if err != nil {
		t.Error(err)
	}

	util.AssertEquals(t, result, "")
}

func Test_MustCompileMultipleEmptyExpressions(t *testing.T) {
	comp := lang.NewCompiler(true)

	result, err := comp.Compile(lang.NewBuffer("<test>", "{{ }} {{ }} {{ }}"))

	if err != nil {
		t.Error(err)
	}

	util.AssertEquals(t, result, "  ")
}

func Test_MustNotCompileNestedExpressions(t *testing.T) {
	comp := lang.NewCompiler(true)

	_, err := comp.Compile(lang.NewBuffer("<test>", "{{ {{ }} }}"))

	util.AssertErrorReturn(t, err)
}
