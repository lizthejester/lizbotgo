package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/Lizthejester/LizianTime/pkg/ltime"
	"github.com/bwmarrin/discordgo"
	"github.com/lizthejester/lizbotgo/pkg/alarm"
	"github.com/lizthejester/lizbotgo/pkg/chanselect"
	"github.com/lizthejester/lizbotgo/pkg/config"
	"github.com/lizthejester/lizbotgo/pkg/inspire"
	"github.com/lizthejester/lizbotgo/pkg/roll"
	"github.com/lizthejester/lizbotgo/pkg/user"
	"github.com/lizthejester/lizbotgo/pkg/vote"
	"golang.org/x/exp/rand"
)

type Lizbot struct {
	Name string `json:"Lizbot"`
}

var UserManager *user.UserManager = user.NewManager()
var ServerManager *chanselect.ServerManager = chanselect.NewServerManager()

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	userMessage := m.Content
	if userMessage == "" {
		return
	}
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
	response := getResponse(s, m, userMessage)
	if response == "" {
		s.ChannelMessageSend(m.ChannelID, "sorry, I don't know that command! :)")
	} else if response == "No response" {
		return
	} else {
		s.ChannelMessageSend(m.ChannelID, response)
	}
}

func getResponse(s *discordgo.Session, m *discordgo.MessageCreate, userInput string) string {
	user := UserManager.GetUser(m.Author.ID, s, ServerManager)
	lowered := strings.ToLower(userInput)
	fmt.Println(lowered)

	if lowered == "set main channel" {
		ServerManager.GetServer(m.GuildID).SetChannel(m.ChannelID)
		fmt.Println(ServerManager.GetServer(m.GuildID).MainChannel)
		return "Main channel set!"
	}

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
		greetings := []string{"Hello there", "Hi", "Greetings", "Hihi", "Howdy", "Yo", "Salutations"}
		return greetings[rand.Intn(len(greetings))]
	case "goodbye", "bye", "see ya", "later", "see ya later", "see you later", "bye bye", "byebye":
		goodbyes := []string{"farewell Traveler!", "farewell!", "later! ^-^", "see ya! :3", "Bye!", "Bye now! ^-^", "byebye! ^-^"}
		return goodbyes[rand.Intn(len(goodbyes))]
	case "magic 8ball", "magic8ball":
		rand.Seed(uint64(time.Now().UnixNano()))
		ballResp := []string{"Yes, definitely",
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
		jokes := []string{"What do kids play when their mom is using the phone? Bored games.",
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
	if strings.HasPrefix(lowered, "roll a d") {
		return roll.RollDice(lowered[8:])
	}
	// coin flip
	if lowered == "flip a coin" {
		possCoinResults := [2]string{"heads", "tails"}
		return possCoinResults[rand.Intn(len(possCoinResults))]
	}
	// hold a vote
	if strings.HasPrefix(lowered, "hold a vote") {
		if len(lowered) > 11 {
			resp, err := vote.HoldAVote(s, m, userInput[12:])
			if err != nil {
				fmt.Println(err)
			}
			return resp
		} else {
			return "Hold a vote on what?"
		}
	}
	// quotes
	if lowered == "inspire" {
		return inspire.GetQuote()
	}

	// tarot
	switch lowered {
	case "shuffle":
		user.TarotDeck.TarotShuffle()
		return "Shuffled!"

	case "draw":
		deck := user.TarotDeck
		card := deck.Draw()

		return card

	case "reset deck":
		user.TarotDeck.ResetDeck()

		return "Deck reset."
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
	if strings.HasPrefix(lowered, "lizdate") {
		if len(lowered) == 7 {
			currentYear, currentMonth, currentDay := time.Now().Local().Date()
			lizMonth, lizDay := ltime.GetDayMonth(currentYear, currentMonth.String(), currentDay)
			fmt.Println("Current date:", lizMonth, lizDay, ltime.GetDayOfWeek(lizDay, lizMonth))
			response := "Current date: " + strconv.Itoa(lizDay) + " " + lizMonth + ", " + ltime.GetDayOfWeek(lizDay, lizMonth)
			return response
		}

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

	//example command: "set alarm for (day) (month) (year) (time + timezone) (name) (comment string)"
	if strings.HasPrefix(lowered, "set alarm for") {
		//TODO: Parse Deadline
		wrongSyntaxMessage := "Syntax is: 01 30 2006 03:04PM -0800 \"Name\" Description"

		if len(lowered) > 14 {
			firstSpaceIndex := 0
			secondSpaceIndex := 0
			thirdSpaceIndex := 0
			fourthSpaceIndex := 0
			fifthSpaceIndex := 0
			var alarmName string
			var alarmComment string
			var dline string
			for i := 15; firstSpaceIndex == 0; i++ {
				if string(lowered[i]) == " " {
					firstSpaceIndex = i
				}
				if i == len(lowered) {
					return wrongSyntaxMessage
				}
			}
			for i := firstSpaceIndex + 1; secondSpaceIndex == 0; i++ {
				if string(lowered[i]) == " " {
					secondSpaceIndex = i
				}
				if i == len(lowered) {
					return wrongSyntaxMessage
				}
			}
			for i := secondSpaceIndex + 1; thirdSpaceIndex == 0; i++ {
				if string(lowered[i]) == " " {
					thirdSpaceIndex = i
				}
				if i == len(lowered) {
					return wrongSyntaxMessage
				}
			}
			for i := thirdSpaceIndex + 1; fourthSpaceIndex == 0; i++ {
				if string(lowered[i]) == " " {
					fourthSpaceIndex = i
				}
				if i == len(lowered) {
					return wrongSyntaxMessage
				}
			}
			for i := fourthSpaceIndex + 1; fifthSpaceIndex == 0; i++ {
				if string(lowered[i]) == " " {
					fifthSpaceIndex = i
				}
				if i == len(lowered) {
					return wrongSyntaxMessage
				}
			}

			secondQuotationMark := 0
			if string(lowered[fifthSpaceIndex+1]) == "\"" {
				for i := fifthSpaceIndex + 2; secondQuotationMark == 0; i++ {
					if string(lowered[i]) == "\"" {
						secondQuotationMark = i
					}
					if i == len(lowered) {
						return wrongSyntaxMessage
					}
				}
			} else {
				return wrongSyntaxMessage + "this one"
			}

			alarmName = userInput[fifthSpaceIndex+2 : secondQuotationMark]
			alarmComment = userInput[secondQuotationMark+1:]
			dline = userInput[14:fifthSpaceIndex]

			/*if firstSpaceIndex == 0 || secondSpaceIndex == 0 || thirdSpaceIndex == 0 {
				return wrongSyntaxMessage
			}*/

			alm := &alarm.Alarm{
				ChannelID: m.ChannelID,
				Deadline:  dline,
				Content:   alarmComment,
				Name:      alarmName,
				UserID:    m.Author.ID,
				ServerID:  m.GuildID,
			}

			s.ChannelMessageSend(m.ChannelID, alarmName+" set for "+dline)
			UserManager.GetUser(m.Author.ID, s, ServerManager).AlarmManager.SetAlarm(alm, s, m.ChannelID)
			return "No response"
		} else {
			return wrongSyntaxMessage
		}
	}

	fmt.Println("response got")
	return ""
}
func SaveServers() {
	db, err := sql.Open("sqlite3", "file:lizbot.db?cache=shared&_timeout=1000")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	for serverid, server := range ServerManager.Servers {
		if server.MainChannel == "" {
			continue
		}
		_, err = db.Exec("insert into servers (serverid, mainchannel) values(?, ?)", serverid, server.MainChannel)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func SaveAlarms() {
	emptiedTime := "01 02 2006 03:04PM -0700"
	db, err := sql.Open("sqlite3", "file:lizbot.db?cache=shared&_timeout=1000")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	for _, user := range UserManager.GetAllUsers() {
		for _, thisAlarm := range user.AlarmManager.GetAlarms() {
			if thisAlarm.Deadline == emptiedTime {
				continue
			} else {
				_, err = db.Exec("insert into alarms (name, time, comment, channelid, userid, serverid) values(?, ?, ?, ?, ?, ?)",
					thisAlarm.Name,
					thisAlarm.Deadline,
					thisAlarm.Content,
					thisAlarm.ChannelID,
					thisAlarm.UserID,
					thisAlarm.ServerID,
				)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}

func LoadServers() {
	db, err := sql.Open("sqlite3", "file:lizbot.db?cache=shared&_timeout=1000")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	rows, err := db.Query("select * from servers")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var serverid string
		var channelid string
		err = rows.Scan(&serverid, &channelid)
		if err != nil {
			fmt.Println(err)
		}
		ServerManager.Servers[serverid] = &chanselect.Server{
			MainChannel: channelid,
		}
	}
}

/*func LoadExpiredAlarms() {
	db, err := sqlx.Open("sqlite3", "file:lizbot.db?cache=shared&_timeout=1000")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	rows, err := db.Query("select * from alarms")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	tx, err := db.Begin()
	for rows.Next() {
		var newalarm alarm.Alarm
		var dbid int

		if erra := rows.Scan(&dbid, &newalarm.Name, &newalarm.Deadline, &newalarm.Content, &newalarm.ChannelID, &newalarm.UserID, &newalarm.ServerID); err != nil {
			fmt.Println(erra, dbid)
		}
		deadlineTime, errd := time.Parse("01 02 2006 03:04PM -0700", newalarm.Deadline)
		if errd != nil {
			fmt.Println(errd)
		}
		if time.Until(deadlineTime) <= 0 {
			ServerManager.GetServer(newalarm.ServerID).ExpiredAlarmManager.Alarms = append(ServerManager.Servers[newalarm.ServerID].ExpiredAlarmManager.Alarms, newalarm)

			_, err = tx.Exec("DELETE FROM alarms WHERE name=$1 AND time=$2", newalarm.Name, newalarm.Deadline)
		}
	}
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
	} else {
		tx.Commit()
	}
	_, err = db.Exec("delete from servers")
	if err != nil {
		fmt.Println(err)
	}
}*/

func SendExpiredAlarms() {
	for _, server := range ServerManager.Servers {
		if server.MainChannel == "" {
			continue
		}
		var newMessage string
		for _, alarm := range server.ExpiredAlarmManager.Alarms {
			newMessage = newMessage + "<@" + alarm.UserID + "> " + alarm.Name + " has gone off!\n"
		}
		s.ChannelMessageSend(server.MainChannel, newMessage)
	}
}

func LoadAlarms() {
	db, err := sql.Open("sqlite3", "file:lizbot.db?cache=shared&_timeout=1000")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	rows, err := db.Query("select * from alarms")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var newalarm alarm.Alarm
		var dbid int
		if err = rows.Scan(&dbid, &newalarm.Name, &newalarm.Deadline, &newalarm.Content, &newalarm.ChannelID, &newalarm.UserID, &newalarm.ServerID); err != nil {
			fmt.Println(err)
		}
		_ = UserManager.GetUser(newalarm.UserID, s, ServerManager)
	}
}

func DeleteServers() {
	db, err := sql.Open("sqlite3", "file:lizbot.db?cache=shared&_timeout=1000")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	_, err = db.Exec("delete from servers")
	if err != nil {
		fmt.Println(err)
	}
}

func DeleteExpiredAlarms() {
	db, err := sql.Open("sqlite3", "file:lizbot.db?cache=shared&_timeout=1000")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	for _, v := range ServerManager.Servers {
		for _, alarm := range v.ExpiredAlarmManager.Alarms {
			_, err = db.Exec("DELETE FROM alarms WHERE name=? AND time=?", alarm.Name, alarm.Deadline)
			if err != nil {
				fmt.Println(err)
			}

		}
	}
}

var s *discordgo.Session

func main() {
	db, err := sql.Open("sqlite3", "file:lizbot.db?cache=shared&_timeout=1000")
	if err != nil {
		fmt.Println(err)
	}

	sqlstatement := "create table if not exists servers (serverid text not null primary key, mainchannel text)"
	_, err = db.Exec(sqlstatement)
	if err != nil {
		fmt.Println(err)
	}

	sqlstatement = "create table if not exists alarms (id integer not null primary key, name text, time text, comment text, channelid text, userid text, serverid text)"
	_, err = db.Exec(sqlstatement)
	if err != nil {
		fmt.Println(err)
	}
	db.Close()

	s, err = discordgo.New("Bot " + config.DISCORD_KEY)
	if err != nil {
		fmt.Printf("Invalid bot parameters: %v", err)
		return
	}
	//BORROWED CODE
	// Register the messageCreate func as a callback for MessageCreate events.
	s.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	s.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages

	// Open a websocket connection to Discord and begin listening.
	err = s.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	LoadServers()
	DeleteServers()
	//LoadExpiredAlarms()
	LoadAlarms()
	SendExpiredAlarms()
	DeleteExpiredAlarms()

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	SaveAlarms()
	SaveServers()
	// Cleanly close down the Discord session.
	s.Close()
	//END
}
