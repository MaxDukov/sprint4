package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

// parseTraining разбирает строку данных о количестве шагов, типе тренировки и продолжительности тренировки
func parseTraining(data string) (int, string, time.Duration, error) {
	vals := strings.Split(data, ",")
	if len(vals) != 3 {
		return 0, "", 0, errors.New("неверный формат данных")
	}
	steps, err := strconv.Atoi(vals[0])
	if err != nil {
		return 0, "", 0, errors.New("неверный формат данных количества шагов")
	}
	if steps <= 0 {
		return 0, "", 0, errors.New("количество шагов должно быть больше 0")
	}
	trainingType := vals[1]
	duration, err := time.ParseDuration(vals[2])
	if duration <= time.Duration(0) {
		return 0, "", 0, errors.New("продолжительность должна быть больше 0")
	}
	if err != nil {
		return 0, "", 0, errors.New("неверный формат данных продолжительности тренировки")
	}
	return steps, trainingType, duration, nil
}

// distance рассчитывает дистанцию, пройденную за тренировку
func distance(steps int, height float64) float64 {
	stepsLength := float64(steps) * stepLengthCoefficient * height
	return stepsLength / mInKm
}

// meanSpeed рассчитывает среднюю скорость, с которой была пройдена дистанция
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	return distance(steps, height) / duration.Hours()
}

// TrainingInfo возвращает информацию о тренировке
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, trainingType, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}
	if duration <= 0 || steps <= 0 || weight <= 0 || height <= 0 {
		return "", errors.New("неверные данные")
	}
	var calories float64
	switch trainingType {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
	speed := meanSpeed(steps, height, duration)
	distance := distance(steps, height)
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", trainingType, duration.Hours(), distance, speed, calories), nil
}

// RunningSpentCalories рассчитывает количество сожженных калорий при беге
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if duration <= 0 || steps <= 0 || weight <= 0 || height <= 0 {
		return 0, errors.New("неверные данные")
	}

	speed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	calories := (weight * speed * durationInMinutes) / minInH
	return calories, nil
}

// WalkingSpentCalories рассчитывает количество сожженных калорий при ходьбе
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if duration <= 0 || steps <= 0 || weight <= 0 || height <= 0 {
		return 0, errors.New("неверные данные")
	}
	durationInMinutes := duration.Minutes()
	speed := meanSpeed(steps, height, duration)
	calories := walkingCaloriesCoefficient * (weight * speed * durationInMinutes) / minInH
	return calories, nil
}
