package prompt

import (
	"strings"
)

type History struct {
	history        []string
	currentIndex   int
	originalTarget string
	session        *searchSession
}

func (h *History) Add(s string) {
	h.history = append(h.history, s)
}

func (h *History) SearchUp(target string) (string, bool) {
	if !h.session.isSameSession(target) {
		h.session = NewSearchSession(target, h.history)
	}

	return h.session.searchUp(target)
}

func (h *History) SearchDown(target string) (string, bool) {
	if !h.session.isSameSession(target) {
		h.session = NewSearchSession(target, h.history)
	}

	return h.session.searchDown(target)
}

func match(s, t string) bool {
	return strings.Contains(s, t)
}

type searchSession struct {
	originalTarget string
	lastMatchIndex int
	candidates     []string
}

func NewSearchSession(target string, history []string) *searchSession {
	candidates := append([]string{target}, history...)
	candidates = append(candidates, target)

	return &searchSession{
		originalTarget: target,
		lastMatchIndex: len(history) + 1,
		candidates:     candidates,
	}
}

func (s *searchSession) isSameSession(target string) bool {
	if s == nil {
		return false
	}

	return s.candidates[s.lastMatchIndex] == target
}

func (s *searchSession) searchUp(target string) (string, bool) {
	if s.lastMatchIndex == 0 {
		return s.candidates[s.lastMatchIndex], false
	}

	for i := s.lastMatchIndex - 1; i > 0; i-- {
		if match(s.candidates[i], s.originalTarget) {
			s.lastMatchIndex = i
			return s.candidates[i], true
		}
	}

	s.lastMatchIndex = 0
	return s.originalTarget, false
}

func (s *searchSession) searchDown(target string) (string, bool) {
	if s.lastMatchIndex == len(s.candidates)-1 {
		return s.candidates[s.lastMatchIndex], false
	}

	for i := s.lastMatchIndex + 1; i < len(s.candidates)-1; i++ {
		if match(s.candidates[i], s.originalTarget) {
			s.lastMatchIndex = i
			return s.candidates[i], true
		}
	}

	s.lastMatchIndex = len(s.candidates) - 1
	return s.originalTarget, false
}
