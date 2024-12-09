package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/scanoss/scanoss.lui/backend/service"
	"github.com/spf13/cobra"
)

type ArgDef struct {
	Name      string
	Shorthand string
	Default   interface{}
	Usage     string
	Type      string
}

var (
	scanArgs = []ArgDef{
		{"wfp", "w", "", "Scan a WFP File instead of a folder (optional)", "string"},
		{"dep", "p", "", "Use a dependency file instead of a folder (optional)", "string"},
		{"stdin", "s", "", "Scan the file contents supplied via STDIN (optional)", "string"},
		{"files", "e", []string{}, "List of files to scan.", "stringSlice"},
		{"identify", "i", "", "Scan and identify components in SBOM file", "string"},
		{"ignore", "n", "", "Ignore components specified in the SBOM file", "string"},
		{"output", "o", "", "Output result file name (optional - default stdout).", "string"},
		{"format", "f", "plain", "Result output format (optional - default: plain)", "string"},
		{"threads", "T", 5, "Number of threads to use while scanning (optional - default 5)", "int"},
		{"flags", "F", 0, "Scanning engine flags", "int"},
		{"post-size", "P", 32, "Number of kilobytes to limit the post to while scanning (optional - default 32)", "int"},
		{"timeout", "M", 180, "Timeout (in seconds) for API communication (optional - default 180)", "int"},
		{"retry", "R", 5, "Retry limit for API communication (optional - default 5)", "int"},
		{"no-wfp-output", "", false, "Skip WFP file generation", "bool"},
		{"dependencies", "D", false, "Add Dependency scanning", "bool"},
		{"dependencies-only", "", false, "Run Dependency scanning only", "bool"},
		{"sc-command", "", "", "Scancode command and path if required (optional - default scancode).", "string"},
		{"sc-timeout", "", 600, "Timeout (in seconds) for scancode to complete (optional - default 600)", "int"},
		{"dep-scope", "", "", "Filter dependencies by scope - default all (options: dev/prod)", "string"},
		{"dep-scope-inc", "", "", "Include dependencies with declared scopes", "string"},
		{"dep-scope-exc", "", "", "Exclude dependencies with declared scopes", "string"},
		{"settings", "", "", "Settings file to use for scanning (optional - default scanoss.json)", "string"},
		{"skip-settings-file", "", false, "Skip default settings file (scanoss.json) if it exists", "bool"},
		{"debug", "d", false, "Enable debug messages", "bool"},
		{"trace", "t", false, "Enable trace messages, including API posts", "bool"},
		{"quiet", "q", true, "Enable quiet mode", "bool"},
	}
)

func NewScanCmd(scanService service.ScanService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scan [scanDirPath]",
		Short: "Run a scan on the specified folder",
		Args: func(cmd *cobra.Command, args []string) error {
			filesFlag := cmd.Flag("files")
			// If no files are specified and no folder is specified, return an error
			if filesFlag == nil || !filesFlag.Changed && len(args) == 0 {
				return fmt.Errorf("you must specify a folder to scan")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := scanService.CheckDependencies(); err != nil {
				return fmt.Errorf("dependency check failed: %w", err)
			}

			scanOptions := make([]string, 0)
			scanOptions = append(scanOptions, "--quiet")

			var scanDirPath string

			// Check if the folder path is specified as an argument
			// Could be that the user specifies only --files flag (e.g scan --files file1.go file2.go)
			if len(args) == 1 && args[0] != "" {
				scanDirPath = args[0]
				scanOptions = append(scanOptions, scanDirPath)
			}

			for _, arg := range scanArgs {
				flag := cmd.Flag(arg.Name)
				if flag == nil || !flag.Changed {
					continue
				}

				switch arg.Type {
				case "string":
					scanOptions = append(scanOptions, fmt.Sprintf("--%s", arg.Name), flag.Value.String())
				case "stringSlice":
					values, err := cmd.Flags().GetStringSlice(arg.Name)
					if err != nil {
						return fmt.Errorf("an error occurred with argument %s: %w", arg.Name, err)
					}
					commaSeparatedValues := strings.Join(values, ",")
					scanOptions = append(scanOptions, fmt.Sprintf("--%s", arg.Name), commaSeparatedValues)
				case "int":
					scanOptions = append(scanOptions, fmt.Sprintf("--%s", arg.Name), flag.Value.String())
				case "bool":
					if flag.Value.String() == "true" {
						scanOptions = append(scanOptions, fmt.Sprintf("--%s", arg.Name))
					}
				}
			}

			return scanService.Scan(scanOptions)
		},
	}

	for _, arg := range scanArgs {
		switch arg.Type {
		case "string":
			cmd.Flags().StringP(arg.Name, arg.Shorthand, arg.Default.(string), arg.Usage)
		case "stringSlice":
			cmd.Flags().StringSliceP(arg.Name, arg.Shorthand, arg.Default.([]string), arg.Usage)
		case "int":
			cmd.Flags().IntP(arg.Name, arg.Shorthand, arg.Default.(int), arg.Usage)
		case "bool":
			cmd.Flags().BoolP(arg.Name, arg.Shorthand, arg.Default.(bool), arg.Usage)
		}
	}

	return cmd
}

func init() {
	service := service.NewScanServicePythonImpl()
	scanCmd := NewScanCmd(service)

	// This is a workaround to prevent the scan command opening the lui when running tests
	if os.Getenv("GO_TEST") != "true" {
		scanCmd.PostRun = func(cmd *cobra.Command, args []string) {
			os.Exit(0)
		}
	}

	rootCmd.AddCommand(scanCmd)
}
