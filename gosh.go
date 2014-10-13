// This is a comment.
package main

import "bufio"
import "io/ioutil"
import "fmt"
import "os"
import "os/exec"
import "path"

func stripFirstLine(script []byte) []byte {
	i := 0
	for script[i] != '\n' {
		i += 1
	}
	return script[i:]
}

func main() {
	if len(os.Args) != 2 {
		println("wrong number of args")
		os.Exit(1)
	}
	scriptPath := os.Args[1]

	tempDirName, err := ioutil.TempDir("", "gosh")
	if err != nil {
		println("unable to create temp directory")
		os.Exit(1)
	}
	defer os.RemoveAll(tempDirName)
	println("tempDirName:", tempDirName)

	// inputBytes, err := ioutil.ReadFile(scriptPath)
	// if err != nil {
	// 	println("error reading input file", err)
	// 	os.Exit(1)
	// }

	tempFileName := path.Join(
		tempDirName, path.Base(scriptPath) + ".go")

	inputFile, err := os.Open(scriptPath)
	if err != nil {
		println("error opening input file")
		os.Exit(1)
	}
	defer inputFile.Close()
	lines := make(chan string)
	go func() {
		scanner := bufio.NewScanner(inputFile)
		for scanner.Scan() {
			lines <- scanner.Text()
		}
		close(lines)
	}()

	outputFile, err := os.Create(tempFileName)
	if err != nil {
		println("error opening output file", err)
		os.Exit(1)
	}
	defer outputFile.Close()
	
	// for {
	// 	line, more := <- lines
	// 	if ! more {
	// 		break
	// 	}
	// 	println("line:", line)
	// }
	lineNumber := 0
	for line := range lines {
		lineNumber++
		if lineNumber == 1 {
			fmt.Fprintln(outputFile, "package main")
		} else {
			fmt.Fprintln(outputFile, line)

		}
	}

	catCmd := exec.Command("nl", tempFileName)
	catCmd.Stdout = os.Stdout
	catCmd.Run()

	// var outputBytes []byte
	// outputBytes = append(outputBytes, "package main\n")
	// outputBytes = stripFirstLine(inputBytes)

	// err = ioutil.WriteFile(tempFileName, outputBytes, 0400)
	// if err != nil {
	// 	println("error writing temp file")
	// 	os.Exit(1)
	// }

	path, err := exec.LookPath("go")
	if err != nil {
		println("error finding 'go'", err)
		os.Exit(1)
	}
	cmd := &exec.Cmd{
		Path: path,
		Args: []string{"go", "run", tempFileName},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	err = cmd.Run()
	if err != nil {
		fmt.Println(err);
		return
	}
}
