package logIPFS

var logSingleton *logger

func Info(msg interface{}) {
	check()

	logSingleton.info(msg)
}

func Warn(msg interface{}) {
	check()

	logSingleton.warn(msg)
}

func Error(msg interface{}) {
	check()

	logSingleton.error(msg)
}

func check() {
	if logSingleton != nil {
		return
	}

	logSingleton = &logger{
		System: newLogger("system"),
	}
}
