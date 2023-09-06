package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sshelll/sinfra/tview/txtview"
	"github.com/sshelll/sinfra/util"
)

const corporate = `Leverage the a the b to c to d.
[yellow]Press Enter, then Tab/Backtab for word selections`

const textZH = `你好世界，我是txtviewer，我是一个文本查看器`

func main() {
	viewer()
}

func viewer() {
	opts := txtview.NewDefaultOpts()
	opts.Title = util.Ptr("Test TxtViewer")
	opts.DoneFunc = func(key tcell.Key, v *txtview.Viewer) {
		if key == tcell.KeyEnter {
			v.Stop()
		}
	}
	viewer := txtview.NewViewer(opts)
	go func() {
		for _, word := range strings.Split(corporate, " ") {
			if word == "the" {
				word = "[red]the[white]"
			}
			fmt.Fprintf(viewer, "%s ", word)
			time.Sleep(50 * time.Millisecond)
		}
	}()
	if err := viewer.Run(); err != nil {
		panic(err)
	}
}

func raw() {
	app := tview.NewApplication()
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	textView.SetBackgroundColor(tcell.ColorDefault)
	textView.SetTitle("TextView")
	textView.SetWrap(false)
	textView.SetBorder(true).SetBorderAttributes(tcell.AttrBold)
	numSelections := 0
	go func() {
		for _, ch := range corporate {
			word := string(ch)
			if word == "t" {
				word = "[red]t[white]"
			}
			fmt.Fprintf(textView, "%s", word)
			time.Sleep(200 * time.Millisecond)
		}
	}()
	textView.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Stop()
		}
		currentSelection := textView.GetHighlights()
		if key == tcell.KeyEnter {
			if len(currentSelection) > 0 {
				textView.Highlight()
			} else {
				textView.Highlight("0").ScrollToHighlight()
			}
		} else if len(currentSelection) > 0 {
			index, _ := strconv.Atoi(currentSelection[0])
			if key == tcell.KeyTab && numSelections > 0 {
				index = (index + 1) % numSelections
			} else if key == tcell.KeyBacktab {
				index = (index - 1 + numSelections) % numSelections
			} else {
				return
			}
			textView.Highlight(strconv.Itoa(index)).ScrollToHighlight()
		}
	})
	textView.SetBorder(true)
	if err := app.SetRoot(textView, true).SetFocus(textView).Run(); err != nil {
		panic(err)
	}
}
