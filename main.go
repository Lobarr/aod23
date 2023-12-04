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

// instruments a routine
func instrument(ctx string, f func()) {
	t := time.Now()
	f()
	d := time.Now().Sub(t)
	log.Printf(
		"[%s] took %d ms (%d ns) to execute", ctx, d.Milliseconds(), d.Nanoseconds())
}

// helper function for traversing a file line by line. line processors accept a variable used
// to maintain global state and the line to process.
func traverseFile[T any](inputFile string, ctx *T, processLine func(*T, string)) {
	// open file
	f, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("unable to open file %s", inputFile)
	}

	scanner := bufio.NewScanner(f)

	// go through the file lin by line
	for scanner.Scan() {
		line := scanner.Text()
		// skip empty lines and comments
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		processLine(ctx, line)
	}

	if err = scanner.Err(); err != nil {
		log.Fatalf("unable to read contents of file %s", inputFile)
	}

	if err = f.Close(); err != nil {
		log.Fatalf("unable to close file %s", f.Name())
	}
}

func day1Part1() {
	type state struct {
		sumOfCalibrationValues int
		lineCount              int
	}

	inputFile := "./day1-input.txt"
	ctx := &state{}

	traverseFile(inputFile, ctx, func(s *state, line string) {
		digits := [2]rune{}
		start := 0
		end := len(line) - 1
		foundFirstNumber := false
		foundLastNumber := false

		// use two pointers to traverse the line. one that start from the fron and
		// one that starts from the back. they can overlap and one can finish
		// before the other.
		for (start < len(line) && end >= 0) && (!foundFirstNumber || !foundLastNumber) {

			// optimization - if the pointers overlap and we havent't found the lhe last number, we
			// can reuse the first one we've found in order to not recheck parts
			// we've already checked
			if start >= end && foundFirstNumber && !foundLastNumber {
				digits[1] = digits[0]
				foundLastNumber = true
				break
			}

			// handle finding the first number
			if !foundFirstNumber && start < len(line) {
				c := rune(line[start])
				if isNumeric(c) {
					digits[0] = c
					// halt interator of this pointer
					foundFirstNumber = true
				} else {
					// next
					start += 1
				}
			}

			// handle finding the last number
			if !foundLastNumber && end > 0 {
				c := rune(line[end])
				if isNumeric(c) {
					digits[1] = c
					// halt interator of this pointer
					foundLastNumber = true
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

		s.sumOfCalibrationValues += valueNum
		s.lineCount += 1

	})

	log.Printf("[part 1] answer = %d", ctx.sumOfCalibrationValues)
}

func day1Part2() {
	type state struct {
		sumOfCalibrationValues int
		lineCount              int
		wordsToNumber          *trie.Trie[rune]
	}

	inputFile := "./day1-input.txt"
	ctx := &state{
		wordsToNumber: trie.BuildFromMap(map[string]rune{
			"one":   '1',
			"two":   '2',
			"three": '3',
			"four":  '4',
			"five":  '5',
			"six":   '6',
			"seven": '7',
			"eight": '8',
			"nine":  '9',
		}),
	}

	traverseFile(inputFile, ctx, func(s *state, line string) {
		digits := [2]rune{}
		start := 0
		end := 0

		// use sliding window to compute subsequences that could make a word. we
		// use two pointers here and they both start from the front.
		for end < len(line) {
			c := rune(line[end])
			if isNumeric(c) {
				// attempt to persist the first digit
				if digits[0] == 0 {
					digits[0] = c
				}

				// always overwrite the last digit
				digits[1] = c
				start += 1
			} else {
				// ensure we are still building towards a valid word
				for {
					_, isSubsequence := s.wordsToNumber.SubTrie([]byte(line[start:end+1]), true)
					if isSubsequence || start >= end {
						break
					}

					start += 1
				}
				// handle word
				num, prefixLen, ok := s.wordsToNumber.SearchPrefixInString(line[start : end+1])
				if ok && prefixLen == len(line[start:end+1]) {
					if digits[0] == 0 {
						digits[0] = num
					}

					// always overwrite the last digit
					digits[1] = num
					start = end
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

		s.sumOfCalibrationValues += valueNum
		s.lineCount += 1
	})

	log.Printf("[part 2] answer = %d", ctx.sumOfCalibrationValues)
}

func day1(dashes string) {
	log.Printf("%sDay 1%s", dashes, dashes)
	instrument("part 1", day1Part1)
	instrument("part 2", day1Part2)
}

func day2Part1() {
	type state struct {
		gamesSum int
	}

	inputFile := "./day2-input.txt"
	ctx := &state{}

	traverseFile(inputFile, ctx, func(s *state, line string) {
		const maxNumRed = 12
		const maxNumGreen = 13
		const maxNumBlue = 14

		seperator := strings.Index(line, ": ")
		if seperator == -1 {
			log.Fatal("unable to fine seperator ':'")
		}

		gameStr := strings.TrimPrefix(line[:seperator], "Game ")
		game, err := strconv.Atoi(gameStr)
		if err != nil {
			log.Fatalf("unable to convert '%s' to number on line \"%s\"", gameStr, line)
		}

		parts := strings.Split(line[seperator+2:], "; ")

		isColorValid := func(color string, suffix string, max int) bool {
			if strings.HasSuffix(color, suffix) {
				numStr := strings.TrimSuffix(color, suffix)
				num, err := strconv.Atoi(numStr)
				if err != nil {
					log.Fatalf("unable to convert '%s' to number on color \"%s\"", numStr, color)
				}

				if num > max {
					return false
				}
				return true
			}
			return true
		}

		valid := true
		for _, part := range parts {
			colors := strings.Split(part, ", ")

			for _, color := range colors {
				if !isColorValid(color, " red", maxNumRed) ||
					!isColorValid(color, " green", maxNumGreen) ||
					!isColorValid(color, " blue", maxNumBlue) {
					valid = false
				}
			}
		}

		if valid {
			s.gamesSum += game
		}
	})

	log.Printf("[part 1] answer = %d", ctx.gamesSum)
}

func day2Part2() {
	type state struct {
		gamesSum int
	}

	inputFile := "./day2-input.txt"
	ctx := &state{}

	traverseFile(inputFile, ctx, func(s *state, line string) {
		seperator := strings.Index(line, ": ")
		if seperator == -1 {
			log.Fatal("unable to fine seperator ':'")
		}

		parts := strings.Split(line[seperator+2:], "; ")

		handleColor := func(color string, suffix string, rbg *[3]int, i int) {
			if strings.HasSuffix(color, suffix) {
				numStr := strings.TrimSuffix(color, suffix)
				num, err := strconv.Atoi(numStr)
				if err != nil {
					log.Fatalf("unable to convert '%s' to number on color \"%s\"", numStr, color)
				}

				if num > rbg[i] {
					rbg[i] = num
				}
			}
		}

		maxRGB := [3]int{}
		for _, part := range parts {
			colors := strings.Split(part, ", ")

			for _, color := range colors {
				handleColor(color, " red", &maxRGB, 0)
				handleColor(color, " green", &maxRGB, 1)
				handleColor(color, " blue", &maxRGB, 2)
			}
		}

		s.gamesSum += maxRGB[0] * maxRGB[1] * maxRGB[2]
	})

	log.Printf("[part 2] answer = %d", ctx.gamesSum)
}

func day2(dashes string) {
	log.Printf("%sDay 2%s", dashes, dashes)
	instrument("part 1", day2Part1)
	instrument("part 2", day2Part2)
}

func day3Part1() {
	type state struct {
		answer            int
		prevLine          string
		prevNumberWindows [][2]int
		prevSymbols       []int
	}

	// inputFile := "./day3-testinput.txt"
	inputFile := "./day3-input.txt"
	ctx := &state{}

	traverseFile(inputFile, ctx, func(s *state, line string) {
		numberWindows := [][2]int{}
		symbols := []int{}
		start := -1
		end := -1
		
		// prerocess line to gather all contxt needed for evaluation
		for i, c := range line {
			if isNumeric(c) {
				// handle number window
				if start == -1 {
					start = i
					end = i
				} else {
					end = i
				}
			} else {
				// handle end of number window
				if start != -1 && end != -1 {
					numberWindows = append(numberWindows, [2]int{start, end})
					start = -1
					end = -1
				}

				// store symbols
				if c != '.' {
					symbols = append(symbols, i)
				}
			}
		}

		seenLines := map[string]bool{}

		handleNumbersWithAdjacentWindows := func(windows [][2]int, symbolIndexes []int, l string) {
			if l != "" {
				if _, seen := seenLines[line]; !seen {
					log.Println(line)
					seenLines[line] = true
				}

				for _, window := range windows {
					prev := window[0] - 1
					next := window[1] + 1

					if prev < 0 {
						prev = 0
					}

					if next >= len(line) {
						next = len(line) - 1
					}

					for _, symbolIndex := range symbolIndexes {
						// found symbol adjacent to a number window
						if prev <= symbolIndex && symbolIndex <= next {
							numStr := l[window[0] : window[1]+1]

							log.Printf("found symbol adjacent to number %s",numStr)

							num, err := strconv.Atoi(numStr)
							if err != nil {
								log.Fatalf("unable to convert '%s' to number on line \"%s\"", numStr, l)
							}

							ctx.answer += num
						}
					}
				}
			}
		}

		// compare the symbols in the current line with the number of the
		// current line
		handleNumbersWithAdjacentWindows(numberWindows, symbols, line)
		// compare the nubmers of the previous line with the symbols in
		// the current line
		handleNumbersWithAdjacentWindows(ctx.prevNumberWindows, symbols, ctx.prevLine)
		// compare the symbols in the previous line with the numbers of the current
		// line
		handleNumbersWithAdjacentWindows(numberWindows, ctx.prevSymbols, line)

		ctx.prevLine = line
		ctx.prevSymbols = symbols
		ctx.prevNumberWindows = numberWindows
	})

	log.Printf("[part 1] answer = %d", ctx.answer)
}

func day3(dashes string) {
	log.Printf("%sDay 1%s", dashes, dashes)
	instrument("part 1", day3Part1)
	// instrument("part 2", day3Part2)
}

func main() {
	dashes := strings.Repeat("-", 20)
	log.Println("Advent of Code 2023!")
	day1(dashes)
	day2(dashes)
	day3(dashes)
}
