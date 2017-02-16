package highlight

import (
	"regexp"
	"strings"
)

func combineLineMatch(src, dst LineMatch) LineMatch {
	for k, v := range src {
		if g, ok := dst[k]; ok {
			if g == "" {
				dst[k] = v
			}
		} else {
			dst[k] = v
		}
	}
	return dst
}

type State *Region

type LineStates interface {
	Lines() [][]byte
	State(lineN int) State
	SetState(lineN int, s State)
}

type Highlighter struct {
	endRegions []*Region
	def        *Def
}

func NewHighlighter(def *Def) *Highlighter {
	h := new(Highlighter)
	h.def = def
	return h
}

type LineMatch map[int]string

func FindIndex(regex *regexp.Regexp, str []byte, canMatchStart, canMatchEnd bool) []int {
	regexStr := regex.String()
	if strings.Contains(regexStr, "^") {
		if !canMatchStart {
			return nil
		}
	}
	if strings.Contains(regexStr, "$") {
		if !canMatchEnd {
			return nil
		}
	}
	return regex.FindIndex(str)
}

func FindAllIndex(regex *regexp.Regexp, str []byte, canMatchStart, canMatchEnd bool) [][]int {
	regexStr := regex.String()
	if strings.Contains(regexStr, "^") {
		if !canMatchStart {
			return nil
		}
	}
	if strings.Contains(regexStr, "$") {
		if !canMatchEnd {
			return nil
		}
	}
	return regex.FindAllIndex(str, -1)
}

func (h *Highlighter) highlightRegion(start int, canMatchEnd bool, lineNum int, line []byte, region *Region) LineMatch {
	highlights := make(LineMatch)

	if len(line) == 0 {
		if canMatchEnd {
			h.endRegions[lineNum] = region
		}

		return highlights
	}

	loc := FindIndex(region.end, line, start == 0, canMatchEnd)
	if loc != nil {
		if region.parent == nil {
			highlights[start+loc[1]] = ""
			return combineLineMatch(highlights,
				combineLineMatch(h.highlightRegion(start, false, lineNum, line[:loc[0]], region),
					h.highlightEmptyRegion(start+loc[1], canMatchEnd, lineNum, line[loc[1]:])))
		}
		highlights[start+loc[1]] = region.parent.group
		return combineLineMatch(highlights,
			combineLineMatch(h.highlightRegion(start, false, lineNum, line[:loc[0]], region),
				h.highlightRegion(start+loc[1], canMatchEnd, lineNum, line[loc[1]:], region.parent)))
	}

	for _, r := range region.rules.regions {
		loc = FindIndex(r.start, line, start == 0, canMatchEnd)
		if loc != nil {
			highlights[start+loc[0]] = r.group
			return combineLineMatch(highlights,
				combineLineMatch(h.highlightRegion(start, false, lineNum, line[:loc[0]], region),
					h.highlightRegion(start+loc[1], canMatchEnd, lineNum, line[loc[1]:], r)))
		}
	}

	for _, p := range region.rules.patterns {
		matches := FindAllIndex(p.regex, line, start == 0, canMatchEnd)
		for _, m := range matches {
			highlights[start+m[0]] = p.group
			if _, ok := highlights[start+m[1]]; !ok {
				highlights[start+m[1]] = region.group
			}
		}
	}

	if canMatchEnd {
		h.endRegions[lineNum] = region
	}

	return highlights
}

func (h *Highlighter) highlightEmptyRegion(start int, canMatchEnd bool, lineNum int, line []byte) LineMatch {
	highlights := make(LineMatch)
	if len(line) == 0 {
		if canMatchEnd {
			h.endRegions[lineNum] = nil
		}
		return highlights
	}

	for _, r := range h.def.rules.regions {
		loc := FindIndex(r.start, line, start == 0, canMatchEnd)
		if loc != nil {
			highlights[start+loc[0]] = r.group
			return combineLineMatch(highlights,
				combineLineMatch(h.highlightEmptyRegion(start, false, lineNum, line[:loc[0]]),
					h.highlightRegion(start+loc[1], canMatchEnd, lineNum, line[loc[1]:], r)))
		}
	}

	for _, p := range h.def.rules.patterns {
		matches := FindAllIndex(p.regex, line, start == 0, canMatchEnd)
		for _, m := range matches {
			highlights[start+m[0]] = p.group
			if _, ok := highlights[start+m[1]]; !ok {
				highlights[start+m[1]] = ""
			}
		}
	}

	if canMatchEnd {
		h.endRegions[lineNum] = nil
	}

	return highlights
}

func (h *Highlighter) Highlight(input string) []LineMatch {
	lines := strings.Split(input, "\n")
	var lineMatches []LineMatch

	h.endRegions = make([]*Region, len(lines))

	for i := 0; i < len(lines); i++ {
		line := []byte(lines[i])

		if i == 0 || h.endRegions[i-1] == nil {
			lineMatches = append(lineMatches, h.highlightEmptyRegion(0, true, i, line))
		} else {
			lineMatches = append(lineMatches, h.highlightRegion(0, true, i, line, h.endRegions[i-1]))
		}
	}

	return lineMatches
}

func (h *Highlighter) ReHighlight(input LineStates, startline int) []LineMatch {
	lines := input.Lines()
	var lineMatches []LineMatch

	h.endRegions = make([]*Region, len(lines))

	for i := startline; i < len(lines); i++ {
		line := []byte(lines[i])

		if i == 0 || h.endRegions[i-1] == nil {
			lineMatches = append(lineMatches, h.highlightEmptyRegion(0, true, i, line))
		} else {
			lineMatches = append(lineMatches, h.highlightRegion(0, true, i, line, h.endRegions[i-1]))
		}

		curState := h.endRegions[i]
		lastState := input.State(i)

		if curState == lastState {
			break
		}

		input.SetState(i, curState)
	}

	return lineMatches
}
