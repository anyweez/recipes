package config

import (
	"testing"
)

/**
 * Check to make sure an error is properly returned if the requested file
 * doesn't exist.
 *
 * Expected:
 * 		- first param is a RecipeConfig with no values set
 * 		- second param is an error
 */
func TestNoConfigFile(t *testing.T) {
	_, err := New("")

	// Check to make sure that an error was returned.
	if err == nil {
		t.Error()
	}
}

/**
 * Check to make sure that if one of the sections doesn't exist in the
 * configuration then we throw an error.
 *
 * Expected:
 * 		- same as TestNoConfigFile
 */

/*
func TestMissingSection(t *testing.T) {
	_, err := New("../../../test/missing-section-recipes.conf")

	if err == nil {
		t.Error("Missing section doesn't generate error.")
	}
}
*/

/**
 * Check to make sure that if a field doesn't exist in the configuration
 * then we throw an error.
 *
 * Expected:
 * 		- same as TestNoConfigFile
 */
/*
func TestMissingField(t *testing.T) {
	_, err := New("../../../test/missing-field-recipes.conf")

	if err == nil {
		t.Error("Missing fields don't generate error.")
	}
}
*/
/**
 * Test to make sure that if the data is invalid (wrong data types in this
 * case) that an error is properly returned.
 *
 * Expected:
 * 		- same as TestNoConfigFile
 */
func TestInvalidData(t *testing.T) {
	_, err := New("../../../test/invalid-data-recipes.conf")

	if err == nil {
		t.Error()
	}
}

/**
 * Test to ensure that we can print the connection string.
 *
 * Expected:
 * 		- Configuration file is read successfully.
 * 		- No error is returned.
 * 		- Output string properly generated without throwing an exception.
 */
func TestCreateConnectionString(t *testing.T) {
	conf, err := New("../../../test/recipes.conf")

	// This configuration is expected to load correctly.
	if err != nil {
		t.Error(err.Error())
	}

	// Some back and forth to appease Go's rigorous compiler warnings.
	str := conf.Mongo.ConnectionString()
	str = str
}
