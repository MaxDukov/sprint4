package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

// parsePackage разбирает строку данных о количестве шагов и продолжительности прогулки
func parsePackage(data string) (int, time.Duration, error) {
	vals := strings.Split(data, ",")
	if len(vals) != 2 {
		return 0, 0, errors.New("неверный формат данных")
	}
	steps, err := strconv.Atoi(vals[0])
	if steps <= 0 {
		return 0, 0, errors.New("количество шагов должно быть больше 0")
	}
	if err != nil {
		return 0, 0, errors.New("неверный формат данных")
	}
	duration, err := time.ParseDuration(vals[1])
	if err != nil {
		return 0, 0, errors.New("неверный формат данных")
	}
	if duration <= time.Duration(0) {
		return 0, 0, errors.New("продолжительность должна быть больше 0")
	}
	return steps, duration, nil
}

// DayActionInfo возвращает информацию о количестве шагов, дистанции и затраченных калориях
func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println("неверный формат данных", err)
		return ""
	}
	if steps <= 0 {
		log.Println("количество шагов должно быть больше 0")
		return ""
	}

	distance := float64(steps) * stepLength / mInKm
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println("неверный формат данных", err)
		return ""
	}
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distance, calories)
}
