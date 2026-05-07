# Manual

## Sumary

- [Variables](#variables)
  - [Variable Syntax](#variable-syntax)
  - [Global Variables](#global-variables)
    - [`__TIMESTAMP__`](#__timestamp__)
    - [`__DATE__`](#__date__)
    - [`__TIME__`](#__time__)
    - [`__DATETIME__`](#__datetime__)
    - [`__YEAR__`](#__year__)
    - [`__MONTH__`](#__month__)
    - [`__DAY__`](#__day__)
    - [`__HOUR__`](#__hour__)
    - [`__MINUTE__`](#__minute__)
    - [`__SECOND__`](#__second__)
    - [`__MILLISECOND__`](#__millisecond__)
    - [`__MICROSECOND__`](#__microsecond__)
    - [`__NANOSECOND__`](#__nanosecond__)

  - [Local Variables](#local-variables)
    - [`__UUID__`](#__uuid__)
    - [`__16_DIGIT__`](#__16_digit__)
    - [`__8_DIGIT__`](#__8_digit__)
    - [`__4_DIGIT__`](#__4_digit__)
    - [`__2_DIGIT__`](#__2_digit__)
    - [`__BASEPATH__`](#__basepath__)
    - [`__BASENAME__`](#__basename__)
    - [`__DIRNAME__`](#__dirname__)
    - [`__FILENAME__`](#__filename__)
    - [`__FILEDATA__`](#__filedata__)

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
| [`__TIMESTAMP__`](#__timestamp__) | Unix timestamp in seconds. | `1714859201` |
| [`__DATE__`](#__date__) | Current date in ISO format. | `2026-05-04` |
| [`__TIME__`](#__time__) | Current time. | `19:42:31` |
| [`__DATETIME__`](#__datetime__) | Current date and time. | `2026-05-04T19:42:31Z` |
| [`__YEAR__`](#__year__) | Current year. | `2026` |
| [`__MONTH__`](#__month__) | Current month. | `05` |
| [`__DAY__`](#__day__) | Current day of month. | `04` |
| [`__HOUR__`](#__hour__) | Current hour. | `19` |
| [`__MINUTE__`](#__minute__) | Current minute. | `42` |
| [`__SECOND__`](#__second__) | Current second. | `31` |
| [`__MILLISECOND__`](#__millisecond__) | Current millisecond component. | `123` |
| [`__MICROSECOND__`](#__microsecond__) | Current microsecond component. | `123456` |
| [`__NANOSECOND__`](#__nanosecond__) | Current nanosecond component. | `123456789` |

---

### Local Variables

| Variable | Description | Example |
|----------|-------------|---------|
| [`__UUID__`](#__uuid__) | Unique identifier for the current entry. | `550e8400-e29b-41d4-a716-446655440000` |
| [`__16_DIGIT__`](#__16_digit__) | Random 16-digit number. | `4839201748291038` |
| [`__8_DIGIT__`](#__8_digit__) | Random 8-digit number. | `04829173` |
| [`__4_DIGIT__`](#__4_digit__) | Random 4-digit number. | `0037` |
| [`__2_DIGIT__`](#__2_digit__) | Random 2-digit number. | `07` |
| [`__BASEPATH__`](#__basepath__) | Parent path of the current entry. | `parent/child` |
| [`__BASENAME__`](#__basename__) | Last path component of the current entry. | `file.ext` |
| [`__DIRNAME__`](#__dirname__) | Full pathname of the current directory entry. Available only for directories. | `parent/child` |
| [`__FILENAME__`](#__filename__) | Full pathname of the current file entry. Available only for files. | `parent/child/file.ext` |
| [`__FILEDATA__`](#__filedata__) | Raw file contents. Available only for files. | `!"#$%&'()*+,-./:;<=>?@[]^_`{|}~` |

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
