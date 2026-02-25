package main

import (
	"flag"
	"fmt"
	"os"

	"pumu/internal/scanner"
)

const version = "v1.0.0-beta.1"

func main() {
	sweepCmd := flag.NewFlagSet("sweep", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	versionFlag := flag.Bool("version", false, "Print version information")
	flag.BoolVar(versionFlag, "v", false, "Print version information (shorthand)")

	if len(os.Args) < 2 {
		fmt.Println("Running refresh in current directory...")
		err := scanner.RefreshCurrentDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	switch os.Args[1] {
	case "version", "--version", "-v":
		fmt.Printf("pumu version %s\n", version)
		return
	case "sweep":
		reinstallFlag := sweepCmd.Bool("reinstall", false, "Reinstall packages after removing their folders")
		if err := sweepCmd.Parse(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
			os.Exit(1)
		}
		err := scanner.SweepDir(".", false, *reinstallFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "list":
		if err := listCmd.Parse(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
			os.Exit(1)
		}
		err := scanner.SweepDir(".", true, false)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "repair":
		repairCmd := flag.NewFlagSet("repair", flag.ExitOnError)
		verboseFlag := repairCmd.Bool("verbose", false, "Show details for all projects, including healthy ones")
		if err := repairCmd.Parse(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
			os.Exit(1)
		}
		err := scanner.RepairDir(".", *verboseFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "prune":
		pruneCmd := flag.NewFlagSet("prune", flag.ExitOnError)
		thresholdFlag := pruneCmd.Int("threshold", 50, "Minimum score to prune (0-100)")
		dryRunFlag := pruneCmd.Bool("dry-run", false, "Only analyze and list, don't delete")
		if err := pruneCmd.Parse(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
			os.Exit(1)
		}
		err := scanner.PruneDir(".", *thresholdFlag, *dryRunFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("Unknown command '%s'. Run 'pumu list', 'pumu sweep', 'pumu repair', 'pumu prune' or just 'pumu'.\n", os.Args[1])
		os.Exit(1)
	}
}
