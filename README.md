# Mimic Templating Library

## Overview

Mimic interpret `.mimic` files in the source path (`./.mimic` directory by default) and generate copies of them in the target path (the current directory by default).

Variables can be used in:

* File contents
* File names
* Directory names

## Variable Syntax

```plaintext
{{ something }}
```

Variables can be modified:

```plaintext
{{ upper(something) }}
```

Modifiers can be nested:

```plaintext
{{ upper(kebab(something)) }}
```

## Usage

Basic usage:

```bash
$ mimic            # Without specifing the source and target path
$ mimic ./.mimic . # Specifing the source and target path
```

## Flags

| Flag | Long Flag | Description | Default |
| :--- | :--- | :--- | :--- |
| `-v` | `--var` | Set a var directly by passing as a `key=value` pair | |
| | `--expr-open` | Set the open expression syntax | `{{` |
| | `--expr-close` | Set the close expression syntax | `}}` |
| `-h` | `--help` | Print Help (this message) and exit
| | `--version` | Print version information and exit

## How It Works

1. Mimic scans the source directory for `.mimic` files
2. It detects variables like `{{ name }}`, `{{ lower(name) }}`, etc.
3. Expressions are evaluated:

   * From `--var` flags if provided
   * Otherwise via interactive prompts
4. Values are modified (optional)
5. Files are generated in the target directory with variables evaluated

## Interactive Mode

If a variable is not provided via CLI, Mimic will prompt:

```txt
{{ pascal(name) }}
```

```bash
$ Please enter a value for "name": Some name or description
```

This will be evaluated as something like:

```txt
SomeNameOrDescription
```

## Non-Interactive Mode

Provide variables directly:

```bash
mimic -v name=value -v "name=value" -v name="value" ...
```

## Behavior Details

* Only files ending with `.mimic` are processed
* The `.mimic` suffix is removed in generated files
* Directory structure is preserved
* Missing directories are created automatically
* Unknown variables remain unchanged

## Limitations

* No conditionals or loops
