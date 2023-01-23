package component

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const (
	AppLabel = "EasierConnect"
	APPID    = "github.com/lyc8503/EasierConnect"
)

var (
	Connected = false
)

// Hello Demo
func Hello() {
	a := app.New()
	w := a.NewWindow("Hello World")

	w.SetContent(widget.NewLabel("Hello World!"))
	w.ShowAndRun()
}

func EasierConnectUI() {
	a := app.NewWithID(APPID)
	pref := a.Preferences()
	w := a.NewWindow(AppLabel)
	w.Resize(fyne.Size{Height: 300, Width: 400})

	validator1 := validation.NewRegexp(`[\w\d.]+`, "Invalid Url")
	validator2 := validation.NewRegexp(`\w+`, "Can't be empty")
	validator3 := validation.NewRegexp(`[:]?\d+`, "Be a number")
	url := widget.NewEntry()
	url.SetPlaceHolder("vpn.domain.com") //Hint
	url.Validator = validator1
	url.SetText(pref.StringWithFallback("Url", "")) //Save user input after exit
	port := widget.NewEntry()
	port.SetPlaceHolder("443")
	port.SetText(pref.StringWithFallback("Port", "443"))
	port.Validator = validator3

	username := widget.NewEntry()
	username.SetPlaceHolder("UserName")
	username.Validator = validator2
	username.SetText(pref.StringWithFallback("UserName", ""))
	passwd := widget.NewPasswordEntry()
	passwd.SetPlaceHolder("Password")
	passwd.Validator = validator2
	passwd.SetText(pref.StringWithFallback("Password", ""))

	socks5 := widget.NewEntry()
	socks5.SetPlaceHolder(":1080")
	socks5.Validator = validator3
	socks5.SetText(pref.StringWithFallback("Socks5", ":1080"))

	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
		},
	}

	form.CancelText = "Exit"
	form.SubmitText = "Connect"
	form.OnSubmit = func() {
		pref.SetString("Url", url.Text)
		pref.SetString("Port", port.Text)
		pref.SetString("UserName", username.Text)
		pref.SetString("Password", passwd.Text)
		pref.SetString("Socks5", socks5.Text)

		if !Connected {
			form.SubmitText = "DisConnect"
			form.Refresh()
			Connected = true
			go Process(url.Text, port.Text, username.Text, passwd.Text, socks5.Text)
		} else {
			form.SubmitText = "Connect"
			form.Refresh()
			Connected = true
		}

	}
	form.OnCancel = func() {
		w.Close()
		a.Quit()
	}

	// we can also append items
	form.Append("Url", url)
	form.Append("Port", port)
	form.Append("UserName", username)
	form.Append("PassWord", passwd)
	form.Append("Socks5 Listen", socks5)
	container1 := container.New(layout.NewVBoxLayout(), form)
	w.SetContent(container1)

	// add systray icon on desktop system;
	if desk, ok := a.(desktop.App); ok {
		m := fyne.NewMenu(AppLabel,
			fyne.NewMenuItem("Show", func() {
				w.Show()
			}))
		desk.SetSystemTrayMenu(m)
		w.SetCloseIntercept(func() {
			w.Hide()
		})
	}
	w.Show()

	a.Run()
}
