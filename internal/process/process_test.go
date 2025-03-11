package process

import (
	"testing"
)

func asssertmessage(t testing.TB, got int, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d but want %d", got, want)
	}
}

func TestParameter(t *testing.T) {

	t.Run("Parameter of 2 and 1", func(t *testing.T) {
		got := getPerimeter(2, 1)
		want := 6
		asssertmessage(t, got, want)
	})

	t.Run("Parameter of 4 and 5", func(t *testing.T) {
		got := getPerimeter(4, 5)
		want := 18
		asssertmessage(t, got, want)
	})

	t.Run("Parameter of 10 and 11", func(t *testing.T) {
		got := getPerimeter(10, 11)
		want := 42
		asssertmessage(t, got, want)
	})

}
