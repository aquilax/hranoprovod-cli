# Configuration file

Often used options can be set in a configuration file to remove the need to
explicitly add them each time.

## Default configuration file

The default configuration file is loaded from $HOME/.hranoprovod/config but
any valid configuration file can be passed as a parameter using the `-c` or
`--config` command line options.

## Configuration file format

The configuration file is following the INI file structure with the following
main sections (all of them are optional):

* `Global` - general configuration
* `Resolver` - resolver specific configuration

Details for each section follow:

### Global

#### Now [String]

Allows the user to overwrite the current time and date. Must contain date time
formatted using `RFC3339` e.g. `2020-01-01T01:00:00Z` . By default the current
date and time will be used.

#### DbFileName [String]

The full path to the database file e.g. `/tmp/db.yaml`

#### LogFileName [String]

The full path to the log journal file e.g. `/tmp/log.yaml`

#### DateFormat [String]

The format in which the dates will be written in the log journal file. Formatted
using the go date layout described [here](https://pkg.go.dev/time#pkg-constants)
e.g. `2006-01-02`

### Resolver

#### MaxDepth [Integer]

Sets the maximum depth to which the resolver will try to resolve recipes
e.g. `10`

## Example configuration file

```ini
[Global]
Now=2020-01-01T01:00:00Z
DbFileName=/tmp/db.yaml
LogFileName=/tmp/log.yaml
DateFormat=2006-01-02
[Resolver]
MaxDepth=10
```
