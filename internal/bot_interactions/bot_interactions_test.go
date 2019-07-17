package bot_interactions_test

import (
	"strings"
	"telegram_bot/internal/bot_interactions"
	"testing"
)

func TestBuildResultMessage(t *testing.T) {
	ch := make(chan string)
	want := "this is a test string"

	go func() {
		words := strings.Split(want, " ")
		for index := range strings.Split(want, " ") {
			word := words[index]
			if index != len(words)-1 {
				word += " "
			}

			ch <- word
		}

		close(ch)
	}()

	got := bot_interactions.BuildResultMessage(ch)
	if got != want {
		t.Errorf("bot_interactions.buildResultMessage() error!  expectd %q, want %q", want, got)
	}
}

func TestParseLimitFromMsg(t *testing.T) {
	t.Log("Given a number in message it should pass")
	{
		msg := "42"
		got, err := bot_interactions.ParseLimitFromMsg(msg)
		want := 42
		if err != nil {
			t.Errorf("bot_interactions.parseLimitFromMsg() error: %v", err)
		}

		if got != want {
			t.Errorf("bot_interactions.parseLimitFromMsg() error! expected %d, gor %d", want, got)
		}
	}

	t.Log("Given a number with a word it should not pass")
	{
		msg := "42 is a meaning of life"
		_, err := bot_interactions.ParseLimitFromMsg(msg)
		if err == nil {
			t.Error("bot_interactions.parseLimitFromMsg() expected to return an error, got nil")
		}
	}
}
