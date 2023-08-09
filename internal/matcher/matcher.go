package matcher

import (
	"encoding/base64"
	"fmt"
	"gmail-unwanted-remover/internal/config"
	"google.golang.org/api/gmail/v1"
	"regexp"
	"strings"
)

type Matcher struct {
	rxForbiddenWords *regexp.Regexp
	rxAllowedWords   *regexp.Regexp
	rxEmail          *regexp.Regexp
	cfg              config.Matcher
}

func New(cfg config.Matcher) *Matcher {
	return &Matcher{
		cfg:              cfg,
		rxForbiddenWords: regexp.MustCompile("(?i)(" + strings.Join(cfg.StopList.Words, "|") + ")"),
		rxAllowedWords:   regexp.MustCompile("(?i)(" + strings.Join(cfg.AllowedWords, "|") + ")"),
		rxEmail:          regexp.MustCompile("(?i)<([^>]+)>"),
	}
}

func (m *Matcher) MatchMessage(msg *gmail.Message) bool {
	// match subject
	if m.rxForbiddenWords.MatchString(m.GetSubject(msg)) {
		return true
	}

	// match body
	for _, p := range msg.Payload.Parts {
		body, _ := base64.URLEncoding.DecodeString(p.Body.Data)
		if m.rxForbiddenWords.MatchString(string(body)) && !m.rxAllowedWords.MatchString(string(body)) {
			return true
		}
	}

	// match email and domain
	if email, err := m.getEmailInFrom(m.getFrom(msg)); err == nil {
		for _, e := range m.cfg.StopList.Emails {
			if e == email {
				return true
			}
		}

		if domain, err := m.getDomainFromEmail(email); err == nil {
			for _, d := range m.cfg.StopList.Domains {
				if d == domain {
					return true
				}
			}
		}
	}

	return false
}

func (m *Matcher) getEmailInFrom(f string) (string, error) {
	if res := m.rxEmail.FindStringSubmatch(f); len(res) == 2 {
		parts := strings.Split(res[1], "@")
		if len(parts) == 2 {
			return res[1], nil
		}
	}
	return "", fmt.Errorf("email not found in from clause")
}

func (m *Matcher) getDomainFromEmail(e string) (string, error) {
	parts := strings.Split(e, "@")
	if len(parts) == 2 {
		return parts[1], nil
	}
	return "", fmt.Errorf("domain not found in given email")
}

func (m *Matcher) GetSubject(msg *gmail.Message) string {
	for _, h := range msg.Payload.Headers {
		if h.Name == "Subject" {
			return h.Value
		}
	}
	return ""
}

func (m *Matcher) getFrom(msg *gmail.Message) string {
	for _, h := range msg.Payload.Headers {
		if h.Name == "From" {
			return h.Value
		}
	}
	return ""
}
