package ucid_test

import (
	"errors"
	"testing"

	"github.com/MatthiasSchild/ucid"
)

func TestMustNew(t *testing.T) {
	// Generate 1000000 entries, all should be generated without error
	for i := 0; i < 1000000; i++ {
		ucid.MustNew("test")
	}
}

func TestDataFromUCID(t *testing.T) {
	_, err := ucid.DataFromUCID("test_0000_0000_0001_1111")
	if err != nil {
		t.FailNow()
	}

	_, err = ucid.DataFromUCID("te_0000_0000_0001_1111")
	if err != nil {
		t.FailNow()
	}

	_, err = ucid.DataFromUCID("testtest_0000_0000_0001_1111")
	if err != nil {
		t.FailNow()
	}

	_, err = ucid.DataFromUCID("testtoolong_0000_0000_0001_1111")
	if err == nil {
		t.FailNow()
	}

	_, err = ucid.DataFromUCID("t_0000_0000_0001_1111")
	if err == nil {
		t.FailNow()
	}

	_, err = ucid.DataFromUCID("test_00000_0000_0001_1111")
	if err == nil {
		t.FailNow()
	}
}

func TestData_ToUCID(t *testing.T) {
	_, err := ucid.Data{
		Context:   "test",
		Timestamp: 0,
		Random:    0,
	}.ToUCID()
	if err != nil {
		t.Logf("Should work, but got error: %v", err)
		t.FailNow()
	}

	_, err = ucid.Data{
		Context:   "",
		Timestamp: 0,
		Random:    0,
	}.ToUCID()
	if !errors.Is(err, ucid.ErrContextTooShort) {
		t.Logf("Should get ErrContextTooShort, but got instead: %v", err)
		t.FailNow()
	}

	_, err = ucid.Data{
		Context:   "testtesttest",
		Timestamp: 0,
		Random:    0,
	}.ToUCID()
	if !errors.Is(err, ucid.ErrContextTooLong) {
		t.Logf("Should get ErrContextTooLong, but got instead: %v", err)
		t.FailNow()
	}

	_, err = ucid.Data{
		Context:   "testA",
		Timestamp: 0,
		Random:    0,
	}.ToUCID()
	if !errors.Is(err, ucid.ErrContextOnlyLowerCaseLetter) {
		t.Logf("Should get ErrContextOnlyLowerCaseLetter, but got instead: %v", err)
		t.FailNow()
	}

	_, err = ucid.Data{
		Context:   "test",
		Timestamp: -1,
		Random:    0,
	}.ToUCID()
	if !errors.Is(err, ucid.ErrTimestampNegative) {
		t.Logf("Should get ErrTimestampNegative, but got instead: %v", err)
		t.FailNow()
	}

	_, err = ucid.Data{
		Context:   "test",
		Timestamp: ucid.MaxTimestamp + 1,
		Random:    0,
	}.ToUCID()
	if !errors.Is(err, ucid.ErrTimestampTooHigh) {
		t.Logf("Should get ErrTimestampTooHigh, but got instead: %v", err)
		t.FailNow()
	}

	_, err = ucid.Data{
		Context:   "test",
		Timestamp: 0,
		Random:    -1,
	}.ToUCID()
	if !errors.Is(err, ucid.ErrRandomNegative) {
		t.Logf("Should get ErrRandomNegative, but got instead: %v", err)
		t.FailNow()
	}

	_, err = ucid.Data{
		Context:   "test",
		Timestamp: 0,
		Random:    ucid.MaxRandom + 1,
	}.ToUCID()
	if !errors.Is(err, ucid.ErrRandomTooHigh) {
		t.Logf("Should get ErrRandomTooHigh, but got instead: %v", err)
		t.FailNow()
	}
}
