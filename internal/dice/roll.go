package dice

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	diceRollRegex = `^(\d*)d(\d+)(([+-]\d+)+)?$`
	modifierRegex = `[+-]?\d+`
)

// Roll parses the input string and returns the result of the dice roll along with the individual rolls
func Roll(input string) (string, int, error) {
	// Remove all spaces from the input string
	input = strings.ReplaceAll(input, " ", "")

	// Define a regular expression to match the dice roll pattern
	matches := regexp.MustCompile(diceRollRegex).FindStringSubmatch(input)

	// Check if input format is correct
	if matches == nil || len(matches) < 1 {
		return "", 0, errors.New("invalid input format")
	}

	// Parse the number of rolls (default to 1 if not present)
	numRolls := 1
	if matches[1] != "" {
		var err error
		numRolls, err = strconv.Atoi(matches[1])
		if err != nil {
			return "", 0, fmt.Errorf("invalid number of rolls: %v", err)
		}
	}

	// Parse the number of sides on the dice
	sides, err := strconv.Atoi(matches[2])
	if err != nil {
		return "", 0, fmt.Errorf("invalid number of sides: %v", err)
	}

	// Parse the modifiers and evaluate it
	modifiers := matches[3]
	modifier, err := evaluateModifiers(modifiers)
	if err != nil {
		return "", 0, fmt.Errorf("invalid modifier: %v", err)
	}

	// Create a new rand.Rand instance with a seed
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Roll the dice and sum the results
	totalRoll := 0
	var rolls []int
	for i := 0; i < numRolls; i++ {
		roll := r.Intn(sides) + 1
		rolls = append(rolls, roll)
		totalRoll += roll
	}

	// Add the modifier to the total roll
	result := totalRoll + modifier

	// Construct the output string showing individual rolls and modifiers
	var rollStrings []string
	for _, roll := range rolls {
		rollStrings = append(rollStrings, "("+strconv.Itoa(roll)+")")
	}
	rollOutput := strings.Join(rollStrings, " + ")
	if modifiers != "" {
		// Split the modifiers into individual terms with spaces
		modifierTerms := regexp.MustCompile(modifierRegex).FindAllString(modifiers, -1)
		for _, modifierTerm := range modifierTerms {
			if string(modifierTerm[0]) == "+" {
				modifierTerm = modifierTerm[1:]
			}
			rollOutput += " + (" + modifierTerm + ")"
		}
	}

	return rollOutput, result, nil
}

// evaluateModifiers evaluates the arithmetic expression in the modifiers
func evaluateModifiers(modifiers string) (int, error) {
	if modifiers == "" {
		return 0, nil
	}

	// Split the modifiers into individual terms
	terms := regexp.MustCompile(modifierRegex).FindAllString(modifiers, -1)
	total := 0

	// Sum up all the terms
	for _, term := range terms {
		value, err := strconv.Atoi(term)
		if err != nil {
			return 0, err
		}
		total += value
	}

	return total, nil
}
