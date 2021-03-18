// Copyright 2021 OpenIndex.de.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"log/syslog"
	"os"
	"path/filepath"
	"strings"
)

// Syslog connection, used in non-debug mode.
var logger *syslog.Writer

// If an error occurred, log it and shutdown the application.
func errorFail(err error) {
	if err != nil {
		logFatal(err)
	}
}

// If an error occurred, log it and shutdown the application.
//goland:noinspection GoUnusedFunction
func errorWarn(err error) {
	if err != nil {
		logWarn(err)
	}
}

// Test, if the application is running in debug mode.
func isDebug() bool {
	if len(os.Args) < 2 {
		return false
	}
	for i := 1; i < len(os.Args); i++ {
		arg := strings.ToLower(strings.TrimSpace(os.Args[i]))
		if arg == "debug" {
			return true
		}
	}
	return false
}

// Test, if a string is empty.
func isEmpty(text string) bool {
	return len(text) == 0
}

// Test, if a path points to a file.
func isFile(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// Log an execution error.
func logExecutionError(err interface{}, javaHome string, javaCommand string, javaArguments []string) {
	logFatal(fmt.Sprintf("Execution failed: %s\nJava home: %s\nJava command: %s\nJava arguments: %s",
		fmt.Sprintf("%s", err), javaHome, javaCommand, javaArguments))
}

// Log a fatal error and shutdown the application.
func logFatal(err interface{}) {
	// Log to syslog.
	if logger != nil {
		e := logger.Err(fmt.Sprintf("%s", err))
		if e == nil {
			os.Exit(1)
		}
	}

	// Otherwise log to console.
	log.Println(fmt.Sprintf("[ERROR] %s", err))
	os.Exit(1)
}

// Log an information.
func logInfo(message interface{}) {
	// Log to syslog.
	if logger != nil {
		e := logger.Info(fmt.Sprintf("%s", message))
		if e == nil {
			return
		}
	}

	// Otherwise log to console.
	log.Println(fmt.Sprintf("[INFO] %s", message))
}

// Log a warning.
func logWarn(message interface{}) {
	// Log to syslog.
	if logger != nil {
		e := logger.Warning(fmt.Sprintf("%s", message))
		if e == nil {
			return
		}
	}

	// Otherwise log to console.
	log.Println(fmt.Sprintf("[WARNING] %s", message))
}

// Close connection to syslog.
func loggerClose() {
	if logger != nil {
		err := logger.Close()
		logger = nil
		if err != nil {
			logFatal(err)
		}
	}
}

// Create syslog logger, if the application is not running debug mode.
func loggerInit() {
	if !isDebug() {
		l, err := syslog.New(syslog.LOG_ERR, filepath.Base(filepath.Dir(filepath.Dir(filepath.Dir(os.Args[0])))))
		errorFail(err)
		logger = l
		defer loggerClose()
	}
}

// Get java home path from JAVA_HOME environment variable.
func getJavaHomeViaEnv() string {
	javaHome, javaHomeFound := os.LookupEnv("JAVA_HOME")
	if !javaHomeFound {
		return ""
	}
	return strings.TrimSpace(javaHome)
}

// Get user language from LANG environment variable.
func getLanguageViaEnv() string {
	language, languageFound := os.LookupEnv("LANG")
	if !languageFound {
		return ""
	}

	language = strings.TrimSpace(language)
	if isEmpty(language) {
		return ""
	}

	// Normalise language - e.g. "de_DE.UTF8" => "de-DE"
	if strings.Contains(language, ".") {
		language = strings.Split(language, ".")[0]
	}
	language = strings.ReplaceAll(language, "_", "-")
	if strings.Contains(language, "-") {
		parts := strings.Split(language, "-")
		language = fmt.Sprintf("%s-%s", parts[0], parts[1])
	}

	return language
}
