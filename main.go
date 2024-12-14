package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/lizthejester/lizbotgo/src/config"
	"github.com/lizthejester/lizbotgo/src/ltime"
	"golang.org/x/exp/rand"
)

type Lizbot struct {
	Name string `json:"Lizbot"`
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	userMessage := m.Content

	// dont respond to self
	if m.Author.ID == s.State.User.ID {
		return
		// dont respond to pk
	}
	if m.Author.ID == "1115685378704277585" {
		return
		// only respond to "?" proxy
	}
	if string(userMessage[0]) != "?" {
		return
	}
	userMessage = userMessage[1:]
	fmt.Println(userMessage)
	fmt.Println("getting response")
	response := getResponse(m, userMessage)
	if response == "" {
		s.ChannelMessageSend(m.ChannelID, "sorry, I don't know that command! :)")
	} else {
		s.ChannelMessageSend(m.ChannelID, response)
	}
}

func getResponse(m *discordgo.MessageCreate, userInput string) string {
	lowered := strings.ToLower(userInput)
	fmt.Println(lowered)

	//Empty message
	if lowered == "" {
		return "well you're awfully silent..."
	}

	//command list
	if lowered == "command list" {
		directory := "?magic8ball\n?flip a coin\n?roll a d4, d6, d8, d10, d12, or d20\n?inspire\n?joke\n?lizdate"
		return directory
	}

	//chat
	switch lowered {
	case "hello", "hi", "hihi", "howdy", "hiya", "hey", "greetings", "yo", "salutations":
		greetings := [...]string{"Hello there", "Hi", "Greetings", "Hihi", "Howdy", "Yo", "Salutations"}
		return greetings[rand.Intn(len(greetings))]
	case "goodbye", "bye", "see ya", "later", "see ya later", "see you later", "bye bye", "byebye":
		goodbyes := [...]string{"farewell Traveler!", "farewell!", "later! ^-^", "see ya! :3", "Bye!", "Bye now! ^-^", "byebye! ^-^"}
		return goodbyes[rand.Intn(len(goodbyes))]
	case "magic 8ball", "magic8ball":
		rand.Seed(uint64(time.Now().UnixNano()))
		ballResp := [...]string{"Yes, definitely",
			"It is certain",
			"Without a doubt",
			"You may rely on it",
			"As I see it, yes",
			"Most likely",
			"Outlook good",
			"Signs point to yes",
			"Yes",
			"Definitely",
			"Don’t count on it",
			"My reply is no",
			"My sources say no",
			"Outlook not so good",
			"Very doubtful",
			"Reply hazy, try again",
			"Ask again later",
			"Better not tell you now",
			"Cannot predict now",
			"Concentrate and ask again"}
		return ballResp[rand.Intn(len(ballResp))]
	case "tell me a joke", "joke", "tell a joke", "what's a good joke", "what's a good joke?", "know any jokes", "know a good joke?":
		jokes := [...]string{"What do kids play when their mom is using the phone? Bored games.",
			"What do you call an ant who fights crime? A vigilANTe!",
			"Why did the teddy bear say no to dessert? Because she was stuffed.",
			"Why did the scarecrow win a Nobel prize? Because she was outstanding in her field.",
			"What kind of shoes do frogs love? Open-toad!",
			"What did the ghost call his Mum and Dad? His transparents.",
			"What was a more useful invention than the first telephone? The second telephone.",
			"What’s a snake’s favorite subject in school? Hiss-tory.",
			"What animal is always at a baseball game? A bat."}
		return jokes[rand.Intn(len(jokes))]
	}

	// rolling
	if len(lowered) > 8 {
		if lowered[0:8] == "roll a d" {
			rand.Seed(uint64(time.Now().UnixNano()))
			switch lowered[8:] {
			case "4":
				return "You rolled:" + strconv.Itoa(rand.Intn(4)+1)
			case "6":
				return "You rolled:" + strconv.Itoa(rand.Intn(6)+1)
			case "8":
				return "You rolled:" + strconv.Itoa(rand.Intn(8)+1)
			case "10":
				return "You rolled:" + strconv.Itoa(rand.Intn(10)+1)
			case "12":
				return "You rolled:" + strconv.Itoa(rand.Intn(12)+1)
			case "20":
				return "You rolled:" + strconv.Itoa(rand.Intn(20)+1)
			default:
				return "I don't have that die! :("
			}
		}
	}
	// coin flip
	if lowered == "flip a coin" {
		coinResults := [2]string{"heads", "tails"}
		return coinResults[rand.Intn(len(coinResults))]
	}

	// tarot
	switch lowered {
	case "shuffle":
		initDeck(m)
		for _, user := range userCards {
			if user.ID == m.Author.ID {
				tarotShuffle(user.deck)
				return "Shuffled!"
			}
		}
	case "draw":
		initDeck(m)
		var inversion string
		if rand.Intn(2) == 1 {
			inversion = " inverted"
		}
		for _, user := range userCards {
			if user.ID == m.Author.ID {
				user.hand = append(user.hand, user.deck[len(user.deck)-1])
				user.deck = user.deck[:len(user.deck)-1]

				return user.hand[len(user.hand)-1] + inversion
			}
		}
	}

	//Miss Amie suggests
	/*sc := bufio.NewScanner(strings.NewReader(userInput))
	sc.Split(bufio.ScanWords)
	sc.Scan()
	fmt.Println(sc.Text())
	sc.Scan()
	fmt.Println(sc.Text())
	sc.Scan()
	fmt.Println(sc.Text())*/

	// calendar
	if len(lowered) > 6 {
		if lowered[0:7] == "lizdate" {
			if len(lowered) == 7 {
				currentYear, currentMonth, currentDay := time.Now().Date()
				lizMonth, lizDay := ltime.GetDayMonth(currentYear, currentMonth.String(), currentDay)
				fmt.Println("Current date:", lizMonth, lizDay, ltime.GetDayOfWeek(lizDay, lizMonth))
				response := "Current date: " + strconv.Itoa(lizDay) + " " + lizMonth + ", " + ltime.GetDayOfWeek(lizDay, lizMonth)
				return response
			} else {
				var firstSpaceIndex int
				var secondSpaceIndex int
				for i := 8; firstSpaceIndex == 0 && i < len(lowered); i++ {
					if string(lowered[i]) == " " {
						firstSpaceIndex = i
					}
					if i == len(lowered) {
						return "That command looks wrong"
					}
				}
				for i := firstSpaceIndex + 1; secondSpaceIndex == 0 && i < len(lowered); i++ {
					if string(lowered[i]) == " " {
						secondSpaceIndex = i
					}
					if i == len(lowered) {
						return "That command looks wrong"
					}
				}
				if firstSpaceIndex == 0 || secondSpaceIndex == 0 {
					return "That command looks wrong"
				}
				gregDay, err := strconv.Atoi(lowered[8:(firstSpaceIndex)])
				if err != nil {
					return "Formatting error; First argument should be a number (Day)"
				}
				fmt.Println("gregday:", gregDay)
				gregMonth := userInput[(firstSpaceIndex + 1):(secondSpaceIndex)]
				switch gregMonth {
				case "January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December":

				default:
					return "Formatting error; Second argument should be a Gregorian month (i.e. January)"
				}
				fmt.Println("gregmonth:", gregMonth)
				gregYear, err2 := strconv.Atoi(lowered[(secondSpaceIndex + 1):])
				if err2 != nil {
					return "Formatting error; Third argument should be a number (Year)"
				}
				fmt.Println("gregyear:", gregYear)
				lizMonth, lizDay := ltime.GetDayMonth(gregYear, gregMonth, gregDay)
				response := strconv.Itoa(lizDay) + " " + lizMonth + ", " + ltime.GetDayOfWeek(lizDay, lizMonth)
				return response
			}
		}
	}
	fmt.Println("response got")
	return ""
}

