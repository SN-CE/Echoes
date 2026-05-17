package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"
)

var stop chan os.Signal

func shutdown(app *tview.Application) {
	<-stop
	app.QueueUpdateDraw(func() {
		app.Stop()
	})
}

func receiver(conn *websocket.Conn, chat *tview.TextView, app *tview.Application) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			app.QueueUpdateDraw(func() {
				chat.Write([]byte("[red::b]Disconnected Or Server Went Offline[-::-]"))
				app.Stop()
			})
			return
		}
		if string(msg) == "[red::b]\n\nWELCOME TO ECHOES[-::-]" {
			app.QueueUpdateDraw(func() {
				chat.Clear()
			})
			app.QueueUpdateDraw(func() {
				chat.Write([]byte(string(msg)))
			})
			time.Sleep(2 * time.Second)
			app.QueueUpdateDraw(func() {
				chat.Clear()
			})
			continue
		}

		app.QueueUpdateDraw(func() {
			chat.Write([]byte(string(msg) + "[-::-]\n\n"))
		})
	}
}

func sender(conn *websocket.Conn, composer *tview.InputField, chat *tview.TextView, app *tview.Application) {
	composer.SetDoneFunc(func (key tcell.Key) {
		if key == tcell.KeyEnter {
			msg := composer.GetText()
			if msg != "" {
				err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
				composer.SetText("")
				if err != nil {
					app.QueueUpdateDraw(func() {
						chat.Write([]byte("[red::b]Disconnected Or Server Went Offline[-::-]"))
						app.Stop()
					})
					return
				}
			}
		}
	})
}

func main() {
	stop = make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	var addr string;
	fmt.Printf("Enter the IP or URL of the Echoes Server you wish to connect to - ")
	fmt.Scanln(&addr);
	addr = "ws://" + addr;

	conn, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		fmt.Printf("Server Offline Or Busy\n")
		return
	}
	defer conn.Close()

	app := tview.NewApplication()

	chat := tview.NewTextView()
	chat.SetDynamicColors(true)
	chat.SetScrollable(true)
	chat.SetChangedFunc(func() {
		chat.ScrollToEnd()
		app.Draw()
	})
	chat.SetBorder(true)
	chat.SetTitle("[::b]   Echoes   ")

	composer := tview.NewInputField()
	composer.SetLabel("[black:gray]You : [-]")
	composer.SetBorder(true)
	composer.SetFieldWidth(0)
	composer.SetFieldTextColor(tcell.ColorBlack)
	composer.SetBackgroundColor(tcell.ColorGray)
	composer.SetFieldBackgroundColor(tcell.ColorGray)

	layout := tview.NewFlex()
	layout.SetDirection(tview.FlexRow)
	layout.AddItem(chat, 0, 1, false)
	layout.AddItem(composer, 3, 0, true)

	app.SetRoot(layout, true)
	app.EnableMouse(true)

	sender(conn, composer, chat, app)
	go receiver(conn, chat, app)

	go shutdown(app)

	err = app.Run()
	if err != nil {
	panic(err)
	}
}
