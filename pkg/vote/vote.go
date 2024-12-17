package vote

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

func HoldAVote(s *discordgo.Session, m *discordgo.MessageCreate, q string) (string, error) {
	// Initial Message
	fmt.Println("Holding a vote!", q)
	var yeses, nos int
	voteDeclare := "VOTE: " + q + " React 👍 or 👎 (Vote will be open for 5 minutes.)"
	voteMessage, err := s.ChannelMessageSend(m.ChannelID, voteDeclare)
	if err != nil {
		return "Couldn't send the vote message!", fmt.Errorf("failed to send initial vote message: %w", err)
	}

	// Timer
	if err := s.MessageReactionAdd(m.ChannelID, voteMessage.ID, "👍"); err != nil {
		return "Couldn't react to vote message!", fmt.Errorf("failed to react to initial vote message: %w", err)
	}
	if err := s.MessageReactionAdd(m.ChannelID, voteMessage.ID, "👎"); err != nil {
		return "Couldn't react to vote message!", fmt.Errorf("failed to react to initial vote message: %w", err)
	}

	timer1 := time.NewTimer(15 * time.Second)
	<-timer1.C
	fmt.Println("timer fired")
	// Counting the vote

	yesUsers, err := s.MessageReactions(m.ChannelID, voteMessage.ID, "👍", int(100), "", "")
	if err != nil {
		return "Uhm, I lost count :(", fmt.Errorf("Couldn't get the upvote reactions from message: %s", err)
	}
	noUsers, err := s.MessageReactions(m.ChannelID, voteMessage.ID, "👎", int(100), "", "")
	if err != nil {
		return "Uhm, I lost count :(", fmt.Errorf("Couldn't get the downvote reactions from message: %s", err)
	}
	yeses = len(yesUsers)
	nos = len(noUsers)
	yeses = yeses - 1
	nos = nos - 1
	/*reactions := voteMessage.Reactions
	for _, v := range reactions {
		fmt.Println(v.Emoji.Name, v.Emoji.ID)
		if v.Emoji.Name == "👍" {
			yeses = v.Count
		} else if v.Emoji.Name == "👎" {
			nos = v.Count
		}
	}*/
	fmt.Println("Vote concluded: Yes", yeses, "No", nos)
	if yeses > nos {
		return q + " Vote concluded! Chat votes yes! " + strconv.Itoa(yeses) + " Yes votes, " + strconv.Itoa(nos) + " No votes.", nil
	} else if nos > yeses {
		return q + " Vote concluded! Chat votes no! " + strconv.Itoa(nos) + " No votes, " + strconv.Itoa(yeses) + " Yes votes.", nil
	} else {
		if nos == 0 && yeses == 0 {
			return q + " Nobody voted! hmm...", nil
		} else if nos == yeses {
			return q + " Vote concluded! It's a tie! " + strconv.Itoa(nos) + " No votes, " + strconv.Itoa(yeses) + " Yes votes.", nil
		}
	}
	return "", nil
}
