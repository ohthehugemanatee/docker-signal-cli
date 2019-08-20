#!/bin/sh

# Sets up a signal listener which will dump content to stdout, and attachments to the mapped directory.

# You have to set the MYPHONE env variable to your phone number, in format +4912345678901
if [ -z ${MYPHONE} ]; then 
  echo "you must set the myphone variable to your phone number, in format +4912345678901."
  exit 1
fi

# You have to choose where to save all of signal's stuff. A reasonable default is ~/.local/share/signal-cli , but whatever man.
if [ -z ${DATADIR} ]; then 
  echo "you must set the DATADIR variable to the local directory where you want to save Signal stuff."
  exit 1
fi
mkdir -p $HOME/.local/share/signal-cli
docker run -v $DATADIR:/root/.local/share/signal-cli ohthehugemanatee/signal-cli signal-cli -u $MYPHONE receive -t -1 --json
