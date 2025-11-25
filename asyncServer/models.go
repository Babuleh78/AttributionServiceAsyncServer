package main

// ComposerAnalysisRequest запрос на расчет совпадения
type ComposerAnalysisRequest struct {
	ComposerAnalysisID int64  `json:"composer_analysis_id"`
	ComposerID         int64  `json:"composer_id"`
	AnalysisID         int64  `json:"analysis_id"`
	SecretKey          string `json:"secret_key"`
}

// ComposerAnalysisUpdate запрос на обновление совпадения
type ComposerAnalysisUpdate struct {
	PotentialCoincidence float64 `json:"potential_coincidence"`
}

// ComposerAnalysisResponse ответ для асинхронного callback
type ComposerAnalysisResponse struct {
	ComposerAnalysisID   int64   `json:"composer_analysis_id"`
	PotentialCoincidence float64 `json:"potential_coincidence"`
	SecretKey            string  `json:"secret_key"`
}

// ComposerAnalysisSyncResponse синхронный ответ с результатом расчета
type ComposerAnalysisSyncResponse struct {
	ComposerAnalysisID   int64   `json:"composer_analysis_id"`
	PotentialCoincidence float64 `json:"potential_coincidence"`
	SecretKey            string  `json:"secret_key"`
	Status               string  `json:"status"`
}

// ComposerAnalysis представляет связь м-м
type ComposerAnalysis struct {
	ID                     int64   `json:"id"`
	ComposerID             int64   `json:"composer_id"`
	AnalysisID             int64   `json:"analysis_id"`
	AnonUnisonsSecondsFreq *string `json:"anon_unisons_seconds_freq"`
	AnonThirdsFreq         *string `json:"anon_thirds_freq"`
	AnonFourthsFifthsFreq  *string `json:"anon_fourths_fifths_freq"`
	AnonSixthsSeventhsFreq *string `json:"anon_sixths_sevenths_freq"`
	AnonOctavesFreq        *string `json:"anon_octaves_freq"`
	PotentialCoincidence   string  `json:"potential_coincidence"`
}

// IntervalStat представляет статистику интервала из interval_stats
type IntervalStat struct {
	IntervalGroup string   `json:"IntervalGroup"`
	Frequency     *float64 `json:"Frequency"`
	StdDev        *float64 `json:"StdDev"`
}

// Composer представляет модель композитора с вложенной структурой
type Composer struct {
	ID             int64          `json:"id"`
	Name           string         `json:"name"`
	Biography      *string        `json:"biography"`
	Image          *string        `json:"image"`
	AnalyzedWorks  int            `json:"analyzed_works"`
	TotalIntervals int            `json:"total_intervals"`
	Period         string         `json:"period"`
	PolyphonyType  string         `json:"polyphony_type"`
	IntervalStats  []IntervalStat `json:"interval_stats"`
}
