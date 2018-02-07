package log

import (
	"os"

	"github.com/abchain/fabric/peerex/logging"
	"github.com/spf13/viper"
)

var (
	modules      []string
	defaultLevel = logging.DEBUG
)

func init() {

	format := logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{level:.4s} [%{module:.6s}] %{shortfile} %{shortfunc} ▶ %{message}%{color:reset}`)

	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)

	modules = make([]string, 1, 1)
}

func InitLogger(module string) *logging.Logger {

	modules = append(modules, module)

	l := logging.MustGetLogger(module)
	logging.SetLevel(defaultLevel, module)

	return l
}

func SetLogLevel() error {

	// Read config
	switch viper.GetString("logging.level") {
	case "debug":
		defaultLevel = logging.DEBUG
	case "info":
		defaultLevel = logging.INFO
	case "notice":
		defaultLevel = logging.NOTICE
	case "warning":
		defaultLevel = logging.WARNING
	case "error":
		defaultLevel = logging.ERROR
	case "critical":
		defaultLevel = logging.CRITICAL
	default:
		defaultLevel = defaultLevel
	}

	// Set log level
	for _, m := range modules {
		logging.SetLevel(defaultLevel, m)
	}

	return nil
}

func SetModuleLogLevel(module string, level logging.Level) error {

	logging.SetLevel(level, module)

	return nil
}
