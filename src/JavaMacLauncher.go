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
	"bytes"
	"fmt"
	"howett.net/plist"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const VERSION = "1.1.0-SNAPSHOT"

// Main structure of a plist file.
type plistDataType struct {
	BundleName        string        `plist:"CFBundleName"`
	BundleDisplayName string        `plist:"CFBundleDisplayName"`
	BundleVersion     string        `plist:"CFBundleVersion"`
	BundleIconFile    string        `plist:"CFBundleIconFile"`
	JavaMacLauncher   plistJavaType `plist:"JavaMacLauncher"`
}

// JavaMacLauncher configuration in the plist file.
type plistJavaType struct {
	JavaHome             string            `plist:"JavaHome"`
	JavaCommand          string            `plist:"JavaCommand"`
	JavaOptions          []string          `plist:"JavaOptions"`
	JavaClassPath        []string          `plist:"JavaClassPath"`
	JavaModulePath       []string          `plist:"JavaModulePath"`
	WorkingDirectory     string            `plist:"WorkingDirectory"`
	ApplicationCommand   string            `plist:"ApplicationCommand"`
	ApplicationArguments []string          `plist:"ApplicationArguments"`
	HeapMinimum          string            `plist:"HeapMinimum"`
	HeapMaximum          string            `plist:"HeapMaximum"`
	SplashImage          string            `plist:"SplashImage"`
	DockName             map[string]string `plist:"DockName"`
	DockIcon             string            `plist:"DockIcon"`
	UseScreenMenuBar     bool              `plist:"UseScreenMenuBar"`
	LaunchInForeground   bool              `plist:"LaunchInForeground"`
}

