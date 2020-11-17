# Alarm

**Модуль "Будильник"**

Пакет представляет собой простой будильник, который отправляет пустую структуру в канал в заданый момент времени, с учетом необходимой временной зоны.

Пример:

```
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	alarm "github.com/NuclearLouse/utilities-alarm"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	timeAlarm, err := time.Parse("15:04:05", "17:58:00")
	if err != nil {
		log.Fatalln("set time alarm:", err)
	}

	a, err := alarm.New(ctx, timeAlarm, "Europe/Chisinau")
	if err != nil {
		log.Fatalln("create alarm:", err)
	}
LOOP:
	for {
		select {
		case ta := <-a.A:
			fmt.Printf("Будильник сработал в : %s\n", ta.Format("15:04:05"))
			// делаю что нужно по будильнику и/или выхожу
			ctx.Done()
			break LOOP
		default:
			// тут что-то делаю до срабатывания будильника

		}
		time.Sleep(1 * time.Second)
	}

}

```
