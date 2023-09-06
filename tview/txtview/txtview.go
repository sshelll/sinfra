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

func (v *viewer) Run() error {
	v.apply()
	return v.app.
		SetRoot(v.view, v.opts.FullScreen).
		SetFocus(v.view).
		Run()
}

func (v *viewer) apply() {
	v.view.SetBackgroundColor(v.opts.BgColor)
	v.view.SetRegions(v.opts.Regions).
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
	if v.opts.Title != nil {
		v.view.SetTitle(*v.opts.Title)
	}
	if !v.opts.FullScreen {
		v.view.SetRect(v.opts.X, v.opts.Y, v.opts.Cols, v.opts.Rows)
	}
}
