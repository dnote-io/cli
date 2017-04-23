# Dnote CLI

![Dnote](assets/main.png)

A command line interface for spontaneously capturing the things you learn while coding

## Installation

On macOS, or Linux, run:

    curl -s https://raw.githubusercontent.com/dnote-io/cli/master/install.sh | sh

In some cases, you might need `sudo`. Feel free to inspect [install.sh](https://github.com/dnote-io/cli/blob/master/install.sh):

    curl -s https://raw.githubusercontent.com/dnote-io/cli/master/install.sh | sudo sh

On Windows, download [binary](https://github.com/dnote-io/cli/releases)

## Overview

Dnote categorizes your **notes** by **books**.

All your books and notes are stored in `$HOME/.dnote` as a JSON file.

You can optionally sync your note with Dnote server. Syncing will allow you to interact with your notes using the web frontend at https://dnote.io, and set up digest notifications.

## Commands

### dnote use [book name]
*alias: u*

Change the book to write your note in.

e.g.

    dnote use linux

### dnote new "[note]"
*alias: n*

Write a new note under the current book with auto-generated note names.

e.g.

    dnote new "set -e instructs bash to exit immediately if any command has non-zero exit status"

**dnote new -t "[note name or UID]" "[note content]"**

Write a new note under the current book with custom note names.

e.g

    dnote new -t "note_001" "set -e instructs bash to exit immediately if any command has non-zero exit status"

**dnote edit "[note name or UID]" "[note content]"**

Edit a note using it's note name or UID.

e.g

    dnote edit "note_001" "chdir changes working directory"

**dnote delete "[note name or UID]"**

Delete a note using it's name or UID.

e.g

    dnote delete "note_001"

**dnote books**

List all the books that you created

e.g.

    $ dnote books
      javascript
    * linux
      tmux
      css

### dnote notes

List all notes in the current book

#### Options

* `-b [book name]`

Specify the name of the book to read from

e.g.

    $ dnote notes
    On note JS
    * .bind() creates a new function
    * arrow function uses less memory than function with .bind()
    * the time passed to setTimeout is minimum, no guaranteed time

### dnote sync

Sync notes with Dnote server

### dnote login

Start a login procedure which will store the APIKey to communicate with the server


## Links

* [Website](https://dnote.io)
* [Making Dnote (blog article)](https://sungwoncho.io/making-dnote/)

## License

MIT

-------

> Made by [sung](https://sungwoncho.io)
