package txtview

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type viewer struct {
	app  *tview.Application
	view *tview.TextView
	opts *Options
}

type Options struct {
	Title *string

	BgColor       tcell.Color
	Wrap          bool
	DynamicColors bool
	Regions       bool
	FullScreen    bool

	X, Y, Rows, Cols int // Position and size of the viewer when FullScreen is false

	Border      bool
	BorderAttr  *tcell.AttrMask
	BorderStyle *tcell.Style

	DoneFunc func(key tcell.Key)
}

func NewDefaultOpts() *Options {
	borderAtter := tcell.AttrBold
	opts := &Options{
		BgColor:       tcell.ColorDefault,
		DynamicColors: true,
		Regions:       true,
		Wrap:          true,
		FullScreen:    true,
		Border:        true,
		BorderAttr:    &borderAtter,
	}
	return opts
}

func NewViewer(opts ...*Options) *viewer {
	tv := &viewer{
		app:  tview.NewApplication(),
		view: tview.NewTextView(),
		opts: NewDefaultOpts(),
	}
	if len(opts) > 0 {
		tv.opts = opts[0]
	}
	return tv
}

func (v *viewer) SetOpts(opts *Options) *viewer {
	v.opts = opts
	return v
}

func (v *viewer) App() *tview.Application {
	return v.app
}

func (v *viewer) View() *tview.TextView {
	return v.view
}

func (v *viewer) Write(p []byte) (n int, err error) {
	return v.view.Write(p)
}

func (v *viewer) Run() error {
	v.apply()
	return v.app.
		SetRoot(v.view, v.opts.FullScreen).
		SetFocus(v.view).
		Run()
}

func (v *viewer) apply() {
	view := v.view
	view.SetBackgroundColor(v.opts.BgColor)
	view.SetRegions(v.opts.Regions).
		SetWrap(v.opts.Wrap).
		SetDynamicColors(v.opts.DynamicColors).
		SetChangedFunc(func() {
			v.app.Draw()
		}).
		SetDoneFunc(func(key tcell.Key) {
			// reserved key for exit
			if key == tcell.KeyEscape {
				v.app.Stop()
			}
			if v.opts.DoneFunc != nil {
				v.opts.DoneFunc(key)
			}
		})
	view.SetBorder(v.opts.Border)
	if v.opts.BorderAttr != nil {
		view.SetBorderAttributes(*v.opts.BorderAttr)
	}
	if v.opts.BorderStyle != nil {
		view.SetBorderStyle(*v.opts.BorderStyle)
	}
	if v.opts.Title != nil {
		view.SetTitle(*v.opts.Title)
	}
	if !v.opts.FullScreen {
		view.SetRect(v.opts.X, v.opts.Y, v.opts.Cols, v.opts.Rows)
	}
}
