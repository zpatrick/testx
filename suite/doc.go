// Package suite is built as a way to provide test suite capabilities without
// interfering with the 'go test' ecosystem. Other solutions typically require
// the suite itself become the test entrypoint, making it more difficult or
// impossible to use existing tooling to do things like run specific tests.
//
// Instead, this package uses lazy loading to solve a number of issues:
//
// • Test entrypoints do not change: the test's suite is retrieved within the test function body.
// Tooling to run individual tests (or groups of tests) will work the same.
//
// • Suite setup & teardown methods are only executed on an as-needed basis.
// Iteration speed should remain as fast as possible during active development.
//
// The flow of the package is pretty simple:
//
// Step 1: Register suites in TestMain by using the suite.Register function.
//
// Step 2: Make sure suites' Teardown methods will be called by using the suite.Run function in TestMain.
//
// Step 3: Retrieve a suite within a test function by using the suite.Get function.
//
// To get into a little more detail: the first time suite.Get[T](tb) is called for a given suite type T,
// T's Setup method will be executed.
// If T's Setup method returns an error, tb.Fatal will be called.
// Any subsequent calls to suite.Get[T](tb) will immediately call tb.Fatal with the same error.
//
// The suite.Run function calls suite.Teardown once the tests have finished running.
// Any suite whose Setup method was executed will have their Teardown method executed.
// While suite.Teardown can technically be called at any time, it's recommended to use suite.Run instead
// of calling suite.Teardown manually. Teardown methods happen on a FILO basis from which they are registered;
// the suite that should be torn down last should be registered first.
//
//  package example
//
//  import (
//    "database/sql"
//    "os"
//    "testing"
//
//    _ "github.com/go-sql-driver/mysql"
//    "github.com/zpatrick/testx/assert"
//    "github.com/zpatrick/testx/suite"
//  )
//
//   type DBSuite struct {
//   	DB *sql.DB
//   }
//
//   func (d *DBSuite) Setup(tb testing.TB) error {
//   	db, err := sql.Open("mysql", "...")
//   	if err != nil {
//   	   return err
//   	}
//
//   	d.DB = db
//   	return nil
//   }
//
//   func (d *DBSuite) Teardown() error {
//   	if d.DB == nil {
//   	   return nil
//   	}
//
//   	return d.DB.Close()
//   }
//
//   func TestMain(m *testing.M) {
//   	suite.Register(&DBSuite{})
//
//   	os.Exit(suite.Run(m))
//   }
//
//   func TestDB_xxx(t *testing.T) {
//   	s := suite.Get[*DBSuite](t)
//
//   	row := s.DB.QueryRow("...")
//   	assert.NilError(t, row.Err())
//   }
//
//   func TestDB_yyy(t *testing.T) {
//   	s := suite.Get[*DBSuite](t)
//
//   	row := s.DB.QueryRow("...")
//   	assert.NilError(t, row.Err())
//   }
package suite
