package files

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

type awsProfile struct {
	name   string
	config map[string]string
}

func GetAWSProfileCredentials(profile string) (map[string]string, error) {
	profiles, err := readAWSCredentials(os.ExpandEnv("$HOME/.aws/credentials"))
	if err != nil {
		return nil, err
	}

	for _, p := range profiles {
		if p.name == profile {
			return p.config, nil
		}
	}

	return nil, errors.New("profile not found")
}

func readAWSCredentials(filepath string) ([]awsProfile, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	var profiles []awsProfile
	var currentProfile *awsProfile

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			if currentProfile != nil {
				profiles = append(profiles, *currentProfile)
			}
			profileName := strings.Trim(line, "[]")
			currentProfile = &awsProfile{
				name:   profileName,
				config: make(map[string]string),
			}
		} else if currentProfile != nil {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				currentProfile.config[key] = value
			}
		}
	}

	if currentProfile != nil {
		profiles = append(profiles, *currentProfile)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return profiles, nil
}
