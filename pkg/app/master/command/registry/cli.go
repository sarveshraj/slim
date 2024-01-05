package registry

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/slimtoolkit/slim/pkg/app"
	"github.com/slimtoolkit/slim/pkg/app/master/command"
)

const (
	Name  = "registry"
	Usage = "Execute registry operations"
	Alias = "r"

	PullCmdName      = "pull"
	PullCmdNameUsage = "Pull a container image from registry"

	PushCmdName      = "push"
	PushCmdNameUsage = "Push a container image to a registry"

	CopyCmdName      = "copy"
	CopyCmdNameUsage = "Copy a container image from one registry to another"

	ImageIndexCreateCmdName      = "image-index-create"
	ImageIndexCreateCmdNameUsage = "Create an image index (aka manifest list) with the referenced images (already in the target registry)"

	ServerCmdName      = "server"
	ServerCmdNameUsage = "Start a registry server"
)

func fullCmdName(subCmdName string) string {
	return fmt.Sprintf("%s.%s", Name, subCmdName)
}

type CommonCommandParams struct {
	UseDockerCreds bool
	CredsAccount   string
	CredsSecret    string
}

func CommonCommandFlagValues(ctx *cli.Context) (*CommonCommandParams, error) {
	values := &CommonCommandParams{
		UseDockerCreds: ctx.Bool(FlagUseDockerCreds),
		CredsAccount:   ctx.String(FlagCredsAccount),
		//prefer env var for secret (todo: add interactive and file read modes)
		CredsSecret: ctx.String(FlagCredsSecret),
	}

	return values, nil
}

type PullCommandParams struct {
	*CommonCommandParams
	TargetRef    string
	SaveToDocker bool
}

func PullCommandFlagValues(ctx *cli.Context) (*PullCommandParams, error) {
	common, err := CommonCommandFlagValues(ctx)
	if err != nil {
		return nil, err
	}

	values := &PullCommandParams{
		CommonCommandParams: common,
		TargetRef:           ctx.String(command.FlagTarget),
		SaveToDocker:        ctx.Bool(FlagSaveToDocker),
	}

	return values, nil
}

type ImageIndexCreateCommandParams struct {
	*CommonCommandParams
	ImageIndexName  string
	ImageNames      []string
	AsManifestList  bool
	InsecureRefs    bool
	DumpRawManifest bool
}

func ImageIndexCreateCommandFlagValues(ctx *cli.Context) (*ImageIndexCreateCommandParams, error) {
	common, err := CommonCommandFlagValues(ctx)
	if err != nil {
		return nil, err
	}

	values := &ImageIndexCreateCommandParams{
		CommonCommandParams: common,
		ImageIndexName:      ctx.String(FlagImageIndexName),
		ImageNames:          ctx.StringSlice(FlagImageName),
		AsManifestList:      ctx.Bool(FlagAsManifestList),
		InsecureRefs:        ctx.Bool(FlagInsecureRefs),
		DumpRawManifest:     ctx.Bool(FlagDumpRawManifest),
	}

	return values, nil
}

type ServerCommandParams struct {
	*CommonCommandParams
	EnableReferrersAPI bool
}

func ServerCommandFlagValues(ctx *cli.Context) (*ServerCommandParams, error) {
	common, err := CommonCommandFlagValues(ctx)
	if err != nil {
		return nil, err
	}

	values := &ServerCommandParams{
		CommonCommandParams: common,
		EnableReferrersAPI:  ctx.Bool(FlagEnableReferrersAPI),
	}

	return values, nil
}

var CLI = &cli.Command{
	Name:    Name,
	Aliases: []string{Alias},
	Usage:   Usage,
	Flags: []cli.Flag{
		cflag(FlagUseDockerCreds),
		cflag(FlagCredsAccount),
		cflag(FlagCredsSecret),
	},
	Subcommands: []*cli.Command{
		{
			Name:  PullCmdName,
			Usage: PullCmdNameUsage,
			Flags: []cli.Flag{
				command.Cflag(command.FlagTarget),
				cflag(FlagSaveToDocker),
			},
			Action: func(ctx *cli.Context) error {
				gcvalues := command.GlobalFlagValues(ctx)
				xc := app.NewExecutionContext(
					fullCmdName(PullCmdName),
					gcvalues.QuietCLIMode,
					gcvalues.OutputFormat)

				cparams, err := PullCommandFlagValues(ctx)
				if err != nil {
					return err
				}

				if cparams.TargetRef == "" {
					if ctx.Args().Len() < 1 {
						xc.Out.Error("param.target", "missing target")
						cli.ShowCommandHelp(ctx, Name)
						return nil
					} else {
						cparams.TargetRef = ctx.Args().First()
					}
				}

				OnPullCommand(xc, gcvalues, cparams)
				return nil
			},
		},
		{
			Name:  PushCmdName,
			Usage: PushCmdNameUsage,
			Action: func(ctx *cli.Context) error {
				gcvalues := command.GlobalFlagValues(ctx)
				xc := app.NewExecutionContext(
					fullCmdName(PushCmdName),
					gcvalues.QuietCLIMode,
					gcvalues.OutputFormat)

				OnPushCommand(xc, gcvalues)
				return nil
			},
		},
		{
			Name:  CopyCmdName,
			Usage: CopyCmdNameUsage,
			Action: func(ctx *cli.Context) error {
				gcvalues := command.GlobalFlagValues(ctx)
				xc := app.NewExecutionContext(
					fullCmdName(CopyCmdName),
					gcvalues.QuietCLIMode,
					gcvalues.OutputFormat)

				OnCopyCommand(xc, gcvalues)
				return nil
			},
		},
		{
			Name:  ImageIndexCreateCmdName,
			Usage: ImageIndexCreateCmdNameUsage,
			Flags: []cli.Flag{
				cflag(FlagImageIndexName),
				cflag(FlagImageName),
				cflag(FlagAsManifestList),
				cflag(FlagInsecureRefs),
				cflag(FlagDumpRawManifest),
			},
			Action: func(ctx *cli.Context) error {
				gcvalues := command.GlobalFlagValues(ctx)
				xc := app.NewExecutionContext(
					fullCmdName(ImageIndexCreateCmdName),
					gcvalues.QuietCLIMode,
					gcvalues.OutputFormat)

				cparams, err := ImageIndexCreateCommandFlagValues(ctx)
				if err != nil {
					return err
				}

				OnImageIndexCreateCommand(xc, gcvalues, cparams)
				return nil
			},
		},
		{
			Name:  ServerCmdName,
			Usage: ServerCmdNameUsage,
			Action: func(ctx *cli.Context) error {
				gcvalues := command.GlobalFlagValues(ctx)

				cparams, err := ServerCommandFlagValues(ctx)
				if err != nil {
					return err
				}

				xc := app.NewExecutionContext(fullCmdName(ServerCmdName),
					gcvalues.QuietCLIMode,
					gcvalues.OutputFormat)
				OnServerCommand(xc, gcvalues, cparams)
				return nil
			},
		},
	},
}
