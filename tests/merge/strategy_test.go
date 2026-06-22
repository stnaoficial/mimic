package merge

import (
	"mimic/internal/merge"
	"mimic/tests/util"
	"strings"
	"testing"
)

func Test_MustMergeTwoJavaClasses(t *testing.T) {
	strategy := merge.NewStrategy()

	old := `
	public final class Example {

		public void methodOne() {
			return;
		}
	}
	`

	new := `
	public final class Example {

		public void methodTwo() {
			return;
		}
	}
	`

	expected := `
	public final class Example {

		public void methodOne() {
			return;
		}

		public void methodTwo() {
			return;
		}
	}
	`

	result, err := strategy.Merge(old, new)

	if err != nil {
		t.Error(err)
	}

	util.AssertEquals(t, result, expected)
}

func Test_MustMergeTwoMarkdownSections(t *testing.T) {
	strategy := merge.NewStrategy()

	old := `
	# Heading 1

	Paragraph of heading 1.
	`

	new := `
	# Heading 2

	Paragraph of heading 2.
	`

	expected := `
	# Heading 1

	Paragraph of heading 1.

	# Heading 2

	Paragraph of heading 2.
	`

	result, err := strategy.Merge(old, new)

	if err != nil {
		t.Error(err)
	}

	util.AssertEquals(t, result, expected)
}

func Test_MustMergeTwoPythonClasses(t *testing.T) {
	strategy := merge.NewStrategy()

	old := `
	class DataProcessor:
		def process_old(self):
			return "old"
	`

	new := `
	class DataProcessor:
		def process_new(self):
			return "new"
	`

	expected := `
	class DataProcessor:
		def process_old(self):
			return "old"

		def process_new(self):
			return "new"
	`

	result, err := strategy.Merge(old, new)

	if err != nil {
		t.Error(err)
	}

	util.AssertEquals(t, result, expected)
}

func Test_MustMergeTwoJSONObjects(t *testing.T) {
	strategy := merge.NewStrategy()

	old := `
	{
		"config": "base",
		"feature_a": true
	}
	`

	new := `
	{
		"config": "base",
		"feature_b": false
	}
	`

	expected := `
	{
		"config": "base",
		"feature_a": true

		"feature_b": false
	}
	`

	result, err := strategy.Merge(old, new)

	if err != nil {
		t.Error(err)
	}

	util.AssertEquals(t, result, expected)
}

func Test_MustMergeTwoHTMLNodes(t *testing.T) {
	strategy := merge.NewStrategy()

	old := `
	<div class="container">
		<p>First paragraph</p>
	</div>
	`

	new := `
	<div class="container">
		<p>Second paragraph</p>
	</div>
	`

	expected := `
	<div class="container">
		<p>First paragraph</p>

		<p>Second paragraph</p>
	</div>
	`

	result, err := strategy.Merge(old, new)

	if err != nil {
		t.Error(err)
	}

	util.AssertEquals(t, strings.TrimSpace(result), strings.TrimSpace(expected))
}

func Test_MustHandleIdenticalBuffers(t *testing.T) {
	strategy := merge.NewStrategy()

	content := `
	same line 1
	same line 2
	`

	result, err := strategy.Merge(content, content)
	if err != nil {
		t.Error(err)
	}

	util.AssertEquals(t, result, content)
}

func Test_MustHandleCompletelyDifferentBuffers(t *testing.T) {
	strategy := merge.NewStrategy()

	old := "alpha\nbeta"
	new := "gamma\ndelta"

	expected := "alpha\nbeta\n\ngamma\ndelta"

	result, err := strategy.Merge(old, new)
	if err != nil {
		t.Error(err)
	}

	util.AssertEquals(t, result, expected)
}
