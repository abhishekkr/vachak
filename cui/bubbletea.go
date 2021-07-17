package cui

import (
	"fmt"
	"os"
	"strings"

	"github.com/abhishekkr/vachak/book"

	viewport "github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	runewidth "github.com/mattn/go-runewidth"
)

/*
below code is from charmbracelet's example
*/

const (
	useHighPerformanceRenderer = false

	headerHeight = 3
	footerHeight = 3

	defaultPagerLogPath = "/tmp/vachak-pager"
)

var (
	pagerLogPath = os.Getenv("PAGER_LOG")
)

func render(page book.Page, content string) {
	if pagerLogPath != "" {
		fyl, err := tea.LogToFile(pagerLogPath, defaultPagerLogPath)
		if err != nil {
			fmt.Printf("Could not open file %s: %v", pagerLogPath, err)
			os.Exit(1)
		}
		defer fyl.Close()
	}

	p := tea.NewProgram(
		model{content: content, page: page},
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if err := p.Start(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}
}

type model struct {
	content  string
	ready    bool
	viewport viewport.Model
	page     book.Page
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		verticalMargins := headerHeight + footerHeight

		if !m.ready {
			// Since this program is using the full size of the viewport we need
			// to wait until we've received the window dimensions before we
			// can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.Model{Width: msg.Width, Height: msg.Height - verticalMargins}
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			m.viewport.SetContent(m.content)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMargins
		}

		if useHighPerformanceRenderer {
			// Render (or re-render) the whole viewport. Necessary both to
			// initialize the viewport and when the window is resized.
			//
			// This is needed for high-performance rendering only.
			cmds = append(cmds, viewport.Sync(m.viewport))
		}
	}

	// Because we're using the viewport's default update function (with pager-
	// style navigation) it's important that the viewport's update function:
	//
	// * Receives messages from the Bubble Tea runtime
	// * Returns commands to the Bubble Tea runtime
	//
	m.viewport, cmd = m.viewport.Update(msg)
	if useHighPerformanceRenderer {
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	headerTop := "╭───────────╮"
	headerMid := "│ Vachak :) ├-[ " + m.page.BookName() + " ]├ (mouse OR arrow-keys OR pg-up/down) to SCROLL | qq OR <ESC> to close chapter"
	headerBot := "╰───────────╯"
	headerMid += strings.Repeat("─", m.viewport.Width-runewidth.StringWidth(headerMid))
	header := fmt.Sprintf("%s\n%s\n%s", headerTop, headerMid, headerBot)

	footerTop := "╭──────╮"
	footerMid := fmt.Sprintf("┤ %3.f%% │", m.viewport.ScrollPercent()*100) + " (" + m.page.Creators() + ") "
	footerBot := "╰──────╯"
	gapSize := m.viewport.Width - runewidth.StringWidth(footerMid)
	footerTop = strings.Repeat(" ", gapSize) + footerTop
	footerMid = strings.Repeat("─", gapSize) + footerMid
	footerBot = strings.Repeat(" ", gapSize) + footerBot
	footer := fmt.Sprintf("%s\n%s\n%s", footerTop, footerMid, footerBot)

	return fmt.Sprintf("%s\n%s\n%s", header, m.viewport.View(), footer)
}
