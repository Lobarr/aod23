package main

import (
	"bufio"
	"github.com/porfirion/trie"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// uses ASCII to determine if a charater is a number
func isNumeric(char rune) bool {
	n := int(char)
	return n >= 48 && n <= 57
}

func day1Part1() {
	inputFile := "./day1-input.txt"

	// open file
	f, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("unable to open file %s", inputFile)
	}

	scanner := bufio.NewScanner(f)
	sumOfCalibrationValues := 0
	lineCount := 0

	// go through the file lin by line
	for scanner.Scan() {
		digits := [2]rune{}
		line := scanner.Text()
		start := 0
		end := len(line) - 1
		foundFirstNumber := false
		foundLastNumber := false

		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		// use two pointers to traverse the line. one that start from the fron and
		// one that starts from The back. they can overlap and one can finish
		// befeore the other.
		for (start < len(line) && end >= 0) && (!foundFirstNumber || !foundLastNumber) {
			// if the pointers overlap and we havent't found the lhe last number, we
			// can r{jeuse the first one we've found in order to not recheck parts
			// we've already checked. slight optimization
			if start >= end && foundFirstNumber && !foundLastNumber {
				// log.Printf("for line %s, deducing the last number using the first since the pointer overlap", line)
				digits[1] = digits[0]
				foundLastNumber = true
				break
			}

			if !foundFirstNumber && start < len(line) {
				c := rune(line[start])
				if isNumeric(c) {
					digits[0] = c
					// hal{bt interator of this pointer
					foundFirstNumber = true
					//log.Printf("found first number: val - %c, head - %d, tail - %d, len - %d", c, start, end, len(line))
				} else {
					// next
					start += 1
				}
			}

			if !foundLastNumber && end > 0 {
				c := rune(line[end])
				if isNumeric(c) {
					digits[1] = c
					// halt interator of this pointer
					foundLastNumber = true
					//log.Printf("found last number: val - %c, head - %d, tail - %d, len - %d", c, start, end, len(line))
				} else {
					// previous
					end -= 1
				}
			}
		}

		value := string(digits[:])
		valueNum, err := strconv.Atoi(value)

		if err != nil {
			log.Fatalf("unable to convert '%s' to number", value)
		}

		sumOfCalibrationValues += valueNum
		lineCount += 1
		//log.Printf("'%s'--> %s", line, value)
	}

	if err = scanner.Err(); err != nil {
		log.Fatalf("unable to read contents of file %s", inputFile)
	}

	if err = f.Close(); err != nil {
		log.Fatalf("unable to close file %s", f.Name())
	}

	log.Printf("[part 1] answer = %d", sumOfCalibrationValues)
}

func day1Part2() {
	inputFile := "./day1-input.txt"
	//inputFile := "./day1-testinput.txt"

	// open file
	f, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("unable to open file %s", inputFile)
	}

	scanner := bufio.NewScanner(f)
	sumOfCalibrationValues := 0
	lineCount := 0
	wordsToNumber := trie.BuildFromMap(map[string]rune{
		"one":   '1',
		"two":   '2',
		"three": '3',
		"four":  '4',
		"five":  '5',
		"six":   '6',
		"seven": '7',
		"eight": '8',
		"nine":  '9',
	})

	// go through the file lin by line
	for scanner.Scan() {
		digits := [2]rune{}
		line := scanner.Text()
		start := 0
		end := 0

		// log.Printf("line: %s", line)

		// handle comments and empty lines
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		for end < len(line) {
			c := rune(line[end])
			if isNumeric(c) {
				// log.Printf("found number: val - %c, head - %d, tail - %d", c, start, end)

				if digits[0] == 0 {
					digits[0] = c
				}

				// always overwrite the last digit
				digits[1] = c
				start += 1
			} else {
				// ensure we are still building towards a valid word
				for {
					_, isSubsequence := wordsToNumber.SubTrie([]byte(line[start:end+1]), true)
					if isSubsequence {
						// log.Printf("is subsequence: %s (start - %d, end - %d, line - %s)", line[start:end+1], start, end, line)
						break
					}

					if start >= end {
						break
					}

					start += 1
				}

				//
				if num, prefixLen, ok := wordsToNumber.SearchPrefixInString(line[start:end+1]); ok {
					if prefixLen == len(line[start:end+1]) {
						// log.Printf("[start trail] match word %s to nubmer %c for line %s", line[start:end+1], num, line)

						if digits[0] == 0 {
							digits[0] = num
						}

						// always overwrite the last digit
						digits[1] = num
						start = end
					}
				}
			}

			// next
			end += 1
		}

		value := string(digits[:])
		valueNum, err := strconv.Atoi(value)

		if err != nil {
			log.Fatalf("unable to convert '%s' (digits = %v) to number on line \"%s\"", value, digits, line)
		}

		sumOfCalibrationValues += valueNum
		lineCount += 1
		// log.Printf("'%s'--> %s", line, value)
	}

	if err = scanner.Err(); err != nil {
		log.Fatalf("unable to read contents of file %s", inputFile)
	}

	if err = f.Close(); err != nil {
		log.Fatalf("unable to close file %s", f.Name())
	}

	log.Printf("[part 2] answer = %d", sumOfCalibrationValues)
}

func main() {
	dashes := strings.Repeat("-", 20)
	log.Println("Advent of Code 2023!")

	// Day 1
	log.Printf("%sDay 1%s", dashes, dashes)

	// Day 1 Part 1
	day1Part1T := time.Now()
	day1Part1()
	day1Part1D := time.Now().Sub(day1Part1T)
	log.Printf(
		"[part 1] took %d ms (%d ns) to execute", day1Part1D.Milliseconds(), day1Part1D.Nanoseconds())

	// Day 1 Part 2
	day1Part2T := time.Now()
	day1Part2()
	day1Part2D := time.Now().Sub(day1Part2T)
	log.Printf(
		"[part 2] took %d ms (%d ns) to execute", day1Part2D.Milliseconds(), day1Part2D.Nanoseconds())

}
