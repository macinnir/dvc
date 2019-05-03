package apitest

import (
	"fmt"
)

/**
 * Strings
 */

// SetString sets a string value at `key` globally
func (a *APITest) SetString(key, value string) string {
	a.stringVals[key] = value
	return value
}

// GetString gets the string value identified by `name` from the global cache
func (a *APITest) GetString(key string) string {
	val, ok := a.stringVals[key]
	if !ok {
		return ""
	}

	return val
}

/**
 * Ints
 */

// SetInt sets an int value at `name` globally
func (a *APITest) SetInt(name string, v int64) int64 {
	a.intVals[name] = v
	return v
}

// GetInt gets the int value at `name` globally
func (a *APITest) GetInt(name string) int64 {
	val, ok := a.intVals[name]
	if !ok {
		return 0
	}

	return val
}

// // SetRandInt generates a random number where 0 <= n <= max,
// // sets it to the number cache and returns the number
// func (a *APITest) SetRandInt(name string, max int) int64 {
// 	randNum := RandInt(max)
// 	a.SetInt(name, randNum)
// 	return randNum
// }

// Increment sets an integer in the number cache if it doesn't exist
// And then decrements it by 1
// The resulting value is returned
func (a *APITest) Increment(key string) int64 {

	val, ok := a.intVals[key]

	if !ok {
		val = 0
	}

	val++

	a.intVals[key] = val

	return val
}

// Decrement sets an integer in the number cache if it doesn't exist
// And then decrements it by 1
// The resulting value is returned
func (a *APITest) Decrement(key string) int64 {

	val, ok := a.intVals[key]

	if !ok {
		val = 0
	}

	val--

	a.intVals[key] = val

	return val
}

/**
 * Objects
 */

// SetObject sets an object at `name` globally
func (a *APITest) SetObject(key string, obj interface{}) {
	a.objectVals[key] = obj
}

// GetObject gets an object at `key` globally
func (a *APITest) GetObject(name string) interface{} {
	val, ok := a.objectVals[name]
	if !ok {
		a.logger.Error(fmt.Sprintf("Invalid object name `%s`", name))
		return nil
	}

	return val
}
