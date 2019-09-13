/*
 * Copyright (C) 2019 The "MysteriumNetwork/feedback" Authors.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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
