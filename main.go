package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Question struct {
	Text    string
	Options []string
	Answer  int
}

type GameState struct {
	Player    string
	Points    int
	Questions []Question
}

func (g *GameState) Init() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Seja bem vindo(a) ao quiz")
	fmt.Print("Digite seu nome: ")

	name, err := reader.ReadString('\n')

	if err != nil {
		panic("Erro ao ler a string")
	}

	g.Player = name

	fmt.Printf("Vamos ao jogo, %s", g.Player)
}

func (g *GameState) ProcceessCSV() {
	f, err := os.Open("quiz.csv")

	if err != nil {
		panic("Erro ao ler arquivo")
	}

	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()

	if err != nil {
		panic("Erro ao ler csv")
	}

	for index, record := range records {
		if index > 0 {
			correctAnswer, _ := stringToInt(record[5])

			question := Question{
				Text:    record[0],
				Options: record[1:5],
				Answer:  correctAnswer,
			}

			g.Questions = append(g.Questions, question)
		}
	}
}

func (g *GameState) Run() {
	for index, question := range g.Questions {
		const YELLOW = "\033[33m"
		const RESET = "\033[0m"

		fmt.Printf("%s %d. %s %s\n", YELLOW, (index + 1), question.Text, RESET)

		for j, option := range question.Options {
			fmt.Printf("[%d] %s\n", (j + 1), option)
		}

		fmt.Print("\n> Digite uma alternativa: ")
		var answer int
		var err error

		for {
			reader := bufio.NewReader(os.Stdin)
			read, _ := reader.ReadString('\n')

			answer, err = stringToInt(read[:len(read)-1])

			if err != nil {
				fmt.Println(err.Error())
				fmt.Print("\n> Digite uma alternativa: ")
				continue
			}

			break
		}

		if answer == question.Answer {
			fmt.Println("Resposta correta! :)")
			g.Points += 10
		} else {
			fmt.Println("Ops! Resposta incorreta! :(")
		}

		fmt.Println("======================================================================")
	}
}

func main() {
	game := &GameState{}
	go game.ProcceessCSV()
	game.Init()
	game.Run()

	fmt.Println("\n                              FIM DE JOGO                             ")

	fmt.Printf("> SUA PONTUAÇÃO: %d pts\n", game.Points)
}

func stringToInt(s string) (int, error) {
	num, err := strconv.Atoi(s)

	if err != nil {
		return 0, errors.New("não é permitido caractere não numérico")
	}

	return num, nil
}
