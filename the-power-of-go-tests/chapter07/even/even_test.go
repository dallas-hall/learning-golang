package even_test

import (
	"even"
	"strconv"
	"testing"
)

func TestIsEvenReturningTrueForEvenNumbers(t *testing.T) {
	t.Parallel()
	for i := 0; i < 100; i += 2 {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if !even.IsEven(i) {
				t.Error(false)
			}
		})
	}
}

func TestIsEvenReturningFalseForOddNumbers(t *testing.T) {
	t.Parallel()
	for i := 1; i < 100; i += 2 {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if even.IsEven(i) {
				t.Error(true)
			}
		})
	}
}

func TestIsEvenCacheReturningTrueForEvenNumbers(t *testing.T) {
	t.Parallel()
	for i := 0; i < 100; i += 2 {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if !even.IsEvenCache(i) {
				t.Error(false)
			}
		})
	}
}

func TestIsEvenCacheReturningFalseForOddNumbers(t *testing.T) {
	t.Parallel()
	for i := 1; i < 100; i += 2 {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if even.IsEvenCache(i) {
				t.Error(true)
			}
		})
	}
}
