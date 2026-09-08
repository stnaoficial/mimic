# Manual

## Sumary

- [Variables](#variables)
  - [Variable Syntax](#variable-syntax)

- [Constants](#constants)
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

## Constants

| Constant | Description | Example |
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

| Name         | Description                                                                                           | Example                  | Diacritics |
| ------------ | ----------------------------------------------------------------------------------------------------- | ------------------------ | ---------- |
| `upper`      | Convert a value to UPPERCASE.                                                                         | `{{ upper(name) }}`      | Yes        |
| `lower`      | Convert a value to lowercase.                                                                         | `{{ lower(name) }}`      | Yes        |
| `proper`     | Capitalize the first character of each word in a value.                                               | `{{ proper(name) }}`     | Yes        |
| `title`      | Capitalize the first character of each word except for articles, short prepositions and conjunctions. | `{{ title(name) }}`      | Yes        |
| `capitalize` | Capitalize the first character of a value.                                                            | `{{ capitalize(name) }}` | Yes        |
| `pascal`     | Convert a value to PascalCase.                                                                        | `{{ pascal(name) }}`     | No         |
| `camel`      | Convert a value to camelCase.                                                                         | `{{ camel(name) }}`      | No         |
| `flat`       | Convert a value to flatcase with no separators.                                                       | `{{ flat(name) }}`       | No         |

### Case separators

| Name    | Description                    | Example             | Diacritics |
| ------- | ------------------------------ | ------------------- | ---------- |
| `kebab` | Convert a value to kebab-case. | `{{ kebab(name) }}` | No         |
| `snake` | Convert a value to snake_case. | `{{ snake(name) }}` | No         |

### String manipulation

| Name        | Description                                                                                                                              | Example                         |
| ----------- | ---------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------- |
| `before`    | Return the part of a value before the first occurrence of a target string. Returns the original value if the target is not found.        | `{{ before(name, ".") }}`       |
| `after`     | Return the part of a value after the first occurrence of a target string. Returns the original value if the target is not found.         | `{{ after(name, ".") }}`        |
| `between`   | Return the part of a value between the first occurrence of two target strings. Returns the original value if either target is not found. | `{{ between(name, "[", "]") }}` |
| `replace`   | Replace all occurrences of a value with another value. Requires three arguments: value, old and new.                                     | `{{ replace(name, " ", "-") }}` |
| `normalize` | Remove diacritics from a value.                                                                                                          | `{{ normalize(name) }}`         |
| `delimit`   | Separate the words in a value using the specified delimiter.                                                                             | `{{ delimit(name, "~") }}`      |

### Padding

| Name        | Description                                                               | Example                          |
| ----------- | ------------------------------------------------------------------------- | -------------------------------- |
| `pad_left`  | Pad a value on the left until it reaches the specified length.            | `{{ pad_left(value, 5, "0") }}`  |
| `pad_right` | Pad a value on the right until it reaches the specified length.           | `{{ pad_right(value, 5, "0") }}` |
| `pad_both`  | Pad a value on both sides until it reaches the specified length.          | `{{ pad_both(value, 5, "0") }}`  |
| `zero_fill` | Pad a value with zeros on the left until it reaches the specified length. | `{{ zero_fill(count, 3) }}`      |
