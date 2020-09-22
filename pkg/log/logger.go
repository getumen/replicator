package log

//go:generate mockgen -source=$GOFILE -destination=mock$GOPACKAGE/mock_$GOFILE -package=mock$GOPACKAGE

// Logger is a logger interface
type Logger interface {
	Debug(message string, fields Fields)
	Info(message string, fields Fields)
	Warning(message string, fields Fields)
	Fatal(message string, fields Fields)
	Panic(message string, fields Fields)
}

// Fields is log fields
type Fields map[string]interface{}
