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
	"github.com/lizthejester/lizbotgo/pkg/explain"
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

// This function will be called (due to AddHandler above) every time a newmessage is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	userMessage := m.Content
	// handles empty messages
	if userMessage == "" {
		return
	}
	// doesnt respond to self
	if m.Author.ID == s.State.User.ID {
		return
	}
	// doesnt respond to pluralkit ((I'm not totally sure this actually works))
	if m.Author.ID == "1115685378704277585" {
		return
	}
	// only responds to "?" proxy
	if string(userMessage[0]) != "?" {
		return
	}
	//removes ? from user input once read
	userMessage = userMessage[1:]
	response := getResponse(s, m, userMessage)
	// handles unknown commands
	if response == "" {
		s.ChannelMessageSend(m.ChannelID, "sorry, I don't know that command! :)")
	} else if response == "No response" {
		return
	} else {
		s.ChannelMessageSend(m.ChannelID, response)
	}
}

// get response contains all the user input commands
func getResponse(s *discordgo.Session, m *discordgo.MessageCreate, userInput string) string {
	user := UserManager.GetUser(m.Author.ID, s, ServerManager)
	lowered := strings.ToLower(userInput)
	fmt.Println(lowered)
	// COMMAND LIST
	//set main channel
	if lowered == "set main channel" {
		perms, err := s.UserChannelPermissions(m.Author.ID, m.ChannelID)
		if err != nil {
			fmt.Println(err)
		}
		if perms&discordgo.PermissionManageMessages == discordgo.PermissionManageMessages {
			ServerManager.GetServer(m.GuildID).SetChannel(m.ChannelID)
			fmt.Println(ServerManager.GetServer(m.GuildID).MainChannel)
			return "Main channel set!"
		} else {
			return "Not Admin."
		}
	}
	// command list
	if lowered == "command list" {
		directory := "?magic8ball\n?flip a coin\n?roll a d4, d6, d8, d10, d12, or d20\n?inspire\n?joke\n?lizdate"
		return directory
	}
	// explain
	if strings.HasPrefix(lowered, "explain") {
		if len(lowered) > 7 {
			resp, err := explain.Explain(s, m, lowered[8:])
			if err != nil {
				fmt.Println(err)
			}
			return resp
		} else {
			return "Hold a vote on what?"
		}
	}
	// chat
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
	// Miss Amie suggests
	/*sc := bufio.NewScanner(strings.NewReader(userInput))
	sc.Split(bufio.ScanWords)
	sc.Scan()
	fmt.Println(sc.Text())
	sc.Scan()
	fmt.Println(sc.Text())
	sc.Scan()
	fmt.Println(sc.Text())*/
	// calendar

	// lizdate
	if strings.HasPrefix(lowered, "lizdate") {

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

		gregMonth := lowered[(firstSpaceIndex + 1):(secondSpaceIndex)]
		if len(gregMonth) == 1 {
			gregMonth = "0" + gregMonth
		}
		switch gregMonth {
		case "january", "jan", "01", "february", "feb", "02", "march", "mar", "03", "april", "apr", "04", "may", "05", "june", "jun", "06", "july", "jul", "07", "august", "aug", "08", "september", "sept", "sep", "09", "october", "oct", "10", "november", "nov", "11", "december", "dec", "12":

		default:
			return "Formatting error; Second argument should be a Gregorian month (i.e. January)"
		}
		fmt.Println("gregmonth:", gregMonth)

		gregYear, err2 := strconv.Atoi(lowered[(secondSpaceIndex + 1):])
		if err2 != nil {
			return "Formatting error; Third argument should be a number (Year)"
		}
		fmt.Println("gregyear:", gregYear)

		if len(lowered) == 7 {
			currentYear, currentMonth, currentDay := time.Now().Local().Date()
			lizMonth, lizDay := ltime.GetDayMonth(currentYear, currentMonth.String(), currentDay)
			fmt.Println("Current date: ", lizMonth, lizDay, ltime.GetDayOfWeek(lizDay, lizMonth))
			response := "Current date: " + ltime.GetDayOfWeek(lizDay, lizMonth) + strconv.Itoa(lizDay) + " " + lizMonth + ", " + strconv.Itoa(gregYear)
			return response
		} else {

			lizMonth, lizDay := ltime.GetDayMonth(gregYear, gregMonth, gregDay)
			response := ltime.GetDayOfWeek(lizDay, lizMonth) + strconv.Itoa(lizDay) + " " + lizMonth + ", " + strconv.Itoa(gregYear)
			return response
		}
	}
	// set alarms
	if strings.HasPrefix(lowered, "set alarm for") {
		// example command: "set alarm for (month) (day) (year) ((time)AM/PM) (timezone) ("name") (loop frequency) ("comment string")"
		// note: use of military time does not require colon in time, declaration of AM/PM, but does require timezones.
		// note: any value of loop frequency that is not "daily", "weekly", "monthly", or "yearly" will prevent an alarm from looping but one must be present or the first word of the comment will be used as the loop frequency and the comment will be printed without the first word.
		wrongSyntaxMessage := "Syntax is: month day year 03:04PM timezone \"Name\" \"Description\" loopFrequency \n Example: April 20th 2024 04:20pm PST \"smokin'\" \"That Jazz Cabbage\" daily \n Example: 04 20 24 0420 -0800 \"smokin\" \"that jazz cabbage\" daily"
		// indexes for parsing user input
		if len(lowered) > 14 {
			firstSpaceIndex := 0
			secondSpaceIndex := 0
			thirdSpaceIndex := 0
			fourthSpaceIndex := 0
			fifthSpaceIndex := 0
			for i := 15; firstSpaceIndex == 0; i++ {
				if string(lowered[i]) == " " {
					firstSpaceIndex = i
				}
				if i == len(lowered) {
					return "Not enough spaces. " + wrongSyntaxMessage
				}
			}
			for i := firstSpaceIndex + 1; secondSpaceIndex == 0; i++ {
				if string(lowered[i]) == " " {
					secondSpaceIndex = i
				}
				if i == len(lowered) {
					return "Not enough spaces. " + wrongSyntaxMessage
				}
			}
			for i := secondSpaceIndex + 1; thirdSpaceIndex == 0; i++ {
				if string(lowered[i]) == " " {
					thirdSpaceIndex = i
				}
				if i == len(lowered) {
					return "Not enough spaces. " + wrongSyntaxMessage
				}
			}
			for i := thirdSpaceIndex + 1; fourthSpaceIndex == 0; i++ {
				if string(lowered[i]) == " " {
					fourthSpaceIndex = i
				}
				if i == len(lowered) {
					return "Not enough spaces. " + wrongSyntaxMessage
				}
			}
			for i := fourthSpaceIndex + 1; fifthSpaceIndex == 0; i++ {
				if string(lowered[i]) == " " {
					fifthSpaceIndex = i
				}
				if i == len(lowered) {
					return "Not enough spaces. " + wrongSyntaxMessage
				}
			}
			// set colon index
			colonIndex := 0
			for i := thirdSpaceIndex + 1; colonIndex == 0; i++ {
				if string(lowered[i]) == ":" {
					colonIndex = i
				}
				if i == len(lowered) {
					break
				}
			}
			// set index of closing quotation mark on field "Name"
			secondQuotationMark := 0
			if string(lowered[fifthSpaceIndex+1]) == "\"" {
				for i := fifthSpaceIndex + 2; secondQuotationMark == 0; i++ {
					if string(lowered[i]) == "\"" {
						secondQuotationMark = i
					}
					if i == len(lowered) {
						return "Missing quotation mark. " + wrongSyntaxMessage
					}
				}
			} else {
				return "Missing quotation mark. " + wrongSyntaxMessage
			}
			// reference: (month)1(day)2(year)3((time)AM/PM)4(timezone)5("name")(2q)(6)("comment string")(4q)(7)(loop frequency)
			fourthQuotationMark := 0
			if string(lowered[secondQuotationMark+2]) == "\"" {
				for i := secondQuotationMark + 3; fourthQuotationMark == 0; i++ {
					if string(lowered[i]) == "\"" {
						fourthQuotationMark = i
					}
					if i == len(lowered) {
						return "Missing quotation mark. " + wrongSyntaxMessage
					}
				}
			} else {
				return "Missing quotation mark. " + wrongSyntaxMessage
			}
			var loopFreq string
			if len(lowered)-1 != fourthQuotationMark {
				/*for i := fourthQuotationMark + 1; seventhSpaceIndex == 0; i++ {
					if string(lowered[i]) == " " {
						seventhSpaceIndex = i
					}
					if i == len(lowered)-1 {
						return "Not enough spaces. " + wrongSyntaxMessage
					}
				}*/
				loopFreq = userInput[fourthQuotationMark+2:]
			}
			// reference: (month)1(day)2(year)3((time)AM/PM)4(timezone)5("name")(2q)(6)("comment string")(4q)(7)(loop frequency)
			// parse user input. switches allow for dynamic input and some typos.
			alarmName := userInput[fifthSpaceIndex+2 : secondQuotationMark]
			alarmComment := userInput[secondQuotationMark+3 : fourthQuotationMark]
			dlDay := lowered[firstSpaceIndex+1 : secondSpaceIndex]
			dlMonth := lowered[14:firstSpaceIndex]
			dlDayInt, err := strconv.Atoi(dlDay)
			if err != nil {
				fmt.Println(err)
			}
			leapYear := false
			dlYear := lowered[secondSpaceIndex+1 : thirdSpaceIndex]
			// dlYear
			if len(dlYear) == 2 {
				dlYear = "20" + dlYear
			}
			if len(dlYear) > 4 {
				return "Year too long. " + wrongSyntaxMessage
			}
			if len(dlYear) < 2 {
				return "Year too short. " + wrongSyntaxMessage
			}
			dlYearInt, err := strconv.Atoi(dlYear)
			if err != nil {
				fmt.Println(err)
			}
			if dlYearInt%4 == 0 {
				leapYear = true
			}
			//dlDay switch
			switch dlDay {
			case "01", "1", "first", "1st":
				dlDay = "01"
			case "02", "2", "second", "2nd":
				dlDay = "02"
			case "03", "3", "third", "3rd":
				dlDay = "03"
			case "04", "4", "fourth", "4th":
				dlDay = "04"
			case "05", "5", "fifth", "5th":
				dlDay = "05"
			case "06", "6", "sixth", "6th":
				dlDay = "06"
			case "07", "7", "seventh", "7th":
				dlDay = "07"
			case "08", "8", "eighth", "8th":
				dlDay = "08"
			case "09", "9", "ninth", "9th":
				dlDay = "09"
			case "10", "tenth", "10th":
				dlDay = "10"
			case "11", "eleventh", "11th":
				dlDay = "11"
			case "12", "twelfth", "twelveth", "twelvth", "12th":
				dlDay = "12"
			case "13", "thirteenth", "13th":
				dlDay = "13"
			case "14", "fourteenth", "14th":
				dlDay = "14"
			case "15", "fifteenth", "15th":
				dlDay = "15"
			case "16", "sixteenth", "16th":
				dlDay = "16"
			case "17", "seventeenth", "17th":
				dlDay = "17"
			case "18", "eighteenth", "18th":
				dlDay = "18"
			case "19", "ninteenth", "nineteenth", "19th":
				dlDay = "19"
			case "20", "twentyth", "twentyeth", "20th":
				dlDay = "20"
			case "21", "twentyfirst", "21st":
				dlDay = "21"
			case "22", "twentysecond", "22nd":
				dlDay = "22"
			case "23", "twentythird", "23rd":
				dlDay = "23"
			case "24", "twentyfourth", "24th":
				dlDay = "24"
			case "25", "twentyfifth", "25th":
				dlDay = "25"
			case "26", "twentysixth", "26th":
				dlDay = "26"
			case "27", "twentyseventh", "27th":
				dlDay = "27"
			case "28", "twentyeighth", "28th":
				dlDay = "28"
			case "29", "twentyninth", "twentynineth", "29th":
				dlDay = "29"
			case "30", "thirtyth", "thirtyeth", "thirtieth", "30th":
				dlDay = "30"
			case "31", "thirtyfirst", "31st":
				dlDay = "31"
			case "32", "thirtysecond", "32nd":
				dlDay = "32"
			case "33", "thirtthird", "33rd":
				dlDay = "33"
			case "34", "thirtyfourth", "34th":
				dlDay = "34"
			case "35", "thirtyfifth", "35th":
				dlDay = "35"
			case "36", "thirtysixth", "36th":
				dlDay = "36"
			case "37", "thirtseventhh", "37th":
				dlDay = "37"
			case "38", "thirtyeighth", "38th":
				dlDay = "38"
			default:
				return "Problem with day. " + wrongSyntaxMessage
			}
			// dlMonth switches
			if leapYear {
				switch dlMonth {
				case "january", "jan", "1", "01":
					dlMonth = "01"
					if dlDayInt > 31 {
						return "Problem with day. " + wrongSyntaxMessage
					}
				case "february", "feb", "2", "02":
					dlMonth = "02"
					if dlDayInt > 29 {
						return "Problem with day. " + wrongSyntaxMessage
					}
				case "march", "mar", "3", "03":
					dlMonth = "03"
					if dlDayInt > 31 {
						return "Problem with day. " + wrongSyntaxMessage
					}
				case "april", "apr", "4", "04":
					dlMonth = "04"
					if dlDayInt > 30 {
						return "Problem with day. " + wrongSyntaxMessage
					}
				case "may", "5", "05":
					dlMonth = "05"
					if dlDayInt > 31 {
						return "Problem with day. " + wrongSyntaxMessage
					}
				case "june", "jun", "6", "06":
					dlMonth = "06"
					if dlDayInt > 30 {
						return "Problem with day. " + wrongSyntaxMessage
					}
				case "july", "jul", "7", "07":
					dlMonth = "07"
					if dlDayInt > 31 {
						return "Problem with day. " + wrongSyntaxMessage
					}
				case "august", "aug", "8", "08":
					dlMonth = "08"
					if dlDayInt > 31 {
						return "Problem with day. " + wrongSyntaxMessage
					}
				case "september", "sep", "9", "09":
					dlMonth = "09"
					if dlDayInt > 30 {
						return "Problem with day. " + wrongSyntaxMessage
					}
				case "october", "oct", "10":
					dlMonth = "10"
					if dlDayInt > 31 {
						return "Problem with day. " + wrongSyntaxMessage
					}
				case "november", "nov", "11":
					dlMonth = "11"
					if dlDayInt > 30 {
						return "Problem with day. " + wrongSyntaxMessage
					}
				case "december", "dec", "12":
					dlMonth = "12"
					if dlDayInt > 31 {
						return "Problem with day. " + wrongSyntaxMessage
					}
				case "menotheen":
					if dlDayInt > 31 {
						dlDay = strconv.Itoa(dlDayInt - 31)
						dlMonth = "02"
					} else {
						dlDay = strconv.Itoa(dlDayInt)
						dlMonth = "01"
					}
				case "lengten":
					if dlDayInt > 24 {
						dlDay = strconv.Itoa(dlDayInt - 24)
						dlMonth = "03"
					} else {
						dlDay = strconv.Itoa(dlDayInt + 5)
						dlMonth = "02"
					}
				case "regen":
					if dlDayInt > 18 {
						dlDay = strconv.Itoa(dlDayInt - 18)
						dlMonth = "04"
					} else {
						dlDay = strconv.Itoa(dlDayInt + 13)
						dlMonth = "03"
					}
				case "leorar":
					if dlDayInt > 12 {
						dlDay = strconv.Itoa(dlDayInt - 12)
						dlMonth = "05"
					} else {
						dlDay = strconv.Itoa(dlDayInt + 18)
						dlMonth = "04"
					}
				case "mysund":
					if dlDayInt > 6 {
						dlDay = strconv.Itoa(dlDayInt - 6)
						dlMonth = "06"
					} else {
						dlDay = strconv.Itoa(dlDayInt + 25)
						dlMonth = "05"
					}
				case "heisswerm":
					if dlDayInt > 31 {
						dlDay = strconv.Itoa(dlDayInt - 31)
						dlMonth = "08"
					} else {
						dlDay = strconv.Itoa(dlDayInt)
						dlMonth = "07"
					}
				case "largaheiss":
					if dlDayInt > 25 {
						dlDay = strconv.Itoa(dlDayInt - 25)
						dlMonth = "09"
					} else {
						dlDay = strconv.Itoa(dlDayInt + 6)
						dlMonth = "08"
					}
				case "pommois":
					if dlDayInt > 19 {
						dlDay = strconv.Itoa(dlDayInt - 19)
						dlMonth = "10"
					} else {
						dlDay = strconv.Itoa(dlDayInt + 11)
						dlMonth = "09"
					}
				case "spinnan":
					if dlDayInt > 13 {
						dlDay = strconv.Itoa(dlDayInt - 13)
						dlMonth = "11"
					} else {
						dlDay = strconv.Itoa(dlDayInt + 18)
						dlMonth = "10"
					}
				case "kalt":
					if dlDayInt > 7 {
						dlDay = strconv.Itoa(dlDayInt - 7)
						dlMonth = "12"
					} else {
						dlDay = strconv.Itoa(dlDayInt + 23)
						dlMonth = "11"
					}
				default:
					return "Problem with month. " + wrongSyntaxMessage
				}
			} else {
				switch dlMonth {
				case "january", "jan", "1", "01":
					dlMonth = "01"
					if dlDayInt > 31 {
						return "Problem with day. " + wrongSyntaxMessage
					}
				case "february", "feb", "2", "02":
					dlMonth = "02"
					if dlDayInt > 28 {
						return "Problem with day. " + wrongSyntaxMessage
					}
				case "march", "mar", "3", "03":
					dlMonth = "03"
					if dlDayInt > 31 {
						return "Problem with day. " + wrongSyntaxMessage
					}
				case "april", "apr", "4", "04":
					dlMonth = "04"
					if dlDayInt > 30 {
						return "Problem with day. " + wrongSyntaxMessage
					}
				case "may", "5", "05":
					dlMonth = "05"
					if dlDayInt > 31 {
						return "Problem with day. " + wrongSyntaxMessage
					}
				case "june", "jun", "6", "06":
					dlMonth = "06"
					if dlDayInt > 30 {
						return "Problem with day. " + wrongSyntaxMessage
					}
				case "july", "jul", "7", "07":
					dlMonth = "07"
					if dlDayInt > 31 {
						return "Problem with day. " + wrongSyntaxMessage
					}
				case "august", "aug", "8", "08":
					dlMonth = "08"
					if dlDayInt > 31 {
						return "Problem with day. " + wrongSyntaxMessage
					}
				case "september", "sep", "9", "09":
					dlMonth = "09"
					if dlDayInt > 30 {
						return "Problem with day. " + wrongSyntaxMessage
					}
				case "october", "oct", "10":
					dlMonth = "10"
					if dlDayInt > 31 {
						return "Problem with day. " + wrongSyntaxMessage
					}
				case "november", "nov", "11":
					dlMonth = "11"
					if dlDayInt > 30 {
						return "Problem with day. " + wrongSyntaxMessage
					}
				case "december", "dec", "12":
					dlMonth = "12"
					if dlDayInt > 31 {
						return "Problem with day. " + wrongSyntaxMessage
					}
				case "menotheen":
					if dlDayInt > 31 {
						dlDay = strconv.Itoa(dlDayInt - 31)
						dlMonth = "02"
					} else {
						dlDay = strconv.Itoa(dlDayInt)
						dlMonth = "01"
					}
				case "lengten":
					if dlDayInt > 23 {
						dlDay = strconv.Itoa(dlDayInt - 23)
						dlMonth = "03"
					} else {
						dlDay = strconv.Itoa(dlDayInt + 5)
						dlMonth = "02"
					}
				case "regen":
					if dlDayInt > 17 {
						dlDay = strconv.Itoa(dlDayInt - 17)
						dlMonth = "04"
					} else {
						dlDay = strconv.Itoa(dlDayInt + 14)
						dlMonth = "03"
					}
				case "leorar":
					if dlDayInt > 11 {
						dlDay = strconv.Itoa(dlDayInt - 11)
						dlMonth = "05"
					} else {
						dlDay = strconv.Itoa(dlDayInt + 19)
						dlMonth = "04"
					}
				case "mysund":
					if dlDayInt > 5 && dlDayInt < 36 {
						dlDay = strconv.Itoa(dlDayInt - 5)
						dlMonth = "06"
					} else if dlDayInt == 36 {
						dlDay = "01"
						dlMonth = "07"
					} else {
						dlDay = strconv.Itoa(dlDayInt + 26)
						dlMonth = "05"
					}
				case "heisswerm":
					if dlDayInt > 30 {
						dlDay = strconv.Itoa(dlDayInt - 30)
						dlMonth = "07"
					} else {
						dlDay = strconv.Itoa(dlDayInt + 1)
						dlMonth = "06"
					}
				case "largaheiss":
					if dlDayInt > 24 {
						dlDay = strconv.Itoa(dlDayInt - 24)
						dlMonth = "09"
					} else {
						dlDay = strconv.Itoa(dlDayInt + 7)
						dlMonth = "08"
					}
				case "pommois":
					if dlDayInt > 18 {
						dlDay = strconv.Itoa(dlDayInt - 18)
						dlMonth = "10"
					} else {
						dlDay = strconv.Itoa(dlDayInt + 12)
						dlMonth = "09"
					}
				case "spinnan":
					if dlDayInt > 12 {
						dlDay = strconv.Itoa(dlDayInt - 12)
						dlMonth = "11"
					} else {
						dlDay = strconv.Itoa(dlDayInt + 19)
						dlMonth = "10"
					}
				case "kalt":
					if dlDayInt > 6 {
						dlDay = strconv.Itoa(dlDayInt - 6)
						dlMonth = "12"
					} else {
						dlDay = strconv.Itoa(dlDayInt + 24)
						dlMonth = "11"
					}
				default:
					return "Problem with month. " + wrongSyntaxMessage
				}
			}

			// Lizian parsing
			//dlDayInt, err := strconv.Atoi(dlDay)
			//dlYearInt, err := strconv.Atoi(dlYear)
			//lmonth, lday := ltime.GetDayMonth(dlYearInt, dlMonth, dlDayInt)

			// military time conversion
			hasColon := true
			if colonIndex == 0 {
				hasColon = false
			}
			var dlTimeHours string
			var dlTimeMins string
			var dlm string
			var dlTZone string
			if hasColon {
				// standard operation
				dlTimeHours = lowered[thirdSpaceIndex+1 : colonIndex]
				dlTimeMins = lowered[colonIndex+1 : colonIndex+3]
				dlm = lowered[fourthSpaceIndex-2 : fourthSpaceIndex]
			} else {
				// military time conversion
				dlTimeHours = lowered[thirdSpaceIndex+1 : thirdSpaceIndex+3]
				dlTimeMins = lowered[thirdSpaceIndex+3 : fourthSpaceIndex]
			}
			isPM := false
			if dlm == "pm" {
				isPM = true
			}
			// PM time conversions
			if isPM {
				switch dlTimeHours {
				case "12":
					dlTimeHours = "12"
					dlm = "PM"
				case "01", "1":
					dlTimeHours = "01"
					dlm = "PM"
				case "02", "2":
					dlTimeHours = "02"
					dlm = "PM"
				case "03", "3":
					dlTimeHours = "03"
					dlm = "PM"
				case "04", "4":
					dlTimeHours = "04"
					dlm = "PM"
				case "05", "5":
					dlTimeHours = "05"
					dlm = "PM"
				case "06", "6":
					dlTimeHours = "06"
					dlm = "PM"
				case "07", "7":
					dlTimeHours = "07"
					dlm = "PM"
				case "08", "8":
					dlTimeHours = "08"
					dlm = "PM"
				case "09", "9":
					dlTimeHours = "09"
					dlm = "PM"
				case "10":
					dlTimeHours = "10"
					dlm = "PM"
				case "11":
					dlTimeHours = "11"
					dlm = "PM"
				default:
					return "Problem with time. " + wrongSyntaxMessage
				}
				// AM time conversion
			} else if !isPM && hasColon {
				switch dlTimeHours {
				case "12":
					dlTimeHours = "12"
					dlm = "AM"
				case "01", "1":
					dlTimeHours = "01"
					dlm = "AM"
				case "02", "2":
					dlTimeHours = "02"
					dlm = "AM"
				case "03", "3":
					dlTimeHours = "03"
					dlm = "AM"
				case "04", "4":
					dlTimeHours = "04"
					dlm = "AM"
				case "05", "5":
					dlTimeHours = "05"
					dlm = "AM"
				case "06", "6":
					dlTimeHours = "06"
					dlm = "AM"
				case "07", "7":
					dlTimeHours = "07"
					dlm = "AM"
				case "08", "8":
					dlTimeHours = "08"
					dlm = "AM"
				case "09", "9":
					dlTimeHours = "09"
					dlm = "AM"
				case "10":
					dlTimeHours = "10"
					dlm = "AM"
				case "11":
					dlTimeHours = "11"
					dlm = "AM"
				default:
					return "Problem with time. " + wrongSyntaxMessage
				}
				// military time conversion
			} else if !isPM && !hasColon {
				switch dlTimeHours {
				case "01":
					dlTimeHours = "01"
					dlm = "AM"
				case "02":
					dlTimeHours = "02"
					dlm = "AM"
				case "03":
					dlTimeHours = "03"
					dlm = "AM"
				case "04":
					dlTimeHours = "04"
					dlm = "AM"
				case "05":
					dlTimeHours = "05"
					dlm = "AM"
				case "06":
					dlTimeHours = "06"
					dlm = "AM"
				case "07":
					dlTimeHours = "07"
					dlm = "AM"
				case "08":
					dlTimeHours = "08"
					dlm = "AM"
				case "09":
					dlTimeHours = "09"
					dlm = "AM"
				case "10":
					dlTimeHours = "10"
					dlm = "AM"
				case "11":
					dlTimeHours = "11"
					dlm = "AM"
				case "12":
					dlTimeHours = "12"
					dlm = "PM"
				case "13":
					dlTimeHours = "01"
					dlm = "PM"
				case "14":
					dlTimeHours = "02"
					dlm = "PM"
				case "15":
					dlTimeHours = "03"
					dlm = "PM"
				case "16":
					dlTimeHours = "04"
					dlm = "PM"
				case "17":
					dlTimeHours = "05"
					dlm = "PM"
				case "18":
					dlTimeHours = "06"
					dlm = "PM"
				case "19":
					dlTimeHours = "07"
					dlm = "PM"
				case "20":
					dlTimeHours = "08"
					dlm = "PM"
				case "21":
					dlTimeHours = "09"
					dlm = "PM"
				case "22":
					dlTimeHours = "10"
					dlm = "PM"
				case "23":
					dlTimeHours = "11"
					dlm = "PM"
				case "24":
					dlTimeHours = "12"
					dlm = "AM"
				default:
					return "Problem with time. " + wrongSyntaxMessage
				}
			}
			// time zone conversions
			dlTZone = lowered[fourthSpaceIndex+1 : fifthSpaceIndex]
			switch dlTZone {
			case "acdt", "+1030":
				dlTZone = "+1030"
			case "acst", "+0930":
				dlTZone = "+0930"
			case "act", "−0500":
				dlTZone = "-0500"
			//case	"ACT"	ASEAN Common Time (proposed)
			//dlTZone = +0800
			case "acwst", "+0845": //Australian Central Western Standard Time (unofficial)	UTC+08:45
				dlTZone = "+0845"
			case "adt", "-0300": //Atlantic Daylight Time	UTC−03:00
				dlTZone = "-0300"
			case "aedt", "+1100": //Australian Eastern Daylight Saving Time	UTC+11:00
				dlTZone = "+1100"
			case "aest", "+1000": //Australian Eastern Standard Time	UTC+10:00
				dlTZone = "+1000"
			case "aft", "+0430": //Afghanistan Time	UTC+04:30
				dlTZone = "+0430"
			case "akdt", "-0800": //Alaska Daylight Time	UTC−08:00
				dlTZone = "-0800"
			case "akst", "-0900": //Alaska Standard Time	UTC−09:00
				dlTZone = "-0900"
			case "almt", "+0600": //Alma-Ata Time[1]	UTC+06:00
				dlTZone = "+0600"
			case "amst": //Amazon Summer Time (Brazil)[2]	UTC−03:00
				dlTZone = "-0300"
			case "amt":
				return "specify \"Amazon\" or \"Armenia\""
			case "amazon", "-0400": //(Brazil)[3]	UTC−04:00
				dlTZone = "-0400"
			case "armenia", "+0400": //UTC+04:00
				dlTZone = "+0400"
			case "anat", "+1200": //Anadyr Time[4]	UTC+12:00
				dlTZone = "+1200"
			case "aqtt", "+0500": //Aqtobe Time[5]	UTC+05:00
				dlTZone = "+0500"
			case "art": //Argentina Time	UTC−03:00
				dlTZone = "-0300"
			case "ast":
				return "please specify \"Arabia-Standard\", or \"Atlantic-Standard\""
			case "arabia-standard", "+0300": //Arabia Standard Time	UTC+03:00
				dlTZone = "+0300"
			case "atlantic-standard": //Atlantic Standard Time	UTC−04:00
				dlTZone = "-0400"
			case "awst", "+0800": //Australian Western Standard Time	UTC+08:00
				dlTZone = "+0800"
			case "azost", "+0000": //Azores Summer Time	UTC+00:00
				dlTZone = "+0000"
			case "azot", "-0100": //Azores Standard Time	UTC−01:00
				dlTZone = "-0100"
			case "azt": //Azerbaijan Time	UTC+04:00
				dlTZone = "0400"
			case "bnt": //Brunei Time	UTC+08:00
				dlTZone = "+0800"
			case "biot": //British Indian Ocean Time	UTC+06:00
				dlTZone = "+0600"
			case "bit", "-1200": //Baker Island Time	UTC−12:00
				dlTZone = "-1200"
			case "bot": //Bolivia Time	UTC−04:00
				dlTZone = "-0400"
			case "brst", "-0200": //Brasília Summer Time	UTC−02:00
				dlTZone = "-0200"
			case "brt": //Brasília Time	UTC−03:00
				dlTZone = "-0300"
			case "bst":
				return "specify \"Bangledesh\" or \"Bougainville\""
			case "bangledesh": //Bangladesh Standard Time	UTC+06:00
				dlTZone = "+0600"
			case "bougainville": //Bougainville Standard Time[6]	UTC+11:00
				dlTZone = "+1100"
			//case	"BST", "":	//British Summer Time (British Standard Time from Mar 1968 to Oct 1971)	UTC+01:00
			//dlTZone = ""
			case "btt": //Bhutan Time	UTC+06:00
				dlTZone = "+0600"
			case "cat", "+0200": //Central Africa Time	UTC+02:00
				dlTZone = "+0200"
			case "cct", "+0630": //Cocos Islands Time	UTC+06:30
				dlTZone = "+0630"
			case "cdt":
				return "please specify \"Central-Daylight\", or \"Cuba-Daylight\""
			case "central-daylight", "-0500": //Central Daylight Time (North America)	UTC−05:00
				dlTZone = "-0500"
			case "cuba-daylight": //Cuba Daylight Time[7]	UTC−04:00
				dlTZone = "-0400"
			case "cest": //Central European Summer Time	UTC+02:00
				dlTZone = "+0200"
			case "cet", "+0100": //Central European Time	UTC+01:00
				dlTZone = "+0100"
			case "chadt", "+1345": //Chatham Daylight Time	UTC+13:45
				dlTZone = "+1345"
			case "chast", "+1245": //Chatham Standard Time	UTC+12:45
				dlTZone = "+1245"
			case "chot": //Choibalsan Standard Time	UTC+08:00
				dlTZone = "+0800"
			case "chost", "+0900": //Choibalsan Summer Time	UTC+09:00
				dlTZone = "+0900"
			case "chst": //Chamorro Standard Time	UTC+10:00
				dlTZone = "+1000"
			case "chut": //Chuuk Time	UTC+10:00
				dlTZone = "+1000"
			case "cist": //Clipperton Island Standard Time	UTC−08:00
				dlTZone = "-0800"
			case "ckt", "-1000": //Cook Island Time	UTC−10:00
				dlTZone = "-1000"
			case "clst": //Chile Summer Time	UTC−03:00
				dlTZone = "-0300"
			case "clt": //Chile Standard Time	UTC−04:00
				dlTZone = "-0400"
			case "cost": //Colombia Summer Time	UTC−04:00
				dlTZone = "-0400"
			case "cot": //Colombia Time	UTC−05:00
				dlTZone = "-0500"
			case "cst":
				return "specify \"Central-Standard\", \"China-Standard\", or \"Cuba-Standard\""
			case "central-standard", "-0600": //Central Standard Time (Central America)	UTC−06:00
				dlTZone = "-0600"
			case "china-standard": //China Standard Time	UTC+08:00
				dlTZone = "+0800"
			case "cuba-standard": //Cuba Standard Time	UTC−05:00
				dlTZone = "-0500"
			case "cvt": //Cape Verde Time	UTC−01:00
				dlTZone = "-0100"
			case "cwst": //Central Western Standard Time (Australia) unofficial	UTC+08:45
				dlTZone = "+0845"
			case "cxt", "+0700": //Christmas Island Time	UTC+07:00
				dlTZone = "+0700"
			case "davt": //Davis Time	UTC+07:00
				dlTZone = "+0700"
			case "ddut": //Dumont d'Urville Time (in French Antarctic station)	UTC+10:00
				dlTZone = "+1000"
			case "dft": //AIX-specific equivalent of Central European Time[NB 1]	UTC+01:00
				dlTZone = "+0100"
			case "easst": //Easter Island Summer Time	UTC−05:00
				dlTZone = "-0500"
			case "east": //Easter Island Standard Time	UTC−06:00
				dlTZone = "-0600"
			case "eat": //East Africa Time	UTC+03:00
				dlTZone = "+0300"
			case "ect":
				return "please specify \"Eastern-Caribbean\", or \"Ecuador\"."
			case "eastern-caribbean": //Eastern Caribbean Time (does not recognise DST)	UTC−04:00
				dlTZone = "-0400"
			case "ecuador": //Ecuador Time	UTC−05:00
				dlTZone = "-0500"
			case "edt": //Eastern Daylight Time (North America)	UTC−04:00
				dlTZone = "-0400"
			case "eest": //Eastern European Summer Time	UTC+03:00
				dlTZone = "+0300"
			case "eet": //Eastern European Time	UTC+02:00
				dlTZone = "+0200"
			case "egst": //Eastern Greenland Summer Time	UTC+00:00
				dlTZone = "+0000"
			case "egt": //Eastern Greenland Time	UTC−01:00
				dlTZone = "-0100"
			case "est": //Eastern Standard Time (North America)	UTC−05:00
				dlTZone = "-0500"
			case "fet": //Further-eastern European Time	UTC+03:00
				dlTZone = "+0300"
			case "fjt": //Fiji Time	UTC+12:00
				dlTZone = "+1200"
			case "fkst": //Falkland Islands Summer Time	UTC−03:00
				dlTZone = "-0300"
			case "fkt": //Falkland Islands Time	UTC−04:00
				dlTZone = "-0400"
			case "fnt": //Fernando de Noronha Time	UTC−02:00
				dlTZone = "-0200"
			case "galt": //Galápagos Time	UTC−06:00
				dlTZone = "-0600"
			case "gamt": //Gambier Islands Time	UTC−09:00
				dlTZone = "-0900"
			case "get": //Georgia Standard Time	UTC+04:00
				dlTZone = "+0400"
			case "gft": //French Guiana Time	UTC−03:00
				dlTZone = "-0300"
			case "Gilt": //Gilbert Island Time	UTC+12:00
				dlTZone = "+1200"
			case "git": //Gambier Island Time	UTC−09:00
				dlTZone = "-0900"
			case "gmt": //Greenwich Mean Time	UTC+00:00
				dlTZone = "+0000"
			case "gst":
				return "please specify \"South-Georgia\", or \"Gulf-Standard\""
			case "south-georgia": //South Georgia and the South Sandwich Islands Time	UTC−02:00
				dlTZone = "-0200"
			case "gulf-standard": //Gulf Standard Time	UTC+04:00
				dlTZone = "+0400"
			case "gyt": //Guyana Time	UTC−04:00
				dlTZone = "-0400"
			case "hdt": //Hawaii–Aleutian Daylight Time	UTC−09:00
				dlTZone = "-0900"
			case "haec": //Heure Avancée d'Europe Centrale French-language name for CEST	UTC+02:00
				dlTZone = "+0200"
			case "hst": //Hawaii–Aleutian Standard Time	UTC−10:00
				dlTZone = "-1000"
			case "hkt": //Hong Kong Time	UTC+08:00
				dlTZone = "+0800"
			case "hmt": //Heard and McDonald Islands Time	UTC+05:00
				dlTZone = "+0500"
			case "hovst": //Hovd Summer Time (not used from 2017–present)	UTC+08:00
				dlTZone = "+0800"
			case "hovt": //Hovd Time	UTC+07:00
				dlTZone = "+0700"
			case "ict": //Indochina Time	UTC+07:00
				dlTZone = "+0700"
			case "idlw": //International Date Line West time zone	UTC−12:00
				dlTZone = "-1200"
			case "idt": //Israel Daylight Time	UTC+03:00
				dlTZone = "+0300"
			case "iot": //Indian Ocean Time	UTC+06:00
				dlTZone = "+0600"
			case "irdt": //Iran Daylight Time	UTC+04:30
				dlTZone = "+0430"
			case "irkt": //Irkutsk Time	UTC+08:00
				dlTZone = "+0800"
			case "irst", "+0330": //Iran Standard Time	UTC+03:30
				dlTZone = "+0330"
			case "ist":
				return "please specify \"Indian-Standard\", \"Irish-Standard\", or \"Israel-Standard\""
			case "indian-standard", "+0530": //Indian Standard Time	UTC+05:30
				dlTZone = "+0530"
			case "irish-standard": //Irish Standard Time[8]	UTC+01:00
				dlTZone = "+0100"
			case "isreal-standard": //Israel Standard Time	UTC+02:00
				dlTZone = "+0200"
			case "jst": //Japan Standard Time	UTC+09:00
				dlTZone = "+0900"
			case "kalt": //Kaliningrad Time	UTC+02:00
				dlTZone = "+0200"
			case "kgt": //Kyrgyzstan Time	UTC+06:00
				dlTZone = "+0600"
			case "kost": //Kosrae Time	UTC+11:00
				dlTZone = "+1100"
			case "krat": //Krasnoyarsk Time	UTC+07:00
				dlTZone = "+0700"
			case "kst": //Korea Standard Time	UTC+09:00
				dlTZone = "+0900"
			case "lhst":
				return "please specify \"Howe-Standard\", or \"Howe-Summer\""
			case "howe-standard": //Lord Howe Standard Time	UTC+10:30
				dlTZone = "+1030"
			case "howe-summer": //Lord Howe Summer Time	UTC+11:00
				dlTZone = "+1100"
			case "lint", "+1400": //Line Islands Time	UTC+14:00
				dlTZone = "+1400"
			case "magt": //Magadan Time	UTC+12:00
				dlTZone = "+1200"
			case "mart", "-0930": //Marquesas Islands Time	UTC−09:30
				dlTZone = "-0930"
			case "mawt": //Mawson Station Time	UTC+05:00
				dlTZone = "+0500"
			case "mdt": //Mountain Daylight Time (North America)	UTC−06:00
				dlTZone = "-0600"
			case "met": //Middle European Time (same zone as CET)	UTC+01:00
				dlTZone = "+0100"
			case "mest": //Middle European Summer Time (same zone as CEST)	UTC+02:00
				dlTZone = "+0200"
			case "mht": //Marshall Islands Time	UTC+12:00
				dlTZone = "+1200"
			case "mist": //Macquarie Island Station Time	UTC+11:00
				dlTZone = "+1100"
			case "mit": //Marquesas Islands Time	UTC−09:30
				dlTZone = "-0930"
			case "mmt": //Myanmar Standard Time	UTC+06:30
				dlTZone = "+0630"
			case "msk": //Moscow Time	UTC+03:00
				dlTZone = "+0300"
			case "mst":
				return "please specify \"Malaysian\", or \"Mountain-Standard\"."
			case "malaysian": //Malaysian Standard Time	UTC+08:00
				dlTZone = "+0800"
			case "mountain-standard", "-0700": //Mountain Standard Time (North America)	UTC−07:00
				dlTZone = "-0700"
			case "mut": //Mauritius Time	UTC+04:00
				dlTZone = "+0400"
			case "mvt": //Maldives Time	UTC+05:00
				dlTZone = "+0500"
			case "myt": //Malaysia Time	UTC+08:00
				dlTZone = "+0800"
			case "nct": //New Caledonia Time	UTC+11:00
				dlTZone = "+1100"
			case "ndt", "-0230": //Newfoundland Daylight Time	UTC−02:30
				dlTZone = "-0230"
			case "nft": //Norfolk Island Time	UTC+11:00
				dlTZone = "+1100"
			case "novt": //Novosibirsk Time [9]	UTC+07:00
				dlTZone = "+0700"
			case "npt", "+0545": //Nepal Time	UTC+05:45
				dlTZone = "+0545"
			case "nst", "-0330": //Newfoundland Standard Time	UTC−03:30
				dlTZone = "-0330"
			case "nt": //Newfoundland Time	UTC−03:30
				dlTZone = "-0330"
			case "nut", "-1100": //Niue Time	UTC−11:00
				dlTZone = "-1100"
			case "nzdt", "+1300": //New Zealand Daylight Time	UTC+13:00
				dlTZone = "+1300"
			case "nzdst": //New Zealand Daylight Saving Time	UTC+13:00
				dlTZone = "+1300"
			case "nzst": //New Zealand Standard Time	UTC+12:00
				dlTZone = "+1200"
			case "omst": //Omsk Time	UTC+06:00
				dlTZone = "+0600"
			case "orat": //Oral Time	UTC+05:00
				dlTZone = "+0500"
			case "pdt": //Pacific Daylight Time (North America)	UTC−07:00
				dlTZone = "-0700"
			case "pet": //Peru Time	UTC−05:00
				dlTZone = "-0500"
			case "pett": //Kamchatka Time	UTC+12:00
				dlTZone = "+1200"
			case "pgt": //Papua New Guinea Time	UTC+10:00
				dlTZone = "+100"
			case "phot": //Phoenix Island Time	UTC+13:00
				dlTZone = "+1300"
			case "pht": //Philippine Time	UTC+08:00
				dlTZone = "+0800"
			case "phst": //Philippine Standard Time	UTC+08:00
				dlTZone = "+0800"
			case "pkt": //Pakistan Standard Time	UTC+05:00
				dlTZone = "+0500"
			case "pmdt": //Saint Pierre and Miquelon Daylight Time	UTC−02:00
				dlTZone = "-0200"
			case "pmst": //Saint Pierre and Miquelon Standard Time	UTC−03:00
				dlTZone = "-0300"
			case "pont": //Pohnpei Standard Time	UTC+11:00
				dlTZone = "+1100"
			case "pst": //Pacific Standard Time (North America)	UTC−08:00
				dlTZone = "-0800"
			case "pwt": //Palau Time[11]	UTC+09:00
				dlTZone = "+0900"
			case "pyst": //Paraguay Summer Time[12]	UTC−03:00
				dlTZone = "-0300"
			case "pyt": //Paraguay Time[13]	UTC−04:00
				dlTZone = "-0400"
			case "ret": //Réunion Time	UTC+04:00
				dlTZone = "+0400"
			case "rott": //Rothera Research Station Time	UTC−03:00
				dlTZone = "-0300"
			case "sakt": //Sakhalin Island Time	UTC+11:00
				dlTZone = "+1100"
			case "samt": //Samara Time	UTC+04:00
				dlTZone = "+0400"
			case "sast": //South African Standard Time	UTC+02:00
				dlTZone = "0200"
			case "sbt": //Solomon Islands Time	UTC+11:00
				dlTZone = "+1100"
			case "sct": //Seychelles Time	UTC+04:00
				dlTZone = "+0400"
			case "sdt": //Samoa Daylight Time	UTC−10:00
				dlTZone = "-1000"
			case "sgt": //Singapore Time	UTC+08:00
				dlTZone = "+0800"
			case "slst": //Sri Lanka Standard Time	UTC+05:30
				dlTZone = "+0530"
			case "sret": //Srednekolymsk Time	UTC+11:00
				dlTZone = "+1100"
			case "srt": //Suriname Time	UTC−03:00
				dlTZone = "-0300"
			case "sst": //Samoa Standard Time	UTC−11:00
				dlTZone = "-1100"
			case "syot": //Showa Station Time	UTC+03:00
				dlTZone = "+0300"
			case "taht": //Tahiti Time	UTC−10:00
				dlTZone = "-1000"
			case "tha": //Thailand Standard Time	UTC+07:00
				dlTZone = "+0700"
			case "tft": //French Southern and Antarctic Time[14]	UTC+05:00
				dlTZone = "+0500"
			case "tjt": //Tajikistan Time	UTC+05:00
				dlTZone = "+0500"
			case "tkt": //Tokelau Time	UTC+13:00
				dlTZone = "+1300"
			case "tlt": //Timor Leste Time	UTC+09:00
				dlTZone = "+0900"
			case "tmt": //Turkmenistan Time	UTC+05:00
				dlTZone = "+0500"
			case "trt": //Turkey Time	UTC+03:00
				dlTZone = "+0300"
			case "tot": //Tonga Time	UTC+13:00
				dlTZone = "+1300"
			case "tst": //Taiwan Standard Time	UTC+08:00
				dlTZone = "+0800"
			case "tvt": //Tuvalu Time	UTC+12:00
				dlTZone = "+1200"
			case "ulast": //Ulaanbaatar Summer Time	UTC+09:00
				dlTZone = "+0900"
			case "ulat": //Ulaanbaatar Standard Time	UTC+08:00
				dlTZone = "+0800"
			case "utc": //Coordinated Universal Time	UTC+00:00
				dlTZone = "+0000"
			case "uyst": //Uruguay Summer Time	UTC−02:00
				dlTZone = "-0200"
			case "uyt": //Uruguay Standard Time	UTC−03:00
				dlTZone = "-0300"
			case "uzt": //Uzbekistan Time	UTC+05:00
				dlTZone = "+0500"
			case "vet": //Venezuelan Standard Time	UTC−04:00
				dlTZone = "-0400"
			case "vlat": //Vladivostok Time	UTC+10:00
				dlTZone = "+1000"
			case "volt": //Volgograd Time	UTC+03:00
				dlTZone = "+0300"
			case "vost": //Vostok Station Time	UTC+06:00
				dlTZone = "+0600"
			case "vut": //Vanuatu Time	UTC+11:00
				dlTZone = "+1100"
			case "wakt": //Wake Island Time	UTC+12:00
				dlTZone = "+1200"
			case "wast": //West Africa Summer Time	UTC+02:00
				dlTZone = "+0200"
			case "wat": //West Africa Time	UTC+01:00
				dlTZone = "+0100"
			case "west": //Western European Summer Time	UTC+01:00
				dlTZone = "+0100"
			case "wet": //Western European Time	UTC+00:00
				dlTZone = "+0000"
			case "wib": //Western Indonesian Time	UTC+07:00
				dlTZone = "+0700"
			case "wit": //Eastern Indonesian Time	UTC+09:00
				dlTZone = "+0900"
			case "wita": //Central Indonesia Time	UTC+08:00
				dlTZone = "+0800"
			case "wgst": //West Greenland Summer Time[15]	UTC−02:00
				dlTZone = "-0200"
			case "wgt": //West Greenland Time[16]	UTC−03:00
				dlTZone = "-0300"
			case "wst": //Western Standard Time	UTC+08:00
				dlTZone = "+0800"
			case "yakt": //Yakutsk Time	UTC+09:00
				dlTZone = "+0900"
			case "yekt": //Yekaterinburg Time	UTC+05:00
				dlTZone = "+0500"
			default:
				return "Problem with time zone. " + wrongSyntaxMessage
			}
			// concatenate dline
			if len(dlDay) == 1 {
				dlDay = "0" + dlDay
			}
			dline := dlMonth + " " + dlDay + " " + dlYear + " " + dlTimeHours + ":" + dlTimeMins + dlm + " " + dlTZone
			// character limits
			if len(alarmComment) > 100 || len(alarmName) > 50 {
				return "Name or comment too long, max 50 for name & 100 for comment"
			}

			alm := &alarm.Alarm{
				ChannelID: m.ChannelID,
				Deadline:  dline,
				Content:   alarmComment,
				Name:      alarmName,
				LoopFreq:  loopFreq,
				UserID:    m.Author.ID,
				ServerID:  m.GuildID,
			}
			// parse deadline
			parsedTime, err := time.Parse("01 02 2006 03:04PM -0700", dline)
			if err != nil {
				fmt.Println(err, "this one")
			}
			//convert deadline to unix
			unixTime := parsedTime.Unix()

			// sends message with unix time formatted for discord when alarm is set
			s.ChannelMessageSend(m.ChannelID, alarmName+" set for "+"<t:"+strconv.FormatInt(unixTime, 10)+":F>")
			UserManager.GetUser(m.Author.ID, s, ServerManager).AlarmManager.SetAlarm(alm, s, m.ChannelID)
			return "No response"
		} else {
			return wrongSyntaxMessage
		}
	}
	// list alarms
	switch lowered {
	case "list my alarms", "list alarms", "see alarms", "see my alarms":
		ListAlarms(m.Author.ID, m.ChannelID)
	default:
		return "No response"
	}

	// delete alarms
	if strings.HasPrefix(lowered, "delete alarm") {
		index, err := strconv.Atoi(lowered[13:])
		if err != nil {
			fmt.Println(err)
		}
		DeleteAlarm(m.Author.ID, m.ChannelID, index)
		return "Deleted."
	}
	// in case of empty message
	if lowered == "" {
		return "well you're awfully silent..."
	}
	// if no commands recognized then return empty string (triggers "I dont know that command" message)
	fmt.Println("response got")
	return ""

}

