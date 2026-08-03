package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strings"
)

type MunouLoki struct {
	memories []string
}

func (m *MunouLoki) call(input string) {
	m.memories = append(m.memories, input)
	var memoryFile, errorI = os.OpenFile("./memoryFile", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if errorI != nil {
		fmt.Println(errorI)
		return
	}
	fmt.Fprintln(memoryFile, input)
	memoryFile.Close()
	var selected string
	for selected == "" {
		var random = rand.IntN(len(m.memories))
		selected = m.memories[random]
	}
	fmt.Println(selected)

}

func memoriesReset() {
	os.WriteFile("./memoryFile", []byte{}, 0644)
	fmt.Println("Reset Complete | 初期化完了")
	linebreak()
	menuMode()
}

func linebreak() {
	fmt.Print("\n")
}

func chatMode() {
	fmt.Println("C====CHAT====C")
	fmt.Println("* [>>] is YOU. [<<] is LOKI")
	fmt.Println("* Exit for /bye | * 閉じるには /bye")
	var memory, errorI = os.ReadFile("./memoryFile")
	if errorI != nil {
		fmt.Println(errorI)
		return
	}
	var fmtdMemories = strings.Split(string(memory), "\n")
	var munouLoki = MunouLoki{memories: fmtdMemories}
	for {
		fmt.Print("[>>] ")
		var input string
		fmt.Scanln(&input)
		switch input {
		case "/bye":
			linebreak()
			menuMode()
			return
		case "":
		default:
			fmt.Print("[<<] ")
			munouLoki.call(input)
		}
	}
}
func configMode() {
	for {
		fmt.Println("====CONFIG====")
		fmt.Println("* Enter the following | 以下のいずれかを入力して")
		fmt.Println("[reset]: Reset of memories")
		fmt.Println("[exit]: Close the config | 設定を閉じる")
		fmt.Print("> ")
		var input string
		fmt.Scanln(&input)
		switch input {
		case "reset":
			linebreak()
			memoriesReset()
			return
		case "exit":
			linebreak()
			menuMode()
			return
		default:
			linebreak()
		}
	}
}
func menuMode() {
	for {
		fmt.Println("M====MENU====M")
		fmt.Println("* Enter the following | 以下のいずれかを入力して")
		fmt.Println("[chat]: Start the chat | 対話を始める")
		fmt.Println("[config]: Tweak the config | 設定をいじる")
		fmt.Println("[exit]: End | 終わる")
		fmt.Print("> ")
		var input string
		fmt.Scanln(&input)
		switch input {
		case "chat":
			linebreak()
			chatMode()
			return
		case "config":
			linebreak()
			configMode()
			return
		case "exit":
			linebreak()
			return
		default:
			linebreak()
		}
	}
}

func main() {
	menuMode()
}
