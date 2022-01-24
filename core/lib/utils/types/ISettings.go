package types

import "strconv"

type SettingsManager struct {
	settings map[string]string
}

func NewSettingsManager(settings map[string]string) *SettingsManager {
	return &SettingsManager{settings}
}

// String returns the string value of a setting. If that setting is not found it returns defaultValue.
func (sm *SettingsManager) String(name, defaultValue string) string {

	if _, ok := sm.settings[name]; !ok {
		return defaultValue
	}

	return sm.settings[name]
}

// Int64 returns the value of a setting converted to int64 if possible.
// If the setting value does not exist, or parsing the string value to int64 returns an error, it returns defaultValue.
func (sm *SettingsManager) Int64(name string, defaultValue int64) int64 {

	if _, ok := sm.settings[name]; !ok {
		return defaultValue
	}

	if intVal, e := strconv.ParseInt(sm.settings[name], 10, 64); e != nil {
		return defaultValue
	} else {
		return intVal
	}

}

// Bool returns true or false based on whether the value is the string "1"
// If the setting does not exist, it returns defaultValue
func (sm *SettingsManager) Bool(name string, defaultValue bool) bool {

	if _, ok := sm.settings[name]; !ok {
		return defaultValue
	}

	return sm.settings[name] == "1"

}

// Float64 returns the value of a setting converted to float64 if possible
// If the setting value does not exist, or parsing the string value to float64 returns an error, it returns defaultValue
func (sm *SettingsManager) Float64(name string, defaultValue float64) float64 {

	if _, ok := sm.settings[name]; !ok {
		return defaultValue
	}

	if floatVal, e := strconv.ParseFloat(sm.settings[name], 64); e != nil {
		return defaultValue
	} else {
		return floatVal
	}

}
