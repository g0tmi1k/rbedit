package rbeditCmd

import (
	"github.com/rakshasa/rbedit/actions"
	"github.com/rakshasa/rbedit/inputs"
	"github.com/rakshasa/rbedit/outputs"
	"github.com/spf13/cobra"
)

// PutCmd:

func newPutCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "put [OPTIONS] KEY-PATH...",
		Short: "Put object",
		Run:   putCmdRun,
	}

	setupDefaultCommand(cmd)
	addInputFlags(cmd)
	addOutputFlags(cmd)
	addAnyValueFlags(cmd)

	return cmd
}

func putCmdRun(cmd *cobra.Command, args []string) {
	if len(args) == 0 && !hasChangedFlags(cmd) {
		printCommandUsageAndExit(cmd)
	}

	metadata, err := metadataFromCommand(cmd, WithInput(), WithOutput(), WithAnyValue())
	if err != nil {
		printCommandErrorAndExit(cmd, err)
	}

	input := inputs.NewSingleInput(inputs.NewDecodeBencode(), inputs.NewFileInput())
	output := outputs.NewSingleOutput(outputs.NewEncodeBencode(), outputs.NewFileOutput())

	if err := input.Execute(metadata, actions.NewPut(output, args)); err != nil {
		printCommandErrorAndExit(cmd, err)
	}
}
