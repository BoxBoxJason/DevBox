package envmanager

import (
	"bufio"
	"devbox/pkg/utils"
	"fmt"
	"maps"
	"os"
	"path"
	"strings"

	"go.uber.org/zap"
)

var (
	DEFAULT_SYS_ENV_FILE             = strings.TrimSpace(utils.Getenv("DEVBOX_ENV_FILE", fmt.Sprintf("%s/00-env-devbox.zsh", utils.Getenv("ZSH_CUSTOM", path.Join(os.Getenv("HOME"), ".oh-my-zsh/custom")))))
	systemEnvManager     *EnvManager = nil
)

func SystemEnvManager(envFile string) *EnvManager {
	if systemEnvManager != nil {
		return systemEnvManager
	}

	// Retrieve environment manager file
	if envFile == "" {
		zap.L().Fatal("DEVBOX_ENV_FILE is set to an empty string, please set it to a valid file path")
	}

	var envManager *EnvManager

	if envFileInfo, err := os.Stat(envFile); os.IsNotExist(err) {
		// If the file does not exist, create it
		if err := os.WriteFile(envFile, []byte{}, 0600); err != nil {
			zap.L().Fatal("Failed to create env file", zap.String("file", envFile), zap.Error(err))
		}
	} else if err != nil {
		zap.L().Fatal("Failed to stat env file", zap.String("file", envFile), zap.Error(err))
	} else if envFileInfo.IsDir() {
		zap.L().Fatal("DEVBOX_ENV_FILE points to a directory, expected a file", zap.String("file", envFile))
	} else if envFileInfo.Mode()&os.ModeType != 0 {
		zap.L().Fatal("DEVBOX_ENV_FILE points to a special file, expected a regular file", zap.String("file", envFile))
	} else if envFileInfo.Mode()&0600 == 0 {
		zap.L().Fatal("DEVBOX_ENV_FILE does not have the correct permissions, expected 0600", zap.String("file", envFile))
	}
	envManager = &EnvManager{
		file: envFile,
	}
	if err := envManager.parseEnvFile(); err != nil {
		zap.L().Fatal("Failed to parse env file", zap.String("file", envFile), zap.Error(err))
	}

	return envManager
}

type EnvManager struct {
	file      string
	variables map[string]string
}

func (em *EnvManager) Set(variables map[string]string) []error {
	unsetVariables := make(map[string]string, len(variables))
	for key, value := range variables {
		variableValue, exists := em.variables[key]
		if !exists || variableValue != value {
			unsetVariables[key] = value
		}
	}
	if len(unsetVariables) > 0 {
		// Append to file
		if errs := em.AppendToEnvFile(unsetVariables); len(errs) > 0 {
			return errs
		}

		// Ensure in-memory map exists
		if em.variables == nil {
			em.variables = make(map[string]string)
		}

		// Update the in-memory map with the new variables
		maps.Copy(em.variables, unsetVariables)

		// Also set them in the process environment
		var setErrs []error
		for k, v := range unsetVariables {
			if err := os.Setenv(k, os.ExpandEnv(v)); err != nil {
				setErrs = append(setErrs, fmt.Errorf("failed to set env %s: %w", k, err))
			}
		}
		if len(setErrs) > 0 {
			return setErrs
		}
	}
	return nil
}

// parseEnvFile reads the environment variables from the bash file.
// It parses each line to check if it contains the format "export KEY=VALUE"
func (em *EnvManager) parseEnvFile() error {
	zap.L().Debug("Parsing environment file", zap.String("file", em.file))
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
				key := utils.TrimSpacesAndQuotes(parts[0])
				value := utils.TrimSpacesAndQuotes(parts[1])

				vars[key] = value
			}
		}
	}
	em.variables = vars
	zap.L().Debug("Found existing devbox environment variables", zap.Int("count", len(vars)))
	return scanner.Err()
}

// AppendToEnvFile appends new environment variables to the specified file.
// It checks if the file ends with a newline and adds one if it doesn't.
// The variables are added in the format "export KEY=VALUE".
func (em *EnvManager) AppendToEnvFile(envVars map[string]string) []error {
	zap.L().Info("Appending environment variables to file", zap.String("file", em.file), zap.Any("variables", envVars))

	fileEndsWithNewline, err := utils.FileEndsWithNewline(em.file)
	if err != nil {
		return []error{fmt.Errorf("failed to check if env file ends with newline: %w", err)}
	}

	// Open the file for appending and add a newline
	wf, err := os.OpenFile(em.file, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return []error{fmt.Errorf("failed to open env file for appending: %w", err)}
	}
	defer wf.Close()

	if !fileEndsWithNewline {
		if _, err := wf.WriteString("\n"); err != nil {
			return []error{fmt.Errorf("failed to write newline to env file: %w", err)}
		}
	}

	errorChan := make(chan error, len(envVars))
	// Write each environment variable in the format "export KEY=VALUE"
	for key, value := range envVars {
		line := fmt.Sprintf("export %s=\"%s\"\n", key, value)
		if _, err := wf.WriteString(line); err != nil {
			errorChan <- fmt.Errorf("failed to write new variable to env file: %w", err)
		}
	}
	close(errorChan)
	return utils.MergeErrors(errorChan)
}
