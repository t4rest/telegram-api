package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/shelomentsevd/mtproto"
)

func main() {
	logfile, err := os.OpenFile("logfile.txt", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer logfile.Close()

	log.SetOutput(logfile)
	log.Println("Program started")

	// LoadContacts
	Mtproto, err := mtproto.NewMTProto(
		41994, "269069e15c81241f5670c397941016a2",
		mtproto.WithAuthFile(os.Getenv("HOME")+"/.telegramgo", false))
	if err != nil {
		log.Fatal(err)
	}
	telegramCLI, err := NewTelegramCLI(Mtproto)
	if err != nil {
		log.Fatal(err)
	}
	if err = telegramCLI.Connect(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Welcome to telegram CLI")
	if err := telegramCLI.CurrentUser(); err != nil {
		var phonenumber string
		fmt.Println("Enter phonenumber number below: ")
		fmt.Scanln(&phonenumber)
		err := telegramCLI.Authorization(phonenumber)
		if err != nil {
			log.Fatal(err)
		}
	}
	if err := telegramCLI.LoadContacts(); err != nil {
		log.Fatalf("Failed to load contacts: %s", err)
	}
	// Show help first time
	help()
	stop := make(chan struct{}, 1)
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
	SignalProcessing:
		for {
			select {
			case <-sigc:
				telegramCLI.Read()
			case <-stop:
				break SignalProcessing
			}
		}
	}()

	err = telegramCLI.Run()
	if err != nil {
		log.Println(err)
		fmt.Println("Telegram CLI exits with error: ", err)
	}
	// Stop SignalProcessing goroutine
	stop <- struct{}{}
}
