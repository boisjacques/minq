package main

import (
	"os"
	"bufio"
	"fmt"
	"strings"
	"strconv"
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
		} else if string(line[0:2]) == "10" && string(line[2]) != "0"{
			tenMbTest = append(tenMbTest, line)
		} else if string(line[0:3]) == "100"{
			hundredMbTest = append(hundredMbTest, line)
		} else {
			result = append(result, line)
		}
	}

	times2mb := make([]string, 0)
	times10mb := make([]string, 0)
	times100mb := make([]string, 0)


	for _,line := range twoMbTest{
		lineSplit := strings.Split(line, " ")
		i, err := strconv.Atoi(lineSplit[4])
		if err != nil {
			fmt.Println(err)
		}
		i -= 3
		times2mb = append(times2mb, strconv.Itoa(i))
	}

	for _,line := range tenMbTest{
		lineSplit := strings.Split(line, " ")
		i, err := strconv.Atoi(lineSplit[4])
		if err != nil {
			fmt.Println(err)
		}
		i -= 3
		times10mb = append(times10mb, strconv.Itoa(i))
	}

	for _,line := range hundredMbTest{
		lineSplit := strings.Split(line, " ")
		i, err := strconv.Atoi(lineSplit[4])
		if err != nil {
			fmt.Println(err)
		}
		i -= 3
		times100mb = append(times100mb, strconv.Itoa(i))
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

	passed2mb := 0
	passed10mb := 0
	passed100mb := 0
	for _,line := range result{
		if strings.Contains(line, "passed"){
			if strings.Contains(line, "2MB"){
				passed2mb++
			} else if strings.Contains(line, "10MB"){
				passed10mb++
			} else if strings.Contains(line, "100MB"){
				passed100mb++
			}
		}
	}

	fmt.Println("2MB Tests passed: ", passed2mb)
	fmt.Println("10MB Tests passed: ", passed10mb)
	fmt.Println("100MB Tests passed: ", passed100mb)
}
