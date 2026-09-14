# Mimic Templating Library

## Overview

The library is a template management and interpretation library that also includes a compiler for processing `.mimic` files and directory structures.

Templates can be made available through the `.mimic/templates` directory, and they can come from either the local filesystem or a remote repository. The remote repository and other settings can be configured through the `.mimic/config` file.

For more control, templates can also be used directly by specifying their source and/or target paths.

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
# Without specifying which template to use
mimic local|remote

# Specifying the template name
mimic local|remote -n js-class

# Specifying the source and target path
mimic local|remote -s ./.mimic/templates/js-class -t .
```

## Example

The html-5 template defines a basic HTML project structure:

```bash
.mimic/templates/html-5/
└── {{name}}
    ├── assets
    ├── index.html.mimic
    ├── script.js
    └── style.css
```

When this template is used, Mimic interprets the directory structure and creates a new project from it. The `{{name}}` directory contains a template variable, so its name is resolved when the template is compiled. For example, using html-5 with `-v name=my-site` would produce:

```bash
my-site
├── assets
├── index.html.mimic
├── script.js
└── style.css
```

Files ending in `.mimic` are interpreted by the compiler and generated as regular files, while other files are copied as they are. This allows a template to combine static files and dynamically generated files in the same directory.
