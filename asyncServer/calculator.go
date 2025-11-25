package main

import (
	"math"
	"math/rand"
	"strconv"
	"time"
)

type CoincidenceCalculator struct {
	random *rand.Rand
}

func NewCoincidenceCalculator() *CoincidenceCalculator {
	return &CoincidenceCalculator{
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// CalculatePotentialCoincidence рассчитывает вероятность совпадения
func (c *CoincidenceCalculator) CalculatePotentialCoincidence(composer *Composer, analysis *ComposerAnalysis) float64 {
	// Имитация задержки 5-10 секунд
	delay := 5 + c.random.Float64()*5
	time.Sleep(time.Duration(delay * float64(time.Second)))

	return c.calculateSimpleSimilarity(composer, analysis)
}

// calculateSimpleSimilarity простая формула схожести на основе частот интервалов
func (c *CoincidenceCalculator) calculateSimpleSimilarity(composer *Composer, analysis *ComposerAnalysis) float64 {
	var totalSimilarity float64
	var totalWeight float64

	// Функция для конвертации string в float64
	parseFloat := func(s *string) *float64 {
		if s == nil {
			return nil
		}
		if f, err := strconv.ParseFloat(*s, 64); err == nil {
			return &f
		}
		return nil
	}

	// Сопоставление интервалов
	intervals := []struct {
		composerStatIndex int
		analysisFreq      *string
		weight            float64
	}{
		{0, analysis.AnonUnisonsSecondsFreq, 0.25},
		{1, analysis.AnonThirdsFreq, 0.20},
		{2, analysis.AnonFourthsFifthsFreq, 0.20},
		{3, analysis.AnonSixthsSeventhsFreq, 0.20},
		{4, analysis.AnonOctavesFreq, 0.15},
	}

	for _, interval := range intervals {
		if len(composer.IntervalStats) > interval.composerStatIndex {
			compStat := composer.IntervalStats[interval.composerStatIndex]
			anaFreq := parseFloat(interval.analysisFreq)

			if compStat.Frequency != nil && anaFreq != nil {
				// Простая формула: 1 - относительное отклонение
				deviation := math.Abs(*anaFreq - *compStat.Frequency)
				maxPossibleDeviation := 100.0 // максимальное возможное отклонение в процентах
				similarity := 1.0 - (deviation / maxPossibleDeviation)

				// Учет веса интервала
				totalSimilarity += similarity * interval.weight
				totalWeight += interval.weight
			}
		}
	}

	if totalWeight > 0 {
		// Нормализуем и преобразуем в проценты
		result := (totalSimilarity / totalWeight) * 100
		// Ограничение и округление
		result = math.Max(0, math.Min(100, result))
		return math.Round(result*100) / 100
	}

	// Если данных нет, используем случайное значение
	return 30 + c.random.Float64()*40
}
