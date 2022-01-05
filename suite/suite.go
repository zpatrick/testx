package suite

import (
	"fmt"
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

// A Suite holds some shared state for multiple tests.
// The state should be instantiated during the Setup method
// and (if neccessary) cleaned up during the Teardown method.
type Suite interface {
	Setup() error
	Teardown() error
}

// Base is a placeholder type which can be embedded into other suites
// if they don't need to implement setup or teardown methods.
type Base struct{}

// Setup is a no-op.
func (Base) Setup() error { return nil }

// Teardown is a no-op.
func (Base) Teardown() error { return nil }

// Register allows a suite of type s to be later retrieved using Get.
// Only one instance of type s should be registered - multiple calls to Register
// using the same type will cause a panic.
func Register(s Suite) {
	m := newSuiteManager(s)
	key := m.Type()

	if _, ok := registry[key]; ok {
		panic(fmt.Sprintf("suite of type %v has already been registered", key))
	}

	registry[key] = m
}

// Get returns the instance of S which must have been previously registered using Register.
// If this is the first time Get is called for type S, the suite's Setup method will be ran.
// If the suite's Setup method failed, t.Fatal will be called.
func Get[S Suite](t testing.TB) (s S) {
	sType := newSuiteManager(s).Type()
	m, ok := registry[sType]
	if !ok {
		t.Fatalf("suite of type %v has not been registered", sType)
	}

	if err := m.Setup(); err != nil {
		t.Fatalf("setup failed for suite %v: %s", sType, err.Error())
	}

	return m.suite.(S)
}

// Teardown runs the Teardown method on registered suites who ran their Setup methods.
// If a suite was registered but never retrieved (by using the Get function), it's
// teardown method will not be run.
func Teardown() error {
	var errs []error
	for _, m := range registry {
		if err := m.Teardown(); err != nil {
			errs = append(errs, errors.Wrapf(err, "--- ERROR: Teardown failed for suite %v", m.Type()))
		}
	}

	return multierr.Combine(errs...)
}

// Run is a helper method which calls m.Run and the Teardown function,
// returning the exit code from m.Run, or 1 if an error occured during Teardown.
func Run(m *testing.M) int {
	code := m.Run()
	if err := Teardown(); err != nil {
		for _, err := range multierr.Errors(err) {
			fmt.Fprintln(os.Stderr, err.Error())
		}

		return 1
	}

	return code
}

type suiteManager struct {
	suite Suite

	once     sync.Once
	setupRan bool
	setupErr error
}

func newSuiteManager(s Suite) *suiteManager {
	return &suiteManager{suite: s}
}

func (s *suiteManager) Type() string {
	return reflect.TypeOf(s.suite).String()
}

func (s *suiteManager) Setup() error {
	if s.setupErr != nil {
		return s.setupErr
	}

	s.once.Do(func() {
		s.setupErr = s.suite.Setup()
		s.setupRan = true
	})

	return s.setupErr
}

func (s *suiteManager) Teardown() error {
	if !s.setupRan {
		return nil
	}

	return s.suite.Teardown()
}

var registry map[string]*suiteManager = map[string]*suiteManager{}
