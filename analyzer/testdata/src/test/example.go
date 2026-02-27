package example

import (
	"context"
	"log/slog"
	"os"
)

type User struct {
	email    string
	password string
}

func SomeFunc() {

	slogLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	ctx := context.Background()

	slogLogger.Info("Debug message with uppercase")   // want "message should start with lowercase"
	slogLogger.Warn("Warning message with uppercase") // want "message should start with lowercase"

	slogLogger.Info("тестовое сообщение на русском") // want "invalid symbol"
	slogLogger.Warn("предупреждение на кириллице")   // want "invalid symbol"
	slogLogger.Error("ошибка с русскими символами")  // want "invalid symbol"

	slogLogger.Info("message with *special* characters!")        // want "invalid symbol"
	slogLogger.Warn("message with @#$%^&* symbols")              // want "invalid symbol"
	slogLogger.Debug("message with _underscores_ and -hyphens-") // want "invalid symbol"

	password := "secret123"
	apiKey := "sk-123456789"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"

	slogLogger.Info("user password: " + password) // want "invalid symbol" "sensetive data" "sensetive data {password}"
	slogLogger.Warn("api_key=" + apiKey)          // want "invalid symbol" "sensetive data" "sensetive data {apiKey}"

	slogLogger.Info("user authenticated",
		"password", password, // want "sensetive data"
		"api_key", apiKey, // want "invalid symbol" "sensetive data"
	)

	u := User{
		email:    "secret@email.com",
		password: "super-secret-123",
	}

	slogLogger.Info("user data",
		"user_email", u.email, // want "invalid symbol" "sensetive data {email}"
		"user_password", u.password, // want "invalid symbol" "sensetive data"
	)

	slogLogger.LogAttrs(ctx, slog.LevelInfo, "request with sensitive data",
		slog.String("authorization", "Bearer "+token), // want "sensetive data {token}"
		slog.String("x-api-key", apiKey),              // want "invalid symbol" "sensetive data" "sensetive data {apiKey}"
	)

	slogLogger.Info("debug message with lowercase")
	slogLogger.Warn("warning in english only")
	slogLogger.Info("message without special characters")
	slogLogger.Error("database connection timeout")

	slogLogger.Info("request processed",
		slog.Int("status_code", 200), // want "invalid symbol"
		slog.String("method", "GET"),
		slog.Float64("response_time_ms", 123.45), // want "invalid symbol"
		slog.Bool("cached", false),
		slog.Duration("processing_time", 1000000000), // want "invalid symbol"
		slog.Group("request_info", // want "invalid symbol"
			slog.String("path", "/api/users"), // want "invalid symbol" "sensetive data"
			slog.Int("content_length", 1024),  // want "invalid symbol"
		),
	)

	pPassword := &password
	logins := []string{"l1", "l2", "l3", "l4"}

	slogLogger.Debug("user activity",
		slog.Int("userid", 12345),                     // want "sensetive data"
		slog.String("action", "login"),                // want "sensetive data"
		slog.Any("test data", (*pPassword+logins[0])), // want "sensetive data {pPassword}" "sensetive data {logins}"
		slog.Any("test data 2", (logins[1:3])),        // want "sensetive data {logins}"
	)

	slogLogger.Info("user count updated",
		slog.Int("previouscount", 99),
		slog.Int("newcount", 100),
		slog.Int("difference", 1),
	)

	logData := func(data any) {}
	logData("test")
}
