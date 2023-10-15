package retry

import (
	"reflect"
	"testing"
)

func TestExponentialBackoff(t *testing.T) {
	backOffs := []uint64{}
	backOff := uint64(1)
	backOffs = append(backOffs, backOff) // backOffs = [1]
	backOffFn := ExponentialBackoff(10)
	// backOffs = [1]
	backOff = backOffFn(backOff)
	backOffs = append(backOffs, backOff) // backOffs = [1, 2]

	backOff = backOffFn(backOff)
	backOffs = append(backOffs, backOff) // backOffs = [1, 2, 4]

	backOff = backOffFn(backOff)
	backOffs = append(backOffs, backOff) // backOffs = [1, 2, 4, 8]

	backOff = backOffFn(backOff)
	backOffs = append(backOffs, backOff) // backOffs = [1, 2, 4, 8, 10] Because maxBackOff expected is 10
	expectedBackOffs := []uint64{1, 2, 4, 8, 10}
	if !reflect.DeepEqual(backOffs, expectedBackOffs) {
		t.Fatalf("Expected BackOff should be %+v but found %+v", expectedBackOffs, backOffs)
	}
}

func TestRandomBackOff(t *testing.T) {
	lower := 2
	upper := 7
	initialBackOff := uint64(1)
	backoffFunc := RandomBackoff(2, 7)
	firstBackOff := backoffFunc(initialBackOff)
	secondBackOff := backoffFunc(firstBackOff)
	thirdBackOff := backoffFunc(firstBackOff)

	backOffs := []uint64{firstBackOff, secondBackOff, thirdBackOff}

	for _, num := range backOffs {
		if int(num) < lower || int(num) > upper {
			t.Fatalf("Random number %d should be between %d and %d", num, lower, upper)
		}
	}
}
