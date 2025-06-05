package main

import (
	"errors"
	"fmt"
	"index_plov/internal/parser"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	telebot "gopkg.in/telebot.v4"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	pref := telebot.Settings{
		Token: os.Getenv("BOT_TOKEN"),
	}

	b, err := telebot.NewBot(pref)

	if err != nil {
		fmt.Println(err)
	}

	b.Handle("/start", func(ctx telebot.Context) error {
		message :=
			`Привет. Это бот для получения индекса нашего любимого плова!
Средняя цена 1 килограмма плова с учетом цен на продукты на рынках по данным Агентства статистики.
Данные доступны с 2021 года. Данные можно получить за год, за год и определенный месяц.
		
Формат должен быть "год месяц", "месяц', "год"`
		return ctx.Send(message)
	})

	b.Handle(telebot.OnText, func(ctx telebot.Context) error {
		message := ctx.Text()
		var rm strings.Builder
		result, err := parser.GetData(message)
		var pe *parser.ParseError
		if err != nil {
			if errors.As(err, &pe) {
				return ctx.Send(err)
			} else {
				return ctx.Send("Формат должен	 быть \"год месяц\", \"месяц\", \"год\"")
			}
		}

		switch v := result.(type) {
		case map[string]interface{}:
			rm.WriteString(message + "\n")
			for key, value := range v {
				fmt.Fprintf(&rm, "\t%s : %f\n", key, value)
			}
		case float64:
			fmt.Fprintf(&rm, "%s : %f", time.Now().Month().String(), v)
		default:
			rm.WriteString("Нет информации об этом месяце")
		}

		return ctx.Send(rm.String())
	})

	b.Start()
}
