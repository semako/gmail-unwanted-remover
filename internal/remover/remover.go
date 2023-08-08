package remover

import (
	"fmt"
	"gmail-unwanted-remover/internal/gmail"
	"gmail-unwanted-remover/internal/matcher"
	"time"
)

type Remover struct {
	checkInterval time.Duration
	gmail         *gmail.Gmail
	matcher       *matcher.Matcher
}

func New(checkInterval time.Duration, gmail *gmail.Gmail, matcher *matcher.Matcher) *Remover {
	return &Remover{
		checkInterval: checkInterval,
		gmail:         gmail,
		matcher:       matcher,
	}
}

func (r *Remover) processInbox() error {
	l, err := r.gmail.ListInbox()
	if err != nil {
		return fmt.Errorf("unable to list inbox messages: %w", err)
	}
	for _, m := range l.Messages {
		if m, err = r.gmail.GetMessage(m.Id); err != nil {
			return fmt.Errorf("unable to load message: %w", err)
		}

		if r.matcher.MatchMessage(m) {
			fmt.Printf("Message w/subject is being deleted: %s\n", r.matcher.GetSubject(m))
			if err = r.gmail.DeleteMessage(m.Id); err != nil {
				return fmt.Errorf("unable to delete message: %w", err)
			}
		}
	}
	return nil
}

func (r *Remover) processSpam() error {
	l, err := r.gmail.ListSpam()
	if err != nil {
		return fmt.Errorf("unable to list spam messages: %w", err)
	}
	for _, m := range l.Messages {
		if m, err = r.gmail.GetMessage(m.Id); err != nil {
			return fmt.Errorf("unable to load message: %w", err)
		}
		fmt.Printf("Spam message w/subject is being deleted: %s\n", r.matcher.GetSubject(m))
		if err = r.gmail.DeleteMessage(m.Id); err != nil {
			return fmt.Errorf("unable to delete message: %w", err)
		}
	}
	return nil
}

func (r *Remover) Daemon() error {
	ticker := time.NewTicker(r.checkInterval)
	for range ticker.C {
		fmt.Printf("tick time\n")
		if err := r.processInbox(); err != nil {
			return fmt.Errorf("remover: daemon: processInbox failed: %w", err)
		}
		if err := r.processSpam(); err != nil {
			return fmt.Errorf("remover: daemon: processSpam failed: %w", err)
		}
	}
	return nil
}
