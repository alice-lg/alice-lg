#!/usr/bin/env sh

######################################################################
# @author      : annika
# @file        : init
# @created     : Tuesday Jan 11, 2022 15:35:20 CET
# @description :  Initialize the database
######################################################################

if [ -z $PSQL ]; then
    PSQL="psql"
fi

if [ -z $PGHOST ]; then
    export PGHOST="localhost"
fi

if [ -z $PGPORT ]; then
    export PGPORT="5432"
fi

if [ -z $PGDATABASE ]; then
    export PGDATABASE="alice"
fi

if [ -z $PGUSER ]; then
    export PGUSER="postgres"
fi

if [ -z $PGPASSWORD ]; then
    export PGPASSWORD="postgres"
fi

## Commandline opts: 
OPT_USAGE=0
OPT_TESTING=0

while [ $# -gt 0 ]; do
  case "$1" in
    -h) OPT_USAGE=1 ;;
    -t) OPT_TESTING=1 ;;
  esac
  shift
done

if [ $OPT_USAGE -eq 1 ]; then
    echo "Options:"
    echo "   -t     Use test database"
    exit
fi

if [ $OPT_TESTING -eq 1 ]; then
    echo "++ using test database"
    NAME="${PGDATABASE}_test"
    export PGDATABASE=$NAME
fi

psql

