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

package logconfig

import (
	"flag"
	"fmt"
	"strings"

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

// RegisterFlags registers logger CLI flags
func RegisterFlags() {
	flag.StringVar(&logLevelFlag, "log-level", seelog.DebugStr, fmt.Sprintf("Service logging level (%s)", strings.Join(allLevels, "|")))
}

// Configure configures options using parsed flag values
func Configure() {
	level := logLevelFlag
	levelInt, found := seelog.LogLevelFromString(level)
	if !found {
		levelInt = seelog.DebugLvl
		level = seelog.DebugStr
	}
	(*CurrentLogOptions).logLevelInt = levelInt
	(*CurrentLogOptions).LogLevel = level
}
