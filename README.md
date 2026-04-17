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

```plaintext
{{ name }}
```

Examples:

```plaintext
{{ name }}
{{ Name }}
{{ This is a name }}
```

Modifiers:

```plaintext
{{ camel(some name) }}  // someName
{{ pascal(some name) }} // SomeName
{{ snake(some name) }}  // some_name
{{ kebab(some name) }}  // some-name
{{ dot(some name) }}    // some.name
{{ flat(some name) }}   // somename
{{ lower(some name) }}  // some name
{{ upper(some name) }}  // SOME NAME
```

Modifiers can be nested in different combinations:

```plaintext
{{ upper(snake(some name)) }} // SOME_NAME
{{ upper(kebab(some name)) }} // SOME-NAME
{{ upper(dot(some name)) }} // SOME.NAME
{{ lower(snake(some name)) }} // some_name
{{ lower(kebab(some name)) }} // some-name
{{ lower(dot(some name)) }} // some.name
```

---

## Usage

Basic usage:

```bash
$ mimic -s ./.mimic -t ./output -v key0=value0 -v "key1=value1" ...
```

---

## Flags

| Flag | Long Flag | Description | Default |
| :--- | :--- | :--- | :--- |
| `-s` | `--source` | Set the source directory path of `.mimic` files | `./.mimic` |
| `-t` | `--target` | Set the target path where all files will be copied | `.` |
| `-v` | `--var` | Set a var directly by passing as a `key=value` pair | |
| | `--var-prefix` | Set the var pattern prefix | `{{` |
| | `--var-sufix` | Set the var pattern sufix | `}}` |

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
