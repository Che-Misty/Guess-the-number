package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

type game struct {
	Id        int
	Mode      string
	Attemps   int
	MaxNumber int
}

var gamemode = []game{
	{1, "Easy", 15, 50},
	{2, "Medium", 10, 100},
	{3, "Hard", 5, 200},
}

func main() {
	defer handelPanic()
	for {
		fmt.Println("Выбери уровень сложности: \n1: ", gamemode[0].Mode, "\n2: ", gamemode[1].Mode, "\n3: ", gamemode[2].Mode)
		fmt.Println("Введите число!")

		var input, maxNumber, attemps int
		checkInput(&input, 1, 3)
		maxNumber, attemps = gameChoice(input)
		aimNumber := rand.Intn(maxNumber) + 1

		loader()

		var lastAttemps []int

		for i := 0; i < attemps; i++ {
			var guess int
			color.Yellow("Попытка %d: Введите число: ", i+1)
			checkInput(&guess, 1, maxNumber)
			lastAttemps = append(lastAttemps, guess)
			fmt.Println(lastAttemps)
			if result, ok := checkGuess(guess, aimNumber, i+1, attemps); ok {
				saveData(i+1, result, input)
				break
			}
		}

		fmt.Println("\nСыграть снова? Для продолжения нажмите клавишу (y)")
		var isAgain string
		fmt.Scanln(&isAgain)
		if isAgain != "y" {
			break
		}
	}
}

func gameChoice(input int) (maxNumber, attemps int) {
	switch input {
	case 1:
		color.Cyan("Вы выбрали легкий уровень сложности")
	case 2:
		color.Cyan("Вы выбрали средний уровень сложности")
	case 3:
		color.Cyan("Вы выбрали сложный уровень сложности")
	default:
		color.Cyan("Неверный ввод, попробуйте снова")
	}
	maxNumber = gamemode[input-1].MaxNumber
	attemps = gamemode[input-1].Attemps
	color.Yellow("У вас есть %d попыток, чтобы угадать число от 1 до %d\n", attemps, maxNumber)
	return maxNumber, attemps
}

func checkGuess(guess int, target int, attempt int, maxAttempts int) (bool, bool) {
	if math.Abs(float64(guess-target)) <= 5 {
		color.Yellow("Горячо")
	} else if math.Abs(float64(guess-target)) > 5 && math.Abs(float64(guess-target)) <= 15 {
		color.Yellow("Тепло")
	} else {
		color.Yellow("Холодно")
	}
	if guess < target {
		color.Yellow("Секретное число больше👆")
	} else if guess > target {
		color.Yellow("Секретное число меньше👇")
	} else {
		color.Green("Поздравляем! Вы угадали число!🎉 q(≧▽≦q)")
		return true, true
	}
	if attempt == maxAttempts {
		color.Red("Вы проиграли ━┳━　━┳━")
		color.Red("Загаданное число: %d", target)
		return false, true
	}
	return false, false
}

func checkInput(input *int, minInt, maxInt int) {
	for {
		_, err := fmt.Scanln(input)
		if err != nil {
			color.Cyan("Введите число")
		}
		if *input <= maxInt && *input >= minInt {
			return
		}
		color.Cyan("Значение находится вне диапазона! Пожалуйста введите другое значение.")
	}
}

func handelPanic() {
	if r := recover(); r != nil {
		fmt.Println("Паника обработана: ", r)
	}
}

func loader() {
	for i := range 5 {
		left := strings.Repeat("_", i+1)
		done := strings.Repeat(" ", 5-i-1)
		percentage := (i + 1) * 20
		fmt.Printf("\r Loading... [%v%v] %d%%", left, done, percentage)
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Println("\nСекретное число загадано! Удачи в угадывании!")
}

func saveData(attempt int, result bool, id int) {
	data, err := os.ReadFile("results.json")
	if err != nil {
		if os.IsNotExist(err) {
			color.Yellow("Файл не найден, будет создан новый файл.")
			os.Create("results.json")
		} else {
			color.Red("Ошибка при чтении файла: %v", err)
			return
		}
	}

	var existingData []map[string]any
	if len(data) > 0 {
		err = json.Unmarshal(data, &existingData)
		if err != nil {
			color.Red("Ошибка при разборе JSON: %v", err)
			return
		}
	}
	var isWin string
	if result {
		isWin = "Победа"
	} else {
		isWin = "Проигрыш"
	}
	existingData = append(existingData, map[string]any{
		"Количество попыток": attempt,
		"Режим сложности":    gamemode[id-1].Mode,
		"Итог":               isWin,
		"Дата":               time.Now().Format("2006-01-02 15:04:05"),
	})
	data, err = json.MarshalIndent(existingData, "", "  ")
	os.WriteFile("results.json", data, 0644)
	if err != nil {
		color.Red("Ошибка при кодировании JSON: %v", err)
		return
	}
}
