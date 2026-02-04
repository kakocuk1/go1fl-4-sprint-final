package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) { // Принимает данные и выводит значения шагов, вида активности, продолжительности
	parts := strings.Split(data, ",") // Step 1: Разделить строку на слайс строк + проверить длинну
	if len(parts) != 3 {
		return 0, "", 0, errors.New("неверный формат данных")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, err
	}
	if steps <= 0 {
		return 0, "", 0, errors.New("количество шагов должно быть больше нуля")
	}

	activity := parts[1] // Step 2: указали откуда брать вид активности

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, err
	}
	if duration <= 0 {
		return 0, "", 0, errors.New("продолжительность должна бать больше нуля")
	}
	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 { // Возвращает дистанцию в Км.
	stepLength := height * stepLengthCoefficient // Step 1: Длинна шага
	distanceM := float64(steps) * stepLength     // Step 2: Дистанциях в метрах
	distanceKm := distanceM / mInKm              // Step 3: Дистанция в Км (можно ли объеденить это со следующим шагом?)
	return distanceKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 { // Возвращает среднюю скорость
	if duration <= 0 { // Step 1: Проверка продолжительности
		return 0
	}
	dist := distance(steps, height) // Step 2: Дистаниця в Км
	durationH := duration.Hours()   // Step 3: Продолжительность в часах(H)
	return dist / durationH         // Step 4: Средняя скорость
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) { // функция расчета калорий при беге
	if steps <= 0 { // Step 1: Проверяем каждый входной параметр
		return 0, errors.New("Значение шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, errors.New("Значение веса должно быть больше 0")
	}
	if height <= 0 {
		return 0, errors.New("Значение роста должно быть больше 0")
	}
	if duration <= 0 {
		return 0, errors.New("Значение продолжительности тренировки должно быть больше 0")
	}
	speed := meanSpeed(steps, height, duration) // Step 2: расчитываем ср.скорость км/ч

	durationInMinutes := duration.Minutes() // Step 3: продолжительно в минутах, для расчета калорий

	calories := (weight * speed * durationInMinutes) / minInH // Step 4: количество калорий, потраченных при беге

	return calories, nil

}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) { // Функция расчета калорий при ходьбе
	if steps <= 0 { // Step 1: Проверяем каждый входной параметр
		return 0, errors.New("Значение шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, errors.New("Значение веса должно быть больше 0")
	}
	if height <= 0 {
		return 0, errors.New("Значение роста должно быть больше 0")
	}
	if duration <= 0 {
		return 0, errors.New("Значение продолжительности тренировки должно быть больше 0")
	}
	speed := meanSpeed(steps, height, duration) // Step 2: расчитываем ср.скорость км/ч

	durationInMinutes := duration.Minutes() // Step 3: продолжительно в минутах, для расчета калорий

	calories := (weight * speed * durationInMinutes) / minInH // Step 4: количество калорий

	calories *= walkingCaloriesCoefficient // Step 5: Корректируем на коэфициент для ходьбы

	return calories, nil
}

func TrainingInfo(data string, weight, height float64) (string, error) { // Информация о тренировке
	steps, activity, duration, err := parseTraining(data) // Парсим входные данные
	if err != nil {
		log.Println(err)
		return "", err
	}
	durationInHours := duration.Hours() // Перевод в часы

	var (
		dist         float64
		averageSpeed float64
		calories     float64
	)

	switch activity { //разбираем для каждой активности
	case "Бег":
		dist = distance(steps, height)
		averageSpeed = meanSpeed(steps, height, duration)
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	case "Ходьба":
		dist = distance(steps, height)
		averageSpeed = meanSpeed(steps, height, duration)
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	result := fmt.Sprintf(
		"Тип тренировки: %s\n"+
			"Длительность: %.2f ч.\n"+
			"Дистанция: %.2f км.\n"+
			"Скорость: %.2f км/ч\n"+
			"Сожгли калорий: %.2f\n",
		activity,
		durationInHours,
		dist,
		averageSpeed,
		calories,
	)

	return result, nil
}
