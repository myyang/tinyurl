package tinyurl

import "fmt"
import "testing"

func TestInit(t *testing.T) {
	_ = NewTinyURL()
}

func TestTinyURL(t *testing.T) {
	tu := NewTinyURL()

	tests := []struct {
		URL string
		Len int
	}{
		{"www.google.com", 10},
		{"www.yahoo.com", 8},
		{"www.facebook.com", 7},
		{"www.github.com", 6},
		{"www.linkedin.com", 5},
	}

	for _, exp := range tests {
		tu.SetHashLen(exp.Len)
		h, err := tu.Shorten(exp.URL)
		if err != nil || len(h) != exp.Len {
			t.Errorf("Error to shorten: %v(exp len: %v, got: %v), err: %v\n", exp.URL, exp.Len, len(h), err)
			t.FailNow()
		}
		if url, err := tu.Recover(h); err != nil || url != exp.URL {
			t.Errorf("Error to recover: %v, exp: %v, got: %v. Or err: %v\n", h, exp.URL, url, err)
			t.FailNow()
		}
	}
}

func BenchmarkShorten(b *testing.B) {
	tu := NewTinyURL()
	tmp := "https://www.google.com/%v\n"
	for i := 0; i < b.N; i++ {
		tu.Shorten(fmt.Sprintf(tmp, i))
	}
}