// Execution starts here.
func main() {
	// Test, if the application is running in debug mode.
	debug := isDebug()

	// Init logger.
	loggerInit()

	logInfo(fmt.Sprintf("Entering JavaMacLauncher v%s...", VERSION))
	//for _, arg := range os.Args {
	//	logInfo(fmt.Sprintf("> %s", arg))
	//}
	if debug {
		logInfo(fmt.Sprintf("Running in debug mode."))
	}

	// Determine application location.
	launcherPath, err := filepath.Abs(os.Args[0])
	errorFail(err)
	launcherPath = filepath.Clean(launcherPath)
	macOSPath := filepath.Dir(launcherPath)
	contentsPath := filepath.Dir(macOSPath)
	appPath := filepath.Dir(contentsPath)
	logInfo(fmt.Sprintf("Bundle: %s", appPath))
	plistPath := filepath.Join(contentsPath, "Info.plist")
	if !isFile(plistPath) {
		logFatal(fmt.Sprintf("Can't find \"Info.plist\" at: %s", plistPath))
	}
	logInfo(fmt.Sprintf("Info.plist: %s", plistPath))

	// Read details about the application bundle from Info.plist.
	var plistData plistDataType
	plistFile, err := os.Open(plistPath)
	errorFail(err)
	plistDecoder := plist.NewDecoder(plistFile)
	err = plistDecoder.Decode(&plistData)
	errorFail(err)
	err = plistFile.Close()
	errorFail(err)

	// Get java home & make relative path absolute.
	javaHome := strings.TrimSpace(plistData.JavaMacLauncher.JavaHome)
	if !isEmpty(javaHome) && !strings.HasPrefix(javaHome, "/") {
		javaHome = filepath.Join(appPath, javaHome)
	}

	// Get java command & make relative path absolute.
	javaCommand := strings.TrimSpace(plistData.JavaMacLauncher.JavaCommand)
	if !isEmpty(javaCommand) && !strings.HasPrefix(javaCommand, "/") {
		javaCommand = filepath.Join(appPath, javaCommand)
	}

	// If Java command & Java home are empty, try to detect it automatically.
	if isEmpty(javaCommand) && isEmpty(javaHome) {
		javaHome = getJavaHome()
		if isEmpty(javaHome) {
			logFatal("Unable to obtain java home.")
		}
	}

	// If Java command is empty, get it from Java home.
	if isEmpty(javaCommand) && !isEmpty(javaHome) {
		javaCommand = filepath.Join(javaHome, "bin", "java")
	}

	// Set Java Home environment variable, if available.
	if !isEmpty(javaHome) {
		logInfo(fmt.Sprintf("Java home: %s", javaHome))
		err = os.Setenv("JAVA_HOME", javaHome)
		errorFail(err)
	}

	// Make sure a Java command is available.
	if isEmpty(javaCommand) {
		logFatal("Can't find a Java command.")
	} else if !isFile(javaCommand) {
		logFatal(fmt.Sprintf("Java command does not point to a file: %s", javaCommand))
	} else {
		logInfo(fmt.Sprintf("Java command: %s", javaCommand))
	}

	// Get working directory & make relative path absolute.
	// Users home directory is used, if empty or not configured.
	workingDirectory := strings.TrimSpace(plistData.JavaMacLauncher.WorkingDirectory)
	if isEmpty(workingDirectory) {
		home, err := os.UserHomeDir()
		errorFail(err)
		workingDirectory = home
	} else if !strings.HasPrefix(workingDirectory, "/") {
		workingDirectory = filepath.Join(appPath, workingDirectory)
	}

	var javaArguments []string

	// Add minimum heap size to arguments.
	heapMinimum := strings.TrimSpace(plistData.JavaMacLauncher.HeapMinimum)
	if !isEmpty(heapMinimum) {
		javaArguments = append(javaArguments, fmt.Sprintf("-Xms%s", heapMinimum))
	}

	// Add maximum heap size to arguments.
	heapMaximum := strings.TrimSpace(plistData.JavaMacLauncher.HeapMaximum)
	if !isEmpty(heapMaximum) {
		javaArguments = append(javaArguments, fmt.Sprintf("-Xmx%s", heapMaximum))
	}

	// Add splash image to arguments.
	splashImage := strings.TrimSpace(plistData.JavaMacLauncher.SplashImage)
	if !isEmpty(splashImage) {
		if !strings.HasPrefix(splashImage, "/") {
			splashImage = filepath.Join(appPath, splashImage)
		}
		javaArguments = append(javaArguments, fmt.Sprintf("-splash:%s", splashImage))
	}

	// Add dock name to arguments.
	dockName := getDockName(plistData)
	if !isEmpty(dockName) {
		javaArguments = append(javaArguments, fmt.Sprintf("-Xdock:name=%s", dockName))
	}

	// Add dock icon to arguments.
	dockIcon := getDockIcon(plistData)
	if !isEmpty(dockIcon) {
		if !strings.HasPrefix(dockIcon, "/") {
			dockIcon = filepath.Join(appPath, dockIcon)
		}
		javaArguments = append(javaArguments, fmt.Sprintf("-Xdock:icon=%s", dockIcon))
	}

	// Enable screen menu bar.
	if plistData.JavaMacLauncher.UseScreenMenuBar {
		javaArguments = append(javaArguments, "-Dapple.laf.useScreenMenuBar=true")
	}

	// Add Java options to arguments.
	for i := 0; i < len(plistData.JavaMacLauncher.JavaOptions); i++ {
		javaOption := strings.TrimSpace(plistData.JavaMacLauncher.JavaOptions[i])
		if !isEmpty(javaOption) {
			javaArguments = append(javaArguments, javaOption)
		}
	}

	// Add Java class path to arguments.
	javaClassPath := ""
	for i := 0; i < len(plistData.JavaMacLauncher.JavaClassPath); i++ {
		path := strings.TrimSpace(plistData.JavaMacLauncher.JavaClassPath[i])
		if isEmpty(path) {
			continue
		}
		if !strings.HasPrefix(path, "/") {
			path = filepath.Join(appPath, path)
		}
		if !isEmpty(javaClassPath) {
			javaClassPath += ":"
		}
		javaClassPath += path
	}
	if !isEmpty(javaClassPath) {
		javaArguments = append(javaArguments, "--class-path", javaClassPath)
	}

	// Add Java module path to arguments.
	javaModulePath := ""
	for i := 0; i < len(plistData.JavaMacLauncher.JavaModulePath); i++ {
		path := strings.TrimSpace(plistData.JavaMacLauncher.JavaModulePath[i])
		if isEmpty(path) {
			continue
		}
		if !strings.HasPrefix(path, "/") {
			path = filepath.Join(appPath, path)
		}
		if !isEmpty(javaModulePath) {
			javaModulePath += ":"
		}
		javaModulePath += path
	}
	if !isEmpty(javaModulePath) {
		javaArguments = append(javaArguments, "--module-path", javaModulePath)
	}

	// Add application command to arguments.
	applicationCommand := strings.TrimSpace(plistData.JavaMacLauncher.ApplicationCommand)
	if !isEmpty(applicationCommand) {
		if strings.HasSuffix(strings.ToLower(applicationCommand), ".jar") {
			if !strings.HasPrefix(applicationCommand, "/") {
				applicationCommand = filepath.Join(appPath, applicationCommand)
			}
			javaArguments = append(javaArguments, "-jar", applicationCommand)
		} else {
			parts := strings.Split(applicationCommand, " ")
			for i := 0; i < len(parts); i++ {
				javaArguments = append(javaArguments, parts[i])
			}
		}
	}

	// Add application arguments to arguments.
	for i := 0; i < len(plistData.JavaMacLauncher.ApplicationArguments); i++ {
		applicationArgument := strings.TrimSpace(plistData.JavaMacLauncher.ApplicationArguments[i])
		if !isEmpty(applicationArgument) {
			javaArguments = append(javaArguments, applicationArgument)
		}
	}

	// Change to working directory.
	err = os.Chdir(workingDirectory)
	errorFail(err)

	// Prepare command.
	logInfo(fmt.Sprintf("Java arguments: %s", javaArguments))
	command := exec.Command(javaCommand, javaArguments...)
	//command := exec.Command(javaCommand, "-version")
	command.Dir = workingDirectory

	if debug || plistData.JavaMacLauncher.LaunchInForeground {
		//
		// Execute command in foreground.
		//

		// Fetch command's STDOUT & STDERR in debug mode.
		var commandOut bytes.Buffer
		var commandErr bytes.Buffer
		if debug {
			command.Stdout = &commandOut
			command.Stderr = &commandErr
		}

		// Execute command and wait.
		err = command.Run()
		if err != nil {
			if debug {
				stdout := strings.TrimSpace(commandOut.String())
				if len(stdout) > 0 {
					logInfo(fmt.Sprintf("STDOUT:\n%s\n", stdout))
				}

				stderr := strings.TrimSpace(commandErr.String())
				if len(stderr) > 0 {
					logInfo(fmt.Sprintf("STDERR:\n%s\n", stderr))
				}
			}

			logExecutionError(err, javaHome, javaCommand, javaArguments)
		}
	} else {
		//
		// Execute command in background.
		//

		err = command.Start()
		if err != nil {
			logExecutionError(err, javaHome, javaCommand, javaArguments)
		}
	}

	logInfo("Exiting JavaMacLauncher. Have a nice day!")
}
