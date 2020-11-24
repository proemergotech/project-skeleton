package main

import (
	"io/ioutil"
	"testing"
)

var testText = `
	"github.com/proemergotech/project-skeleton/app/client" DELETE
	"github.com/proemergotech/project-skeleton/app/config"
	"github.com/proemergotech/project-skeleton/app/event" DELETE
	"github.com/proemergotech/project-skeleton/app/rest"
	"github.com/proemergotech/project-skeleton/app/service"
	"github.com/proemergotech/project-skeleton/app/storage" DELETE
	"github.com/proemergotech/project-skeleton/app/validation"
`

func TestRegex(t *testing.T) {
	result := importCleanupRegex.ReplaceAllString(testText, "")
	t.Log(result)
}

func TestGoimports(t *testing.T) {
	f, err := ioutil.ReadFile("../output/app/di/container.go")
	checkErr(t, err)
	f2, err := goimport(f, "project-skeleton")
	checkErr(t, err)
	t.Log(string(f2))
}

func checkErr(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("%+v", err)
	}
}
