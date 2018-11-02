package config

import (
	"github.com/op/go-logging"
)

// We mimic the level.go in op-logging but use a light-weight implements
// (no default formatters)

type ModuleLeveled struct {
	levels map[string]logging.Level
	logging.Backend
}

// duplicate from the default one and replace with another backend
func DuplicateLevelBackend(backend logging.Backend) *ModuleLeveled {

	return &ModuleLeveled{
		levels:  DefaultBackend.levels,
		Backend: backend,
	}
}

// GetLevel returns the log level for the given module.
func (l *ModuleLeveled) GetLevel(module string) logging.Level {
	level, exists := l.levels[module]
	if !exists {
		level, exists = l.levels[""]
		// no configuration exists, default to debug
		if !exists {
			level = logging.DEBUG
		}
	}
	return level
}

// SetLevel sets the log level for the given module.
func (l *ModuleLeveled) SetLevel(level logging.Level, module string) {
	l.levels[module] = level
}

// IsEnabledFor will return true if logging is enabled for the given module.
func (l *ModuleLeveled) IsEnabledFor(level logging.Level, module string) bool {
	return level <= l.GetLevel(module)
}

func (l *ModuleLeveled) Log(level logging.Level, calldepth int, rec *logging.Record) error {
	return l.Backend.Log(level, calldepth+1, rec)
}
