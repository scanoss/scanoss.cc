package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var (
	scanDir          string
	wfp              string
	dep              string
	stdin            string
	files            []string
	identify         string
	ignore           string
	output           string
	format           string
	threads          int
	flags            int
	postSize         int
	timeout          int
	retry            int
	noWfpOutput      bool
	dependencies     bool
	dependenciesOnly bool
	scCommand        string
	scTimeout        int
	depScope         string
	depScopeInc      string
	depScopeExc      string
	settings         string
	skipSettingsFile bool
	debug            bool = false
	trace            bool = false
	quiet            bool = false
)

var scanCmd = &cobra.Command{
	Use:   "scan [scanDirPath]",
	Short: "Run a scan on the specified folder",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := checkPythonDependencies(); err != nil {
			fmt.Fprintf(os.Stderr, "dependency check failed: %v\n", err)
			os.Exit(1)
		}

		scanDir = args[0]

		scanossPyCmd := exec.Command("scanoss-py", "scan", scanDir)
		var stderr bytes.Buffer
		scanossPyCmd.Stderr = &stderr

		output, err := scanossPyCmd.Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error running python script: %v", stderr.String())
			os.Exit(1)
		}
		fmt.Println(string(output))
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	scanCmd.Flags().StringVarP(&wfp, "wfp", "w", "", "Scan a WFP File instead of a folder (optional)")
	scanCmd.Flags().StringVarP(&dep, "dep", "p", "", "Use a dependency file instead of a folder (optional)")
	scanCmd.Flags().StringVarP(&stdin, "stdin", "s", "", "Scan the file contents supplied via STDIN (optional)")
	scanCmd.Flags().StringSliceVarP(&files, "files", "e", nil, "List of files to scan.")
	scanCmd.Flags().StringVarP(&identify, "identify", "i", "", "Scan and identify components in SBOM file")
	scanCmd.Flags().StringVarP(&ignore, "ignore", "n", "", "Ignore components specified in the SBOM file")
	scanCmd.Flags().StringVarP(&output, "output", "o", "", "Output result file name (optional - default stdout).")
	scanCmd.Flags().StringVarP(&format, "format", "f", "plain", "Result output format (optional - default: plain)")
	scanCmd.Flags().IntVarP(&threads, "threads", "T", 5, "Number of threads to use while scanning (optional - default 5)")
	scanCmd.Flags().IntVarP(&flags, "flags", "F", 0, "Scanning engine flags")
	scanCmd.Flags().IntVarP(&postSize, "post-size", "P", 32, "Number of kilobytes to limit the post to while scanning (optional - default 32)")
	scanCmd.Flags().IntVarP(&timeout, "timeout", "M", 180, "Timeout (in seconds) for API communication (optional - default 180)")
	scanCmd.Flags().IntVarP(&retry, "retry", "R", 5, "Retry limit for API communication (optional - default 5)")
	scanCmd.Flags().BoolVar(&noWfpOutput, "no-wfp-output", false, "Skip WFP file generation")
	scanCmd.Flags().BoolVarP(&dependencies, "dependencies", "D", false, "Add Dependency scanning")
	scanCmd.Flags().BoolVar(&dependenciesOnly, "dependencies-only", false, "Run Dependency scanning only")
	scanCmd.Flags().StringVar(&scCommand, "sc-command", "", "Scancode command and path if required (optional - default scancode).")
	scanCmd.Flags().IntVar(&scTimeout, "sc-timeout", 600, "Timeout (in seconds) for scancode to complete (optional - default 600)")
	scanCmd.Flags().StringVar(&depScope, "dep-scope", "", "Filter dependencies by scope - default all (options: dev/prod)")
	scanCmd.Flags().StringVar(&depScopeInc, "dep-scope-inc", "", "Include dependencies with declared scopes")
	scanCmd.Flags().StringVar(&depScopeExc, "dep-scope-exc", "", "Exclude dependencies with declared scopes")
	scanCmd.Flags().StringVar(&settings, "settings", "", "Settings file to use for scanning (optional - default scanoss.json)")
	scanCmd.Flags().BoolVar(&skipSettingsFile, "skip-settings-file", false, "Skip default settings file (scanoss.json) if it exists")
	scanCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug messages")
	scanCmd.Flags().BoolVarP(&trace, "trace", "t", false, "Enable trace messages, including API posts")
	scanCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Enable quiet mode")
}

func checkPythonDependencies() error {
	pythonCmd := exec.Command("python", "--version")
	if err := pythonCmd.Run(); err != nil {
		return fmt.Errorf("python is not installed: %w", err)
	}

	pipCmd := exec.Command("scanoss-py", "--version")
	if err := pipCmd.Run(); err != nil {
		return fmt.Errorf("scanoss-py is not installed or not in PATH: %w", err)
	}

	return nil
}
