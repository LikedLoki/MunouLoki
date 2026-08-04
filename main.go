package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"os"
)

const memoryFilePath string = "./memoryFile"

var emptyMemory = Memory{User: "", Loki: ""}
var emptyMemories = []Memory{}

type Memory struct {
	Loki string `json:"loki"`
	User string `json:"user"`
}

type NMunouLoki struct {
	Memories []Memory `json:"memories"`
	buffer   Memory
}

func getCharNgrams(text []rune, n int) []string {
	var ngrams = make([]string, 0, len(text)-n+1)
	for i := 0; i+n <= len(text); i++ {
		ngrams = append(ngrams, string(text[i:i+n]))
	}
	return ngrams
}
func calcNgrams(textI string, textII string, n int) float64 {
	var arrayI []string = getCharNgrams([]rune(textI), n)
	var arrayII []string = getCharNgrams([]rune(textII), n)
	if len(arrayI) == 0 || len(arrayII) == 0 {
		return 0
	}

	var setI = map[string]bool{}
	var setII = map[string]bool{}
	for _, value := range arrayI {
		setI[value] = true
	}
	for _, value := range arrayII {
		setII[value] = true
	}

	var count int
	for value := range setI {
		if setII[value] {
			count++
		}
	}

	return (2.0 * float64(count)) / float64(len(setI)+len(setII)) * 100
}
func selectElement(array []float64) []int {
	var currentMax float64
	var equalList []int
	for _, value := range array {
		if value > currentMax {
			currentMax = value
		}
	}
	for index, value := range array {
		if value == currentMax {
			equalList = append(equalList, index)
		}
	}
	return equalList
}

func readMemories() ([]Memory, error) {
	var memoriesJSON, errorI = os.ReadFile(memoryFilePath)
	if errorI != nil {
		fmt.Println(errorI)
		return []Memory{}, errorI
	}
	var memoriesStruct []Memory
	var errorII = json.Unmarshal(memoriesJSON, &memoriesStruct)
	if errorII != nil {
		fmt.Println(errorII)
		return emptyMemories, errorII
	}
	return memoriesStruct, nil
}
func writeMemories(memories []Memory) error {
	var MLJSON, errorI = json.Marshal(memories)
	if errorI != nil {
		return errorI
	}
	var errorII = os.WriteFile(memoryFilePath, MLJSON, 0644)
	if errorII != nil {
		return errorII
	}
	return nil
}

func (m *NMunouLoki) call(input string) string {
	if m.buffer.Loki == "" {
		m.buffer.Loki = input
	} else {
		m.buffer.User = input
	}
	if m.buffer.Loki != "" && m.buffer.User != "" {
		m.Memories = append(m.Memories, m.buffer)
		m.buffer = emptyMemory
	}
	if len(m.Memories) == 0 {
		return input
	} else {
		var listA []string
		for _, value := range m.Memories {
			listA = append(listA, value.Loki)
		}
		var output string
		if rand.IntN(15) < 1 {
			output = listA[rand.IntN(len(listA))]
		} else {
			var listB []float64
			for _, value := range listA {
				listB = append(listB, calcNgrams(value, input, 2))
			}
			var candidateIndexList = selectElement(listB)
			var selectedIndex = candidateIndexList[rand.IntN(len(candidateIndexList))]
			output = m.Memories[selectedIndex].User

		}
		return output
	}
}

func memoriesReset() {
	os.WriteFile(memoryFilePath, []byte("[]"), 0644)
	fmt.Println("Reset Complete | 初期化完了")
	linebreak()
}

func linebreak() {
	fmt.Print("\n")
}

func menuMode() {
	var scanner = bufio.NewScanner(os.Stdin)
menuLoop:
	for {
		fmt.Println("M====MENU====M")
		fmt.Println("* Enter the following | 以下のいずれかを入力して")
		fmt.Println("[chat]: Start the chat | 対話を始める")
		fmt.Println("[config]: Tweak the config | 設定をいじる")
		fmt.Println("[exit]: End | 終わる")
		fmt.Print("> ")
		if !scanner.Scan() {
			break menuLoop
		}
		var inputI = scanner.Text()
		switch inputI {
		case "chat":
			linebreak()
			fmt.Println("C====CHAT====C")
			fmt.Println("* [>>] is YOU. [<<] is LOKI")
			fmt.Println("* Exit for /bye | * 閉じるには /bye")
			var memory, errorI = readMemories()
			if errorI != nil {
				fmt.Println(errorI)
				break menuLoop
			}
			var munouLoki = NMunouLoki{Memories: memory, buffer: emptyMemory}
		chatLoop:
			for {
				fmt.Print("[>>] ")
				if !scanner.Scan() {
					break menuLoop
				}
				var input = scanner.Text()
				switch input {
				case "/bye":
					linebreak()
					var errorI = writeMemories(munouLoki.Memories)
					if errorI != nil {
						fmt.Println(errorI)
						break menuLoop
					}
					break chatLoop
				case "":
				default:
					fmt.Print("[<<] ")
					var response = munouLoki.call(input)
					fmt.Println(response)
				}
			}

		case "config":
			linebreak()
		configLoop:
			for {
				fmt.Println("====CONFIG====")
				fmt.Println("* Enter the following | 以下のいずれかを入力して")
				fmt.Println("[reset]: Reset of memories")
				fmt.Println("[exit]: Close the config | 設定を閉じる")
				fmt.Print("> ")
				if !scanner.Scan() {
					break menuLoop
				}
				var input = scanner.Text()
				switch input {
				case "reset":
					linebreak()
					memoriesReset()
					break configLoop
				case "exit":
					linebreak()
					break configLoop
				default:
					linebreak()
				}
			}

		case "exit":
			linebreak()
			break menuLoop
		default:
			linebreak()
		}
	}
	if scanner.Err() != nil {
		fmt.Println("ERROR EXIT...")
		return
	}
	fmt.Println("EXIT...")
}

func main() {
	menuMode()
}
