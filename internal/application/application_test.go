package application_test

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"strings"
	"telegram_bot/internal/application"
	"testing"
)

func TestApplication_InfoPrintF(t *testing.T) {
	buf := new(bytes.Buffer)
	testWriter := bufio.NewWriter(buf)
	testPrefix := "INFO\t"
	logger := log.New(testWriter, testPrefix, log.Ldate|log.Ltime)

	formatStr := "This string must be written inside the buffer: %s"
	formatArg := "info"
	app := &application.Application{InfoLog: logger}
	app.InfoPrintF(formatStr, formatArg)
	testWriter.Flush()

	got := buf.String()
	if !strings.Contains(got, testPrefix) {
		t.Error("Application_InfoPrintF() does not have a prefix")
	}
	want := fmt.Sprintf(formatStr, formatArg)
	if !strings.Contains(got, want) {
		t.Errorf("Application_InfoprintF() performed incorrectly! expectd %q, got %q", want, got)
	}
}

func TestApplication_ErrorPrintF(t *testing.T) {
	buf := new(bytes.Buffer)
	testWriter := bufio.NewWriter(buf)
	testPrefix := "ERROR\t"
	logger := log.New(testWriter, testPrefix, log.Ldate|log.Ltime)

	formatStr := "This string must be written inside the buffer: %s"
	formatArg := "error"
	app := &application.Application{ErrorLog: logger}
	app.ErrorPrintF(formatStr, formatArg)
	testWriter.Flush()

	got := buf.String()
	if !strings.Contains(got, testPrefix) {
		t.Error("Application_ErrorPrintF() does not have a prefix")
	}

	want := fmt.Sprintf(formatStr, formatArg)
	if !strings.Contains(got, want) {
		t.Errorf("Application_ErrorPrintF() performed incorrectly! expectd %q, got %q", want, got)
	}
}
