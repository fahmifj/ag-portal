package logger

import "github.com/hashicorp/go-hclog"

// Shared logger
var Log = hclog.New(&hclog.LoggerOptions{
	Level: hclog.DefaultLevel,
	Color: hclog.AutoColor,
})
