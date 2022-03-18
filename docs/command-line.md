# NAME

hranoprovod-cli - A new cli application

# SYNOPSIS

hranoprovod-cli

```
[--begin|-b]=[value]
[--config|-c]=[value]
[--database|-d]=[value]
[--date-format]=[value]
[--end|-e]=[value]
[--help|-h]
[--logfile|-l]=[value]
[--maxdepth]=[value]
[--no-color]
[--no-database]
[--version|-v]
```

# DESCRIPTION

Diet tracker for the command line

**Usage**:

```
hranoprovod-cli [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--begin, -b**="": Beginning of period `DATE`

**--config, -c**="": Configuration file `FILE` (default: /home/aquilax/.hranoprovod/config)

**--database, -d**="": optional database file name `FILE` (default: food.yaml)

**--date-format**="": Date format for parsing and printing dates `DATE_FORMAT` (default: 2006/01/02)

**--end, -e**="": End of period `DATE`

**--help, -h**: show help

**--logfile, -l**="": log file name `FILE` (default: log.yaml)

**--maxdepth**="": Resolve depth `DEPTH` (default: 10)

**--no-color**: Disable color output

**--no-database**: Disables loading the database (even if database filename is set)

**--version, -v**: print the version


# COMMANDS

## register, reg

Shows the log register report

**--begin, -b**="": Beginning of period `DATE`

**--csv**: Export as CSV

**--end, -e**="": End of period `DATE`

**--group-food, -g**: Single element grouped by food

**--internal-template-name**="": Name of the internal demplate to use: [default, left-aligned] (default: default)

**--no-color**: Disable color output

**--no-totals**: Disable totals

**--shorten**: Shorten longer strings

**--single-element, -s**="": Show only single element

**--single-food, -f**="": Show only single food

**--totals-only**: Show only totals

**--unresolved**: Deprecated: Show unresolved elements only (moved to 'report unresolved')

**--use-old-reg-reporter**: Use the old reg reporter

## balance, bal

Shows food balance as tree

**--begin, -b**="": Beginning of period

**--collapse, -c**: Collapses sole branches

**--collapse-last**: Collapses last dimension

**--end, -e**="": End of period

**--single-element, -s**="": Show only single element

## lint

Lints file for parsing errors

**--silent, -s**: stay silent if no errors are found

## report

Generates various reports

### element-total

Generates total sum for element grouped by food

**--desc**: Descending order

### unresolved

Print list of unresolved elements

### quantity

Total quantities per food

**--desc**: Descending order

## csv

Generates csv exports

### log

Exports the log file as CSV

**--begin, -b**="": Beginning of period `DATE`

**--end, -e**="": End of period `DATE`

### database

Exports the database file as CSV

### database-resolved

Exports the resolved database as CSV

## stats

Provide stats information

## summary

Show summary

### today

Show summary for today

### yesterday

Show summary for yesterday

## gen

Generate documentation

### man

Generate man page

### markdown

Generate markdown page

**--help, -h**: show help

## print

Print log

**--begin, -b**="": Beginning of period `DATE`

**--end, -e**="": End of period `DATE`

## help, h

Shows a list of commands or help for one command
