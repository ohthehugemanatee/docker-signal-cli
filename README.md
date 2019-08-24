# Signal Nixplay bridge

This project creates a bridge between my Signal messenger Groups and my Nixplay account. It simply logs into my signal, and when it sees an image in the target group, it automatically emails it to Nixplay. At least, it WILL email it. For now it just calls it out in a friendly CLI note.

## Working

* Signal-CLI client can connect and get messages
* Golang wrapper around signal-cli to pick out the attachment filenames we care about

## ToDo
* Send email with the files.


## Installation
* Download and install signal-cli on your host machine somewhere.
* Register your new device with signal-cli on the host machine (see instructions in that repo)
* Use run.sh to actually run the container. Make sure to set environment variables first:
  * `$MYPHONE`: phone number of the receiving Signal account
  * `$GROUPID`: The Signal group ID of the group which receives pictures
  * `$DATADIR`: The directory on the host machine where you want to store persistent data