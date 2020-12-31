#!/usr/bin/env bash

dftPortZtf=8848
dftPortZd=8849
interval=7
nowTime=`date +"%Y-%m-%d %H:%M:%S"`
nowDate=`date +"%Y-%m-%d"`

PARAM_NAME=$1
PARAM_PORT=$2

if [ -z "$PARAM_NAME" ]; then
  echo "first parameter - name can't be empty"
  exit 1
fi

DIR="$( cd "$( dirname "$0"  )" && pwd  )"
PORT=`ps -ef | grep "$PARAM_NAME" | grep -v "grep" | grep -v ".sh" | awk '{print $10}'`
echo name: $PARAM_NAME, dir: $DIR, port $PORT

if [ -z "$PARAM_PORT" ]; then
    if [ -z "$PORT" ]; then
      if [ "$PARAM_NAME" = "ztf" ]; then
        PARAM_PORT="$dftPortZtf"
      else
        PARAM_PORT="$dftPortZd"
      fi
    else
      PARAM_PORT="$PORT"
    fi
fi

for var in 1 2

	do

    if [ -n "$PORT" ]; then

      if [ $PORT = $PARAM_PORT ]; then
        echo service is still alive

        if [ ! -f "$DIR/.upgraded"]; then
          echo sleep $interval second the $var time.
          sleep $interval
        else
          echo upgraded, force to restart.
          PORT= # next round, will cause restart
        fi

      else
        echo kill service on port $PORT.
        ps -ef | grep "$PARAM_NAME" | grep -v "grep" | grep -v ".sh" | awk '{print $2}' | xargs kill -9
        PORT=
      fi

    else

      echo $nowTime start service on port $PARAM_PORT in dir $DIR.
      cd $DIR
      nohup ./ztf -P $PARAM_PORT > nohup.log 2&>zenops-agent-$nowDate.log &

      rm -f "$DIR/.upgraded"
      echo ""
      break

    fi

  done