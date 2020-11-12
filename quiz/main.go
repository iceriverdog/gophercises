package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFileName := flag.String("csv", "problems.csv", "a csv file in the format "+
		"of `question, answer`")
	limit := flag.Int("limit", 30, "the time limit for the quiz in seconds, default 30s")
	flag.Parse()

	file, err := os.Open(*csvFileName)
	if err != nil {
		exit(fmt.Sprintf("faild open %v", *csvFileName))
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit(fmt.Sprintf("%v read falied", csvFileName))
	}

	problems := parseProblems(lines)
	timer := time.NewTimer(time.Duration(*limit) * time.Second)

	scores := 0
problemsLoop:
	for i, p := range problems {
		fmt.Printf("problem #%d: %s\n", i+1, p.q)
		answerCh := make(chan string)
		var answer string
		go func() {
			fmt.Scanf("%s", &answer)
			answerCh <- answer
		}()
		select {
		case <-timer.C:
			break problemsLoop
		case answer = <-answerCh:
			if answer == p.a {
				scores++
			}
		}
	}
	fmt.Printf("you scored %d out of %d\n", scores, len(problems))
}
func parseProblems(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
