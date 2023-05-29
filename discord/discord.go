package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	dg *discordgo.Session

	//identify map[uuid.UUID]chan discordgo.Guild
)

func StartDiscord(token string) error {
	var err error
	dg, err = discordgo.New("Bot " + token)
	if err != nil {
		fmt.Printf("error creating Discord session: %s", err)
		return err
	}

	dg.Identify.Intents = discordgo.IntentsGuildMessages // | discordgo.IntentsGuildMembers

	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionMessageComponent:
			if h, ok := componentsHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}
		}
	})

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return err
	}

	fmt.Println("Bot is now running.")
	fmt.Printf("invite with https://discord.com/oauth2/authorize?client_id=%s&permissions=1024&scope=bot\n", dg.State.User.ID)
	return nil
}

func Shutdown() {
	if dg != nil {
		_ = dg.Close()
	}
}

func VerifyUser(userID, code string) error {
	return verifyUser(dg, userID, code)
}

func FindUser(username, discriminator string) (string, error) {
	return findUser(dg, username, discriminator)
}

var componentsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"verified-rsi": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: i.Message.Embeds,
			},
		})
		if err != nil {
			fmt.Printf("Error responding: %s\n", err)
		}

		request, _ := json.Marshal(map[string]string{
			"verified": i.User.ID,
			//				"code":     "maybe put the code in --",
		})

		_, err = http.Post("http://foo.com", "application/json", bytes.NewBuffer(request)) // todo website url
		if err != nil {
			fmt.Printf("Error sending refresh request: %s\n", err)
		}
	},
}

func findUser(s *discordgo.Session, userName, discriminator string) (string, error) {
	guilds := make([]string, len(s.State.Guilds))
	for i, g := range s.State.Guilds {
		guilds[i] = g.ID
	}

	response := make(chan string)

	dg.AddHandlerOnce(func(s *discordgo.Session, c *discordgo.GuildMembersChunk) {
		for _, m := range c.Members {
			if m.User.Username == userName && m.User.Discriminator == discriminator {
				response <- m.User.ID
				return
			}
		}
	}) // should probably make a general handler with channels

	err := s.RequestGuildMembersBatch(guilds, userName, 0, "", false)
	if err != nil {
		return "", err
	}
	select {
	case id := <-response:
		return id, nil
	case <-time.After(1 * time.Second):
		return "", fmt.Errorf("timed out")
	}

}

func verifyUser(s *discordgo.Session, userID, code string) error {
	ch, err := s.UserChannelCreate(userID)
	if err != nil {
		return err
	}

	_, err = s.ChannelMessageSendComplex(ch.ID, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{{
			Color: 0xff00ff,
			Fields: []*discordgo.MessageEmbedField{{
				Value: fmt.Sprintf("Please add the `%s` to your [RSI accounts Short Bio](https://robertsspaceindustries.com/account/profile) to veriy your account, then click done", code),
			}},
		}},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Emoji: discordgo.ComponentEmoji{
							Name: "✔️",
						},
						Label:    "Done",
						Style:    discordgo.PrimaryButton,
						CustomID: "verified-rsi",
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}
