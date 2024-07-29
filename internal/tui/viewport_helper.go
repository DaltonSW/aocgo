package tui

import "strings"

// Public funcs
func (m viewportModel) IsAtTop() bool {
	return m.YOffset <= 0
}

func (m viewportModel) IsAtBottom() bool {
	return m.YOffset >= m.getMaxYOffset()
}

func (m viewportModel) IsPastBottom() bool {
	return m.YOffset > m.getMaxYOffset()
}

func (m viewportModel) GetScrollPercent() float64 {
	if m.Height >= len(m.lines) {
		return 1.0
	}
	curY := float64(m.YOffset)
	height := float64(m.Height)
	total := float64(len(m.lines))
	v := curY / (total - height)
	return max(0, min(1.0, v))
}

func (m *viewportModel) SetContent(s string) {
	s = strings.ReplaceAll(s, "\r\n", "\n") // Normalize line endings
	m.lines = strings.Split(s, "\n")

	if m.YOffset > len(m.lines)-1 {
		m.GoToBottom()
	}
}

func (m *viewportModel) GoToTop() {

}

func (m *viewportModel) GoToBottom() {

}

func (m *viewportModel) SetYOffset(offset int) {
	m.YOffset = clamp(offset, 0, m.getMaxYOffset())
}

// Private funcs
func (m viewportModel) getMaxYOffset() int {
	return max(0, len(m.lines)-m.Height)
}

func (m viewportModel) getVisibleLines() []string {
	var lines []string
	if len(m.lines) > 0 {
		top := max(0, m.YOffset)
		bottom := clamp(m.YOffset+m.Height, top, len(m.lines))
		lines = m.lines[top:bottom]
	}
	return lines
}

// Utils

func clamp(v, low, high int) int {
	if high < low {
		low, high = high, low
	}
	return min(high, max(low, v))
}
