// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Prompter is an autogenerated mock type for the Prompter type
type Prompter struct {
	mock.Mock
}

// Choose provides a mock function with given fields: _a0, _a1
func (_m *Prompter) Choose(_a0 string, _a1 []string) int {
	ret := _m.Called(_a0, _a1)

	var r0 int
	if rf, ok := ret.Get(0).(func(string, []string) int); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// ChooseWithDefault provides a mock function with given fields: _a0, _a1, _a2
func (_m *Prompter) ChooseWithDefault(_a0 string, _a1 string, _a2 []string) (string, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string, []string) string); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, []string) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Password provides a mock function with given fields: _a0
func (_m *Prompter) Password(_a0 string) string {
	ret := _m.Called(_a0)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// RequestSecurityCode provides a mock function with given fields: _a0
func (_m *Prompter) RequestSecurityCode(_a0 string, b bool) string {
	ret := _m.Called(_a0)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// String provides a mock function with given fields: _a0, _a1
func (_m *Prompter) String(_a0 string, _a1 string) string {
	ret := _m.Called(_a0, _a1)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// StringRequired provides a mock function with given fields: _a0
func (_m *Prompter) StringRequired(_a0 string) string {
	ret := _m.Called(_a0)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
