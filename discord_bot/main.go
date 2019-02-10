package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type guildMemberWithAccentsRemoved struct {
	GuildMember *discordgo.Member
	Nickname    string
}

var (
	commandPrefix string
	botID         string
)

var prefixCheckNickname = "~check_nicknames"
var prefixReloadNicknames = "~reload_nicknames"
var prefixHelp = "~help"
var voiceChannelName = map[string]string{
	"187229035758223360": "General",   //Paranoids Gaming
	"470921656865521665": "Raid Chat", //Life in the Math Lane
}
var allowedChannelsToSendCommands = map[string]string{
	"187229035758223360": "187229035758223360", //Paranoids Gaming
	"470921656865521665": "514964509954146308", //Life in the Math Lane
}

var t = transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r)
}

func main() {
	discord, err := discordgo.New(fmt.Sprintf("Bot %s", os.Getenv("DISCORD_API_TOKEN")))

	errCheck("error creating discord session", err)
	user, err := discord.User("@me")
	errCheck("error retrieving account", err)

	botID = user.ID
	discord.AddHandler(commandHandler)
	discord.AddHandler(func(discord *discordgo.Session, ready *discordgo.Ready) {
		err = discord.UpdateStatus(0, "Nickname manager!")
		if err != nil {
			fmt.Println("Error attempting to set my status")
		}
		servers := discord.State.Guilds
		fmt.Printf("SuperAwesomeOmegaTutorBot has started on %d servers\n", len(servers))
	})

	err = discord.Open()
	errCheck("Error opening connection to Discord", err)
	defer discord.Close()

	commandPrefix = "!"

	<-make(chan struct{})

}

func errCheck(msg string, err error) {
	if err != nil {
		fmt.Printf("%s: %+v", msg, err)
		panic(err)
	}
}

func commandHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
	printMessage := false

	if message.Author.ID == botID || message.Author.Bot {
		return
	}

	if strings.Compare(allowedChannelsToSendCommands[message.GuildID], message.ChannelID) != 0 {
		return
	}

	if strings.HasPrefix(message.Content, prefixCheckNickname) {
		fmt.Printf("Message: %+v || From: %s\n", message.Message, message.Author)
		checkNicknames(discord, message)
		printMessage = true
	}

	if strings.HasPrefix(message.Content, prefixReloadNicknames) {
		reloadNicknames(discord, message.ChannelID)
		printMessage = true
	}

	if strings.HasPrefix(message.Content, prefixHelp) {
		showHelp(discord, message.ChannelID, voiceChannelName[message.GuildID])
		printMessage = true
	}

	if printMessage {
		fmt.Printf("Message: %+v || From: %s\n", message.Message, message.Author)
	}

}

func checkNicknames(discord *discordgo.Session, message *discordgo.MessageCreate) {
	usernamesToCheck := strings.Split(message.Content, " ")[1:]
	channel := getChannelByName(discord, message.GuildID, voiceChannelName[message.GuildID])
	if channel == nil {
		return
	}
	users := getUsersConnectedToVoiceChannel(discord, channel, message.GuildID)
	nicknamesFound, nicknamesNotFound, discordMembersWhitoutCorrectNickname := compareUsersToNicknames(discord, users, usernamesToCheck, message.GuildID)

	_, _ = discord.ChannelMessageSendEmbed(message.ChannelID, &discordgo.MessageEmbed{
		Title: "Results of Nickname Check",
		Color: 16642983,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  fmt.Sprintf("Found nicknames (%d/%d)", len(nicknamesFound), len(usernamesToCheck)),
				Value: fmt.Sprintf("```\n%+v\n```", strings.Join(nicknamesFound, "\n")),
			},
			{
				Name:  fmt.Sprintf("Not found nicknames (%d/%d)", len(nicknamesNotFound), len(usernamesToCheck)),
				Value: fmt.Sprintf("```\n%+v\n```", strings.Join(nicknamesNotFound, "\n")),
			},
			{
				Name:  fmt.Sprintf("Discord members whitout correct nickname (%d/%d)", len(discordMembersWhitoutCorrectNickname), len(usernamesToCheck)),
				Value: fmt.Sprintf("```\n%+v\n```", strings.Join(discordMembersWhitoutCorrectNickname, "\n")),
			},
		},
	})
}

func getChannelByName(discord *discordgo.Session, guildID string, channelName string) *discordgo.Channel {
	channels, _ := discord.GuildChannels(guildID)
	for _, channel := range channels {
		if channel.Type != discordgo.ChannelTypeGuildVoice {
			continue
		}
		if strings.Compare(channelName, channel.Name) == 0 {
			return channel
		}
	}
	return nil
}

func getUsersConnectedToVoiceChannel(discord *discordgo.Session, channel *discordgo.Channel, guildID string) []*discordgo.User {
	var users []*discordgo.User
	guild, _ := discord.Guild(guildID)
	for _, vs := range guild.VoiceStates {
		if strings.Compare(vs.ChannelID, channel.ID) == 0 {
			user, err := discord.User(vs.UserID)
			if err != nil {
				continue
			}
			users = append(users, user)
		}
	}
	return users
}

