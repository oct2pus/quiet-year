package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/oct2pus/bocto"
)

var (
	decks   *([]deck)
	session bool
)

func init() {
	// TODO: this only works with a 1 server assumption,
	// this should be replaced with a database at some point if this is
	// expanded upon.
	session = false
}

func begin(bot bocto.Bot, mC *discordgo.MessageCreate, in []string) {
	if !session {
		decks = PrepareDecks()
		session = true
		bot.Session.ChannelMessageSend(mC.ChannelID, "**==Game Begins==**")
	}
}

func drawCard(bot bocto.Bot, mC *discordgo.MessageCreate, in []string) {
	if session {
		var (
			card      card
			err       error
			season    string
			remaining string
		)
		switch {
		case len((*decks)[0].Cards) >= 1:
			card, err = (*decks)[0].Draw()
			season = "spring"
			remaining = (string)(len((*decks)[0].Cards))
		case len((*decks)[1].Cards) >= 1:
			card, err = (*decks)[1].Draw()
			season = "summer"
			remaining = (string)(len((*decks)[1].Cards))
		case len((*decks)[2].Cards) >= 1:
			card, err = (*decks)[2].Draw()
			season = "fall"
			remaining = (string)(len((*decks)[2].Cards))
		case len((*decks)[3].Cards) >= 1:
			card, err = (*decks)[3].Draw()
			season = "winter"
			remaining = (string)(len((*decks)[3].Cards))
		}
		emojiFooter, emojiThumb := getEmoji(season)
		if err != nil {
			bot.Session.ChannelMessageSend(mC.ChannelID, "Somehow drew too many cards. :(\n, ending game prematurely.")
			// TODO: There should be a way to save a session
			// if this bot is expanded upon.
			session = false
		}
		bot.Session.ChannelMessageSendEmbed(mC.ChannelID, &discordgo.MessageEmbed{
			Title: lengthen(card.Face) + " of " + season,
			Fields: []*discordgo.MessageEmbedField{
				{Value: card.Text, Inline: false},
				{Name: "A", Value: card.Option1, Inline: true},
				{Name: "B", Value: card.Option2, Inline: true},
			},
			Footer:    &discordgo.MessageEmbedFooter{Text: remaining + " cards left.", IconURL: "https://raw.github.com/oct2pus/quiet-year/emoji/" + emojiFooter},
			Thumbnail: &discordgo.MessageEmbedThumbnail{URL: "https://raw.github.com/oct2pus/quiet-year/emoji/" + emojiThumb},
		})
	}
}

func attribute(bot bocto.Bot, mC *discordgo.MessageCreate, in []string) {
	bot.Session.ChannelMessageSend(mC.ChannelID, "This bot uses Mutant Standard emoji (https://mutant.tech) which are licensed under a Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International License (https://creativecommons.org/licenses/by-nc-sa/4.0/).")
}

func lengthen(face string) string {
	switch strings.ToLower(string(face[0])) {
	case "a":
		return "Ace"
	case "2":
		return "Two"
	case "3":
		return "Three"
	case "4":
		return "Four"
	case "5":
		return "Five"
	case "6":
		return "Six"
	case "7":
		return "Seven"
	case "8":
		return "Eight"
	case "9":
		return "Nine"
	case "1":
		return "Ten"
	case "j":
		return "Jack"
	case "q":
		return "Queen"
	case "k":
		return "King"
	default:
		return "error"
	}
}
func getEmoji(season string) (footer, thumbnail string) {
	switch season {
	case "spring":
		return "heart_suit.png", "rose.png"
	case "summer":
		return "diamond_suit.png", "sun.png"
	case "fall":
		return "club_suit.png", "maple_leaf.png"
	case "winter":
		return "spade_suit.png", "snowflake.png"
	}
	return "crt_test_pattern.png", "crt_blue_screen.png"
}