type User struct {
	ID   string
	hand []string
	deck []string
}

var userCards []*User

// Initialize deck
func initDeck(m *discordgo.MessageCreate) {
	for _, user := range userCards {
		if user.ID == m.Author.ID {
			return
		}
	}
	var newUser User
	newUser.ID = m.Author.ID
	newUser.hand = []string{}
	newUser.deck = []string{
		"The Fool(0)",
		"The Magician(I)",
		"The High Priestess(II)",
		"The Empress(III)",
		"The Emporer(IV)",
		"The Heirophant(V)",
		"The Lovers(VI)",
		"The Chariot(VII)",
		"Strength(VIII)",
		"The Hermit(IX)",
		"Wheel of Fortune(X)",
		"Justice(XI)",
		"The Hanged Man(XII)",
		"Death(XIII)",
		"Temperance(XIV)",
		"The Devil(XV)",
		"The Tower(XVI)",
		"The Star(XVII)",
		"The Moon(XVIII)",
		"The Sun(XIX)",
		"Judgement(XX)",
		"The World(XXI)",
		"King of Swords",
		"Queen of Swords",
		"Knight of Swords",
		"Page of Swords",
		"One of Swords",
		"Two of Swords",
		"Three of Swords",
		"Four of Swords",
		"Five of Swords",
		"Six of Swords",
		"Seven of Swords",
		"Eight of Swords",
		"Nine of Swords",
		"Ten of Swords",
		"King of Batons",
		"Queen of Batons",
		"Knight of Batons",
		"Page of Batons",
		"One of Batons",
		"Two of Batons",
		"Three of Batons",
		"Four of Batons",
		"Five of Batons",
		"Six of Batons",
		"Seven of Batons",
		"Eight of Batons",
		"Nine of Batons",
		"Ten of Batons",
		"King of Coins",
		"Queen of Coins",
		"Knight of Coins",
		"Page of Coins",
		"One of Coins",
		"Two of Coins",
		"Three of Coins",
		"Four of Coins",
		"Five of Coins",
		"Six of Coins",
		"Seven of Coins",
		"Eight of Coins",
		"Nine of Coins",
		"Ten of Coins",
		"King of Cups",
		"Queen of Cups",
		"Knight of Cups",
		"Page of Cups",
		"One of Cups",
		"Two of Cups",
		"Three of Cups",
		"Four of Cups",
		"Five of Cups",
		"Six of Cups",
		"Seven of Cups",
		"Eight of Cups",
		"Nine of Cups",
		"Ten of Cups",
	}
	tarotShuffle(newUser.deck)
	userCards = append(userCards, &newUser)

	//Rust suggests
	//78uuuuuuuuuuuuuuuuu

	fmt.Println("Appended new user to usercards", newUser.ID)
	fmt.Print(userCards)
}

func tarotShuffle(deck []string) {
	rand.Seed(uint64(time.Now().UnixNano()))
	for i := len(deck) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		deck[i], deck[j] = deck[j], deck[i]
	}
}

var s *discordgo.Session

func main() {
	var err error
	s, err = discordgo.New("Bot " + config.DISCORD_KEY)
	if err != nil {
		fmt.Printf("Invalid bot parameters: %v", err)
		return
	}
	//BORROWED CODE
	// Register the messageCreate func as a callback for MessageCreate events.
	s.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	s.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = s.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	s.Close()
	//END
}
