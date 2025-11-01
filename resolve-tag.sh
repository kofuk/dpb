#!/usr/bin/env bash

app=$1

base_tag=$(head -1 ${app}/Dockerfile | sed -En 's/^FROM [^:]+:([-0-9a-z.]+).+$/\1/p')
cat ${app}/TAG | sed "s/^%base%/${base_tag}/"
