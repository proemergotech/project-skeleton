package main

import (
	"io/ioutil"
	"testing"
)

var testText = `
	"gitlab.com/proemergotech/dliver-project-skeleton/app/client" DELETE
	"gitlab.com/proemergotech/dliver-project-skeleton/app/config"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/event" DELETE
	"gitlab.com/proemergotech/dliver-project-skeleton/app/rest"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/service"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/storage" DELETE
	"gitlab.com/proemergotech/dliver-project-skeleton/app/validation"
`

func TestRegex(t *testing.T) {
	result := importCleanupRegex.ReplaceAllString(testText, "")
	t.Log(result)
}

func TestGoimports(t *testing.T) {
	f, err := ioutil.ReadFile("../output/app/di/container.go")
	checkErr(t, err)
	f2, err := goimport(f, "dliver-project-skeleton")
	checkErr(t, err)
	t.Log(string(f2))
}

func checkErr(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("%+v", err)
	}
}
