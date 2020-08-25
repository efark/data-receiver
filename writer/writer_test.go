/*
Example testing package for writer interface and MemoryWriter.
*/
package writer_test

import (
	"fmt"
	"github.com/efark/data-receiver/writer"
	"os"
	"testing"
)

func TestMemoryWriter_Write(t *testing.T) {
	var w writer.Writer
	memoryWriter, err := writer.NewMemoryWriter()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	w = memoryWriter
	err = w.Write("Test message.")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ms := memoryWriter.GetMessages()
	if len(ms) == 0 {
		t.Log(fmt.Sprintf("Expected len(ms): %d, received: %d", 1, len(ms)))
		t.FailNow()
	}

	if ms[0] != "Test message." {
		t.Log(fmt.Sprintf("Expected ms[0]: %q, received: %q", "Test message.", ms[0]))
		t.FailNow()
	}
}

func TestFileWriter_Write(t *testing.T) {
	filepath := "./test.txt"
	if fileExists(filepath) {
		os.Remove(filepath)
	}

	var w writer.Writer
	fileWriter, err := writer.NewFileWriter(filepath)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	w = fileWriter
	err = w.Write("Test message 1.")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	err = w.Write("Test message 2.")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	t.Log("Closing FileWriter")
	w.Close()

	if fileExists(filepath) {
		os.Remove(filepath)
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func TestFileWriter_OpenFileTwice(t *testing.T) {
	filepath := "./test.txt"
	if fileExists(filepath) {
		os.Remove(filepath)
	}

	var w writer.Writer
	fileWriter, err := writer.NewFileWriter(filepath)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	w = fileWriter
	err = w.Write("Test message 1.")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	err = w.Write("Test message 2.")
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	fileWriter2, err := writer.NewFileWriter(filepath)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	w = fileWriter2
	err = w.Write("Test message 1-2.")
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	t.Log("Closing FileWriter")
	w.Close()
	if fileExists(filepath) {
		os.Remove(filepath)
	}
}
