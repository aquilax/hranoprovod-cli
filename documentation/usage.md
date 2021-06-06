# Usage

Given the following log file:

`$ cat examples/log.yaml`

and the following recipe book file:

`$ cat examples/food.yaml`

## Register report

### Register view (default template)

`$ ./hranoprovod-cli -d examples/food.yaml -l examples/log.yaml --no-color reg`

### Register view (left aligned template)

`$ ./hranoprovod-cli -d examples/food.yaml -l examples/log.yaml --no-color reg --internal-template-name=left-aligned`

### Register report for single food per day

`$ ./hranoprovod-cli -d examples/food.yaml -l examples/log.yaml --no-color reg --single-food coffee/cup`

### Register report for total of ingredients for a day grouped by food

`$ ./hranoprovod-cli -d examples/food.yaml -l examples/log.yaml -b 2021/01/24 -e 2021/01/24 --no-color reg -g -s calories`

## Balance

### Default balance report

`$ ./hranoprovod-cli -d examples/food.yaml -l examples/log.yaml --no-color bal`

### Default balance report (collapse)

`$ ./hranoprovod-cli -d examples/food.yaml -l examples/log.yaml --no-color bal -c`

### Default balance report (collapse last)

`$ ./hranoprovod-cli -d examples/food.yaml -l examples/log.yaml --no-color bal --collapse-last`

### Default balance report (single element)

`$ ./hranoprovod-cli -d examples/food.yaml -l examples/log.yaml --no-color bal -s fat`


## Other reports

### Element-total

Calories for each defined food:

`$ ./hranoprovod-cli -d examples/food.yaml -l examples/log.yaml --no-color report element-total calories`

### Unresolved

Prints list of records which were not found in the recipe book

`$ ./hranoprovod-cli -d examples/food.yaml -l examples/log.yaml --no-color report unresolved`

### CSV Export of the logged items

`$ ./hranoprovod-cli -d examples/food.yaml -l examples/log.yaml --no-color csv log`

### Data statistics

`$ ./hranoprovod-cli -d examples/food.yaml -l examples/log.yaml --no-color stats`

### Total quantities per food

`$ ./hranoprovod-cli -d examples/food.yaml -l examples/log.yaml --no-color report quantity --desc`
