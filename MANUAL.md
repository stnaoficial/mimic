# Manual

## Sumary

- [Variables](#variables)
  - [Variable Syntax](#variable-syntax)
  - [Global Variables](#global-variables)
    - [`__SOURCE_PATH__`](#__target_path__)
    - [`__TARGET_PATH__`](#__source_path__)

  - [Local Variables](#local-variables)
    - [`__COUNT__`](#__count__)
    - [`__PREV_COUNT__`](#__prev_count__)
    - [`__NEXT_COUNT__`](#__next_count__)
    - [`__PATHNAME__`](#__pathname__)
    - [`__DIRNAME__`](#__dirname__)
    - [`__BASENAME__`](#__basename__)

- [Functions](#functions)
  - [Case formatters](#case-formatters)
  - [Case separators](#case-separators)
  - [Miscellaneous](#miscellaneous)

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
| [`__SOURCE_PATH__`](#__source_path__) | Source entry. | `/usr/local` |
| [`__TARGET_PATH__`](#__target_path__) | Target entry. | `/usr/local` |

---

### Local Variables

| Variable | Description | Example |
|----------|-------------|---------|
| [`__COUNT__`](#__count__) | Count of files in the current directory entry. | 0..9 |
| [`__PREV_COUNT__`](#__prev_count__) | Previous count of files in the current directory entry. | 0..9 |
| [`__NEXT_COUNT__`](#__next_count__) | Next count of files in the current directory entry. | 0..9 |
| [`__PATHNAME__`](#__pathname__) | Full pathname of the current file entry. | `parent/child/file.ext` |
| [`__DIRNAME__`](#__dirname__) | Full pathname of the current directory entry. | `parent/child` |
| [`__BASENAME__`](#__basename__) | Last path component of the current entry. | `file.ext` |

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
| `delimit` | Delimit a value by a given space delimiter. | `{{ delimit(name, "~") }}` |
