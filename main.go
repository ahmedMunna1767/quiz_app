package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

/* Shuffles the Problem Set */
func shuffleProblemList(problems []problem) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(problems), func(i, j int) { problems[i], problems[j] = problems[j], problems[i] })
}

func main() {
	/* Declaring flags */
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	shuffle := flag.Bool("shuffle", false, "Shuffles the question order")

	flag.Parse()

	/* Open Csv file */
	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilename))
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the provided CSV file.")
	}

	/* Generate and shuffle the problem Set */
	problems := parseLines(lines)

	if *shuffle == true {
		shuffleProblemList(problems)
	}

	fmt.Printf("Answer %d questions in %d seconds. Lets GO...\n", len(problems), *timeLimit)
	fmt.Println("Press 'Enter' to start")

	var start string
	fmt.Scanf("%s", &start)

	/* Start Timer */
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	correct := 0

	/* loop over problems */
	/* go routine for both timer and the user */
problemloop:
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)

		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- strings.TrimSpace(answer)
		}()

		select {
		case <-timer.C:
			fmt.Println()
			break problemloop
		case answer := <-answerCh:
			if answer == p.a {
				correct++
			}
		}
	}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}

/* Create problem structs from the raw csv data */
func parseLines(lines [][]string) []problem {
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
