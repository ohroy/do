package do

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProvide(t *testing.T) {
	is := assert.New(t)

	type test struct{}

	i := New()

	Provide(i, func(i Injector) (*test, error) {
		return &test{}, nil
	})

	Provide(i, func(i Injector) (test, error) {
		return test{}, fmt.Errorf("error")
	})

	is.Panics(func() {
		// try to erase previous instance
		Provide(i, func(i Injector) (test, error) {
			return test{}, fmt.Errorf("error")
		})
	})

	is.Len(i.self.services, 2)

	s1, ok1 := i.self.services["*github.com/samber/do.test"]
	is.True(ok1)
	if ok1 {
		s, ok := s1.(Service[*test])
		is.True(ok)
		if ok {
			is.Equal("*github.com/samber/do.test", s.getName())
		}
	}

	s2, ok2 := i.self.services["github.com/samber/do.test"]
	is.True(ok2)
	if ok2 {
		s, ok := s2.(Service[test])
		is.True(ok)
		if ok {
			is.Equal("github.com/samber/do.test", s.getName())
		}
	}

	_, ok3 := i.self.services["github.com/samber/do.*plop"]
	is.False(ok3)

	// @TODO: check that all services share the same references
}

func TestProvideNamed(t *testing.T) {
	is := assert.New(t)

	type test struct{}

	i := New()

	ProvideNamed(i, "*foobar", func(i Injector) (*test, error) {
		return &test{}, nil
	})

	ProvideNamed(i, "foobar", func(i Injector) (test, error) {
		return test{}, fmt.Errorf("error")
	})

	is.Panics(func() {
		// try to erase previous instance
		ProvideNamed(i, "foobar", func(i Injector) (test, error) {
			return test{}, fmt.Errorf("error")
		})
	})

	is.Len(i.self.services, 2)

	s1, ok1 := i.self.services["*foobar"]
	is.True(ok1)
	if ok1 {
		s, ok := s1.(Service[*test])
		is.True(ok)
		if ok {
			is.Equal("*foobar", s.getName())
		}
	}

	s2, ok2 := i.self.services["foobar"]
	is.True(ok2)
	if ok2 {
		s, ok := s2.(Service[test])
		is.True(ok)
		if ok {
			is.Equal("foobar", s.getName())
		}
	}

	_, ok3 := i.self.services["*do.plop"]
	is.False(ok3)

	// @TODO: check that all services share the same references
}

func TestProvideValue(t *testing.T) {
	is := assert.New(t)

	i := New()

	type test struct {
		foobar string
	}
	_test := test{foobar: "foobar"}

	ProvideValue(i, 42)
	ProvideValue(i, _test)

	is.Len(i.self.services, 2)

	s1, ok1 := i.self.services["int"]
	is.True(ok1)
	if ok1 {
		s, ok := s1.(Service[int])
		is.True(ok)
		if ok {
			is.Equal("int", s.getName())
			instance, err := s.getInstance(i)
			is.EqualValues(42, instance)
			is.Nil(err)
		}
	}

	s2, ok2 := i.self.services["github.com/samber/do.test"]
	is.True(ok2)
	if ok2 {
		s, ok := s2.(Service[test])
		is.True(ok)
		if ok {
			is.Equal("github.com/samber/do.test", s.getName())
			instance, err := s.getInstance(i)
			is.EqualValues(_test, instance)
			is.Nil(err)
		}
	}

	// @TODO: check that all services share the same references
}

func TestProvideNamedValue(t *testing.T) {
	is := assert.New(t)

	i := New()

	type test struct {
		foobar string
	}
	_test := test{foobar: "foobar"}

	ProvideNamedValue(i, "foobar", 42)
	ProvideNamedValue(i, "hello", _test)

	is.Len(i.self.services, 2)

	s1, ok1 := i.self.services["foobar"]
	is.True(ok1)
	if ok1 {
		s, ok := s1.(Service[int])
		is.True(ok)
		if ok {
			is.Equal("foobar", s.getName())
			instance, err := s.getInstance(i)
			is.EqualValues(42, instance)
			is.Nil(err)
		}
	}

	s2, ok2 := i.self.services["hello"]
	is.True(ok2)
	if ok2 {
		s, ok := s2.(Service[test])
		is.True(ok)
		if ok {
			is.Equal("hello", s.getName())
			instance, err := s.getInstance(i)
			is.EqualValues(_test, instance)
			is.Nil(err)
		}
	}

	// @TODO: check that all services share the same references
}

