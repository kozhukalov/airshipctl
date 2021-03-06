/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"errors"
)

// RunSetContext validates the given command line options and invokes AddContext/ModifyContext
func RunSetContext(o *ContextOptions, airconfig *Config, writeToStorage bool) (bool, error) {
	modified := false
	err := o.Validate()
	if err != nil {
		return modified, err
	}
	if o.Current {
		if airconfig.CurrentContext == "" {
			return modified, ErrMissingCurrentContext{}
		}
		// when --current flag is passed, use current context
		o.Name = airconfig.CurrentContext
	}

	context, err := airconfig.GetContext(o.Name)
	if err != nil {
		var cerr ErrMissingConfig
		if !errors.As(err, &cerr) {
			// An error occurred, but it wasn't a "missing" config error.
			return modified, err
		}

		if o.CurrentContext {
			return modified, ErrMissingConfig{}
		}
		// context didn't exist, create it
		// ignoring the returned added context
		airconfig.AddContext(o)
	} else {
		// Found the desired Current Context
		// Lets update it and be done.
		if o.CurrentContext {
			airconfig.CurrentContext = o.Name
		} else {
			// Context exists, lets update
			airconfig.ModifyContext(context, o)
		}
		modified = true
	}
	// Update configuration file just in time persistence approach
	if writeToStorage {
		if err := airconfig.PersistConfig(true); err != nil {
			// Error that it didnt persist the changes
			return modified, ErrConfigFailed{}
		}
	}

	return modified, nil
}

// RunUseContext validates the given context name and updates it as current context
func RunUseContext(desiredContext string, airconfig *Config) error {
	if _, err := airconfig.GetContext(desiredContext); err != nil {
		return err
	}

	if airconfig.CurrentContext != desiredContext {
		airconfig.CurrentContext = desiredContext
		if err := airconfig.PersistConfig(true); err != nil {
			return err
		}
	}
	return nil
}

// RunSetManifest validates the given command line options and invokes AddManifest/ModifyManifest
func RunSetManifest(o *ManifestOptions, airconfig *Config, writeToStorage bool) (bool, error) {
	modified := false
	err := o.Validate()
	if err != nil {
		return modified, err
	}

	manifest, exists := airconfig.Manifests[o.Name]
	if !exists {
		// manifest didn't exist, create it
		// ignoring the returned added manifest
		airconfig.AddManifest(o)
	} else {
		// manifest exists, lets update
		err = airconfig.ModifyManifest(manifest, o)
		if err != nil {
			return modified, err
		}
		modified = true
	}
	// Update configuration file just in time persistence approach
	if writeToStorage {
		if err := airconfig.PersistConfig(true); err != nil {
			// Error that it didnt persist the changes
			return modified, ErrConfigFailed{}
		}
	}

	return modified, nil
}

// RunSetEncryptionConfig validates the given command line options
// and invokes AddEncryptionConfig/ModifyEncryptionConfig
func RunSetEncryptionConfig(o *EncryptionConfigOptions, airconfig *Config, writeToStorage bool) (bool, error) {
	modified := false
	err := o.Validate()
	if err != nil {
		return modified, err
	}

	encryptionConfig, exists := airconfig.EncryptionConfigs[o.Name]
	if !exists {
		// encryption config didn't exist, create it
		// ignoring the returned added encryption config
		airconfig.AddEncryptionConfig(o)
		modified = true
	} else {
		// encryption config exists, lets update
		airconfig.ModifyEncryptionConfig(encryptionConfig, o)
		modified = true
	}
	// Update configuration file just in time persistence approach
	if writeToStorage {
		if err := airconfig.PersistConfig(true); err != nil {
			// Error that it didnt persist the changes
			return modified, ErrConfigFailed{}
		}
	}

	return modified, nil
}
