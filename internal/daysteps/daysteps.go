package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	// Step 1: Разделить строку на слайс строк. Ориентрируемся на запятую
	part := strings.Split(data, ",")
	if len(part) != 2 { // Step 2: Проверяем на полноту ввода продолжительность + шаги
		return 0, 0, errors.New("Неправильынй формат даты")
	}
	steps, err := strconv.Atoi(part[0]) // Step 3: Преобразовать первый элемент слайса (количество шагов) в тип int, обработать ошибки, убедившись что шагов больше 0
	if err != nil {
		return 0, 0, err
	}
	if steps <= 0 {
		return 0, 0, errors.New("Значение шагов должно быть больше 0")
	}
	duration, err := time.ParseDuration(part[1])
	if err != nil {
		return 0, 0, err
	}
	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data) // Step 1: Получить данные о количестве шагов и продолжительности прогулки с помощью функции parsePackage()
	if err != nil {
		log.Printf("Ошибка парсинга данных: %v", err)
		return ""
	}
	if steps <= 0 { // Step 2: Проверка на шаги
		return ""
	}
	distanceM := float64(steps) * stepLength // Step 3: Дистанция в метрах

	distanceKm := distanceM / mInKm // Step 4: Перевод в км

	calories, err := spentcalories.WalkingSpentCalories(
		steps,
		weight,
		height,
		duration,
	)
	return fmt.Sprintf(
		"Количество шагов: %d. \nДистанция составила %.2f км. \nВы сожгли %.2f ккал.",
		steps,
		distanceKm,
		calories,
	)
}
