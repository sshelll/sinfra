package txtview

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Viewer struct {
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

	DoneFunc func(tcell.Key, *Viewer)
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

func NewViewer(opts ...*Options) *Viewer {
	tv := &Viewer{
		app:  tview.NewApplication(),
		view: tview.NewTextView(),
		opts: NewDefaultOpts(),
	}
	if len(opts) > 0 {
		tv.opts = opts[0]
	}
	return tv
}

func (v *Viewer) SetOpts(opts *Options) *Viewer {
	v.opts = opts
	return v
}

func (v *Viewer) App() *tview.Application {
	return v.app
}

func (v *Viewer) View() *tview.TextView {
	return v.view
}

func (v *Viewer) Write(p []byte) (n int, err error) {
	return v.view.Write(p)
}

func (v *Viewer) Run() error {
	v.apply()
	return v.app.
		SetRoot(v.view, v.opts.FullScreen).
		SetFocus(v.view).
		Run()
}

func (v *Viewer) Stop() {
	v.app.Stop()
}

func (v *Viewer) apply() {
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
				v.opts.DoneFunc(key, v)
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
