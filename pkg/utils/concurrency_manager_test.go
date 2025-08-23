// go
package utils

import (
	"errors"
	"strings"
	"testing"
)

func Test_MergeErrors_various(t *testing.T) {
	t.Parallel()

	e1 := errors.New("err1")
	e2 := errors.New("err2")

	t.Run("nil input returns nil", func(t *testing.T) {
		t.Parallel()
		if got := MergeErrors(nil); got != nil {
			t.Fatalf("expected nil, got: %v", got)
		}
	})

	t.Run("[]error with mixed nils returns only non-nil in order", func(t *testing.T) {
		t.Parallel()
		in := []error{nil, e1, nil, e2}
		got := MergeErrors(in)
		if got == nil || len(got) != 2 {
			t.Fatalf("expected 2 errors, got: %v", got)
		}
		if got[0].Error() != "err1" || got[1].Error() != "err2" {
			t.Fatalf("unexpected errors order/content: %v", got)
		}
	})

	t.Run("*[]error nil pointer returns nil", func(t *testing.T) {
		t.Parallel()
		var p *[]error = nil
		if got := MergeErrors(p); got != nil {
			t.Fatalf("expected nil for nil *[]error, got: %v", got)
		}
	})

	t.Run("*[]error with only nils returns nil", func(t *testing.T) {
		t.Parallel()
		s := []error{nil, nil}
		if got := MergeErrors(&s); got != nil {
			t.Fatalf("expected nil for pointer to slice of only nils, got: %v", got)
		}
	})

	t.Run("*[]error with errors returns those errors", func(t *testing.T) {
		t.Parallel()
		s := []error{e1, nil, e2}
		got := MergeErrors(&s)
		if got == nil || len(got) != 2 {
			t.Fatalf("expected 2 errors, got: %v", got)
		}
		if got[0].Error() != "err1" || got[1].Error() != "err2" {
			t.Fatalf("unexpected errors: %v", got)
		}
	})

	t.Run("error typed nil and non-nil", func(t *testing.T) {
		t.Parallel()
		var enil error = nil
		if got := MergeErrors(enil); got != nil {
			t.Fatalf("expected nil for (error)(nil), got: %v", got)
		}
		if got := MergeErrors(e1); got == nil || len(got) != 1 || got[0].Error() != "err1" {
			t.Fatalf("expected single err1, got: %v", got)
		}
	})

	t.Run("*error pointer cases", func(t *testing.T) {
		t.Parallel()
		var enil error = nil
		// pointer to nil error value
		if got := MergeErrors(&enil); got != nil {
			t.Fatalf("expected nil for pointer to nil error value, got: %v", got)
		}
		// pointer to non-nil
		e := errors.New("ptrerr")
		got := MergeErrors(&e)
		if got == nil || len(got) != 1 || got[0].Error() != "ptrerr" {
			t.Fatalf("expected single ptrerr, got: %v", got)
		}
	})

	t.Run("chan error with mixed sends", func(t *testing.T) {
		// not parallel to avoid scheduler surprises with channels in some runners
		ch := make(chan error, 3)
		ch <- e1
		ch <- nil
		ch <- e2
		close(ch)

		got := MergeErrors(ch)
		if got == nil || len(got) != 2 {
			t.Fatalf("expected 2 errors from chan error, got: %v", got)
		}
		if got[0].Error() != "err1" || got[1].Error() != "err2" {
			t.Fatalf("unexpected errors from chan: %v", got)
		}
	})

	t.Run("chan []error with mixed slices", func(t *testing.T) {
		ch := make(chan []error, 3)
		ch <- []error{e1, nil}
		ch <- []error{}
		ch <- []error{nil, e2}
		close(ch)

		got := MergeErrors(ch)
		if got == nil || len(got) != 2 {
			t.Fatalf("expected 2 errors from chan []error, got: %v", got)
		}
		if got[0].Error() != "err1" || got[1].Error() != "err2" {
			t.Fatalf("unexpected errors from chan []error: %v", got)
		}
	})

	t.Run("unsupported type returns descriptive error", func(t *testing.T) {
		t.Parallel()
		got := MergeErrors(123)
		if got == nil || len(got) != 1 {
			t.Fatalf("expected single error for unsupported type, got: %v", got)
		}
		if !strings.Contains(got[0].Error(), "invalid input type in MergeErrors") {
			t.Fatalf("unexpected error message: %v", got[0].Error())
		}
	})
}