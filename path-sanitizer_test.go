package main

import (
	"reflect"
	"testing"
)

func TestExtendPath(t *testing.T) {
	for _, test := range []struct {
		path    string
		useDot  bool
		parts   []string
		prepend bool
		want    string
	}{
		{
			// append /bin and /sbin
			path:    "/where/ever:/where/ever/else",
			useDot:  false,
			parts:   []string{"/"},
			prepend: false,
			want:    "/where/ever:/where/ever/else:/bin:/sbin",
		},
		{
			// append . and /bin and /sbin
			path:    "/where/ever:/where/ever/else",
			useDot:  true,
			parts:   []string{"/"},
			prepend: false,
			want:    "/where/ever:/where/ever/else:.:/bin:/sbin",
		},
		{
			// prepend /bin and /sbin
			path:    "/where/ever:/where/ever/else",
			useDot:  false,
			parts:   []string{"/"},
			prepend: true,
			want:    "/bin:/sbin:/where/ever:/where/ever/else",
		},
		{
			// prepend . and /bin and /sbin
			path:    "/where/ever:/where/ever/else",
			useDot:  true,
			parts:   []string{"/"},
			prepend: true,
			want:    ".:/bin:/sbin:/where/ever:/where/ever/else",
		},
		{
			// sanitize and append /bin and /sbin
			path:    ":::///where///ever:::///where///ever///else:::",
			useDot:  false,
			parts:   []string{"/"},
			prepend: false,
			want:    "/where/ever:/where/ever/else:/bin:/sbin",
		},
	} {
		if got := extendPath(test.path, test.useDot, test.parts, test.prepend); got != test.want {
			t.Errorf("extendPath(%q,%v,%v,%v) = %q, want %q", test.path, test.useDot, test.parts, test.prepend, got, test.want)
		}
	}
}

func TestSplitPath(t *testing.T) {
	for _, test := range []struct {
		path string
		want []string
	}{
		{
			// normal
			path: "/bin:/usr/bin",
			want: []string{"/bin", "/usr/bin"},
		},
		{
			// empty parts get skipped
			path: ":::/bin:::/usr/bin:::",
			want: []string{"/bin", "/usr/bin"},
		},
		{
			// deduplication
			path: "/bin:/usr/bin:/bin:/usr/bin:/bin:/usr/bin",
			want: []string{"/bin", "/usr/bin"},
		},
	} {
		if got := splitPath(test.path); !reflect.DeepEqual(got, test.want) {
			t.Errorf("splitPath(%q) = %v, want %v", test.path, got, test.want)
		}
	}
}

func TestEvalString(t *testing.T) {
	for _, test := range []struct {
		parts []string
		shell string
		want  string
	}{
		{
			parts: []string{"a", "b", "c"},
			shell: "bash",
			want:  `export PATH="a:b:c"`,
		},
		{
			parts: []string{"a", "b", "c"},
			shell: "zsh",
			want:  `export PATH="a:b:c"`,
		},
		{
			parts: []string{"a", "b", "c"},
			shell: "fish",
			want:  `set -gx PATH "a:b:c"`,
		},
	} {
		if got := evalString(test.shell, test.parts); got != test.want {
			t.Errorf(`evalString(%q, %q) = %q, want %q`, test.shell, test.parts, got, test.want)
		}
	}
}

func TestIsDir(t *testing.T) {
	for _, test := range []struct {
		name string
		want bool
	}{
		{
			// existing dir
			name: "/etc",
			want: true,
		},
		{
			// non-existing dir
			name: "/a/b/c/d/this/does/not/exist",
			want: false,
		},
		{
			// existing, but a file, not a dir
			name: "path-sanitizer.go",
			want: false,
		},
	} {
		if got := isDir(test.name); got != test.want {
			t.Errorf("isdir(%q) = %v, want %v", test.name, got, test.want)
		}
	}
}
