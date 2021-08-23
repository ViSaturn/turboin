package main

import (
	// Standard Libraries
	"fmt"
	"time"
	"strings"
	"os"
	"os/exec"

	// Keybinds
	// "github.com/jezek/xgb/xproto"
	"github.com/jezek/xgbutil"
	keybinds "github.com/jezek/xgbutil/keybind"
	"github.com/jezek/xgbutil/xevent"

	// Process Communication
	// zmq "github.com/pebbe/zmq4"
)

func run(currentMode string) {
	// Get the X Server
	X, err := xgbutil.NewConn()
	if err != nil { panic(err) }

	// Recive runtime requests (e.g. change mode)


	// Get the home directory
	home, err := os.UserHomeDir()
	if err != nil { panic(err) }

	// Get config
	modes := setupModes() // setupModes located in config.go

	// Initialize Keybindings
	keybinds.Initialize(X)

	// Nest
	convertStringToShellCmd := func(home string, str string) func() {
		// Allow usage of home directory
		replacer := strings.NewReplacer("~", home)
		str = replacer.Replace(str)

		// Make readable to exec.Command,
		// then invert to function & return
		splitStr := strings.Split(str, " ")
		execString := func() { exec.Command(splitStr[0], splitStr[1:]...).Run() }

		return execString
	}

	// Cycle through modes & set the default if its there,
	// after that, cycle through the mode's keybindings
	// and setup all the keybindings
	go func() {
		for modeName, mode := range modes {
			/*
			if modeName == "default" {
				currentMode = "default"
			}
			*/

			for keybind, toRun := range mode.(map[string]interface{}) {
				// Attach settings will be changed depending on the command as the code executes
				attachSettings := map[string]interface{} {
					"doubleClick": false,
				}

				// Get command the user is looking to run (from toRun)
				var run func()
				var run2 func() // run2 is ONLY used for double clicks

				switch toRun.(type) {
				// Double Clicks & One Clicks:
				case []interface{}:
					attachSettings["doubleClick"] = true
					obj1 := toRun.([]interface{})[0]
					obj2 := toRun.([]interface{})[1]

					// One Click:
					switch obj1.(type) {
					case string:
						run = convertStringToShellCmd(home, obj1.(string))
					case func():
						run = obj1.(func())
					}

					// Double Click:
					switch obj2.(type) {
					case string:
						run2 = convertStringToShellCmd(home, obj2.(string))
					case func():
						run2 = obj2.(func())
					}
					// Normal Clicks:
				case string:
					attachSettings["doubleClick"] = false
					run = convertStringToShellCmd(home, toRun.(string))
				case func():
					attachSettings["doubleClick"] = false
					run = toRun.(func())
				}

				// Attach
				alwaysModeName := modeName
				timesSent := 0

				keybinds.KeyPressFun(func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
					// Double Click Handler
					if attachSettings["doubleClick"] == true {
						if timesSent == 0 {
							timesSent = timesSent + 1

							go func() {
								time.Sleep(190 * time.Millisecond)

								if timesSent > 1 { // Double Click:
									timesSent = 0
									// this is just here to reset, the actual double click function
									// runs a few lines down (the advantage to putting it there
									// instead of here is to avoid the delay on running the double
									// click function, however there is still no way to completely
									// remove the delay for the one click function at the moment)
								} else { // One Click:
									// fmt.Println("one click")
									timesSent = 0
									if currentMode == alwaysModeName { go run() }
								}
							}()
						} else if timesSent == 1 {
							timesSent = timesSent + 1
							// fmt.Println("double click")
							if currentMode == alwaysModeName { go run2() }
						}
						fmt.Println(timesSent)

					// Normal Click Handler
					} else {
						if currentMode == alwaysModeName { go run() }
					}
				}).Connect(X, X.RootWin(), keybind, true)
			}
		}
	}()

	// Start main event loop
	xevent.Main(X)
	/*
	go xevent.Main(X)
	fmt.Scanln()
	*/
}

/*
func restart(currentMode string) {
	stop()
	run(currentMode)
}
*/

/*
func switchMode(server *zmq.Socket, mode string) {
	fmt.Println(mode)
	fmt.Println(server)

	server.Send(mode, 0)
}
*/

// CLI Stuff
func main() {
	// Connect to server to send runtime requests
	/*
	zctx, _ := zmq.NewContext()
	server, _ := zctx.NewSocket(zmq.REQ)
	server.Connect("tcp://localhost:5005")
	*/

	// Variables
	currentMode := "default"

	// Dumb way to get args
	args := os.Args
	arg1 := ""
	arg2 := ""

	switch len(args) {
	case 3:
		arg1 = args[1]
		arg2 = args[2]
	case 2:
		arg1 = args[1]
	default:
		return
	}

	// Connect args to functions
	if arg1 != "" {
		switch args[1] {
		case "run": run(currentMode)
		case "stop": exec.Command("killall", "turboin").Run()
		case "mode": if arg2 != "" { /*server.Send(args[2], 1)*/ }
		}
	}
}
