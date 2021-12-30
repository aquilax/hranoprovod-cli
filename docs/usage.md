# Usage

Given the following log file:

```sh
$ cat examples/log.yaml
2021/01/24:
  day/nonworking: 1
  coffee/cup: 1
  sandwich/egg/lettuce/100g: 1.20
  candy/snickers/bar: 1

2021/01/25:
  day/nonworking: 1
  coffee/cup: 1
  sandwich/egg/lettuce/100g: 1.50
  coffee/cup: 1

```

and the following recipe book file:

```sh
$ cat examples/food.yaml
# daily nutrition budget
day/nonworking:
  calories: -1200
  fat: -124
  carbohydrate: -50
  protein: -104

bread/rye/100g:
  # barcode: 0000000000000
  calories: 259
  fat: 3.3
  carbohydrate: 48
  protein: 9

egg/boiled/100g:
  # boiling time: 12 min
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

## Register report

### Register view (default template)

```sh
$ ./hranoprovod-cli -d examples/food.yaml -l examples/log.yaml --no-color reg
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
2021/01/25
	day/nonworking              :      1.00
		            calories   -1200.00
		        carbohydrate     -50.00
		                 fat    -124.00
		             protein    -104.00
	coffee/cup                  :      2.00
		          coffee/cup       2.00
	sandwich/egg/lettuce/100g   :      1.50
		            calories     410.40
		        carbohydrate      29.82
		                 fat       7.68
		             protein       9.87
	-- TOTAL  ----------------------------------------------------
		            calories     410.40   -1200.00 =   -789.60
		        carbohydrate      29.82     -50.00 =    -20.18
		          coffee/cup       2.00       0.00 =      2.00
		                 fat       7.68    -124.00 =   -116.32
		             protein       9.87    -104.00 =    -94.13

```

### Register view (left aligned template)

```sh
$ ./hranoprovod-cli -d examples/food.yaml -l examples/log.yaml --no-color reg --internal-template-name=left-aligned
2021/01/24
        1.00  day/nonworking
    -1200.00    calories
      -50.00    carbohydrate
     -124.00    fat
     -104.00    protein
        1.00  coffee/cup
        1.00    coffee/cup
        1.20  sandwich/egg/lettuce/100g
      328.32    calories
       23.86    carbohydrate
        6.14    fat
        7.90    protein
        1.00  candy/snickers/bar
      280.00    calories
       35.10    carbohydrate
       13.60    fat
        4.29    protein
------------------------------------------------------- TOTAL --
      608.32   -1200.00 =    -591.68  calories
       58.96     -50.00 =       8.96  carbohydrate
        1.00       0.00 =       1.00  coffee/cup
       19.74    -124.00 =    -104.26  fat
       12.19    -104.00 =     -91.81  protein
2021/01/25
        1.00  day/nonworking
    -1200.00    calories
      -50.00    carbohydrate
     -124.00    fat
     -104.00    protein
        2.00  coffee/cup
        2.00    coffee/cup
        1.50  sandwich/egg/lettuce/100g
      410.40    calories
       29.82    carbohydrate
        7.68    fat
        9.87    protein
------------------------------------------------------- TOTAL --
      410.40   -1200.00 =    -789.60  calories
       29.82     -50.00 =     -20.18  carbohydrate
        2.00       0.00 =       2.00  coffee/cup
        7.68    -124.00 =    -116.32  fat
        9.87    -104.00 =     -94.13  protein

```

### Register report for single food per day

```sh
$ ./hranoprovod-cli -d examples/food.yaml -l examples/log.yaml --no-color reg --single-food coffee/cup
2021/01/24	coffee/cup	1.00
2021/01/25	coffee/cup	2.00

```

### Register report for total of ingredients for a day grouped by food

```sh
$ ./hranoprovod-cli -d examples/food.yaml -l examples/log.yaml -b 2021/01/24 -e 2021/01/24 --no-color reg -g -s calories
  -1200.00	day/nonworking
    328.32	sandwich/egg/lettuce/100g
    280.00	candy/snickers/bar

```

## Balance

### Default balance report

```sh
$ ./hranoprovod-cli -d examples/food.yaml -l examples/log.yaml --no-color bal
           | candy
           |   snickers
      1.00 |     bar
           | coffee
      3.00 |   cup
           | day
      2.00 |   nonworking
           | sandwich
           |   egg
           |     lettuce
      2.70 |       100g

```

### Default balance report (collapse)

```sh
$ ./hranoprovod-cli -d examples/food.yaml -l examples/log.yaml --no-color bal -c
      1.00 | candy/snickers/bar
      3.00 | coffee/cup
      2.00 | day/nonworking
      2.70 | sandwich/egg/lettuce/100g

```

### Default balance report (collapse last)

```sh
$ ./hranoprovod-cli -d examples/food.yaml -l examples/log.yaml --no-color bal --collapse-last
           | candy
      1.00 |   snickers/bar
      3.00 | coffee/cup
      2.00 | day/nonworking
           | sandwich
           |   egg
      2.70 |     lettuce/100g

```

### Default balance report (single element)

```sh
$ ./hranoprovod-cli -d examples/food.yaml -l examples/log.yaml --no-color bal -s fat
     13.60 | candy
     13.60 |   snickers
     13.60 |     bar
   -248.00 | day
   -248.00 |   nonworking
     13.82 | sandwich
     13.82 |   egg
     13.82 |     lettuce
     13.82 |       100g
-----------|
   -220.58 | fat

```


## Other reports

### Element-total

Calories for each defined food:

```sh
$ ./hranoprovod-cli -d examples/food.yaml -l examples/log.yaml --no-color report element-total calories
-1200.00	day/nonworking
15.00	vegetables/lettuce/romaine/100g
155.00	egg/boiled/100g
259.00	bread/rye/100g
273.60	sandwich/egg/lettuce/100g
280.00	candy/snickers/bar
680.00	sauce/mayonnaise/100g

```

### Unresolved

Prints list of records which were not found in the recipe book

```sh
$ ./hranoprovod-cli -d examples/food.yaml -l examples/log.yaml --no-color report unresolved
coffee/cup

```

### CSV Export of the logged items

```sh
$ ./hranoprovod-cli -d examples/food.yaml -l examples/log.yaml --no-color csv log
2021-01-24;day/nonworking;1.00
2021-01-24;coffee/cup;1.00
2021-01-24;sandwich/egg/lettuce/100g;1.20
2021-01-24;candy/snickers/bar;1.00
2021-01-25;day/nonworking;1.00
2021-01-25;coffee/cup;2.00
2021-01-25;sandwich/egg/lettuce/100g;1.50

```

### Data statistics

```sh
$ ./hranoprovod-cli -d examples/food.yaml -l examples/log.yaml --no-color stats
  Database file:      examples/food.yaml
  Database records:   7

  Log file:           examples/log.yaml
  Log records:        2
  First record:       2021/01/24 (340 days ago)
  Last record:        2021/01/25 (339 days ago)

```

### Total quantities per food

```sh
$ ./hranoprovod-cli -d examples/food.yaml -l examples/log.yaml --no-color report quantity --desc
3.00	coffee/cup
2.70	sandwich/egg/lettuce/100g
2.00	day/nonworking
1.00	candy/snickers/bar

```
