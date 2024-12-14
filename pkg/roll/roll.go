package roll

import (
	"strconv"
	"time"

	"golang.org/x/exp/rand"
)

func RollDice(sides string) string {
	rand.Seed(uint64(time.Now().UnixNano()))

	switch sides {
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
