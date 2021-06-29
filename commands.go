package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/oct2pus/bocto"
)

const IMAGE_URL = "https://raw.githubusercontent.com/oct2pus/quiet-year/main/emoji/"

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
		for i := range *decks {
			(*decks)[i].Shuffle()
		}
		session = true
		bot.Session.ChannelMessageSend(mC.ChannelID, "**==GAME START==**")
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
			season = "Spring"
			remaining = fmt.Sprintf("%v", (len((*decks)[0].Cards)))
		case len((*decks)[1].Cards) >= 1:
			card, err = (*decks)[1].Draw()
			season = "Summer"
			remaining = fmt.Sprintf("%v", (len((*decks)[1].Cards)))
		case len((*decks)[2].Cards) >= 1:
			card, err = (*decks)[2].Draw()
			season = "Fall"
			remaining = fmt.Sprintf("%v", (len((*decks)[2].Cards)))
		case len((*decks)[3].Cards) >= 1:
			card, err = (*decks)[3].Draw()
			season = "Winter"
			remaining = fmt.Sprintf("%v", (len((*decks)[3].Cards)))
			if card.Face == "K♠" {
				session = false
			}
		}
		emojiFooter, emojiThumb := getEmoji(season)
		if err != nil {
			bot.Session.ChannelMessageSend(mC.ChannelID, "Somehow drew too many cards. :(\n, ending game prematurely.")
			// TODO: There should be a way to save a session
			// if this bot is expanded upon.
			session = false
			return
		}
		embed := &discordgo.MessageEmbed{
			Title:     lengthen(card.Face) + " of " + season,
			Footer:    &discordgo.MessageEmbedFooter{Text: remaining + " cards left.", IconURL: IMAGE_URL + emojiFooter},
			Thumbnail: &discordgo.MessageEmbedThumbnail{URL: IMAGE_URL + emojiThumb},
		}
		if card.Text != "" {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "Description",
				Value:  card.Text,
				Inline: false,
			})
		}
		if card.Option1 != "" {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "A",
				Value:  card.Option1,
				Inline: true,
			})
		}
		if card.Option2 != "" {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "B",
				Value:  card.Option2,
				Inline: true,
			})
		}
		bot.Session.ChannelMessageSendEmbed(mC.ChannelID, embed)
		if !session {
			bot.Session.ChannelMessageSend(mC.ChannelID, "**==GAME OVER==**")
		}
	}
}

func discard(bot bocto.Bot, mC *discordgo.MessageCreate, in []string) {
	var err error
	if session {
		switch {
		case len((*decks)[0].Cards) >= 1:
			_, err = (*decks)[0].Draw()
		case len((*decks)[1].Cards) >= 1:
			_, err = (*decks)[1].Draw()
		case len((*decks)[2].Cards) >= 1:
			_, err = (*decks)[2].Draw()
		case len((*decks)[3].Cards) >= 1:
			_, err = (*decks)[3].Draw()
		}
		if err != nil {
			bot.Session.ChannelMessageSend(mC.ChannelID, "Somehow drew too many cards. :(\n, ending game prematurely.")
			// TODO: There should be a way to save a session
			// if this bot is expanded upon.
			session = false
			return
		}
	}
}

func end(bot bocto.Bot, mC *discordgo.MessageCreate, in []string) {
	if session {
		session = false
		bot.Session.ChannelMessageSend(mC.ChannelID, "**==GAME OVER==**")
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
	case "Spring":
		return "heart_suit.png", "rose.png"
	case "Summer":
		return "diamond_suit.png", "sun.png"
	case "Fall":
		return "club_suit.png", "maple_leaf.png"
	case "Winter":
		return "spade_suit.png", "snowflake.png"
	}
	return "crt_test_pattern.png", "crt_blue_screen.png"
}
