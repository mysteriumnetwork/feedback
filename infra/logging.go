package infra

import (
	"bytes"
	"flag"
	"fmt"
	"strings"
	"text/template"

	"github.com/cihub/seelog"
)

// LogOptions seelog options
type LogOptions struct {
	logLevelInt seelog.LogLevel
	LogLevel    string
}

// CurrentLogOptions current seelog options
var CurrentLogOptions = &LogOptions{
	logLevelInt: seelog.DebugLvl,
	LogLevel:    seelog.DebugStr,
}

var (
	allLevels    = []string{seelog.TraceStr, seelog.DebugStr, seelog.InfoStr, seelog.WarnStr, seelog.ErrorStr, seelog.CriticalStr, seelog.OffStr}
	logLevelFlag string
)

// RegisterLoggerFlags registers logger CLI flags
func RegisterLoggerFlags() {
	flag.StringVar(&logLevelFlag, "log-level", seelog.DebugStr, fmt.Sprintf("Service logging level (%s)", strings.Join(allLevels, "|")))
}

// ConfigureLogger configures options using parsed flag values
func ConfigureLogger() {
	level := logLevelFlag
	levelInt, found := seelog.LogLevelFromString(level)
	if !found {
		levelInt = seelog.DebugLvl
		level = seelog.DebugStr
	}
	(*CurrentLogOptions).logLevelInt = levelInt
	(*CurrentLogOptions).LogLevel = level
}

const seewayLogXMLConfigTemplate = `
<seelog minlevel="{{.LogLevel}}">
	<outputs formatid="main">
		<console/>
	</outputs>
	<formats>
		<format id="main" format="%UTCDate(2006-01-02T15:04:05.999999999) [%Level] [%Func] %Msg%n"/>
	</formats>
</seelog>
`

func buildSeelogConfig(opts LogOptions) string {
	tmpl := template.Must(template.New("seelogcfg").Parse(seewayLogXMLConfigTemplate))

	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, opts)
	if err != nil {
		panic(err)
	}

	return tpl.String()
}

// BootstrapLogger loads seelog package into the overall system
func BootstrapLogger(opts *LogOptions) {
	if opts != nil {
		CurrentLogOptions = opts
	}
	newLogger, err := seelog.LoggerFromConfigAsString(buildSeelogConfig(*CurrentLogOptions))
	if err != nil {
		_ = seelog.Warn("Error parsing seelog configuration", err)
		return
	}
	err = seelog.UseLogger(newLogger)
	if err != nil {
		_ = seelog.Warn("Error setting new logger for seelog", err)
	}
	seelog.Infof("Log level: %s", CurrentLogOptions.LogLevel)
}
