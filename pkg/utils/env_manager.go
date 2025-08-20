package utils

import (
	"bufio"
	"fmt"
	"maps"
	"os"
	"path"
	"strings"
)

var (
	SystemEnvManager *EnvManager = nil
)

func init() {
	// Retrieve environment manager file
	sysEnvFile := Getenv("DEVBOX_ENV_FILE", fmt.Sprintf("%s/00-env-devbox.zsh", Getenv("ZSH_CUSTOM", path.Join(os.Getenv("HOME"), ".oh-my-zsh/custom"))))
	if sysEnvFile == "" {
		panic("DEVBOX_ENV_FILE is not set and the default path does not exist")
	}

	var envManager *EnvManager

	if envFileInfo, err := os.Stat(sysEnvFile); os.IsNotExist(err) {
		// If the file does not exist, create it
		if err := os.WriteFile(sysEnvFile, []byte{}, 0600); err != nil {
			panic(fmt.Sprintf("Failed to create env file: %v", err))
		}
	} else if err != nil {
		panic(fmt.Sprintf("Failed to check env file: %v", err))
	} else if envFileInfo.IsDir() {
		panic(fmt.Sprintf("Expected a file but found a directory: %s", sysEnvFile))
	}
	envManager = &EnvManager{
		file: sysEnvFile,
	}
	if err := envManager.parseEnvFile(); err != nil {
		panic(fmt.Sprintf("Failed to parse env file: %v", err))
	}

	SystemEnvManager = envManager
}

type EnvManager struct {
	file      string
	variables map[string]string
}

func (em *EnvManager) Set(variables map[string]string) error {
	unsetVariables := make(map[string]string, len(variables))
	for key, value := range variables {
		variableValue, exists := em.variables[key]
		if !(exists && variableValue == value) {
			unsetVariables[key] = value
		}
	}
	if err := em.AppendToEnvFile(unsetVariables); err != nil {
		return fmt.Errorf("failed to add environment variables to env file: %w", err)
	}
	// Update the in-memory map with the new variables
	maps.Copy(em.variables, unsetVariables)
	return nil
}

// parseEnvFile reads the environment variables from the bash file.
// It parses each line to check if it contains the format "export KEY=VALUE"
func (em *EnvManager) parseEnvFile() error {
	vars := make(map[string]string)

	// Open the file
	f, err := os.Open(em.file)
	if err != nil {
		return fmt.Errorf("failed to open env file: %w", err)
	}
	defer f.Close()

	// Read the file line by line
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if after, ok := strings.CutPrefix(line, "export "); ok {
			// Remove the "export " prefix
			line = after
			// Split the line into key and value
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				// Trim whitespace AND quotes (simple or double) from key and value.
				key := trimSpacesAndQuotes(parts[0])
				value := trimSpacesAndQuotes(parts[1])

				vars[key] = value
			}
		}
	}
	em.variables = vars
	return scanner.Err()
}

// AppendToEnvFile appends new environment variables to the specified file.
// It checks if the file ends with a newline and adds one if it doesn't.
// The variables are added in the format "export KEY=VALUE".
func (em *EnvManager) AppendToEnvFile(envVars map[string]string) error {
	f, err := os.OpenFile(em.file, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("failed to open env file for appending: %w", err)
	}
	defer f.Close()

	// Check if the file does not end with a newline
	stat, err := f.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat env file: %w", err)
	}
	if stat.Size() > 0 {
		buf := make([]byte, 1)
		_, err := f.ReadAt(buf, stat.Size()-1)
		if err != nil {
			return fmt.Errorf("failed to read last byte of env file: %w", err)
		}
		if buf[0] != '\n' {
			// Write a newline if the last character is not '\n'
			if _, err := f.WriteString("\n"); err != nil {
				return fmt.Errorf("failed to write newline to env file: %w", err)
			}
		}
	}

	// Write each environment variable in the format "export KEY=VALUE"
	for key, value := range envVars {
		line := fmt.Sprintf("export %s=\"%s\"\n", key, value)
		if _, err := f.WriteString(line); err != nil {
			return fmt.Errorf("failed to write new variable to env file: %w", err)
		}
	}
	return nil
}
