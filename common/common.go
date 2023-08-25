package common

import (
	"path/filepath"
	"runtime"
)

func FindProjectPath() string {

	var (
		_, b, _, _  = runtime.Caller(0)
		projectbase = filepath.Join(filepath.Dir(b), "..")
	)

	// m1 := regexp.MustCompile(`\\`)
	// projectbaseclean := m1.ReplaceAllString(projectbase, "/")

	return projectbase
}

func GetLine() int {
	_, _, line, _ := runtime.Caller(1)
	return line
}
