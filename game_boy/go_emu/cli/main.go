package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const hostsFile = "/etc/hosts"
const markerStart = "# ==== FOCUS MODE START ===="
const markerEnd = "# ==== FOCUS MODE END ===="

var blockedDomains = []string{
	// YouTube
	"youtube.com",
	"www.youtube.com",
	"m.youtube.com",
	"youtubei.googleapis.com",
	"youtube-ui.l.google.com",
	"s.ytimg.com",
	"i.ytimg.com",
	"yt3.ggpht.com",
	"googlevideo.com",
	"www.googlevideo.com",

	// Reddit
	"reddit.com",
	"www.reddit.com",
	"reddit.com",

	// Discord
	"discord.com",
	"www.discord.com",
	"www.instagram.com",
	"i.instagram.com",

	// Intagram
	"graph.instagram.com",
	"scontent-cdg4-3.cdninstagram.com",
	"edge-chat.instagram.com",
	"static.cdninstagram.com",
	"cdninstagram.com",
	"scontent.cdninstagram.com",

	// Linkedin
	"www.linkedin.com",
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "block":
		block()
	case "unblock":
		unblock()
	case "status":
		status()
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: focus <command>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  block     Block distracting websites")
	fmt.Println("  unblock   Unblock websites")
	fmt.Println("  status    Show current blocking status")
}

func block() {
	data, err := os.ReadFile(hostsFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", hostsFile, err)
		fmt.Fprintln(os.Stderr, "Try running with sudo.")
		os.Exit(1)
	}

	content := string(data)
	if strings.Contains(content, markerStart) {
		fmt.Println("Already blocking. Use 'unblock' first to reset.")
		return
	}

	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(markerStart + "\n")
	for _, domain := range blockedDomains {
		b.WriteString(fmt.Sprintf("0.0.0.0 %s\n", domain))
		b.WriteString(fmt.Sprintf("::0 %s\n", domain))
	}
	b.WriteString(markerEnd + "\n")

	err = os.WriteFile(hostsFile, []byte(content+b.String()), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", hostsFile, err)
		fmt.Fprintln(os.Stderr, "Try running with sudo.")
		os.Exit(1)
	}

	flushDNS()
	fmt.Println("Blocked. Focus mode ON.")
}

func unblock() {
	data, err := os.ReadFile(hostsFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", hostsFile, err)
		os.Exit(1)
	}

	content := string(data)
	startIdx := strings.Index(content, markerStart)
	endIdx := strings.Index(content, markerEnd)

	if startIdx == -1 || endIdx == -1 {
		fmt.Println("Nothing to unblock.")
		return
	}

	cleaned := content[:startIdx] + content[endIdx+len(markerEnd)+1:]
	cleaned = strings.TrimRight(cleaned, "\n") + "\n"

	err = os.WriteFile(hostsFile, []byte(cleaned), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", hostsFile, err)
		fmt.Fprintln(os.Stderr, "Try running with sudo.")
		os.Exit(1)
	}

	flushDNS()
	fmt.Println("Unblocked. Focus mode OFF.")
}

func status() {
	data, err := os.ReadFile(hostsFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", hostsFile, err)
		os.Exit(1)
	}

	if strings.Contains(string(data), markerStart) {
		fmt.Println("Focus mode: ON")
		fmt.Println("Blocked domains:")
		for _, d := range blockedDomains {
			fmt.Printf("  - %s\n", d)
		}
	} else {
		fmt.Println("Focus mode: OFF")
	}
}

func flushDNS() {
	if runtime.GOOS == "darwin" {
		exec.Command("dscacheutil", "-flushcache").Run()
		exec.Command("killall", "-HUP", "mDNSResponder").Run()
	} else if runtime.GOOS == "linux" {
		exec.Command("systemctl", "restart", "systemd-resolved").Run()
	}
}
