package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

const memoryFilePath string = "./memoryFile"

var emptyMemory = Memory{User: "", Loki: ""}
var emptyMemories = []Memory{}

type Memory struct {
	Loki string
	User string
}
type NMunouLoki struct {
	Memories []Memory `json:"memories"`
	buffer   Memory
}

func getNgrams(text string, n int) []string {
	var fmtedTxt = []rune(text)
	var length = len(fmtedTxt)
	if length < n {
		return []string{text}
	}
	var ngrams = make([]string, 0, length-n+1)
	for i := 0; i+n <= length; i++ {
		ngrams = append(ngrams, string(fmtedTxt[i:i+n]))
	}
	return ngrams
}
func calcCosine(textI string, textII string) float64 {
	// 2, 4, 6, 8のN-gramを使い二次元配列を作成
	var arrayI [][]string
	var arrayII [][]string
	for i := 2; i <= 8; i += 2 {
		arrayI = append(arrayI, getNgrams(textI, i))
		arrayII = append(arrayII, getNgrams(textII, i))
	}

	// 何もなかったなら切り返し
	if len(arrayI) == 0 || len(arrayII) == 0 {
		return 0
	}

	// つくった二次元配列をそれぞれ比較
	var list []float64
	for i := 0; i < 4; i++ {
		var mapI = map[string]float64{}
		var mapII = map[string]float64{}
		for _, value := range arrayI[i] {
			mapI[value]++
		}
		for _, value := range arrayII[i] {
			mapII[value]++
		}
		var dot float64
		for value, countI := range mapI {
			if countII, ok := mapII[value]; ok {
				dot += countI * countII
			}
		}
		var numberI, numberII float64
		for _, n := range mapI {
			numberI += math.Pow(n, 2)
		}
		for _, n := range mapII {
			numberII += math.Pow(n, 2)
		}
		var resultN = math.Sqrt(numberI) * math.Sqrt(numberII)
		if resultN == 0 {
			list = append(list, 0.0)
		}
		list = append(list, dot/resultN)
	}
	var sum float64
	for _, value := range list {
		sum += value
	}
	var average = sum / 4
	return average
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

/*
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
*/

func readMemories(database *sql.DB) ([]Memory, error) {
	var createTableSQL = `CREATE TABLE IF NOT EXISTS memories ("id" INTEGER PRIMARY KEY AUTOINCREMENT, "loki" TEXT, "user" TEXT, "date" TEXT);`
	var _, errorI = database.Exec(createTableSQL)
	if errorI != nil {
		return nil, errorI
	}

	var selectSQL = `SELECT loki, user FROM memories`
	var rows, errorII = database.Query(selectSQL)
	if errorII != nil {
		return nil, errorII
	}
	defer rows.Close()

	var memories []Memory
	for rows.Next() {
		var m Memory
		var errorII = rows.Scan(&m.Loki, &m.User)
		if errorII != nil {
			return nil, errorII
		}
		memories = append(memories, m)
	}

	if errorII := rows.Err(); errorII != nil {
		return nil, errorII
	}

	return memories, nil
}
func writeMemories(database *sql.DB, memory Memory) error {
	date := time.Now().Format("2006-01-02|15:04:05")
	var insertSQL = `INSERT INTO memories (loki, user, date) VALUES (?, ?, ?)`
	var _, errorII = database.Exec(insertSQL, memory.Loki, memory.User, date)
	if errorII != nil {
		return errorII
	}
	return nil
}

func decisionFn(m *NMunouLoki, input string, output *string) {
	if len(m.Memories) == 0 {
		*output = input
	} else {
		var listA []string
		for _, value := range m.Memories {
			listA = append(listA, value.Loki)
		}
		if rand.IntN(15) < 1 {
			*output = listA[rand.IntN(len(listA))]
		} else {
			var listB []float64
			for _, value := range listA {
				listB = append(listB, calcCosine(value, input))
			}
			var candidateIndexList = selectElement(listB)
			var selectedIndex = candidateIndexList[rand.IntN(len(candidateIndexList))]
			*output = m.Memories[selectedIndex].User

		}
	}
}

func (m *NMunouLoki) call(database *sql.DB, input string) string {
	var output string
	decisionFn(m, input, &output)
	if m.buffer.Loki != "" {
		m.buffer.User = input
		m.Memories = append(m.Memories, m.buffer)
		writeMemories(database, m.buffer)
	}
	m.buffer = Memory{
		Loki: output,
		User: "",
	}
	return output
}

func memoriesReset(database *sql.DB) {
	database.Exec("DROP TABLE memories")
	fmt.Println("Reset Complete | 初期化完了")
	linebreak()
}

func linebreak() {
	fmt.Print("\n")
}

func menuMode() {
	var database, errorI = sql.Open("sqlite", "./memoryFile")
	if errorI != nil {
		fmt.Println(errorI)
		return
	}
	defer database.Close()
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
			fmt.Println("* \x1b[38;2;255;192;192m[>>] is YOU.\x1b[0m [<<] is LOKI")
			fmt.Println("* Exit for /bye | * 閉じるには /bye")
			var memory, errorII = readMemories(database)
			if errorII != nil {
				fmt.Println(errorII)
				break menuLoop
			}
			var munouLoki = NMunouLoki{Memories: memory, buffer: emptyMemory}
		chatLoop:
			for {
				fmt.Print("\x1b[38;2;255;192;192m[>>] ")
				if !scanner.Scan() {
					break menuLoop
				}
				fmt.Print("\x1b[0m")
				var input = scanner.Text()
				switch input {
				case "/bye":
					linebreak()
					break chatLoop
				case "":
				default:
					fmt.Print("[<<] ")
					var response = munouLoki.call(database, input)
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
					memoriesReset(database)
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
