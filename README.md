# Mimic

Mimic interprets `.mimic` files in the source path (`./.mimic` directory by default) and creates copies of them in the target path (the current directory by default).

---

## Overview

Mimic scans a directory for `.mimic` files, replaces variables defined in templates, and writes the resulting files into a target directory while preserving the structure.

Variables can be used in:

* File contents
* File names
* Directory names

---

## Variable Syntax

Variables use double curly braces and are named in PascalCase format:

```txt
{{Variable}}
```

Examples:

```txt
{{Name}}
{{Domain}}
{{Class}}
```

Whitespace is allowed for modifiers:

```txt
{{camel Name}}  // fooBarBaz
{{pascal Name}} // FooBarBaz
{{snake Name}}  // foo_bar_baz
{{kebab Name}}  // foo-bar-baz
{{lower Name}}  // foobarbaz
{{upper Name}}  // FOOBARBAZ
{{dot Name}}    // foo.bar.baz
```

---

## Usage

Basic usage:

```bash
mimic -s ./.mimic -t ./output -v key0=value0 -v key1=value1 ...
```

---

## Flags

| Flag             | Description                                                      |
| ---------------- | ---------------------------------------------------------------- |
| `-s`, `--source` | Source directory containing `.mimic` files (default: `./.mimic`) |
| `-t`, `--target` | Target directory where files will be generated (default: `.`)    |
| `-v`, `--var`    | Define variables manually (`key=value`)                          |

---

## How It Works

1. Mimic scans the source directory for `.mimic` files
2. It detects variables like `{{Package}}`, `{{Class}}`, etc.
3. Values are resolved:

   * From `--var` flags if provided
   * Otherwise via interactive prompts
4. Values are modified (optional)
5. Files are generated in the target directory with variables replaced

---

## Interactive Mode

If a variable is not provided via CLI, Mimic will prompt:

```bash
$ Please enter a value for "Domain":
```

---

## Non-Interactive Mode

Provide variables directly:

```bash
./mimic \
  -s ./templates/.mimic \
  -t ./output \
  -v Domain=User \
  -v Version=1.0.0
```

## Behavior Details

* Only files ending with `.mimic` are processed
* The `.mimic` suffix is removed in generated files
* Directory structure is preserved
* Missing directories are created automatically
* Unknown variables remain unchanged

## Limitations

* Only supports simple `{{Name}}` and `{{modfier Name}}` variable replacement
* No conditionals or loops
* No built-in validation for variable values
