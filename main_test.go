package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func Test_isPrime(t *testing.T) {
	primeTests := []struct {
		name     string
		testNum  int
		expected bool
		msg      string
	}{
		{"prime", 7, true, "7 is a prime number!"},
		{"not prime", 8, false, "8 is not a prime number because it is divisible by 2!"},
		{"zero", 0, false, "0 is not prime, by definition!"},
		{"one", 1, false, "1 is not prime, by definition!"},
		{"negative number", -11, false, "Negative numbers are not prime, by definition!"},
	}

	for _, e := range primeTests {
		result, msg := isPrime(e.testNum)
		if e.expected && !result {
			t.Errorf("%s: expected true but got false", e.name)
		}

		if !e.expected && result {
			t.Errorf("%s: expected false but got true", e.name)
		}

		if e.msg != msg {
			t.Errorf("%s: expected %s but got %s", e.name, e.msg, msg)
		}
	}
}

func Test_prompt(t *testing.T) {
	// Replace Stdout with buffer
	old := os.Stdout
	defer func() {
		os.Stdout = old
	}()

	r, w, _ := os.Pipe()
	os.Stdout = w

	prompt()

	// Read stdout into buffer
	var buf bytes.Buffer
	w.Close()
	io.Copy(&buf, r)
	r.Close()

	expected := "-> "
	actual := buf.String()
	if expected != actual {
		t.Errorf("Expected '%s', but got '%s'", expected, actual)
	}
}

func Test_intro(t *testing.T) {
	// Replace Stdout with buffer
	old := os.Stdout
	defer func() {
		os.Stdout = old
	}()

	r, w, _ := os.Pipe()
	os.Stdout = w

	intro()

	// Read stdout into buffer
	var buf bytes.Buffer
	w.Close()
	io.Copy(&buf, r)
	r.Close()

	expectedOutput := "Is it Prime?\n------------\nEnter a whole number, and we'll tell you if it is a prime number or not. Enter q to quit.\n-> "
	actualOutput := buf.String()

	if expectedOutput != actualOutput {
		t.Errorf("Expected output:\n%s\nBut got:\n%s", expectedOutput, actualOutput)
	}
}

func Test_checkNumbers(t *testing.T) {

	inputStr := "7\nq\nabc\n"
	scanner := bufio.NewScanner(strings.NewReader(inputStr))
	scanner.Split(bufio.ScanLines)

	expectedOutputs := []string{"7 is a prime number!", "", "Please enter a whole number!"}
	expectedQuits := []bool{false, true, false}

	for i := 0; i < len(expectedOutputs); i++ {
		output, quit := checkNumbers(scanner)
		if output != expectedOutputs[i] {
			t.Errorf("Expected output '%s', but got '%s'", expectedOutputs[i], output)
		}
		if quit != expectedQuits[i] {
			t.Errorf("Expected quit value '%t', but got '%t'", expectedQuits[i], quit)
		}
	}

	if scanner.Scan() {
		t.Error("Scanner did not read all lines from input")
	}
}

func Test_readUserInput(t *testing.T) {

	inputStr := "q\n"
	buf := bytes.NewBufferString(inputStr)

	doneChan := make(chan bool)
	go readUserInput(buf, doneChan)

	expectedQuit := true

	quit := <-doneChan
	if quit != expectedQuit {
		t.Errorf("Expected quit value '%t', but got '%t'", expectedQuit, quit)
	}
	close(doneChan)
}