func compareUsersToNicknames(discord *discordgo.Session, users []*discordgo.User, nicknames []string, guildID string) ([]string, []string, []string) {
	var foundNicknames []string
	var foundNicknamesWhitoutAccents []string
	var notFoundNicknames []string
	var discordMembersWhitoutCorrectNickname []string

	guildMembersWhitoutAccents := removeAccentsFromGuildMemberList(getGuildMembersFromUserList(discord, users, guildID))

NICKNAME_LOOP:
	for _, nickname := range nicknames {
		nicknameDeDicatedCaseSensetive, _, _ := transform.String(t, nickname)
		nicknameDeDicatedLowerCase := strings.ToLower(nicknameDeDicatedCaseSensetive)
		for _, guildMemberWhitoutAccents := range guildMembersWhitoutAccents {
			if strings.Compare(guildMemberWhitoutAccents.Nickname, nicknameDeDicatedLowerCase) == 0 {
				foundNicknames = append(foundNicknames, nickname)
				foundNicknamesWhitoutAccents = append(foundNicknamesWhitoutAccents, nicknameDeDicatedLowerCase)
				continue NICKNAME_LOOP
			}
		}
		notFoundNicknames = append(notFoundNicknames, nickname)
	}

	for _, guildMemberWhitoutAccents := range guildMembersWhitoutAccents {
		if isGuildMemberWithAccentsInNicknamesWhitoutAccentsList(foundNicknamesWhitoutAccents, guildMemberWhitoutAccents) {
			continue
		}
		discordMembersWhitoutCorrectNickname = append(discordMembersWhitoutCorrectNickname, guildMemberWhitoutAccents.Nickname)
	}

	sort.Strings(foundNicknames)
	sort.Strings(notFoundNicknames)
	sort.Strings(discordMembersWhitoutCorrectNickname)

	return foundNicknames, notFoundNicknames, discordMembersWhitoutCorrectNickname
}

func removeAccentsFromGuildMemberList(guildMembers []*discordgo.Member) []guildMemberWithAccentsRemoved {
	var guildMembersWithAccentsRemoved []guildMemberWithAccentsRemoved

	for _, guildMember := range guildMembers {
		if len(guildMember.Nick) > 0 {
			guildMemberNickWhitoutAccents, _, _ := transform.String(t, guildMember.Nick)
			guildMembersWithAccentsRemoved = append(guildMembersWithAccentsRemoved, guildMemberWithAccentsRemoved{
				GuildMember: guildMember,
				Nickname:    strings.ToLower(guildMemberNickWhitoutAccents),
			})
		} else {
			usernameWhitoutAccents, _, _ := transform.String(t, guildMember.User.Username)
			guildMembersWithAccentsRemoved = append(guildMembersWithAccentsRemoved, guildMemberWithAccentsRemoved{
				GuildMember: guildMember,
				Nickname:    strings.ToLower(usernameWhitoutAccents),
			})
		}
	}

	return guildMembersWithAccentsRemoved
}

func getGuildMembersFromUserList(discord *discordgo.Session, users []*discordgo.User, guildID string) []*discordgo.Member {
	var foundGuildMembers []*discordgo.Member
	guild, _ := discord.Guild(guildID)
	for _, guildMember := range guild.Members {
		if isDiscordUserInDiscordUserList(guildMember.User, users) {
			if isDiscordUserInGuildMembersList(guildMember.User, foundGuildMembers) {
				continue
			}
			foundGuildMembers = append(foundGuildMembers, guildMember)
		}
	}
	return foundGuildMembers
}

func isDiscordUserInDiscordUserList(user *discordgo.User, users []*discordgo.User) bool {
	for _, userFromList := range users {
		if strings.Compare(userFromList.ID, user.ID) == 0 {
			return true
		}
	}
	return false
}

func isDiscordUserInGuildMembersList(user *discordgo.User, guildMembers []*discordgo.Member) bool {
	for _, guildMember := range guildMembers {
		if strings.Compare(guildMember.User.ID, user.ID) == 0 {
			return true
		}
	}
	return false
}

func isGuildMemberWithAccentsInNicknamesWhitoutAccentsList(nicknames []string, guildMemberWithAccentsRemoved guildMemberWithAccentsRemoved) bool {
	for _, nickname := range nicknames {
		if strings.Compare(nickname, guildMemberWithAccentsRemoved.Nickname) == 0 {
			return true
		}
	}
	return false
}

func reloadNicknames(discord *discordgo.Session, channelID string) {
	discord.Close()
	err := discord.Open()
	errCheck("Error opening connection to Discord", err)
	_, _ = discord.ChannelMessageSendEmbed(channelID, &discordgo.MessageEmbed{
		Title: "Reload nicknames done",
		Color: 16642983,
	})
}

func showHelp(discord *discordgo.Session, channelID string, channelName string) {
	_, _ = discord.ChannelMessageSendEmbed(channelID, &discordgo.MessageEmbed{
		Title: "Jeeves Bot Help",
		Color: 16642983,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "~check_nicknames <List of nicknames>",
				Value: fmt.Sprintf("```\nChecks in %s voice chat if the nicknames are set correctly.\n```", channelName),
			},
			{
				Name:  "~reload_nicknames",
				Value: "```\nReloades the nicknames in the server\n```",
			},
			{
				Name:  "~help",
				Value: "```\nShows this help page\n```",
			},
		},
	})
}
