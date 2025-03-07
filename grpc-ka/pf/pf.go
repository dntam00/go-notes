package pf

import (
	"fmt"
	"os"
	"os/exec"
)

const pfCustomFile = "/etc/pf.custom"
const pfFile = "/etc/pf.conf"
const BlockRule = "block out proto tcp from any to any port 5577\n"
const DropBlockRule = "#block out proto tcp from any to any port 5577\n"

func main() {
	// Open pf.custom in append mode
	ApplyRule(BlockRule)
}

func ApplyRule(rule string) {
	askpassPath := "/Users/dntam/Projects/golang/play-around/grpc-ka/pf/askpass.sh" // Update this with the correct path
	os.Setenv("SUDO_ASKPASS", askpassPath)
	f, err := os.OpenFile(pfCustomFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed to open pf.custom:", err)
		return
	}
	defer f.Close()

	// Append the new rule
	_, err = f.WriteString(rule)
	if err != nil {
		fmt.Println("Failed to write rule:", err)
		return
	}

	// Reload the custom anchor
	cmd := exec.Command("sudo", "-A", "pfctl", "-f", pfFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error loading rule into pf:", err)
		return
	}

	fmt.Println("Successfully added rule and reloaded pf.")
}
