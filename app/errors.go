package app

import "errors"

var COMMAND_NOT_PROVIDED = errors.New("Command not provided in the command line")
var COMMAND_NOT_FOUND = errors.New("Command not found in the cache")
