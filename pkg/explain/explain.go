package explain

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func Explain(s *discordgo.Session, m *discordgo.MessageCreate, topic string) (string, error) {
	fmt.Println("Sure, I'll explain ", topic)
	switch topic {
	case "alarms", "alarm", "set alarm", "setting alarms", "setting alarm":
		header := "## Setting Alarms \n"
		cmdIntro := "**Alarms are set by using the command:** \n"
		cmdLN1 := "`?set alarm for ` \n"
		cmdLN1half := "followed by month, day, year, time, am or pm, timezone, alarm name, alarm description, and loop frequency (optional). \n"
		cmdLn2 := "- -# **Note: User input is not case sensitive and supports dynamic input for each required field, making it easier for you to remember commands! While our fields do take dynamic input the order of the information is non negotiable, with the exception of \"loop frequency\" as explained in the following section.**\n"
		almExample0 := "	  - -# For example this command: \n"
		almExample1 := "	```?set alarm for april 10 2025 4:20PM pst \"Alarm Name\" \"Description\"``` \n"
		almExample2 := "	  - -# will produce the same results as: \n"
		almExample3 := "	```?set alarm for 04 10th 25 4:20pm -0800 \"Alarm Name\" \"Description\"``` \n"
		header0 := "## Seeing Your Alarms \n"
		cmdLn3 := "**Alarms set by the user can be seen by using the command:** \n"
		cmdLn3half := "		`?list my alarms` \n"
		header1 := "## Deleting Your Alarms \n"
		cmdLn4 := "**Alarms may be deleted by using the command:** \n"
		cmdLn4half := "		`?Delete Alarm` \n"
		cmdln4close := "followed by the number of the alarm, which can be found by using the alarm list command above. Alarms are listed in the order they were created \n"
		emptyline := "\n"
		header2 := "## Loops \n"
		loopexp1 := "Loop frequency is an optional value, meaning you can ignore this field if you want, like the examples from **Alarm Basics** above- Loops can be set to \"daily\", \"weekly\", \"monthly\", \"yearly\", or \"none\". \n"
		loopexp2 := "- -# **Note: Any value of loop frequency that is not listed here will prevent an alarm from looping. You may also leave the loop value blank when setting an alarm.** \n"
		almExample4 := " 	- -# Example command with loop: \n"
		almExamp4half := "		```?set alarm for april 10 2025 4:20PM pst \"Alarm Name\" \"Description\" daily``` \n"
		header3 := "## Military Time \n"
		milTime1 := "Alarms may be set using military time. \n"
		milTime2 := "- -# **Note: Use of military time does not require colon in time or declaration of AM/PM, but does require timezones.** \n"
		almExample5 := " 	- -# Example command using military time: \n"
		almExample5half := "	 ```?set alarm for april 10 2025 0420 pst \"Alarm Name\" \"Description\" daily``` \n"

		explanation := header + cmdIntro + cmdLN1 + cmdLN1half + cmdLn2 + almExample0 + almExample1 + almExample2 + almExample3 + header0 + cmdLn3 + cmdLn3half + emptyline + header1 + cmdLn4 + cmdLn4half + cmdln4close + emptyline + header2 + loopexp1 + loopexp2 + almExample4 + almExamp4half + emptyline + header3 + milTime1 + milTime2 + almExample5 + almExample5half
		return explanation, nil
	/*case "lizdate", "liztime", "liz time", "lizian time", "ltime", "lizian date", "lizian calendar", "liz calendar":
	header := "## The Lizian Calendar \n"
	intro := "The Lizian calendar is effectively a resturcturing of the gregorian calendar where there are 10 months, with 36 and 37 days, alternating.\nEach month has 3 ten-day weeks, and one shorter week with the remaining 6 or 7 days.\nLeap years are handled at the end of the year, in Kalt, where the last week of the year has 8 days."
	compintro := "## Lizian months compared to Gregorian dates"
	lybanner := "### LEAP YEAR\n"
	lymenotheen := "Menotheen: Jan 1 - 31 + Feb 1-5\n"
	lylengten := "Lengten: Feb 6-29 + Mar 1-13\n"
	lyregen := "Regen: Mar 14-31 + Apr 1-18\n"
	lyleorar := "Leorar: Apr 19-30 + May 1-25\n"
	lymysund := "Mysund: May 26-31 + Jun 1-30\n"
	lymeisswerm := "Heisswerm: Jul 1-31 + Aug 1-6\n"
	lylargaheiss := "Largaheiss: Aug 7-31 + Sept 1-11\n"
	lypommois := "Pommois: Sept 12-30 + Oct 1-18\n"
	lyspinnan := "Spinnan: Oct 19-31 + Nov 1-23\n"
	lykalt := "Kalt: Nov 24-30 + Dec 1-31 \n"
	nlybanner := "### NON LEAP YEAR\n"
	nlymenotheen := "Menotheen: Jan 1 - 31 + Feb 1-5\n"
	nlylengten := "Lengten: Feb 6-28 + Mar 1-14\n"
	nlyregen := "Regen: Mar 15-31 + Apr 1-19\n"
	nlyleorar := "Leorar: Apr 20-30 + May 1-26\n"
	nlymysund := "Mysund: May 27-31 + Jun 1-30 + Jul 1\n"
	nlymeisswerm := "Heisswerm: Jul 2-31 + Aug 1-7\n"
	nlylargaheiss := "Largaheiss: Aug 8-31 + Sept 1-12\n"
	nlypommois := "Pommois: Sept 13-30 + Oct 1-19\n"
	nlyspinnan := "Spinnan: Oct 20-31 + Nov 1-24\n"
	nlykalt := "Kalt: Nov 25-30 + Dec 1-31\n"
	emptyline := "\n"
	cmdIntro := "**Gregorian dates can be converted to Lizian dates using the command:** \n"
	cmdLN1 := "`?lizdate ` \n"
	cmdLN1half := "followed by day, month, and year.\n"
	cmdLn2 := "- -# **Note: User input is not case sensitive and supports dynamic input for each required field, making it easier for you to remember commands! While our fields do take dynamic input the order of the information is non negotiable, with the exception of \"loop frequency\" as explained in the following section.**\n"
	almExample0 := "	  - -# For example this command: \n"
	almExample1 := "	```?set alarm for april 10 2025 4:20PM pst \"Alarm Name\" \"Description\"``` \n"
	almExample2 := "	  - -# will produce the same results as: \n"
	almExample3 := "	```?set alarm for 04 10th 25 4:20pm -0800 \"Alarm Name\" \"Description\"``` \n"
	header0 := "## Seeing Your Alarms \n"
	cmdLn3 := "**Alarms set by the user can be seen by using the command:** \n"
	cmdLn3half := "		`?list my alarms` \n"
	header1 := "## Deleting Your Alarms \n"
	cmdLn4 := "**Alarms may be deleted by using the command:** \n"
	cmdLn4half := "		`?Delete Alarm` \n" */
	default:
		return "I didn't understand the subject :/ try again?", fmt.Errorf("topic not found")
	}
	return "I didn't understand the subject :/ try again?", fmt.Errorf("")
}
