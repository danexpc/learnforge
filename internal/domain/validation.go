package domain

var (
	ValidModes  = []string{"lesson", "flashcards", "quiz"}
	ValidLevels = []string{"beginner", "intermediate", "advanced"}
)

func ValidateMode(mode string) bool {
	if mode == "" {
		return true
	}
	for _, m := range ValidModes {
		if mode == m {
			return true
		}
	}
	return false
}

func ValidateLevel(level *string) bool {
	if level == nil || *level == "" {
		return true
	}
	for _, l := range ValidLevels {
		if *level == l {
			return true
		}
	}
	return false
}

