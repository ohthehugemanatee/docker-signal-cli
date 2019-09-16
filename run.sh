#!/bin/sh

# Signal Photo Mailer start script.
# Add your configuration to `.env` (you can use the included env.example as a
# starter), and run this script to start the mailer.
#
# You will probably want to configure the client first. You can do that by
# appending commands to `signal-cli` as arguments to this script. e.g.:
# `./run.sh signal-cli -u +49123456789 register`
# `./run.sh signal-cli -u +49123456789 listGroups`

DIR="$(cd -P "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Sets up a signal listener which will dump content to stdout, and attachments to the mapped directory.
. $DIR/.env

# You have to set the MYPHONE env variable to your phone number, in format +4912345678901
if [ -z ${MYPHONE} ]; then
  echo "you must set the myphone variable to your phone number, in format +4912345678901."
  exit 1
fi

# You have to set the target Group ID.
if [ -z ${GROUPID} ]; then
  echo "you must set the GROUPID variable to the Group which contains your pictures."
  exit 1
fi

if [ -z ${DESTMAIL} ]; then
  echo "you must set the DESTMAIL variable to the destination Nixplay email address."
  exit 1
fi

if [ -z ${SMTPUSER} ]; then
  echo "you must set the SMTPUSER variable to the SMTP username for sending."
  exit 1
fi
if [ -z ${SMTPPASS} ]; then
  echo "you must set the SMTPPASS variable to the SMTP password for sending."
  exit 1
fi
if [ -z ${SMTPSERVER} ]; then
  echo "you must set the SMTPSERVER variable to the SMTP server for sending."
  exit 1
fi
if [ -z ${SMTPFROM} ]; then
  echo "you must set the SMTPFROM variable to the SMTP 'from' email for sending."
  exit 1
fi

# You have to choose where to save all of signal's stuff. A reasonable default is ~/.local/share/signal-cli , but whatever man.
if [ -z ${DATADIR} ]; then
  echo "you must set the DATADIR variable to the local directory where you want to save Signal stuff."
  exit 1
fi

mkdir -p $HOME/.local/share/signal-cli
docker run -v $DATADIR:/root/.local/share/signal-cli -e MYPHONE=${MYPHONE} -e GROUPID=${GROUPID} -e DESTMAIL=${DESTMAIL} -e SMTPUSER=${SMTPUSER} -e SMTPPASS=${SMTPPASS} -e SMTPSERVER=${SMTPSERVER} -e SMTPFROM=${SMTPFROM} -e SMTPPORT=${SMTPPORT}  ohthehugemanatee/signal-photo-mail "$@"