func TestProvideTransiant(t *testing.T) {
	is := assert.New(t)

	type test struct{}

	i := New()

	ProvideTransiant(i, func(i Injector) (*test, error) {
		return &test{}, nil
	})

	ProvideTransiant(i, func(i Injector) (test, error) {
		return test{}, fmt.Errorf("error")
	})

	is.Panics(func() {
		// try to erase previous instance
		Provide(i, func(i Injector) (test, error) {
			return test{}, fmt.Errorf("error")
		})
	})

	is.Len(i.self.services, 2)

	s1, ok1 := i.self.services["*github.com/samber/do.test"]
	is.True(ok1)
	if ok1 {
		s, ok := s1.(Service[*test])
		is.True(ok)
		if ok {
			is.Equal("*github.com/samber/do.test", s.getName())
		}
	}

	s2, ok2 := i.self.services["github.com/samber/do.test"]
	is.True(ok2)
	if ok2 {
		s, ok := s2.(Service[test])
		is.True(ok)
		if ok {
			is.Equal("github.com/samber/do.test", s.getName())
		}
	}

	_, ok3 := i.self.services["github.com/samber/do.*plop"]
	is.False(ok3)

	// @TODO: check that all services share the same references
}
func TestProvideNamedTransiant(t *testing.T) {
	is := assert.New(t)

	type test struct{}

	i := New()

	ProvideNamed(i, "*foobar", func(i Injector) (*test, error) {
		return &test{}, nil
	})

	ProvideNamed(i, "foobar", func(i Injector) (test, error) {
		return test{}, fmt.Errorf("error")
	})

	is.Panics(func() {
		// try to erase previous instance
		ProvideNamed(i, "foobar", func(i Injector) (test, error) {
			return test{}, fmt.Errorf("error")
		})
	})

	is.Len(i.self.services, 2)

	s1, ok1 := i.self.services["*foobar"]
	is.True(ok1)
	if ok1 {
		s, ok := s1.(Service[*test])
		is.True(ok)
		if ok {
			is.Equal("*foobar", s.getName())
		}
	}

	s2, ok2 := i.self.services["foobar"]
	is.True(ok2)
	if ok2 {
		s, ok := s2.(Service[test])
		is.True(ok)
		if ok {
			is.Equal("foobar", s.getName())
		}
	}

	_, ok3 := i.self.services["*do.plop"]
	is.False(ok3)

	// @TODO: check that all services share the same references
}

func TestOverride(t *testing.T) {
	is := assert.New(t)

	type test struct {
		foobar int
	}

	i := New()

	is.NotPanics(func() {
		Provide(i, func(i Injector) (*test, error) {
			return &test{42}, nil
		})
		is.Equal(42, MustInvoke[*test](i).foobar)

		Override(i, func(i Injector) (*test, error) {
			return &test{1}, nil
		})
		is.Equal(1, MustInvoke[*test](i).foobar)

		// OverrideNamed(i, "*github.com/samber/do.test", func(i Injector) (*test, error) {
		// 	return &test{2}, nil
		// })
		// is.Equal(2, MustInvoke[*test](i).foobar)

		// OverrideValue(i, &test{3})
		// is.Equal(3, MustInvoke[*test](i).foobar)

		// OverrideNamedValue(i, "*github.com/samber/do.test", &test{4})
		// is.Equal(4, MustInvoke[*test](i).foobar)
	})
}

