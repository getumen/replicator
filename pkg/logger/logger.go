package logger

//go:generate mockgen -source=$GOFILE -destination=mock$GOPACKAGE/mock_$GOFILE -package=mock$GOPACKAGE

// Logger is a logger interface
type Logger interface {
	Debug(message string)
	Info(message string)
	Warning(message string)
	Fatal(message string)
	WithFields(Fields) Logger
}

// Fields is log fields
type Fields map[string]interface{}
