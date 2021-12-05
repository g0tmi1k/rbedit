package rbeditCmd

import (
	"context"
	"fmt"

	"github.com/rakshasa/rbedit/objects"
	"github.com/spf13/cobra"
)

var (
	announcePath = []string{"announce"}
)

// AnnounceCmd:

func newAnnounceCommand(ctx context.Context) (*cobra.Command, context.Context) {
	cmd := &cobra.Command{
		Use:   "announce [OPTIONS] COMMAND",
		Short: "BitTorrent announce related commands",
		Args:  cobra.ExactArgs(0),
		Run:   func(cmd *cobra.Command, args []string) { printCommandUsage(cmd) },
	}

	setupDefaultCommand(cmd, "rbedit-announce")

	ctx = addCommand(ctx, cmd, newAnnounceGetCommand)
	ctx = addCommand(ctx, cmd, newAnnouncePutCommand)

	return cmd, ctx
}

// AnnounceGetCmd:

func newAnnounceGetCommand(ctx context.Context) (*cobra.Command, context.Context) {
	cmd := &cobra.Command{
		Use:   "get [OPTIONS]",
		Short: "Get announce url",
		Args:  cobra.ExactArgs(0),
		Run:   announceGetCmdRun,
	}

	setupDefaultCommand(cmd, "rbedit-announce-get-state")

	addInputFlags(ctx, cmd)

	return cmd, ctx
}

func announceGetCmdRun(cmd *cobra.Command, args []string) {
	metadata, err := metadataFromCommand(cmd, WithInput())
	if err != nil {
		printCommandErrorAndExit(cmd, err)
	}

	input := objects.NewSingleInput(objects.NewDecodeBencode(), objects.NewFileInput())

	if err := input.Execute(metadata, func(rootObj interface{}, metadata objects.IOMetadata) error {
		obj, err := objects.LookupKeyPath(rootObj, announcePath)
		if err != nil {
			printCommandErrorAndExit(cmd, err)
		}

		if _, ok := objects.AsAbsoluteURI(obj); !ok {
			printCommandErrorAndExit(cmd, fmt.Errorf("announce not a valid URI string"))
		}

		objects.PrintObject(obj)
		return nil

	}); err != nil {
		printCommandErrorAndExit(cmd, err)
	}
}

// AnnouncesPutCmd:

func newAnnouncePutCommand(ctx context.Context) (*cobra.Command, context.Context) {
	cmd := &cobra.Command{
		Use:   "put [OPTIONS] URI",
		Short: "Set tracker announce URI",
		Args:  cobra.ExactArgs(1),
		Run:   announcePutCmdRun,
	}

	setupDefaultCommand(cmd, "rbedit-announce-put-state")

	addInputFlags(ctx, cmd)
	addOutputFlags(ctx, cmd)

	return cmd, ctx
}

func announcePutCmdRun(cmd *cobra.Command, args []string) {
	tracker := args[0]
	if !objects.VerifyAbsoluteURI(tracker) {
		printCommandErrorAndExit(cmd, fmt.Errorf("failed to validate URI"))
	}

	metadata, err := metadataFromCommand(cmd, WithInput(), WithOutput())
	if err != nil {
		printCommandErrorAndExit(cmd, err)
	}

	input := objects.NewSingleInput(objects.NewDecodeBencode(), objects.NewFileInput())
	output := objects.NewSingleOutput(objects.NewEncodeBencode(), objects.NewFileOutput())

	if err := input.Execute(metadata, func(rootObj interface{}, metadata objects.IOMetadata) error {
		rootObj, err := objects.SetObject(rootObj, tracker, announcePath)
		if err != nil {
			printCommandErrorAndExit(cmd, err)
		}

		if err := output.Execute(rootObj, metadata); err != nil {
			printCommandErrorAndExit(cmd, err)
		}

		return nil

	}); err != nil {
		printCommandErrorAndExit(cmd, err)
	}
}
