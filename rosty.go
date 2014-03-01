/*
   Small app to manage entries in your host file on Mac OSX.
*/
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// Path of the hosts file
// const hostFile string = "hosts"

const hostFile string = "/etc/hosts"

// Print the help for the app
func printHelp() {
	fmt.Println("Example Usage")
	fmt.Println("rosty [action] [options]")
	fmt.Println("actions:")
	fmt.Println("\tadd\t[ip] [host]\tAdd an entry in the file")
	fmt.Println("\tget\tReturns the content of the host file")
	fmt.Println("\tdel\tDelete an entry from the host file")
}

// Check for an error and panic is needed
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Read the hosts file and dump the content as an array of strings
func getHostFileContent() []string {
	dat, err := ioutil.ReadFile(hostFile)
	check(err)
	return strings.Split(string(dat), "\n")
}

// Parse the options from the command line
func parseOptions() map[string]string {
	argsWithoutProg := os.Args[1:]
	actions := map[string]bool{"add": true, "del": true, "get": true}
	options := make(map[string]string)
	optionSize := len(argsWithoutProg)

	if optionSize != 0 {

		if val := argsWithoutProg[0]; actions[val] {
			options["action"] = val
		}

		if optionSize == 2 {
			options["index"] = argsWithoutProg[1]

		} else if optionSize == 3 {
			options["ip"] = argsWithoutProg[1]
			options["host"] = argsWithoutProg[2]
		}

	}
	return options
}

func displayItems() {
	lines := getHostFileContent()

	for key, value := range lines {
		if value != "" {
			fmt.Printf("%d -> %s\n", key, value)
		}
	}
}

func writeTofile(newline string) (int, error) {
	f, err := os.OpenFile(hostFile, os.O_WRONLY|os.O_APPEND, 0644)
	check(err)

	defer f.Close()
	n, err := f.WriteString(newline)

	return n, err
}

func addItem(item string) error {
	_, err := writeTofile(item)
	return err
}

func getItem() int {
	var i int
	fmt.Println("Pick an entry to delete")
	displayItems()
	fmt.Print("-> ")
	fmt.Scanf("%d", &i)
	return i
}

func delItem() error {
	lines := getHostFileContent()
	item := getItem()
	f, err1 := os.OpenFile(hostFile, os.O_WRONLY|os.O_TRUNC, 0644)
	check(err1)
	defer f.Close()

	var err2 error

	for key, value := range lines {
		if key != item {
			_, err2 = f.WriteString(fmt.Sprintf("%s\n", value))
		}
	}
	return err2
}

func printError(err error) {
	fmt.Println("ERROR:")
	fmt.Println(err)
}

func makeBackup() {
	backup := fmt.Sprintf("%s.bk", hostFile)
	if _, err := os.Stat(backup); os.IsNotExist(err) {
		cmd := exec.Command("cp", "/etc/hosts", "/etc/hosts.bk")
		err := cmd.Run()
		if err != nil {
			printError(err)
		}
	}
}

func main() {
	options := parseOptions()
	if len(options) == 0 {
		printHelp()
		os.Exit(1)
	}
	makeBackup()
	switch action := options["action"]; {
	case action == "get":
		displayItems()
	case action == "add":
		if err := addItem(fmt.Sprintf("%s %s", options["ip"], options["host"])); err != nil {
			printError(err)
			os.Exit(1)
		}
		fmt.Println("Entry added successfully")
	case action == "del":
		if err := delItem(); err != nil {
			printError(err)
			os.Exit(1)
		}
		fmt.Println("Entry deleted successfully")
	default:
		printHelp()
	}
}
