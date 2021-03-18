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

// +build darwin

package main

import (
	"github.com/progrium/macdriver/objc"
	"os/exec"
	"strings"
)

// Returns the dock name based on current language.
func getDockName(plistData plistDataType) string {
	// Don't use a dock name, if not configured.
	if len(plistData.JavaMacLauncher.DockName) == 0 {
		return ""
	}

	// Get default dock name.
	defaultName, isMapContainsKey := plistData.JavaMacLauncher.DockName["default"]
	if isMapContainsKey {
		defaultName = strings.TrimSpace(defaultName)

		// Return the default name, if no other dock names are configured.
		if len(plistData.JavaMacLauncher.DockName) == 1 {
			return defaultName
		}
	}

	// Use BundleDisplayName as default, if no default was specified.
	if isEmpty(defaultName) {
		defaultName = strings.TrimSpace(plistData.BundleDisplayName)
	}

	// Use BundleName as default, if no default was specified.
	if isEmpty(defaultName) {
		defaultName = strings.TrimSpace(plistData.BundleName)
	}

	// Get user language.
	userLanguage := getLanguage()

	// Use default name, if no user language was found.
	if isEmpty(userLanguage) {
		return defaultName
	}

	// Convert language to lowercase to make lookups case insensitive.
	userLanguage = strings.ToLower(userLanguage)

	// Lookup for explicit language - e.g. "de-DE".
	for lang, dockName := range plistData.JavaMacLauncher.DockName {
		l := strings.ToLower(strings.TrimSpace(lang))
		if l == userLanguage {
			return strings.TrimSpace(dockName)
		}
	}

	// Lookup for simple language - e.g. "de".
	if strings.Contains(userLanguage, "-") {
		userLanguage = strings.TrimSpace(strings.Split(userLanguage, "-")[0])
		for lang, dockName := range plistData.JavaMacLauncher.DockName {
			l := strings.ToLower(strings.TrimSpace(lang))
			if l == userLanguage {
				return strings.TrimSpace(dockName)
			}
		}
	}

	return defaultName
}

// Returns the configured dock icon.
func getDockIcon(plistData plistDataType) string {
	return strings.TrimSpace(plistData.JavaMacLauncher.DockIcon)
}

// Get java home path on Mac systems.
func getJavaHome() string {
	var javaHome string

	javaHome = getJavaHomeViaEnv()
	if !isEmpty(javaHome) {
		return javaHome
	}

	javaHome = getJavaHomeViaLibexec()
	if !isEmpty(javaHome) {
		return javaHome
	}

	return ""
}

// Get java home path by executing "/usr/libexec/java_home".
// TODO: Might be improved by selecting a certain Java Version.
//goland:noinspection SpellCheckingInspection
func getJavaHomeViaLibexec() string {
	out, err := exec.Command("/usr/libexec/java_home").CombinedOutput()
	if (err != nil || len(out) == 0) {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// Get user language on Mac systems.
func getLanguage() string {
	var language string

	language = getLanguageViaEnv()
	if !isEmpty(language) {
		return language
	}

	language = getLanguageViaObjectiveC()
	if !isEmpty(language) {
		return language
	}

	language = getLanguageViaOsascript()
	if !isEmpty(language) {
		return language
	}

	return ""
}

// Get user language by calling macOS Foundation API via Objective C.
// https://developer.apple.com/documentation/foundation/nslocale
// https://developer.apple.com/documentation/foundation/nslocale/1415614-preferredlanguages
// https://developer.apple.com/documentation/foundation/nsarray/1412852-firstobject
//goland:noinspection SpellCheckingInspection
func getLanguageViaObjectiveC() string {
	language := objc.Get("NSLocale").Get("preferredLanguages").Get("firstObject").String()
	return strings.TrimSpace(language)
}

// Get user language by executing "/usr/bin/osascript".
//goland:noinspection SpellCheckingInspection
func getLanguageViaOsascript() string {
	out, err := exec.Command("/usr/bin/osascript", "-e", "user locale of (get system info)").CombinedOutput()
	if err != nil || len(out) == 0 {
		return ""
	}
	return strings.TrimSpace(string(out))
}
