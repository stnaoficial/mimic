<h1 align="center">
   <img src="./assets/icon.jpg" width="80px" height="80px" />
   <br>
   Mimic Templating Library
</h1>

## Overview

Mimic interprets files and directories from a source path (`.mimic` directory by default) and generates copies in a target path (the current directory by default).

## Usage

Basic usage:

```bash
$ mimic copy            # Without specifing the source and target path
$ mimic copy ./.mimic . # Specifing the source and target path
```

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
$ Please enter a value for "name": My variable name
```

This will be evaluated as:

```txt
MyVariableName
```

## Non-Interactive Mode

Provide variables directly:

```bash
mimic copy -v name0=value -v "name1=value" -v name2="value" ...
```

Customize prompt messages:

```bash
mimic copy -p name0="My custom prompt message: " ...
```

## Behavior Details

* Only files ending with `.mimic` are processed
* The `.mimic` suffix is removed in generated files
* Directory structure is preserved, as are non-mimic files
* Missing directories are created automatically