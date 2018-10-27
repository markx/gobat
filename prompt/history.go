package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const historyFile = "history.txt"

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
		h.session = newSearchSession(target, h.history)
	}

	return h.session.searchUp(target)
}

func (h *History) SearchDown(target string) (string, bool) {
	if !h.session.isSameSession(target) {
		h.session = newSearchSession(target, h.history)
	}

	return h.session.searchDown(target)
}

func match(s, t string) bool {
	return strings.Contains(s, t)
}

func (h *History) Save() error {
	file, err := os.OpenFile(historyFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	last := h.history[len(h.history)-1]
	_, err = file.WriteString(last + "\n")
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	return nil
}

func (h *History) Load() error {
	file, err := os.Open(historyFile)
	if err != nil {
		return fmt.Errorf("could not open history file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		h.history = append(h.history, scanner.Text())
	}

	return nil
}

type searchSession struct {
	originalTarget string
	lastMatchIndex int
	candidates     []string
}

func newSearchSession(target string, history []string) *searchSession {
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
