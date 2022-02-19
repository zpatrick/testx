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

// A Suite holds shared dependencies for multiple tests.
type Suite interface {
	// Setup should create any dependencies required for the suite to run.
	// The tb parameter is passed only as a means to setup dependent suites; it
	// should not be used as a control-flow mechanism (e.g. by calling tb.Fatal).
	// Any error(s) encountered during this method should be returned.
	Setup(tb testing.TB) error

	// Teardown should cleanup any resources created by Setup.
	// If Setup returns an error, this method will still be executed in order
	// to cleanup partially-created dependencies. This requires that Teardown
	// methods be idempotent.
	Teardown() error
}

// Base is a placeholder type which can be embedded into types
// which don't need to implement the Setup or Teardown methods.
type Base struct{}

// Setup is a no-op.
func (Base) Setup(tb testing.TB) error { return nil }

// Teardown is a no-op.
func (Base) Teardown() error { return nil }

// Register allows a suite of s's concrete type to be later retrieved using Get.
// Only one instance of type s's concrete type should be registered.
// The order in which suites are registered determines the order teardown methods are called:
// Teardowns happen on a FILO (first in, last out) basis.
func Register(s Suite) {
	m := newSuiteManager(s)

	if err := defaultRegistry.Insert(m.Type(), m); err != nil {
		panic(err)
	}
}

// Get returns the instance of S which must have been previously registered using Register.
// If this is the first time Get is called for type S, the suite's Setup method will be ran.
// If the suite's Setup method fails, tb.Fatal will be called.
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
// If a suite was registered but never retrieved (by using the Get function), its
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
// returning the exit code from m.Run or 1 if an error occured during Teardown.
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

	mux      sync.Mutex
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
	s.mux.Lock()
	defer s.mux.Unlock()

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
	// Execute teardowns on a FILO basis.
	r.teardownOrder = append([]string{key}, r.teardownOrder...)
	return nil
}

func (r *registry) Get(key string) (*suiteManager, bool) {
	s, ok := r.suites[key]
	return s, ok
}
