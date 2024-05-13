// 版权 @2023 凹语言 作者。保留所有权利。

package appwat2wasm

import (
	"fmt"
	"os"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/wazero"
)

var CmdWat2wasm = &cli.Command{
	Hidden:    true,
	Name:      "wat2wasm",
	Usage:     "convert wat format to wasm binary format",
	ArgsUsage: "<file.wat>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "set output file",
		},
	},
	Action: func(c *cli.Context) error {
		if c.NArg() == 0 {
			fmt.Fprintf(os.Stderr, "no input file")
			os.Exit(1)
		}

		infile := c.Args().First()
		outfile := c.String("output")
		if outfile == "" {
			outfile = infile + ".wasm"
		}

		source, err := os.ReadFile(infile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		wasmBytes, err := wazero.Wat2Wasm(source)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = os.WriteFile(outfile, wasmBytes, 0666)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		return nil
	},
}
