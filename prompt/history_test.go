package prompt

import "testing"

func TestSearchUp(t *testing.T) {
	history := &History{
		history: []string{
			"ab",
			"b",
			"abc",
		},
	}

	result, hit := history.SearchUp("a")
	if !hit {
		t.Error("should hit")
	}
	expected := "abc"
	if result != expected {
		t.Fatalf("expected: %v; got: %v\n", expected, result)
	}

	result, hit = history.SearchUp("abc")
	if !hit {
		t.Error("should hit")
	}
	expected = "ab"
	if result != expected {
		t.Fatalf("expected: %v; got: %v\n", expected, result)
	}

	result, hit = history.SearchUp("ab")
	if hit {
		t.Error("should not hit")
	}
	expected = "a"
	if result != expected {
		t.Fatalf("expected: %v; got: %v\n", expected, result)
	}

	result, hit = history.SearchUp("a")
	if hit {
		t.Error("should not hit")
	}
	expected = "a"
	if result != expected {
		t.Fatalf("expected: %v; got: %v\n", expected, result)
	}

}

func TestSearchDown(t *testing.T) {
	history := &History{
		history: []string{
			"ab",
			"b",
			"abc",
		},
	}

	result, hit := history.SearchUp("a")
	if !hit {
		t.Error("should hit")
	}
	expected := "abc"
	if result != expected {
		t.Fatalf("expected: %v; got: %v\n", expected, result)
	}

	result, hit = history.SearchUp("abc")
	if !hit {
		t.Error("should hit")
	}
	expected = "ab"
	if result != expected {
		t.Fatalf("expected: %v; got: %v\n", expected, result)
	}

	result, hit = history.SearchDown("ab")
	if !hit {
		t.Error("should hit")
	}
	expected = "abc"
	if result != expected {
		t.Fatalf("expected: %v; got: %v\n", expected, result)
	}

	result, hit = history.SearchDown("abc")
	if hit {
		t.Error("should not hit")
	}
	expected = "a"
	if result != expected {
		t.Fatalf("expected: %v; got: %v\n", expected, result)
	}
}
