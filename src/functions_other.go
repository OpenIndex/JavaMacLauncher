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

// +build !darwin

package main

// Always return an empty dock name on non-Mac systems.
//goland:noinspection GoUnusedParameter
func getDockName(plistData plistDataType) string {
	return ""
}

// Always return an empty dock icon on non-Mac systems.
//goland:noinspection GoUnusedParameter
func getDockIcon(plistData plistDataType) string {
	return ""
}

// Get java home path on non-Mac systems.
func getJavaHome() string {
	return getJavaHomeViaEnv()
}

// Get user language on non-Mac systems.
//goland:noinspection GoUnusedFunction
func getLanguage() string {
	return getLanguageViaEnv()
}
