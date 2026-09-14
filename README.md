# Mimic Templating Library

## Overview

Mimic is a template management library that interprets templates from source paths (`.mimic/templates` by default) to target paths (the current directory by default), allowing them to serve as a foundation for new projects.

## Installation

Mimic supports Unix-like systems, including Linux and MacOS.

Install the latest version with:

```bash
curl -fsSL https://raw.githubusercontent.com/stnaoficial/mimic/main/install.sh | sh
```

Verify the installation:

```bash
mimic --version
```

Alternatively, download the appropriate archive from the [latest GitHub release](https://github.com/stnaoficial/mimic/releases/latest) and install the binary manually:

```bash
sudo install mimic /usr/local/bin/mimic
```

## Usage

Basic usage:

```bash
$ mimic local                  # Without specifying the source and target path
$ mimic local -n js-class      # Specifying the template name
$ mimic local -s ./.mimic -t . # Specifying the source and target path
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
mimic local -v name0=value -v "name1=value" -v name2="value" ...
```

Customize prompt messages:

```bash
mimic local -p name0="My custom prompt message: " ...
```

## Behavior Details

* Only files ending with `.mimic` are processed
* The `.mimic` suffix is removed in generated files
* Directory structure is preserved, as are non-mimic files
* Missing directories are created automatically
