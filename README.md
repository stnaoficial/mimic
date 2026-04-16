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

Variables use a double-brace format:

```
{{variable}}
```

Examples:

```
{{name}}
{{domain}}
{{class}}
```

Whitespace is allowed but treated literally:

```
{{ name }}
```

---

## Usage

Basic usage:

```
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
2. It detects variables like `{{name}}`, `{{class}}`, etc.
3. Values are resolved:

   * From `--var` flags if provided
   * Otherwise via interactive prompts
4. Files are generated in the target directory with variables replaced

---

## Interactive Mode

If a variable is not provided via CLI, Mimic will prompt:

```
Please enter a value for {{name}}:
```

---

## Non-Interactive Mode

Provide variables directly:

```
./mimic \
  -s ./templates/.mimic \
  -t ./output \
  -v domain=user \
  -v class=UserEntity \
  -v interface=UserModel
```

---

## Example

### Template Structure

```
.mimic/
└── com/java/app/{{domain}}/
    ├── {{class}}.java.mimic
    └── {{interface}}.java.mimic
```

### Template File

```
package app.{{domain}};

public class {{class}} {
}
```

### Command

```
./mimic \
  -s ./.mimic \
  -t ./src \
  -v domain=user \
  -v class=UserEntity \
  -v interface=UserModel
```

### Output

```
src/com/java/app/user/UserEntity.java
```

```
package app.user;

public class UserEntity {
}
```

---

## Behavior Details

* Only files ending with `.mimic` are processed
* The `.mimic` suffix is removed in generated files
* Directory structure is preserved
* Missing directories are created automatically
* Unknown variables remain unchanged

## Limitations

* Only supports simple `{{variable}}` replacement
* No conditionals or loops
* No built-in validation for variable values
