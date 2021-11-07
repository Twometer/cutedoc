package utils

import "testing"

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

func TestGetFileName(t *testing.T) {
	assertEqual(t, "test_file", GetFileName("../a/b/c/test_file.txt"))
	assertEqual(t, "test_file", GetFileName("../a/b/test_file"))
	assertEqual(t, "test_file", GetFileName("test_file.txt"))
	assertEqual(t, "test_file", GetFileName("test_file"))
	assertEqual(t, "test_file", GetFileName("test_file.txt/"))
	assertEqual(t, "test_file", GetFileName("../:/:.:/:/::/..../test_file.txt/"))
	assertEqual(t, "test_file", GetFileName("\\a\\b\\test_file.txt"))
}

func TestStripParentDirectories(t *testing.T) {
	assertEqual(t, "test/a/b.txt", StripParentDirectories("../../test/./a/b.txt"))
	assertEqual(t, "test/a/b.txt", StripParentDirectories("..\\..\\test\\.\\a\\b.txt"))
	assertEqual(t, "b.txt", StripParentDirectories("./b.txt"))
	assertEqual(t, "b.txt", StripParentDirectories("b.txt"))
}
