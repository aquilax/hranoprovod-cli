#!/usr/bin/env bash

########################
# include the magic
########################
. demo-magic.sh

########################
# Configure the options
########################

# hide the evidence
clear

p "# Welcome"
p "hranoprovod-cli is a command line tool which uses plain text files to keep a record of your diet and excercise"
p ""

p "# Data"
p "hranoprovod-cli uses two files as a database. Both files use a subset of yaml with some modifications to improve the usability"
p "food.yaml contains the ingredients/recipes"

## list the database file
pei "cat examples/food.yaml"

p "log.yaml contans your journal of consumption"

## list the log journalfile
pei "cat examples/log.yaml"
p ""

## Register report
p "Register report"
p ""
p "here is how you generate a report for all days in the journal"
pei "./hranoprovod-cli -d examples/food.yaml -l examples/log.yaml reg"

p "positive and negative contributions are tracked separately"
p "in that way eating a candy can add to the calories balance"
p "while excersising can remove from it"
p ""

## Balance report
p "Balance report"
p ""
p "you can aslo generate a balance report running the following command"
pei "./hranoprovod-cli -d examples/food.yaml -l examples/log.yaml bal"