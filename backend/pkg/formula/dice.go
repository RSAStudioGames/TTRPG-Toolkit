package formula

import (
	"fmt"
	"regexp"
	"strconv"
)

var diceTokenRE = regexp.MustCompile(`^(\d+)d(\d+)((kh|kl|dhl)(\d+)|([rxct])(\d+))?$`)

// ParseDiceToken parses a lexer Dice token (e.g. 2d6, 4d6kh3).
func ParseDiceToken(raw string) (*DiceRoll, error) {
	m := diceTokenRE.FindStringSubmatch(raw)
	if m == nil {
		return nil, fmt.Errorf("invalid dice notation %q", raw)
	}
	count, err := strconv.Atoi(m[1])
	if err != nil {
		return nil, fmt.Errorf("invalid dice count in %q", raw)
	}
	sides, err := strconv.Atoi(m[2])
	if err != nil {
		return nil, fmt.Errorf("invalid dice sides in %q", raw)
	}
	d := &DiceRoll{Count: count, Sides: sides}
	if m[3] == "" {
		return d, nil
	}
	if m[4] != "" {
		arg, err := strconv.Atoi(m[5])
		if err != nil {
			return nil, fmt.Errorf("invalid dice modifier in %q", raw)
		}
		d.Mod = &DiceMod{Kind: m[4], Arg: arg}
		return d, nil
	}
	arg, err := strconv.Atoi(m[7])
	if err != nil {
		return nil, fmt.Errorf("invalid dice modifier in %q", raw)
	}
	d.Mod = &DiceMod{Kind: m[6], Arg: arg}
	return d, nil
}
