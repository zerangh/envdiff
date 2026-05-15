package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/user/envdiff/internal/encoder"
	"github.com/user/envdiff/internal/parser"
)

func runEncode(args []string, out io.Writer) error {
	fs := flag.NewFlagSet("encode", flag.ContinueOnError)
	format := fs.String("format", "dotenv", "output format: dotenv, exports, docker")
	omitEmpty := fs.Bool("omit-empty", false, "omit keys with empty values")
	quoteAll := fs.Bool("quote-all", false, "quote all values")

	if err := fs.Parse(args); err != nil {
		return err
	}

	positional := fs.Args()
	if len(positional) < 1 {
		return fmt.Errorf("usage: envdiff encode [flags] <file>")
	}

	env, err := parser.ParseFile(positional[0])
	if err != nil {
		return fmt.Errorf("encode: failed to parse %s: %w", positional[0], err)
	}

	opts := encoder.Options{
		Format:    encoder.Format(*format),
		OmitEmpty: *omitEmpty,
		QuoteAll:  *quoteAll,
		SortKeys:  true,
	}

	result, err := encoder.Encode(env, opts)
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	fmt.Fprint(out, result)
	return nil
}

func init() {
	_ = os.Args // ensure os is used for potential future flag defaults
}
