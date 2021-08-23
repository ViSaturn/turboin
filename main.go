package main

import (
	// Standard Libraries
	// "fmt"
	"time"
	"strings"
	"os"
	"os/exec"

	// Keybinds
	"github.com/jezek/xgbutil"
	keybinds "github.com/jezek/xgbutil/keybind"
	"github.com/jezek/xgbutil/xevent"
)

func run(currentMode string) {
	// Get the X Server
	X, err := xgbutil.NewConn()
	if err != nil { panic(err) }

	// Get the home directory
	home, err := os.UserHomeDir()
	if err != nil { panic(err) }

	// Get config
	modes := setupModes() // setupModes located in config.go

	// Create special list for special keys to use
	// TODO: I probably shouldnt use a map for this since I only need one value,
	// not a key and a value
	special := map[string]interface{} {}

	// Initialize Keybindings
	keybinds.Initialize(X)

	// Nest
	convertStringToShellCmd := func(home string, str string) func() {
		// Allow usage of home directory
		replacer := strings.NewReplacer("~", home)
		str = replacer.Replace(str)

		// Make readable to exec.Command,
		// then convert to function & return
		splitStr := strings.Split(str, " ")
		execString := func() { exec.Command(splitStr[0], splitStr[1:]...).Run() }

		return execString
	}

	makeSpecialKeySet := func(specialKey string, normalKey string, doubleClick bool, doubleClickDelay int) {
		// example:
		// "{mod4:z}-mod4-x": func() { fmt.Println("mod4-z = special key, mod4-x = normal key") }

		// Note that due to limitations the special key's seperator has to be different than
		// the normal key's seperator

		// How this works:
		// pressConn and releaseConn place/remove the special key from the specials map
		// when you press or release the special key on your keyboard
		//
		// while the special key is in the specials map it can be detected when linkConn
		// is ran, linkConn is just there to remove the specialKey from the map in case
		// releaseConn can't find it, the actual code checking the specials map is in the
		// conditionsMet function

		// Add function to specialFuncs map upon press
		pressConn := keybinds.KeyPressFun(func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
			special[normalKey] = ""
		})

		// Remove function from specialFuncs map upon release
		releaseConn := keybinds.KeyReleaseFun(func(X *xgbutil.XUtil, e xevent.KeyReleaseEvent) {
			if _, ok := special[normalKey]; ok {
				delete(special, normalKey)
				_ = 1 // avoid "declared but not used"
			}
		})

		// Link to the normal key, when this is pressed the special key has to be pressed
		// again for this to work again, because of limitations in xgbutil
		linkConn := keybinds.KeyPressFun(func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
			if _, ok := special[normalKey]; ok {
				go func() {
					if doubleClick == true {
						// 5 ms window for the run function to detect the special key
						// This is probably not the best way to do this but i'll probs find a
						// better way later
						time.Sleep(time.Duration(doubleClickDelay + 5) * time.Millisecond)
					}
					delete(special, normalKey)
					_ = 1 // avoid "declared but not used"
				}()
			}
		})

		// Connect to bindings
		errPress := pressConn.Connect(X, X.RootWin(), specialKey, true)
		errRelease := releaseConn.Connect(X, X.RootWin(), specialKey, true)
		errLink := linkConn.Connect(X, X.RootWin(), normalKey, true)

		// Handle Errors
		if errPress != nil { panic(errPress) }
		if errRelease != nil { panic(errRelease) }
		if errLink != nil { panic(errLink) }
	}

	// Cycle through modes & set the default if its there,
	// after that, cycle through the mode's keybindings
	// and setup all the keybindings
	go func() {
		for modeName, mode := range modes {
			// when currentMode is passed to run it's already
			// set to default
			/*
			if modeName == "default" {
				currentMode = "default"
			}
			*/

			for keybind, toRun := range mode.(map[string]interface{}) {
				// Attach settings will be changed depending on the command as the code executes
				attachSettings := map[string]interface{} {
					"doubleClick": false,
						"doubleClickDelay": 190,
					"specialKeySituation": false,
						"specialKey": "",
						"normalKey": "",
				}

				// Function to check if conditions are met to run the keybind's command
				conditionsMet := func(run func(), currentMode string, modeName string) {
					// checks whethr there is a special key situation, if there is
					// it'll also check if it can run, if there is no special key situation
					// it will check if something else has a special key situation with this
					// key as a part of the set
					if attachSettings["specialKeySituation"].(bool) == true {
						if _, ok := special[attachSettings["normalKey"].(string)]; ok {
							go run()
							_ = 1 // avoid "declared but not used"
						}
					} else {
						go run()
						// FIXME
						/*
						if canRun, ok := specialFuncs[attachSettings[keybind].(string)]; ok {
							fmt.Println(canRun)
						} else {
							go run()
						}
						*/
					}
				}

				// Get command the user is looking to run (from toRun)
				var run func()
				var run2 func() // run2 is ONLY used for double clicks

				switch toRun.(type) {
				// Double Clicks & One Clicks:
				case []interface{}:
					attachSettings["doubleClick"] = true
					mapItem1 := toRun.([]interface{})[0]
					mapItem2 := toRun.([]interface{})[1]

					// One Click:
					switch mapItem1.(type) {
					case string:
						run = convertStringToShellCmd(home, mapItem1.(string))
					case func():
						run = mapItem1.(func())
					}

					// Double Click:
					switch mapItem2.(type) {
					case string:
						run2 = convertStringToShellCmd(home, mapItem2.(string))
					case func():
						run2 = mapItem2.(func())
					}
					// Normal Clicks:
				case string:
					attachSettings["doubleClick"] = false
					run = convertStringToShellCmd(home, toRun.(string))
				case func():
					attachSettings["doubleClick"] = false
					run = toRun.(func())
				}

				// Split keybind to take a look at each key (to look for special keys)
				splitKeybind := strings.Split(keybind, "-") // seperator right now is -

				for _, key := range splitKeybind {
					if strings.HasPrefix(key, "{") && strings.HasSuffix(key, "}") {
						// This means there is a special key situation that must be added
						// to attachSettings
						specialKey := strings.Replace(key, "{", "", 1) // remove prefix
						specialKey = strings.Replace(specialKey, "}", "", 1) // remove suffix
						specialKey = strings.Replace(specialKey, ":", "-", 1) // replace seperator

						normalKey := strings.Replace(keybind, key + "-", "", 1) // remove special key from keybind

						attachSettings["specialKeySituation"] = true
						attachSettings["specialKey"] = specialKey
						attachSettings["normalKey"] = normalKey
						keybind = normalKey

						makeSpecialKeySet(specialKey, normalKey, attachSettings["doubleClick"].(bool), attachSettings["doubleClickDelay"].(int))

						_ = 1 // avoid "declared but not used"
					}
				}

				// Attach
				alwaysModeName := modeName
				timesSent := 0

				// Any other checks should be moved to the conditionsMet function
				keybinds.KeyPressFun(func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
					// Check mode
					if currentMode == alwaysModeName {
						// Double Click Handler
						if attachSettings["doubleClick"] == true {
							if timesSent == 0 {
								timesSent = timesSent + 1

								go func() {
									time.Sleep(time.Duration(attachSettings["doubleClickDelay"].(int)) * time.Millisecond)

									if timesSent > 1 { // Double Click:
										timesSent = 0
										// this is just here to reset, the actual double click function
										// runs a few lines down (the advantage to putting it there
										// instead of here is to avoid the delay on running the double
										// click function, however there is still no way to completely
										// remove the delay for the one click function at the moment)
									} else { // One Click:
										timesSent = 0
										conditionsMet(run, currentMode, alwaysModeName)
									}
								}()
							} else if timesSent == 1 {
								timesSent = timesSent + 1
								conditionsMet(run2, currentMode, alwaysModeName)
							}

							// Normal Click Handler
						} else {
							conditionsMet(run, currentMode, alwaysModeName)
						}
					}
				}).Connect(X, X.RootWin(), keybind, true)
			}
		}
	}()

	// Start main event loop
	xevent.Main(X)
}

// CLI Stuff
func main() {
	// Variables
	currentMode := "default"

	// Dumb way to get args but honestly i cant think of anything else
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
		case "mode": if arg2 != "" {  }
		}
	}
}
