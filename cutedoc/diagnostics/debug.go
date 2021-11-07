package diagnostics

const debugEnabled = false

func Debug(fn func()) {
	if debugEnabled {
		fn()
	}
}
