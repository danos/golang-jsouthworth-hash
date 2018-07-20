package hash

import (
	"testing"
)

func TestHashString(t *testing.T) {
	hashkeyInit(1)
	t.Log(hash("sss", 0))
}

func TestHashInt(t *testing.T) {
	hashkeyInit(1)
	t.Log(hash(int(111), 0))
}

func TestHashUint(t *testing.T) {
	hashkeyInit(1)
	t.Log(hash(uint(111), 0))
}

type something struct {
	a, b, c, d, e, f, g, h int
}

func TestHashStruct(t *testing.T) {
	hashkeyInit(1)
	t.Log(hash(something{a: 1}, 0))
}

func TestHashStructPtr(t *testing.T) {
	hashkeyInit(1)
	t.Log(hash(&something{a: 1}, 0))
}

func TestHashStructPtr2(t *testing.T) {
	hashkeyInit(1)
	t.Log(hash(&something{a: 2}, 0))
}

func TestHashSlice(t *testing.T) {
	hashkeyInit(1)
	arr := [5]int{1, 0, 0, 0, 0}
	s1 := arr[:]
	s2 := arr[1:]
	s3 := arr[1:]
	h1 := hash(s1, 0)
	h2 := hash(s2, 0)
	h3 := hash(s3, 0)
	if h2 != h3 {
		t.Fatalf("%q(%d) != %q(%d)\n", s2, h2, s3, h3)
	}
	if h1 == h2 || h1 == h3 {
		t.Fatalf("unexpected equivalent hashes %q(%d) == %q(%d)\n",
			s1, h1, s2, h2)
	}
}

func TestHashStructWithSlice(t *testing.T) {
	hashkeyInit(1)
	t.Log(hash(struct {
		a int
		b []int
	}{a: 1, b: []int{1, 0, 0, 0, 0}}, 0))
}
