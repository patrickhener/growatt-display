package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/patrickhener/growatt-display/api"
	"github.com/patrickhener/growatt-display/utils"

	"github.com/howeyc/gopass"
	"github.com/inancgumus/screen"
)

func cleanup() {
	screen.Clear()
	fmt.Println("Caught CTRL+C, exiting...")
	fmt.Print("\x1b[?25h")
}

func main() {
	mode := flag.String("mode", "login", "[login/genhash]")
	username := flag.String("username", "", "username")
	password := flag.String("password", "", "hashed password")
	server := flag.String("server", "https://server.growatt.com/", "growatt api endpoint")
	loop := flag.Bool("loop", false, "Keep display open and update every minute")
	timeout := flag.Int("timeout", 30000, "Timeout in milliseconds for loop mode (default 30s)")
	flag.Parse()

	switch *mode {
	case "login":
		if *username == "" || *password == "" {
			fmt.Println("You need to provide -username and -password")
			os.Exit(-1)
		}
		api, err := api.New(*server, *username, *password)
		if err != nil {
			panic(err)
		}

		if err := api.Login(); err != nil {
			panic(err)
		}
		fmt.Println("Login successful")
		fmt.Println("")
		if !*loop {
			if err := api.Display(); err != nil {
				panic(err)
			}
		} else {
			error_count := 0
			fmt.Print("\x1b[?25l")
			screen.Clear()

			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)
			go func() {
				<-c
				cleanup()
				os.Exit(1)
			}()

			for {
				screen.MoveTopLeft()
				if err := api.Display(); err != nil {
					screen.Clear()
					error_count = error_count + 1
					fmt.Printf("There was an error: %+v\n", err)
					fmt.Printf("Wait one cycle and see if it resolves. Errors in row: %d\n", error_count)
				}
				if err == nil {
					screen.Clear()
					error_count = 0
				}
				time.Sleep(time.Duration(*timeout * int(time.Millisecond)))
			}
		}
	case "genhash":
		fmt.Print("Enter your password: ")
		input, _ := gopass.GetPasswdMasked()
		fmt.Printf("Provide this hash with -password <hash>: %s\n", utils.HashPassword(string(input)))
	default:
		panic("You need to select either mode login or mode genhash")
	}
}
