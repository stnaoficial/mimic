# Manual

## Sumary

- [Variables](#variables)
  - [Variable Syntax](#variable-syntax)
  - [Global Variables](#global-variables)
    - [`__SOURCEPATH__`](#__targetpath__)
    - [`__TARGETPATH__`](#__sourcepath__)

  - [Local Variables](#local-variables)
    - [`__BASEPATH__`](#__basepath__)
    - [`__BASENAME__`](#__basename__)
    - [`__DIRNAME__`](#__dirname__)
    - [`__FILENAME__`](#__filename__)
    - [`__FILEDATA__`](#__filedata__)

- [Functions](#functions)
  - [Case formatters](#case-formatters)
  - [Case separators](#case-separators)
  - [Miscellaneous](#miscellaneous)

- [Shell Support](#shell-support)

## Variables

### Variable Syntax

```plaintext
{{ var }}
```

Variables can be used in:

* Directory names
* File names
* File contents

### Global Variables

| Variable | Description | Example |
|----------|-------------|---------|
| [`__SOURCEPATH__`](#__sourcepath__) | Source entry. | `/usr/local` |
| [`__TARGETPATH__`](#__targetpath__) | Target entry. | `/usr/local` |

---

### Local Variables

| Variable | Description | Example |
|----------|-------------|---------|
| [`__BASEPATH__`](#__basepath__) | Parent path of the current entry. | `parent/child` |
| [`__BASENAME__`](#__basename__) | Last path component of the current entry. | `file.ext` |
| [`__DIRNAME__`](#__dirname__) | Full pathname of the current directory entry. Available only for directories. | `parent/child` |
| [`__FILENAME__`](#__filename__) | Full pathname of the current file entry. Available only for files. | `parent/child/file.ext` |
| [`__FILEDATA__`](#__filedata__) | Raw file contents. Available only for files. | |

## Functions

Functions may be used in expressions like `{{ upper(name) }}` or `{{ lower(replace(name, " ", "-")) }}`.

### Case formatters

| Name | Description | Example | Diacritics |
|-|-|-|-|
| `upper` | Convert a value to UPPERCASE. | `{{ upper(name) }}` | Yes |
| `lower` | Convert a value to lowercase. | `{{ lower(name) }}` | Yes |
| `proper` | Capitalize the first character of each word in a value. | `{{ proper(name) }}` | Yes |
| `title` | Capitalize the first character of each word except for articles, short prepositions and conjunctions. | `{{ title(name) }}` | Yes |
| `capitalize` | Capitalize the first character of a value. | `{{ capitalize(name) }}` | Yes |
| `pascal` | Convert a value to PascalCase. | `{{ pascal(name) }}` | No |
| `camel` | Convert a value to camelCase. | `{{ camel(name) }}` | No |
| `flat` | Convert to a value flatcase (no separators). | `{{ flat(name) }}` | No |

### Case separators

| Name | Description | Example | Diacritics |
|-|-|-|-|
| `kebab` | Convert a value to kebab-case. | `{{ kebab(name) }}` | No |
| `snake` | Convert a value to snake_case. | `{{ snake(name) }}` | No |

### Miscellaneous

| Name | Description | Example |
|-|-|-|
| `replace` | Replaces all occurrences of a value with another. Requires three arguments: value, old and new. | `{{ replace(name, " ", "-") }}` |
| `normalize` | Removing diacritics of a value. | `{{ normalize(name) }}` |
| `join` | Join a value by a given space separator. | `{{ join(name, "~") }}` |

## Shell Support

Shell support is provided via the `$()` function.