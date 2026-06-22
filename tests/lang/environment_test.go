package tests

import (
	"mimic/internal/lang"
	"mimic/tests/util"
	"testing"
)

func Test_MustCompilePadLeftCallable(t *testing.T) {
	comp := lang.NewCompiler(true)

	result1, err := comp.Compile(lang.NewBuffer("<test>", "{{ pad_left(1, 3, '-') }}"))

	if err != nil {
		t.Error(err)
	}

	result2, err := comp.Compile(lang.NewBuffer("<test>", "{{ pad_left('1', 3, '-') }}"))

	if err != nil {
		t.Error(err)
	}

	util.AssertEquals(t, result1, "--1")
	util.AssertEquals(t, result2, "--1")
}

func Test_MustCompilePadRightCallable(t *testing.T) {
	comp := lang.NewCompiler(true)

	result1, err := comp.Compile(lang.NewBuffer("<test>", "{{ pad_right(1, 3, '-') }}"))

	if err != nil {
		t.Error(err)
	}

	result2, err := comp.Compile(lang.NewBuffer("<test>", "{{ pad_right('1', 3, '-') }}"))

	if err != nil {
		t.Error(err)
	}

	util.AssertEquals(t, result1, "1--")
	util.AssertEquals(t, result2, "1--")
}

func Test_MustCompilePadCallable(t *testing.T) {
	comp := lang.NewCompiler(true)

	result1, err := comp.Compile(lang.NewBuffer("<test>", "{{ pad(1, 3, '-') }}"))

	if err != nil {
		t.Error(err)
	}

	result2, err := comp.Compile(lang.NewBuffer("<test>", "{{ pad('1', 3, '-') }}"))

	if err != nil {
		t.Error(err)
	}

	util.AssertEquals(t, result1, "-1-")
	util.AssertEquals(t, result2, "-1-")
}

func Test_MustCompileZeroFillCallable(t *testing.T) {
	comp := lang.NewCompiler(true)

	result1, err := comp.Compile(lang.NewBuffer("<test>", "{{ zero_fill(1, 3) }}"))

	if err != nil {
		t.Error(err)
	}

	result2, err := comp.Compile(lang.NewBuffer("<test>", "{{ zero_fill('1', 3) }}"))

	if err != nil {
		t.Error(err)
	}

	util.AssertEquals(t, result1, "001")
	util.AssertEquals(t, result2, "001")
}