func DeleteAlarm(userid string, channnelid string, i int) {
	UserManager.GetUser(userid, s, ServerManager).AlarmManager.Alarms[i-1].Deadline = "01 02 2006 03:04PM -0700"
}
func ListAlarms(userid string, channelid string) {
	// alarmlist empty slice of strings
	alarmlist := []string{}
	Lalarmlist := []string{}
	//range alarms per user (message author)
	for i, v := range UserManager.GetUser(userid, s, ServerManager).AlarmManager.Alarms {
		// skip entries with emptied time
		if v.Deadline == "01 02 2006 03:04PM -0700" {
			continue
		}
		// parse deadline of current alarm in to Time Package format
		parsedTime, err := time.Parse("01 02 2006 03:04PM -0700", v.Deadline)
		if err != nil {
			fmt.Println(err)
		}
		// turns parsed time in to unix timestamp
		unixTime := parsedTime.Unix()
		// parse to Lizian Time
		dlYear := v.Deadline[6:10]
		gYear, err := strconv.Atoi(dlYear)
		if err != nil {
			fmt.Println(err)
		}
		dlMonth := v.Deadline[:2]
		var gMonth string
		switch dlMonth {
		case "01":
			gMonth = "January"
		case "02":
			gMonth = "February"
		case "03":
			gMonth = "March"
		case "04":
			gMonth = "April"
		case "05":
			gMonth = "May"
		case "06":
			gMonth = "June"
		case "07":
			gMonth = "July"
		case "08":
			gMonth = "August"
		case "09":
			gMonth = "September"
		case "10":
			gMonth = "October"
		case "11":
			gMonth = "November"
		case "12":
			gMonth = "December"
		}
		dlDay := v.Deadline[3:5]
		gDay, err := strconv.Atoi(dlDay)
		if err != nil {
			fmt.Println(err)
		}

		lizianMonth, lizianDay := ltime.GetDayMonth(gYear, gMonth, gDay)
		// concatenate alarmlist entry and append to alarmlist
		alarmlist = append(alarmlist, "- "+"-# **Alarm "+strconv.Itoa(i+1)+": "+v.Name+"** \n"+"<t:"+strconv.FormatInt(unixTime, 10)+":F>\n"+"-# \""+v.Content+"\" \n "+v.LoopFreq+"\n")
		Lalarmlist = append(Lalarmlist, "- "+"-# **Alarm "+strconv.Itoa(i+1)+": "+v.Name+"** \n"+"`"+ltime.GetDayOfWeek(lizianDay, lizianMonth)+", "+lizianMonth+" "+strconv.Itoa(lizianDay)+", "+dlYear+"` <t:"+strconv.FormatInt(unixTime, 10)+":t> \n"+"-# \""+v.Content+"\" \n "+v.LoopFreq+" \n")
	}
	// splits message in to 8 alarms per message
	hasheader := true
	var splitMessage func(aList []string)
	splitMessage = func(aList []string) {
		header := "## Alarms(Gregorian Time): \n"
		if len(aList) > 8 {
			var divMessage string
			for _, v := range aList[:8] {
				if hasheader {
					divMessage = header + divMessage + v
					hasheader = false
				} else {
					divMessage = divMessage + v
				}
			}
			s.ChannelMessageSend(channelid, divMessage)
			splitMessage(aList[8:])
		} else if len(aList) == 0 {
			/*divMessage := "No response"
			s.ChannelMessageSend(channelid, divMessage)*/
		} else {
			var divMessage string
			for _, v := range aList {
				if hasheader {
					divMessage = header + divMessage + v
					hasheader = false
				} else {
					divMessage = divMessage + v
				}
			}
			s.ChannelMessageSend(channelid, divMessage)
		}
	}
	// separate splitter for Lizian times (necessary for formatting)
	placeheader := true
	var splitMessageL func(aList []string)
	splitMessageL = func(aList []string) {
		header := "## Alarms(Lizian Time): \n"
		if len(aList) > 8 {
			var divMessage string
			for _, v := range aList[:8] {
				if placeheader {
					divMessage = header + divMessage + v
					placeheader = false
				} else {
					divMessage = divMessage + v
				}
			}
			s.ChannelMessageSend(channelid, divMessage)
			splitMessageL(aList[8:])
		} else if len(aList) == 0 {
			divMessage := "No alarms! :)"
			s.ChannelMessageSend(channelid, divMessage)
		} else {
			var divMessage string
			for _, v := range aList {
				if placeheader {
					divMessage = header + divMessage + v
					placeheader = false
				} else {
					divMessage = divMessage + v
				}
			}
			s.ChannelMessageSend(channelid, divMessage)
		}
	}

	splitMessageL(Lalarmlist)
	splitMessage(alarmlist)

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
				_, err = db.Exec("insert into alarms (name, time, comment, channelid, userid, serverid, loopfreq) values(?, ?, ?, ?, ?, ?, ?)",
					thisAlarm.Name,
					thisAlarm.Deadline,
					thisAlarm.Content,
					thisAlarm.ChannelID,
					thisAlarm.UserID,
					thisAlarm.ServerID,
					thisAlarm.LoopFreq,
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
func SendExpiredAlarms() {
	for _, server := range ServerManager.Servers {
		if server.MainChannel == "" {
			continue
		}
		alarmList := []string{}
		for _, alarm := range server.ExpiredAlarmManager.Alarms {
			alarmList = append(alarmList, "<@"+alarm.UserID+"> "+alarm.Name+" has gone off!\n")
			deadlineTime, err := time.Parse("01 02 2006 03:04PM -0700", alarm.Deadline)
			if err != nil {
				fmt.Println(err)
			}
			switch alarm.LoopFreq {
			case "daily":
				alarm.ChannelID = server.MainChannel
				alarm.Deadline = deadlineTime.AddDate(0, 0, 1).Format("01 02 2006 03:04PM -0700")
				go UserManager.GetUser(alarm.UserID, s, ServerManager).AlarmManager.SetAlarm(&alarm, s, server.MainChannel)
			case "weekly":
				alarm.ChannelID = server.MainChannel
				alarm.Deadline = deadlineTime.AddDate(0, 0, 7).Format("01 02 2006 03:04PM -0700")
				go UserManager.GetUser(alarm.UserID, s, ServerManager).AlarmManager.SetAlarm(&alarm, s, server.MainChannel)
			case "monthly":
				alarm.ChannelID = server.MainChannel
				alarm.Deadline = deadlineTime.AddDate(0, 1, 0).Format("01 02 2006 03:04PM -0700")
				go UserManager.GetUser(alarm.UserID, s, ServerManager).AlarmManager.SetAlarm(&alarm, s, server.MainChannel)
			case "yearly":
				alarm.ChannelID = server.MainChannel
				alarm.Deadline = deadlineTime.AddDate(1, 0, 0).Format("01 02 2006 03:04PM -0700")
				go UserManager.GetUser(alarm.UserID, s, ServerManager).AlarmManager.SetAlarm(&alarm, s, server.MainChannel)
			}
		}
		var splitMessage func(aList []string)

		splitMessage = func(aList []string) {
			if len(aList) > 8 {
				var divMessage string
				for _, v := range aList[:8] {
					divMessage = divMessage + v
				}
				s.ChannelMessageSend(server.MainChannel, divMessage)
				splitMessage(aList[8:])
			} else {
				var divMessage string
				for _, v := range aList {
					divMessage = divMessage + v
				}
				s.ChannelMessageSend(server.MainChannel, divMessage)
			}
		}

		splitMessage(alarmList)
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
		if err = rows.Scan(&dbid, &newalarm.Name, &newalarm.Deadline, &newalarm.Content, &newalarm.ChannelID, &newalarm.UserID, &newalarm.ServerID, &newalarm.LoopFreq); err != nil {
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
func DeleteAlarms() {
	db, err := sql.Open("sqlite3", "file:lizbot.db?cache=shared&_timeout=1000")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	_, err = db.Exec("delete from alarms")
	if err != nil {
		fmt.Println(err)
	}
}

var s *discordgo.Session

func main() {
	// open database
	db, err := sql.Open("sqlite3", "file:lizbot.db?cache=shared&_timeout=1000")
	if err != nil {
		fmt.Println(err)
	}
	// create table "servers" if it doesnt exist
	sqlstatement := "create table if not exists servers (serverid text not null primary key, mainchannel text)"
	_, err = db.Exec(sqlstatement)
	if err != nil {
		fmt.Println(err)
	}
	// create table "alarms" if it doesnt exist
	sqlstatement = "create table if not exists alarms (id integer not null primary key, name text, time text, comment text, channelid text, userid text, serverid text, loopfreq text)"
	_, err = db.Exec(sqlstatement)
	if err != nil {
		fmt.Println(err)
	}
	db.Close()

	// create session
	s, err = discordgo.New("Bot " + config.DISCORD_KEY)
	if err != nil {
		fmt.Printf("Invalid bot parameters: %v", err)
		return
	}

	//BORROWED CODE
	// Register the messageCreate func as a callback for MessageCreate events.
	s.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	s.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages | discordgo.IntentsGuildMembers | discordgo.IntentsAll

	// Open a websocket connection to Discord and begin listening.
	err = s.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	LoadServers()
	DeleteServers()
	LoadAlarms()
	SendExpiredAlarms()
	DeleteExpiredAlarms()
	DeleteAlarms()
	defer SaveAlarms()
	defer SaveServers()

	// borrowed code
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	// Cleanly close down the Discord session.
	s.Close()
	//END
}
