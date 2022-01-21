#!/bin/bash

######################################################################
# @author      : Annika Hannig
# @file        : entrypoint
# @created     : Thursday Sep 23, 2021 19:24:51 CEST
#
# @description : Start a command or start the dev server
#   and ensure all dependencies are installed.
######################################################################

cd /ui

CMD=$1
shift
ARGS=$@

case "$CMD" in
    devserver)
        yarn install
        yarn start $ARGS
        ;;
    test)
        yarn install
        yarn test $ARGS
        ;;
    *)
        # Just run the command
        $CMD $ARGS
esac

