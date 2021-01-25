# hranoprovod-cli

Hranoprovod is command line tracking tool. It supports nested recipies and custom defined tracking elements, which makes it perfect for tracking calories, nutrition data, exercises and other accumulative data.

[![asciicast](https://asciinema.org/a/257200.svg)](https://asciinema.org/a/257200)

## General idea

hranoprovod-cli borrows the idea of keeping transactions in plain text files from plain text accounting and applies the same principles to tracking nutrition, exercise or any other quantifiable metric.

The data is stored in two text files:

* journal - a list of multiple entries for each date

Example:

```
2021/01/24:
  day/nonworking: 1
  coffee/cup: 1
  sandwich/egg/lettuce/100g: 1.20
  candy/snickers/bar: 1
```

* recipe book - a file where new recipes can be created for previously defined ingredients

Example:

```
# daily nutrition budget
day/nonworking:
  calories: -1200
  fat: -124
  carbohydrate: -50
  protein: -104

bread/rye/100g:
  calories: 259
  fat: 3.3
  carbohydrate: 48
  protein: 9

egg/boiled/100g:
  calories: 155
  fat: 11
  carbohydrate: 1.1
  protein: 13

vegetables/lettuce/romaine/100g:
  calories: 15
  fat: 0.5
  carbohydrate: 1.7
  protein: 0.9

sauce/mayonnaise/100g:
  calories: 680
  fat: 7.5
  carbohydrate: 0.6
  protein: 1

sandwich/egg/lettuce/100g:
  bread/rye/100g: 0.40
  egg/boiled/100g: 0.20
  vegetables/lettuce/romaine/100g: 0.20
  sauce/mayonnaise/100g: 0.20

candy/snickers/bar:
  calories: 280
  fat: 13.6
  carbohydrate: 35.1
  protein: 4.29

```

Running a daily `reg` report using the above examples produces the following result:

```
2021/01/24
        day/nonworking              :      1.00
                            calories   -1200.00
                        carbohydrate     -50.00
                                 fat    -124.00
                             protein    -104.00
        coffee/cup                  :      1.00
                          coffee/cup       1.00
        sandwich/egg/lettuce/100g   :      1.20
                            calories     328.32
                carbohydrate      23.86
                                 fat       6.14
                             protein       7.90
        candy/snickers/bar          :      1.00
                            calories     280.00
                        carbohydrate      35.10
                                 fat      13.60
                     protein       4.29
        -- TOTAL  ----------------------------------------------------
                            calories     608.32   -1200.00 =   -591.68
                        carbohydrate      58.96     -50.00 =      8.96
                          coffee/cup       1.00       0.00 =      1.00
                                 fat      19.74    -124.00 =   -104.26
                             protein      12.19    -104.00 =    -91.81
```

And this is the daily summary:

```
2021/01/24 :
    608.32 : calories
     58.96 : carbohydrate
      1.00 : coffee/cup
     19.74 : fat
     12.19 : protein
------------
      1.00 : day/nonworking
      1.00 : coffee/cup
      1.20 : sandwich/egg/lettuce/100g
      1.00 : candy/snickers/bar
```

## How does it work

When you run a report `hranoprovod-cli` will look up each entry in the journal in the recipe book file. Each ingredient in the recipe book will also try to lookup itself up until no further reference is found. In the examples above `calories`, `carbohydrate`, `couffee/cup`, `fat` and `protein` are returned as aggregate numbers in the totals since they were never defined. This allows for a great tracking flexibility and also making the tool metric agnostic.

## Keeping a journal

TODO

## Keeping a recipe book

TODO

## Reporting

TODO

## Exporting data

TODO

### Filtering

TODO

### Agregation

TODO