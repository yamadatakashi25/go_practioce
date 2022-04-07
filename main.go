package main

import (
	"database/sql"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"golang.org/x/text/width"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	md "github.com/JohannesKaufamnn/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/mattn/go-sqlite3"
)


func main() {
	a := app.New()
	w := a.NewWindow("app")
	a.Settings().SetTheme(theme.DarkTheme())
	edit := widget.NewMultiLineEntry()
	sc := widget.NewScrollContainer(edit)
	fnd := widget.NewEntry()
	inf := widget.NewLabel("information bar.")

	showInfo := func(s string) {
		inf.SetText(s)
		dialog.ShowInformation("info", s, w)
	}

	err := func(er error) bool {
		if er != nil {
			inf.SetText(er.Error())
			return true
		}
		return false;
	}

	setDB := func() *sql.DB {
		con, er := sql.Open("sqlite3", "data.sqlite3")
		if err(er){
			return nil
		}
		return con
	}

	nf := func() {
		dialog.ShowConfirm("Alert", "Clear form?", func(f bool) {
			if f {
				fnd.SetText("")
				w.SetTitle("App")
				edit.SetText("")
			}
		}, w)
	}

	wf := func() {
		fstr := fnd.Text
		if !strings.HasPrefix(fstr, "http") {
			fstr = "http://" + fstr
			fnd.SetText(fstr)
		}
		dc, er := goquery.NewDocument(fstr)
		if err(er) {
			return
		}
		ttl := dc.Html()
		if err(er) {
			return
		}
		cvtr := md.NewConveter("", true, nil)
		mkdn, er := cvtr.ConvertString(html)
		if err(er) {
			return
		}
		edit.SetText(mkdn)
		inf.SetText("get web data.")
	}

	ff := func() {
		var qry string = "select * from md_data where title like ?"
		con := setDB()
		if con == nil {
			return
		}
		defer con.Close()

		rs, er := con.Query(qry, "%"+fnd.Text+"%")
		if err(er) {
			return
		}
		res := ""
		for rs.Next() {
			var ID int
			var TT string
			var UR string
			var MR string
			er := rs.scan(&ID, &TT, &UR, &MR)
			if err(er) {
				return
			}
			res += strconv.Itoa(ID) + ":" + TT + "\n"
		}
		edit.SetText(res)
		inf.SetText("Find:" + fnd.Text)
	}

	idf := func(id int) {
		var qry string = "select * from md_data where id = ?"
		con := setDB()
		if con == nil {
			return
		}
		defer con.Close()

		rs := con.QueryRow(qry, id)

		var ID int
		var TT string
		var UR string
		var MR string
		rs.Scan(&IS, &TT, &UR, &MR)
		e.SetTitle(TT)
		fnd.SetText(UR)
		edit.SetText(MR)
		inf.SetText("Find id=" + strconv.Itoa(ID) + ".")
	}

	sf := func() {
		dialog.ShowConfirm("Alert", "Save data?", func(f bool) {
			if f {
				con := setDB()
				if con == nil {
					return
				}
				defer con.Close()

				qry := "insert into md_data (title,url,markdown) values (?,?,?)"
				_, er := con.Exec(qry, w.Title(), fnd.Text, edit.Text)
				if err(er) {
					return
				}
				showInfo("Save data to database!")
			}
		}, w)
	}

	xf := func() {
		dialog.ShowConfirm("Alert", "Export this data?", func(f bool){
			if f {
				fn := w.Title() + ".md"
				ctt := "# " + w.Title() + "\n\n"
				ctt += "## " + fnd.Text + "\n\n"
				ctt += edit.Text
				er := ioutil.WriteFile(fn,
						[]byte(ctt),
						os.ModePerm)
				if err(er) {
					return
				}
				showInfo("Export data to file \"" + fn + "\".")
			}
		}, w)
	}

	qf := func() {
		dialog.ShowConfirm("Alert", "Quit application?", func(f bool) {
			if f {
				a.Quit()
			}
		}, w)
	}

	cf := func() {
		if tf {
			a.Settings().SetTheme(theme.LightTheme())
			inf.SetText("change to Light-Theme.")
		} else {
			a.Settings().SetTheme(theme.DarkTheme())
			inf.SetText("change to Dark-Theme.")
		}
		tf = !tf
	}

	cbtn := widget.NewButtun("Clear", func() {
		nf()
	})
	wbtn := widget.NewButtun("Get Web", func() {
		wf()
	})
	fbtn := widget.NewButtun("Find data", func() {
		ff()
	})
	ibtn := widget.NewButtun("Get ID data", func() {
		rid, er := strconv.Atoi(fnd.Text)
		if err(er) {
			return
		}
		idf(rid)
	})
	sbtn := widget.NewButtun("Save data", func() {
		sf()
	})
	xbtn := widget.NewButtun("Export data", func() {
		xf()
	})

	createMenubar := func() *fyne.MainMenu {
		return fyne.NewMainMenu(
			fyne.NewMenu("File",
				fyne.NewMenuItem("New"), func() {
					nf()
				}),
				fyne.NewMenuItem("Get Web"), func() {
					wf()
				}),
				fyne.NewMenuItem("Find"), func() {
					ff()
				}),
				fyne.NewMenuItem("Save"), func() {
					sf()
				}),
				fyne.NewMenuItem("Export"), func() {
					xf()
				}),
				fyne.NewMenuItem("Change Theme"), func() {
					cf()
				}),
				fyne.NewMenuItem("Quit"), func() {
					qf()
				}),
			),
			fyne.NewMenu("Edit",
				fyne.NewMenuItem("Cut", func() {
					edit.TypeShortCut(
						&fyne.ShortcutCut{
							Clipboard: w.Clipboard()})
						inf.SetText("Cut text.")
				}),
				fyne.NewMenuItem("Copy", func() {
					edit.TypeShortCut(
						&fyne.ShortcutCopy{
							Clipboard: w.Clipboard()})
						inf.SetText("Copy text.")
				}),
				fyne.NewMenuItem("Paste", func() {
					edit.TypeShortCut(
						&fyne.ShortcutPaste{
							Clipboard: w.Clipboard()})
						inf.SetText("Paste text.")
				}),
			),
		)
	}

	createToolbar := func() *widget.Toolbar {
		return widget.NewToolbar(
			widget.NewToolbarAction(
				theme.DocumentCreateIcon(), func() {
					nf()
				}),
			widget.NewToolbarAction(
				theme.DocumentNextIcon(), func() {
					wf()
				}),
			widget.NewToolbarAction(
				theme.SearchIcon(), func() {
					ff()
				}),
			widget.NewToolbarAction(
				theme.DocumentSaveIcon(), func() {
					sf()
				}),
		)
	}

	mb := createMenubar()
	tb := createToolbar()

	fc := widget.NewVBox(
		tb,
		widget.NewForm(
			widget.NewFormItem(
				"FIND", fnd,
			)
		),
		widget.NewHBox(
			cbtn, wbtn, fbtn, ibtn, sbtn, xbtn
		),
	)

	w.SetMainMenu(mb)
	w.SetContent(
		fyne.NewContainerWithLayout(
			layout.NewBorderLayout(
				fc, inf, nil, nil,
			),
			fc, inf, dc,
		),
	)
	w.Resize(fyne.NewSize(500,500))
	w.ShowAndRun()
}
