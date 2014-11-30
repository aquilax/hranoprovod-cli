## hranoprovod-cli [![Build Status](https://travis-ci.org/Hranoprovod/hranoprovod-cli.svg?branch=master)](https://travis-ci.org/Hranoprovod/hranoprovod-cli)

## Description

Hranoprovod is command line tracking tool. It supports nested recipies and custom defined tracking elements, which makes it perfect for tracking calories, nutionin data, excercises and other accumulative data.

## Installation

First make sure you have go (golang) installed.

    http://golang.org/

Download the source code.
  
    go get github.com/Hranoprovod/hranoprovod-cli

## Help

Running the `hranoprovod-cli` command will show you the command line options

## Usage

Hranoprovod uses two files with similar format to operate.

### Database file (food.yaml)

Contains all the "recipes" in the following format:

```
fish/tuna/canned/100g:
  calories: 184
  fat: 6
  carbohydrate: 0
  protein: 0

bread/white/100g:
  calories: 265
  fat: 3.2
  carbohydrate: 49
  protein: 9
```


Let's say you love tuna sandwiches then you can combine these two ingredients into one:

```
sandwich/tuna/100g:
  fish/tuna/canned/100g: .6
  bread/white/100g: .4
  
sandwich/tuna/pc:
  sandwich/tuna/100g: 1.5
```

This means that the sandwich is composed of 60% tuna and 40% bread and a sandwich weights arount 150g.

Hranoprovod is measure agnostic and it's up to the user to use or state the measurements.

### Log file (log.yaml)

The log file contains dated usage of the recipes, defined in the database file.

```
2014/12/17:
  tea/cup: 1
  sandwich/tuna/pc: 2
  calories: 300
  biking/km: 10
```

Note: it's not mandatory to have the elements in the database file. Elements which are not found will be represented as they are. They can always be added later to the database.

Given this example, the result will look like:

```
$hranoprovod-cli -d food.yaml -l log.yaml  reg
2014/12/17
	tea/cup                     :      1.00
		             tea/cup       1.00
	sandwich/tuna/pc            :      2.00
		            calories     649.20
		        carbohydrate      58.80
		                 fat      14.64
		             protein      10.80
	calories                    :    300.00
		            calories     300.00
	biking/km                   :     10.00
		           biking/km      10.00
	-- TOTAL  ----------------------------------------------------
		           biking/km      10.00       0.00 =     10.00
		            calories     949.20       0.00 =    949.20
		        carbohydrate      58.80       0.00 =     58.80
		                 fat      14.64       0.00 =     14.64
		             protein      10.80       0.00 =     10.80
		             tea/cup       1.00       0.00 =      1.00
```