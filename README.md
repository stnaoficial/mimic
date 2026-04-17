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

Variables use double curly braces:

```txt
{{ name }}
```

Examples:

```txt
{{ name }}
{{ Name }}
{{ This is a name }}
```

Modifiers:

```txt
{{ camel(some name) }}  // someName
{{ pascal(some name) }} // SomeName
{{ snake(some name) }}  // some_name
{{ kebab(some name) }}  // some-name
{{ dot(some name) }}    // some.name
{{ flat(some name) }}   // somename
{{ lower(some name) }}  // some name
{{ upper(some name) }}  // SOME NAME
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
2. It detects variables like `{{ name }}`, `{{ lower(name) }}`, etc.
3. Values are resolved:

   * From `--var` flags if provided
   * Otherwise via interactive prompts
4. Values are modified (optional)
5. Files are generated in the target directory with variables replaced

---

## Interactive Mode

If a variable is not provided via CLI, Mimic will prompt:

```bash
$ Please enter a value for "name":
$ Please enter a value for "This is a name":
```

---

## Non-Interactive Mode

Provide variables directly:

```bash
mimic -v name=value -v "This is a name = This is a value"
```

## Behavior Details

* Only files ending with `.mimic` are processed
* The `.mimic` suffix is removed in generated files
* Directory structure is preserved
* Missing directories are created automatically
* Unknown variables remain unchanged

## Limitations

* No conditionals or loops
