package main

import (
	"bufio"
	"bytes"
	"log/slog"
	"regexp"
	"strings"
)

type TableOfContents struct {
	Sections []*Section
}

type Section struct {
	Level     int
	ID        string
	LineStart int
	LineEnd   int
}

var headingPrefixPattern = regexp.MustCompile(`^#{1,6}\b`)

func BuildTableOfContents(md []byte, ids []string) (*TableOfContents, error) {
	toc := &TableOfContents{
		Sections: make([]*Section, len(ids)),
	}
	scanner := bufio.NewScanner(bytes.NewReader(md))

	i := 0
	lineno := 0
	for scanner.Scan() {
		lineno++
		line := scanner.Text()

		// Skip non-heading lines.
		if !strings.HasPrefix(line, "#") {
			continue
		}

		// Count leading '#' characters.
		var lvl int
		for line[lvl] == '#' {
			lvl++
		}

		// Skip H1s.
		if lvl == 1 {
			continue
		}

		s := &Section{
			Level:     lvl,
			ID:        ids[i],
			LineStart: lineno,
			LineEnd:   -1,
		}

		// if i > 0 {
		// backtrack:
		// 	for j := i - 1; j >= 0; j-- {
		// 		prev := toc.Sections[j]
		//
		// 		if prev.LineEnd < 0 && prev.Level >= s.Level {
		// 			s.LineEnd = s.LineStart - 1
		// 		}
		//
		// 		if prev.Level == s.Level {
		// 			break backtrack
		// 		}
		// 	}
		// }

		slog.Debug("found section", "level", s.Level, "id", s.ID, "line", s.LineStart)

		toc.Sections[i] = s
		i++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// // Close any dangling headings.
	// for _, s := range toc.Sections {
	// 	if s.LineEnd < 0 {
	// 		s.LineEnd = lineno
	// 	}
	// }

	stack := newStack[*Section]()
	for _, cur := range toc.Sections {
	inner:
		for stack.Len() > 0 {
			prev := stack.Top()
			if prev.Level >= cur.Level {
				stack.Pop()
				prev.LineEnd = cur.LineStart - 1
			} else {
				break inner
			}
		}

		stack.Push(cur)
	}

	for stack.Len() > 0 {
		rem, _ := stack.Pop()
		rem.LineEnd = lineno
	}

	return toc, nil
}

type stack[T any] struct {
	slice []T
}

func newStack[T any]() *stack[T] {
	return &stack[T]{
		slice: make([]T, 0),
	}
}

func (s *stack[T]) Len() int {
	return len(s.slice)
}

func (s *stack[T]) Top() T {
	return s.slice[len(s.slice)-1]
}

func (s *stack[T]) Push(value T) int {
	s.slice = append(s.slice, value)
	return s.Len()
}

func (s *stack[T]) Pop() (value T, ok bool) {
	if len(s.slice) == 0 {
		return *new(T), false
	}

	value = s.slice[len(s.slice)-1]
	s.slice = s.slice[:len(s.slice)-1]
	return value, true
}
