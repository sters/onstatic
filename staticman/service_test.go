package staticman

import "testing"

func Test_generateNewUniqueName(t *testing.T) {
	m := map[string]struct{}{}
	for i := 1; i < 20; i++ {
		key := generateNewUniqueName()
		if _, ok := m[key]; ok {
			t.Fatal("it's not unique name")
		}
	}
}
