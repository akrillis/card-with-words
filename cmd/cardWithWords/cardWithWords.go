package main

import (
	"cardWithWords/internal/pkg/storage"
	"cardWithWords/internal/pkg/telegram"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/jessevdk/go-flags"
)

type Opts struct {
	Token         string `long:"token" description:"telegram api token" required:"true"`
	WordsQuantity int    `long:"wordsQuantity" description:"how many words do we have on each card?" default:"8"`
}

var (
	opts         Opts
	chanOsSignal = make(chan os.Signal)
	chanDone     = make(chan struct{})
	wg           = new(sync.WaitGroup)
)

func dispatcher() {
	log.Println("[dispatcher] starts")

	signal.Notify(chanOsSignal, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-chanOsSignal

	log.Printf("[dispatcher] got %v", sig)
	chanDone <- struct{}{}

	log.Println("[dispatcher] done")
}

func main() {
	log.Println("start")

	// parse command line flags
	if _, err := flags.Parse(&opts); err != nil {
		log.Fatalf("Parse flags: %v\n", err)
	}

	words, err := storage.GetAccessToWords("./wordsDb")
	if err != nil {
		log.Fatalf("couldn't initialize new words storage: %v\n", err)
	}

	bot, err := telegram.GetAccessToTelegramApi(opts.Token, words, opts.WordsQuantity)
	if err != nil {
		log.Fatalf("couldn't initialize telegram bot api: %v\n", err)
	}

	go dispatcher()

	wg.Add(1)

	go func(wg *sync.WaitGroup, done chan struct{}) {
		if err := bot.ListenAndServeForWords(wg, done); err != nil {
			log.Fatalf("couldn't listen and serve words via telegram bot api: %v\n", err)
		}
	}(wg, chanDone)

	wg.Wait()
	log.Println("done")
}
