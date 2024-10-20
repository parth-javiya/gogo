package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide a project name as an argument.")
	}
	projectName := os.Args[1]

	// Create base project directory
	err := os.Mkdir(projectName, 0755)
	if err != nil {
		log.Fatalf("Failed to create project directory: %v", err)
	}

	// Folder structure to create
	dirs := []string{
		filepath.Join("cmd", projectName), // Project name in cmd folder
		"internal/handlers",
		"internal/services",
		"internal/repository",
		"internal/models/api",
		"internal/models/db",
		"internal/middlewares",
		"internal/utils",
		"pkg/logger", // Logger folder in pkg
		"pkg/config", // Config folder in pkg
		"tests/unit",
		"tests/integration",
		"migrations",
		"docs",
	}

	// Create the directories
	for _, dir := range dirs {
		dirPath := filepath.Join(projectName, dir)
		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			log.Fatalf("Failed to create directory %s: %v", dirPath, err)
		}
	}

	// Create initial files
	createFile(filepath.Join(projectName, filepath.Join("cmd", projectName, "main.go")), mainGoContent(projectName))
	createFile(filepath.Join(projectName, ".env"), envFileContent()) // .env file
	createFile(filepath.Join(projectName, ".gitignore"), gitignoreContent())
	createFile(filepath.Join(projectName, "Makefile"), makefileContent(projectName))

	// Add logger package files
	createFile(filepath.Join(projectName, filepath.Join("pkg", "logger", "logger.go")), loggerGoContent())

	// Add config package files
	createFile(filepath.Join(projectName, filepath.Join("pkg", "config", "config.go")), configGoContent())

	// Initialize Git
	initGit(projectName)

	fmt.Printf("Project %s has been created successfully!\n", projectName)
}

// Function to create a file with given content
func createFile(filePath, content string) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Failed to create file %s: %v", filePath, err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		log.Fatalf("Failed to write to file %s: %v", filePath, err)
	}
}

// Initialize Git (but no commit or add)
func initGit(projectDir string) {
	cmd := exec.Command("git", "init")
	cmd.Dir = projectDir
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to initialize Git: %v", err)
	}
}

// Returns the content for .gitignore
func gitignoreContent() string {
	return `# Binaries for programs and plugins
*.exe
*.dll
*.so
*.dylib

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (vendor)
vendor/

# IDE/editor configurations
.idea/
.vscode/
*.swp
`
}

// Returns the content for main.go
func mainGoContent(projectName string) string {
	return fmt.Sprintf(`package main

import (
	"fmt"
	"log"

	"%s/pkg/config"
	"%s/pkg/logger"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize logger
	log, err := logger.NewLogger(cfg.LogFile)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %%v", err)
	}

	log.Info().Msg("Starting the application")
	fmt.Println("Server Port:", cfg.ServerPort)
}
`, projectName, projectName)
}

// Returns the content for .env file
func envFileContent() string {
	return `APP_NAME=myapi
SERVER_PORT=8080
LOG_FILE=logs/myapi.log
DB_USER=root
DB_PASSWORD=password
DB_HOST=localhost
DB_PORT=5432
DB_NAME=mydatabase
`
}

// Returns the content for Makefile
func makefileContent(projectName string) string {
	return fmt.Sprintf(`run:
	go run cmd/%s/main.go

test:
	go test ./...

migrate:
	migrate -path ./migrations -database $(DB_URL) up
`, projectName)
}

// Returns the content for pkg/logger/logger.go
func loggerGoContent() string {
	return `package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// NewLogger creates a new logger that logs to the specified file
func NewLogger(logFile string) (*zerolog.Logger, error) {
	file, err := os.Create(logFile)
	if err != nil {
		return nil, err
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	multi := zerolog.MultiLevelWriter(os.Stdout, file)

	logger := zerolog.New(multi).With().Timestamp().Logger()
	log.Logger = logger
	return &logger, nil
}
`
}

// Returns the content for pkg/config/config.go
func configGoContent() string {
	return `package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

// Config holds the configuration for the application
type Config struct {
	AppName    string ` + "`" + `mapstructure:"APP_NAME"` + "`" + `
	ServerPort string ` + "`" + `mapstructure:"SERVER_PORT"` + "`" + `
	LogFile    string ` + "`" + `mapstructure:"LOG_FILE"` + "`" + `
	DBUser     string ` + "`" + `mapstructure:"DB_USER"` + "`" + `
	DBPassword string ` + "`" + `mapstructure:"DB_PASSWORD"` + "`" + `
	DBHost     string ` + "`" + `mapstructure:"DB_HOST"` + "`" + `
	DBPort     string ` + "`" + `mapstructure:"DB_PORT"` + "`" + `
	DBName     string ` + "`" + `mapstructure:"DB_NAME"` + "`" + `
}

// LoadConfig reads the .env file and returns the application configuration
func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error loading .env file: %%v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Error unmarshalling config: %%v", err)
	}

	return &cfg
}
`
}
