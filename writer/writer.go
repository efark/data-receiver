/*
Package writer defines the Writer interface and has a memory writer for testing purposes.
To use this application, you should implement a writer that suits your needs.
*/
package writer

import (
	"os"
	"strings"
	"time"
)

// Writer interface has Write and Close methods.
type Writer interface {
	Write(content string) error
	Close()
}

// CreateWriter generates the right writer based on the received parameters.
func CreateWriter(class string, params map[string]string) (Writer, error) {
	var w Writer
	var err error
	switch class {
	case "MemoryWriter":
		w, err = NewMemoryWriter()
	case "FileWriter":
		w, err = NewFileWriter(params["filepath"])
	default:
		w, err = NewConsoleWriter()
	}
	return w, err
}

// ConsoleWriter is an empty struct, this one just logs messages to the apps output.
type ConsoleWriter struct {
}

// NewConsoleWriter creates a ConsoleWriter.
func NewConsoleWriter() (*ConsoleWriter, error) {
	return &ConsoleWriter{}, nil
}

// Write logs the message.
func (w *ConsoleWriter) Write(content string) error {
	log.Info("Message received: " + content)
	return nil
}

// Close just logs a closing message.
func (w *ConsoleWriter) Close() {
	log.Info("Closing ConsoleWriter.")
}

// MemoryWriter stores messages in an internal []string.
type MemoryWriter struct {
	Messages []string
}

// NewMemoryWriter generates an empty struct.
func NewMemoryWriter() (*MemoryWriter, error) {
	log.Info("Starting MemoryWriter.")
	return &MemoryWriter{}, nil
}

// Write appends the message in the MemoryWriter.
func (w *MemoryWriter) Write(content string) error {
	log.Info("Storing message in MemoryWriter.")
	w.Messages = append(w.Messages, content)
	return nil
}

// GetMessages returns a []string with all the messages.
func (w *MemoryWriter) GetMessages() []string {
	return w.Messages
}

// Close deletes all the messages from the MemoryWriter.
func (w *MemoryWriter) Close() {
	log.Info("Closing MemoryWriter.")
	w.Messages = []string{}
}

// FileWriter has all the fields necessary to write the messages into a local file.
type FileWriter struct {
	file     *os.File
	deadline time.Time
	size     int64
	mchan    chan string
	done     chan bool
}

// NewFileWriter stores the filepath in an inner field.
func NewFileWriter(filepath string) (*FileWriter, error) {
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		return &FileWriter{}, err
	}

	f := &FileWriter{file: file, mchan: make(chan string), done: make(chan bool)}
	go f.fileWrite()

	return f, err
}

// Write sends the message into an inner channel.
func (w *FileWriter) Write(content string) error {
	log.Info("Storing message in FileWriter.")
	w.mchan <- content
	return nil
}

// Close closes the inner channel and sends a done message through another channel to finish the writing method.
func (w *FileWriter) Close() {
	log.Info("Closing FileWriter.")
	close(w.mchan)
	<-w.done
	err := w.file.Close()
	if err != nil {
		slog.Error(err)
	}
}

// fileWrite receives the messages from the channel and waits for the done channel to receive a message.
func (w *FileWriter) fileWrite() {
	log.Info("Start writing to file.")

	for m := range w.mchan {
		var newline string
		if !strings.HasSuffix(m, "\n") {
			newline = "\n"
		}
		n, err := w.file.WriteString(m + newline)
		if err != nil {
			slog.Error(err.Error())
			return
		}
		w.size += int64(n)
	}

	w.done <- true

	return
}
