package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/amar-jay/minibuf/generate"
	"github.com/amar-jay/minibuf/parser"
	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
)

/*
*
minibufc schema.mb -o types/minibuf --c --ts # generates the types/minibuf.c and types/minibuf.h for C and types/minibuf.ts for TypeScript
*/
func generateCommand() cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		if cmd.NArg() == 0 {
			color.New(color.FgRed).Fprintf(os.Stderr, "Error: No input files provided\n")
			return cli.Exit("", 1)
		}

		inputFiles := cmd.Args().Slice()
		outputDir := cmd.String("output")
		if outputDir == "" {
			outputDir = "generated"
		}

		generateC := cmd.Bool("c")
		generateTS := cmd.Bool("ts")
		if !generateC && !generateTS {
			color.New(color.FgRed).Fprintf(os.Stderr, "Error: At least one of --c or --ts must be specified\n")
			return cli.Exit("", 1)
		}

		// create the output directory if it doesn't exist mkdir -p
		if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
			color.New(color.FgRed).Fprintf(os.Stderr, "Error: Failed to create output directory: %v\n", err)
			return cli.Exit("", 1)
		}

		p := parser.ParseInputFiles(inputFiles)
		if p == nil {
			return cli.Exit("", 1)
		}
		if generateC {
			header, ccode := generate.C(p)
			headerPath := filepath.Join(outputDir, "minibuf.h")
			if err := os.WriteFile(headerPath, []byte(header), 0644); err != nil {
				color.New(color.FgRed).Fprintf(os.Stderr, "Error writing header file: %v\n", err)
				return cli.Exit("", 1)
			}
			cPath := filepath.Join(outputDir, "minibuf.c")
			if err := os.WriteFile(cPath, []byte(ccode), 0644); err != nil {
				color.New(color.FgRed).Fprintf(os.Stderr, "Error writing C file: %v\n", err)
				return cli.Exit("", 1)
			}
		}

		if generateTS {
			ts := generate.TS(p)
			tsPath := filepath.Join(outputDir, "minibuf.ts")
			if err := os.WriteFile(tsPath, []byte(ts), 0644); err != nil {
				color.New(color.FgRed).Fprintf(os.Stderr, "Error writing TS file: %v\n", err)
				return cli.Exit("", 1)
			}
		}

		// Placeholder for actual code generation logic
		log.Printf("Parsed %d schemas from %v to %s (C: %v, TS: %v)", len(p.Schemas), inputFiles, outputDir, generateC, generateTS)
		// print parsed schemas for debugging
		for _, schema := range p.Schemas {
			log.Printf("Schema: %s", schema.Name)
			for _, field := range schema.Fields {
				log.Printf("  Field: %s (%s)", field.Name, field.DataType)
			}
			for defName, defValue := range schema.Defaults {
				log.Printf("  Default: %s = %v", defName, defValue)
			}
		}

		return nil
	}
}

func main() {
	app := &cli.Command{
		Name:  "minibufc",
		Usage: "A minibuf compiler that generates C and TypeScript code from minibuf schema files.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "output, o",
				Aliases: []string{"o"},
				Usage:   "Output directory for generated code",
			},

			// language flags
			&cli.BoolFlag{
				Name:  "c",
				Usage: "Generate C code",
			},
			&cli.BoolFlag{
				Name:  "ts",
				Usage: "Generate TypeScript code",
			},
		},
		Action: generateCommand(),
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
