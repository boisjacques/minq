package main

import (
	"os"
	"bufio"
	"fmt"
	"strings"
)

func main(){
	filePath := os.Args[1]
	twoMbTest := make([]string, 0)
	tenMbTest := make([]string, 0)
	hundredMbTest := make([]string, 0)
	result := make([]string, 0)

	// Open the file.
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Create a new Scanner for the file.
	scanner := bufio.NewScanner(f)
	// Loop over all lines in the file and print them.
	for scanner.Scan() {
		line := scanner.Text()
		if string(line[0]) == "2"{
			twoMbTest = append(twoMbTest, line)
		} else if string(line[:1]) == "10"{
			tenMbTest = append(tenMbTest, line)
		} else if string(line[:2]) == "100"{
			hundredMbTest = append(hundredMbTest, line)
		} else {
			result = append(result, line)
		}
	}

	times2mb := make([]string, 0)
	times10mb := make([]string, 0)
	times100mb := make([]string, 0)
	results := make([]string, 0)


	for _,line := range twoMbTest{
		lineSplit := strings.Split("  ", line)
		times2mb = append(times2mb, lineSplit[1])
	}

	for _,line := range tenMbTest{
		lineSplit := strings.Split("  ", line)
		times10mb = append(times10mb, lineSplit[1])
	}

	for _,line := range hundredMbTest{
		lineSplit := strings.Split("  ", line)
		times100mb = append(times100mb, lineSplit[1])
	}

	for _,line := range results{
		lineSplit := strings.Split(" ", line)
		times2mb = append(times2mb, lineSplit[3])
	}

	file, err := os.Create("2mb.result")
	if err != nil {
		return
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range times2mb {
		fmt.Fprintln(w, line)
	}


	file, err = os.Create("10mb.result")
	if err != nil {
		return
	}
	defer file.Close()

	w = bufio.NewWriter(file)
	for _, line := range times2mb {
		fmt.Fprintln(w, line)
	}


	file, err = os.Create("100mb.result")
	if err != nil {
		return
	}
	defer file.Close()

	w = bufio.NewWriter(file)
	for _, line := range times2mb {
		fmt.Fprintln(w, line)
	}

	passed := 0
	for _,line := range results{
		if strings.Contains(line, "passed"){
			passed++
		}
	}

	fmt.Println("Tests passed: ", passed)
}
