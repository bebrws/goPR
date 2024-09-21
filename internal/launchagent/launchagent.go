package launchagent

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/bebrws/goPR/config"
	"github.com/bebrws/goPR/internal/di"
)

const plistTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Label</key>
	<string>{{.Label}}</string>
	<key>ProgramArguments</key>
	<array>
		<string>{{.Executable}}</string>
	</array>
	<key>StartInterval</key>
	<integer>1200</integer> <!-- 1200 seconds = 20 minutes -->
	<key>RunAtLoad</key>
	<true/>
	<key>KeepAlive</key>
	<false/>
</dict>
</plist>
`

// PlistData holds the necessary data for plist generation
type PlistData struct {
	Label      string
	Executable string
}

// CreateLaunchAgent creates the plist file in ~/Library/LaunchAgents/ if it doesn't exist
func CreateLaunchAgent(deps *di.Deps) error {
	launchAgentsDir := filepath.Join(deps.HomeDir, "Library", "LaunchAgents")
	err := os.MkdirAll(launchAgentsDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create LaunchAgents directory: %w", err)
	}

	plistPath := filepath.Join(launchAgentsDir, config.LaunchAgentLabel+".plist")

	// Check if the plist already exists
	if _, err := os.Stat(plistPath); err == nil {
		fmt.Printf("plist %s already exists, won't create a one\n", plistPath)
		return nil
	}

	// Gather data for plist template
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	plistData := PlistData{
		Label:      config.LaunchAgentLabel,
		Executable: executable,
	}

	// Generate plist content
	var plistContent bytes.Buffer
	tmpl, err := template.New("plist").Parse(plistTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse plist template: %w", err)
	}

	err = tmpl.Execute(&plistContent, plistData)
	if err != nil {
		return fmt.Errorf("failed to execute plist template: %w", err)
	}

	// Write the plist file
	err = os.WriteFile(plistPath, plistContent.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("failed to write plist file: %w", err)
	}

	// Load the plist using launchctl
	cmd := exec.Command("launchctl", "load", plistPath)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to load plist with launchctl: %w", err)
	}

	return nil
}