func TestOverrideNamed(t *testing.T) {
	is := assert.New(t)

	type test struct {
		foobar int
	}

	i := New()

	Provide(i, func(i Injector) (*test, error) {
		return &test{42}, nil
	})
	is.Equal(42, MustInvoke[*test](i).foobar)

	OverrideNamed(i, "*github.com/samber/do.test", func(i Injector) (*test, error) {
		return &test{2}, nil
	})
	is.Equal(2, MustInvoke[*test](i).foobar)
}

func TestOverrideValue(t *testing.T) {
	is := assert.New(t)

	type test struct {
		foobar int
	}

	i := New()

	Provide(i, func(i Injector) (*test, error) {
		return &test{42}, nil
	})
	is.Equal(42, MustInvoke[*test](i).foobar)

	OverrideNamed(i, "*github.com/samber/do.test", func(i Injector) (*test, error) {
		return &test{2}, nil
	})
	is.Equal(2, MustInvoke[*test](i).foobar)
}

func TestOverrideNamedValue(t *testing.T) {
	is := assert.New(t)

	type test struct {
		foobar int
	}

	i := New()

	Provide(i, func(i Injector) (*test, error) {
		return &test{42}, nil
	})
	is.Equal(42, MustInvoke[*test](i).foobar)

	OverrideNamedValue(i, "*github.com/samber/do.test", &test{4})
	is.Equal(4, MustInvoke[*test](i).foobar)
}

func TestOverrideTransiant(t *testing.T) {
	// @TODO
}

func TestOverrideNamedTransiant(t *testing.T) {
	// @TODO
}

func TestInvoke(t *testing.T) {
	is := assert.New(t)

	type test struct {
		foobar string
	}

	i := New()

	Provide(i, func(i Injector) (test, error) {
		return test{foobar: "foobar"}, nil
	})

	is.Len(i.self.services, 1)

	s0a, ok0a := i.self.services["github.com/samber/do.test"]
	is.True(ok0a)

	s0b, ok0b := s0a.(*ServiceLazy[test])
	is.True(ok0b)
	is.False(s0b.built)

	s1, err1 := Invoke[test](i)
	is.Nil(err1)
	if err1 == nil {
		is.Equal("foobar", s1.foobar)
	}

	is.True(s0b.built)

	_, err2 := Invoke[*test](i)
	is.NotNil(err2)
	is.Errorf(err2, "do: service not found")
}

func TestMustInvoke(t *testing.T) {
	is := assert.New(t)

	i := New()

	type test struct {
		foobar string
	}
	_test := test{foobar: "foobar"}

	Provide(i, func(i Injector) (test, error) {
		return _test, nil
	})

	is.Len(i.self.services, 1)

	is.Panics(func() {
		_ = MustInvoke[string](i)
	})

	is.NotPanics(func() {
		instance1 := MustInvoke[test](i)
		is.EqualValues(_test, instance1)
	})
}

func TestInvokeNamed(t *testing.T) {
	is := assert.New(t)

	i := New()

	type test struct {
		foobar string
	}
	_test := test{foobar: "foobar"}

	ProvideNamedValue(i, "foobar", 42)
	ProvideNamedValue(i, "hello", _test)

	is.Len(i.self.services, 2)

	service0, err0 := InvokeNamed[string](i, "plop")
	is.NotNil(err0)
	is.Empty(service0)

	instance1, err1 := InvokeNamed[test](i, "hello")
	is.Nil(err1)
	is.EqualValues(_test, instance1)
	is.EqualValues("foobar", instance1.foobar)

	instance2, err2 := InvokeNamed[int](i, "foobar")
	is.Nil(err2)
	is.EqualValues(42, instance2)
}

func TestMustInvokeNamed(t *testing.T) {
	is := assert.New(t)

	i := New()

	ProvideNamedValue(i, "foobar", 42)

	is.Len(i.self.services, 1)

	is.Panics(func() {
		_ = MustInvokeNamed[string](i, "hello")
	})

	is.Panics(func() {
		_ = MustInvokeNamed[string](i, "foobar")
	})

	is.NotPanics(func() {
		instance1 := MustInvokeNamed[int](i, "foobar")
		is.EqualValues(42, instance1)
	})
}
