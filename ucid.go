package ucid

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
	MinContextLen = 2
	MaxContextLen = 8
	MaxTimestamp  = int64(0x0000_0FFF_FFFF_FFFF)
	MaxRandom     = int32(0x000F_FFFF)
	ErrorUCID     = "err_0000_0000_0000_0000"
)

var (
	ErrContextTooShort            = errors.New("context too short")
	ErrContextTooLong             = errors.New("context too long")
	ErrContextOnlyLowerCaseLetter = errors.New("only lower-case letters in context allowed")
	ErrTimestampNegative          = errors.New("the timestamp should not be negative")
	ErrTimestampTooHigh           = errors.New("the timestamp is too high")
	ErrRandomNegative             = errors.New("the random value should not be negative")
	ErrRandomTooHigh              = errors.New("the random value is too high")
	ErrInvalidUCID                = errors.New("the ucid format is not valid")
)

var (
	// RandomGenerator is the function, which is used to generate the random part of the UCID
	RandomGenerator = rand.Int31

	ucidParsePattern = regexp.MustCompile("^([a-z]{2,8})_([0-9a-f]{4})_([0-9a-f]{4})_([0-9a-f]{3})([0-9a-f])_([0-9a-f]{4})$")
)

type UCID string

// Data contains the destructured information which are contained in a UCID
type Data struct {
	Context   string
	Timestamp int64
	Random    int32
}

// DataFromUCID takes a ucid formatted string and returns the parsed data, or error if one is occurred
func DataFromUCID(ucid string) (Data, error) {
	result := ucidParsePattern.FindStringSubmatch(ucid)
	if result == nil {
		return Data{}, ErrInvalidUCID
	}

	context := result[1]
	timestampString := result[2] + result[3] + result[4]
	randomString := result[5] + result[6]

	timestamp, err := strconv.ParseInt(timestampString, 16, 64)
	if err != nil {
		return Data{}, err
	}
	random, err := strconv.ParseInt(randomString, 16, 32)
	if err != nil {
		return Data{}, err
	}

	return Data{
		Context:   context,
		Timestamp: timestamp,
		Random:    int32(random),
	}, nil
}

// ToUCID taks the data from the data structure and generates the UCID string
func (d Data) ToUCID() (UCID, error) {
	// Validate
	if len(d.Context) < MinContextLen {
		return ErrorUCID, ErrContextTooShort
	}
	if len(d.Context) > MaxContextLen {
		return ErrorUCID, ErrContextTooLong
	}
	for _, r := range d.Context {
		if r < 'a' || r > 'z' {
			return ErrorUCID, ErrContextOnlyLowerCaseLetter
		}
	}
	if d.Timestamp < 0 {
		return ErrorUCID, ErrTimestampNegative
	}
	if d.Timestamp > MaxTimestamp {
		return ErrorUCID, ErrTimestampTooHigh
	}
	if d.Random < 0 {
		return ErrorUCID, ErrRandomNegative
	}
	if d.Random > MaxRandom {
		return ErrorUCID, ErrRandomTooHigh
	}

	// Generate string value
	hexPart := fmt.Sprintf("%011x%05x", d.Timestamp, d.Random)
	ucid := strings.Join([]string{
		d.Context,
		hexPart[0:4],
		hexPart[4:8],
		hexPart[8:12],
		hexPart[12:16],
	}, " ")

	return UCID(ucid), nil
}

// NewData generates a new UCID data structure using a given context.
// The timestamp will be set to now and a random number will be generated using RandomGenerator.
// The context should have at least 2 characters and have a maximum length of 8 characters.
// It should consist only of small letters.
func NewData(context string) Data {
	return Data{
		Context:   context,
		Timestamp: time.Now().UnixMilli() % (MaxTimestamp + 1),
		Random:    RandomGenerator() % (MaxRandom + 1),
	}
}

// New generates a UCID string using a given context
// The timestamp will be set to now and a random number will be generated using RandomGenerator.
// The context should have at least 2 characters and have a maximum length of 8 characters.
// It should consist only of small letters.
func New(context string) (UCID, error) {
	return NewData(context).ToUCID()
}

// MustNew generates a UCID string using a given context
// The timestamp will be set to now and a random number will be generated using RandomGenerator.
// The context should have at least 2 characters and have a maximum length of 8 characters.
// It should consist only of small letters.
//
// In case an error occurs, this function will panic.
// Only use this function for cases where it will definitely not fail.
// e.g. when you statically set the context (MustNew("user"))
func MustNew(context string) UCID {
	ucid, err := New(context)
	if err != nil {
		panic(err)
	}
	return ucid
}
