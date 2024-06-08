package actions

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
)

// Contains the literal JSON-string to send to Firebase.
// It only contains its own part so actions can be combined later.
// For example:
// `"state": "ON"` or
// `"color": "#ff0000"`
type Action struct {
	Type  string
	Path  string
	Value interface{}
}

func Parse(s string) ([]Action, error) {
	if len(s) == 0 {
		return []Action{}, errors.New("s cannot be empty")
	}

	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "")
	split := strings.Split(s, ",")

	var allErrors error
	actions := make([]Action, 0, len(split))

	for _, actionString := range split {
		action, err := parseAction(actionString)
		if err != nil {
			allErrors = multierror.Append(allErrors, err)
			continue
		}
		actions = append(actions, action)
	}

	return actions, allErrors
}

var UnknownAction = errors.New("unknown action")

var colors = map[string]string{
	"red":        "#ff0000",
	"orange":     "#f79719",
	"yellow":     "#ffff00",
	"green":      "#00ff00",
	"blue":       "#0000ff",
	"lightblue":  "#00eeff",
	"light_blue": "#00eeff",
	"purple":     "#a600ff",
	"magenta":    "#ff00e6",
}

// Parses one action.
// Actions can look like either one of these:
// - 'on' (turns the light on)
// - 'off' (turns the light off)
// - '[color]' (sets the color)
// - '[n]' (sets the brightness to n%)
func parseAction(a string) (Action, error) {
	switch a {
	case "on":
		return Action{Type: "on", Path: "state", Value: "ON"}, nil
	case "off":
		return Action{Type: "off", Path: "state", Value: "OFF"}, nil
	default:
		hexColor, ok := colors[a]
		if ok {
			// This is a [color] action.
			return Action{Type: "color", Path: "color", Value: hexColor}, nil
		}
		brightness, err := strconv.Atoi(a)
		if err == nil {
			decimalBrightness := float32(brightness) / 100. * 255 // between 0 and 255
			return Action{Type: "brightness", Path: "brightness", Value: decimalBrightness}, nil
		}
	}

	return Action{}, fmt.Errorf("%q: %w", a, UnknownAction)
}
