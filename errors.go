// Copyright 2022 Kami
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package errors

import (
	"reflect"

	"github.com/goentf/runpoint"
)

type errorChain struct {
	pc   runpoint.PCounter
	text string
	next error
}

func (e *errorChain) Error() string {
	return e.text
}

// New returns an error in the format of the given text.
// Each call to New returns a different error value, even if the text is the same.
//
// underlay: An underlying error can be provided via the optional parameter.
func New(text string, cause ...error) error {
	var e error
	if len(cause) == 1 {
		e = cause[0]
	}
	return &errorChain{
		text: text,
		next: e,
		pc:   runpoint.PC(1),
	}
}

// File returns the file path where the error occurred.
func File(err error) string {
	if err == nil {
		return ""
	}
	e, ok := err.(*errorChain)
	if !ok {
		return ""
	}
	return e.pc.File()
}

// File returns the line number where the error occurred.
func Line(err error) int {
	if err == nil {
		return 0
	}
	e, ok := err.(*errorChain)
	if !ok {
		return 0
	}
	return e.pc.Line()
}

// PC returns the PCounter of the error.
func PC(err error) (pc runpoint.PCounter) {
	if err == nil {
		return
	}
	e, ok := err.(*errorChain)
	if !ok {
		return
	}
	return e.pc
}

// Cause returns the underlay error of err, if err's
// type is *errorChain returning error. Otherwise, returns nil.
func Cause(err error) error {
	if err == nil {
		return nil
	}
	ec, ok := err.(*errorChain)
	if !ok {
		return nil
	}
	return ec.next
}

// ForCauses gets all errors on the error chain.
func ForCauses(err error, fun func(error)) {
	if err == nil {
		return
	}
	for {
		fun(err)
		// Find the next error
		if e, ok := err.(*errorChain); ok && e.next != nil {
			err = e.next
		} else {
			return
		}
	}
}

// OneCauseOf is used to determine whether
// the cause of the specified error is the target error.
func OneCauseOf(err error, target error) bool {
	if target == nil {
		return err == target
	}

	isComparable := reflect.TypeOf(target).Comparable()
	for {
		if isComparable && err == target {
			return true
		}
		if ec, ok := err.(*errorChain); ok && ec.next != nil {
			err = ec.next
		} else {
			return false
		}
	}
}
