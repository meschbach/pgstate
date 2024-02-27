package pgstate

import (
	"math/rand"
	"slices"
	"strings"
	"time"
)

// charClass represents a set of runes to choose for password generation
type charClass []rune

// choose selects a rune from the character class given the entropy
func (c charClass) choose(entropy *rand.Rand) rune {
	index := entropy.Intn(len(c))
	return c[index]
}

var lowerCassClass = charClass("abcdefghijklmnopqrstuvwxyz")
var upperCaseClass = charClass("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
var digitClass = charClass("01234567890")
var specialClass = charClass("!@#$%^&*()_+-= ")

type classSpec struct {
	class     charClass
	remaining int
	maxRun    int
}

func (c classSpec) choose(entropy *rand.Rand, target *strings.Builder, base, max int) (int, bool) {
	total := entropy.Intn(min(max-base, c.remaining, c.maxRun))
	c.remaining = c.remaining - total
	for count := total; count > 0; count-- {
		target.WriteRune(c.class.choose(entropy))
	}
	return total, c.remaining == 0
}

func GeneratePassword() string {
	cfg := GenPasswordConfig{AllowSpecial: true}
	return GeneratePasswordWithConfig(cfg)
}

type GenPasswordConfig struct {
	AllowSpecial bool
}

func GeneratePasswordWithConfig(config GenPasswordConfig) string {
	alphaLower := classSpec{class: lowerCassClass, maxRun: 16, remaining: 64}
	alphaUpper := classSpec{class: upperCaseClass, maxRun: 16, remaining: 64}
	digits := classSpec{class: digitClass, maxRun: 4, remaining: 32}
	classes := []*classSpec{&alphaLower, &alphaUpper, &digits}
	if config.AllowSpecial {
		special := classSpec{class: specialClass, maxRun: 4, remaining: 16}
		classes = append(classes, &special)
	}

	entropy := rand.New(rand.NewSource(time.Now().Unix()))
	size := 64
	limit := size
	output := &strings.Builder{}
	output.Grow(size)
	base := 0
	for base < limit {
		classIndex := entropy.Intn(len(classes))
		c := classes[classIndex]
		count, remove := c.choose(entropy, output, base, size+1)
		base += count
		if remove {
			classes = slices.Delete(classes, classIndex, classIndex)
			if len(classes) == 0 {
				panic("ran out of characters")
			}
		}
	}
	return output.String()
}
