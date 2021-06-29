package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/oct2pus/bocto"
)

const (
	pre = "!!"
)

func main() {
	var token string
	var bot bocto.Bot
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()
	// create bot
	err := bot.New("QYBot", pre, token, 0xee5d6c)
	if err != nil {
		fmt.Printf("%v can't login\nerror: %v\n", bot.Name, err)
		return
	}

	// add commands and responses
	bot = addCommands(bot)
	bot.Confused = confused
	bot.Mentioned = mentioned
	// Event Handlers
	bot.Session.AddHandler(bot.ReadyEvent)
	bot.Session.AddHandler(bot.MessageCreate)
	// Open Bot
	err = bot.Session.Open()
	if err != nil {
		fmt.Printf("Error openning connection: %v\nDump bot info %v\n",
			err,
			bot.String())
	}
	// wait for ctrl+c to close.

	signalClose := make(chan os.Signal, 1)

	signal.Notify(signalClose,
		syscall.SIGINT,
		syscall.SIGTERM,
		os.Interrupt,
		os.Kill)
	<-signalClose

	bot.Session.Close()
}

func addCommands(bot bocto.Bot) bocto.Bot {
	bot.AddCommand("start", begin, true)
	bot.AddCommand("draw", drawCard, true)
	bot.AddCommand("about", attribute, true)
	return bot
}

func confused(bot bocto.Bot, mC *discordgo.MessageCreate, in []string) {
	bot.Session.ChannelMessageSend(mC.ChannelID, "?")
}

func mentioned(bot bocto.Bot, mC *discordgo.MessageCreate, in []string) {
	bot.Session.ChannelMessageSend(mC.ChannelID, "My prefix is `"+pre+"`.")
}
