package main

import (
	"fmt"
	"strings"
	"unicode"
	"os"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/text/transform"
  "golang.org/x/text/unicode/norm"
)

var (
	commandPrefix string
	botID         string
)

var prefixUsernameCheck = "~check_nicknames"
var voiceChannelName = "Raid Chat"
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
		fmt.Printf("SuperAwesomeOmegaTutorBot has started on %d servers", len(servers))
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
	if message.Author.ID == botID || message.Author.Bot {
		return
	}
	if strings.HasPrefix(message.Content, prefixUsernameCheck) {
		fmt.Printf("Message: %+v || From: %s\n", message.Message, message.Author)
		checkUsernames(discord, message)
	}

}

func checkUsernames(discord *discordgo.Session, message *discordgo.MessageCreate) {
	usernamesToCheck := strings.Split(message.Content, " ")[1:]
	channel := getChannelByName(discord, message.GuildID, voiceChannelName)
	if channel == nil {
		return
	}
	users := getUsersConnectedToVoiceChannel(discord, channel, message.GuildID)
	nicknamesFound, nicknamesNotFound := compareUsersToNicknames(discord, users, usernamesToCheck, message.GuildID)
	_, _ = discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Found nicknames: %+v", nicknamesFound))
	_, _ = discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Not found nicknames: %+v", nicknamesNotFound))
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

func compareUsersToNicknames(discord *discordgo.Session, users []*discordgo.User, nicknames []string, guildID string) ([]string, []string) {
	var foundNicknames []string
	var notFoundNicknames []string

	//This is used to speedup the loop
	var interestingGuildMembers []*discordgo.Member
	guild, _ := discord.Guild(guildID)
	for _, guildMember := range guild.Members {
		if isDiscordUserInDiscordUserList(guildMember.User, users) && len(guildMember.Nick) > 0 {
			if isDiscordUserInGuildMembersList(guildMember.User, interestingGuildMembers) {
				continue
			}
			interestingGuildMembers = append(interestingGuildMembers, guildMember)
		}
	}

NICKNAME_LOOP:
	for _, nickname := range nicknames {
		nicknameDeDicatedCaseSensetive, _, _ := transform.String(t, nickname)
		nicknameDeDicatedLowerCase := strings.ToLower(nicknameDeDicatedCaseSensetive)
		for _, guildMember := range interestingGuildMembers {
			guildMemberNickDeDicated, _, _ := transform.String(t, guildMember.Nick)
			if strings.Compare(strings.ToLower(guildMemberNickDeDicated), nicknameDeDicatedLowerCase) == 0 {
				foundNicknames = append(foundNicknames, nickname)
				continue NICKNAME_LOOP
			}
		}
		for _, user := range users {
			userUsernameDeDicated, _, _ := transform.String(t, user.Username)
			if strings.Compare(strings.ToLower(userUsernameDeDicated), strings.ToLower(nicknameDeDicatedLowerCase)) == 0 {
				foundNicknames = append(foundNicknames, nickname)
				continue NICKNAME_LOOP
			}
		}
		notFoundNicknames = append(notFoundNicknames, nickname)
	}

	return foundNicknames, notFoundNicknames
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
