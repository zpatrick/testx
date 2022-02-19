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
	Setup(t testing.TB) error
	Teardown() error
}

// Base is a placeholder type which can be embedded into other suites
// if they don't need to implement setup or teardown methods.
type Base struct{}

// Setup is a no-op.
func (Base) Setup(tb testing.TB) error { return nil }

// Teardown is a no-op.
func (Base) Teardown() error { return nil }

// Register allows a suite of type s to be later retrieved using Get.
// Only one instance of type s should be registered - multiple calls to Register
// using the same type will cause a panic.
func Register(s Suite) {
	m := newSuiteManager(s)

	if err := defaultRegistry.Insert(m.Type(), m); err != nil {
		panic(err)
	}
}

// Get returns the instance of S which must have been previously registered using Register.
// If this is the first time Get is called for type S, the suite's Setup method will be ran.
// If the suite's Setup method failed, t.Fatal will be called.
func Get[S Suite](tb testing.TB) (s S) {
	sType := newSuiteManager(s).Type()
	m, ok := defaultRegistry.Get(sType)
	if !ok {
		tb.Fatalf("suite of type %v has not been registered", sType)
	}

	if err := m.Setup(tb); err != nil {
		tb.Fatalf("setup failed for suite %v: %s", sType, err.Error())
	}

	return m.suite.(S)
}

// Teardown runs the Teardown method on registered suites who ran their Setup methods.
// If a suite was registered but never retrieved (by using the Get function), it's
// teardown method will not be run.
func Teardown() error {
	var errs []error
	for _, key := range defaultRegistry.teardownOrder {
		m, ok := defaultRegistry.Get(key)
		if !ok {
			// This should never happen in theory - just making life easier in case is a bug is introduced.
			return fmt.Errorf("suite %s specified by teardownOrder missing from defaultRegistry", key)
		}

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

func (s *suiteManager) Setup(tb testing.TB) error {
	if s.setupErr != nil {
		return s.setupErr
	}

	s.once.Do(func() {
		s.setupErr = s.suite.Setup(tb)
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

var defaultRegistry = &registry{
	suites:        map[string]*suiteManager{},
	teardownOrder: []string{},
}

type registry struct {
	suites        map[string]*suiteManager
	teardownOrder []string
}

func (r *registry) Insert(key string, s *suiteManager) error {
	if _, ok := r.suites[key]; ok {
		return fmt.Errorf("suite %s has already been registered", key)
	}

	r.suites[key] = s

	// Order happens in reverse order of registration
	r.teardownOrder = append([]string{key}, r.teardownOrder...)
	return nil
}

func (r *registry) Get(key string) (*suiteManager, bool) {
	s, ok := r.suites[key]
	return s, ok
}
