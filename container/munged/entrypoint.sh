#!/bin/bash

[ -n "${DEBUG}" ] && set -x

exec gosu munge /usr/sbin/munged -F
