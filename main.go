// Command almen-int evaluates shot-peening Almen arc height from a case file.
//
// Usage:
//
//	almen-int height <case.json>
//
// The "height" subcommand loads a JSON case file describing the shot stream,
// the Almen strip, the residual layer and the process coverage, applies the
// pinned bending model and the exponential coverage law, and prints the arc
// height, the coverage, the saturation state and the recommended Almen strip.
// Invalid input produces a diagnostic on stderr and a non-zero exit code.
package main

import (
	"flag"
	"fmt"
	"os"

	"almen-int/internal/coverage"
	"almen-int/internal/grade"
	"almen-int/internal/iofmt"
	"almen-int/internal/model"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "almen-int: "+err.Error())
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing subcommand; use 'almen-int height <case.json>'")
	}
	switch args[0] {
	case "height":
		return runHeight(args[1:])
	case "help", "-h", "--help":
		printHelp()
		return nil
	default:
		return fmt.Errorf("unknown subcommand %q; only 'height' is supported", args[0])
	}
}

func runHeight(args []string) error {
	fs := flag.NewFlagSet("height", flag.ContinueOnError)
	verbose := fs.Bool("v", false, "print the intermediate bending and coverage quantities")
	asJSON := fs.Bool("json", false, "emit the result as JSON")
	if err := fs.Parse(args); err != nil {
		return err
	}
	rest := fs.Args()
	if len(rest) != 1 {
		return fmt.Errorf("height expects exactly one case file argument")
	}
	path := rest[0]

	doc, err := iofmt.LoadFile(path)
	if err != nil {
		return err
	}
	if issues := iofmt.AllIssues(doc); len(issues) > 0 {
		return iofmt.JoinIssues("invalid case file:", issues)
	}

	mp := iofmt.BuildModelParams(doc)
	cp := iofmt.BuildCoverageParams(doc)

	mRep := model.BuildReport(mp)
	cRep, err := coverage.BuildReport(cp)
	if err != nil {
		return err
	}

	res := grade.Assemble(mRep.PlateauArcHeight, cp.Coverage, cRep.Gain, cRep.Saturation)

	if *asJSON {
		return iofmt.PrintJSON(os.Stdout, res)
	}
	out := iofmt.RenderResult(res, mRep, cRep, iofmt.RenderOptions{Verbose: *verbose})
	fmt.Print(out)
	return nil
}

func printHelp() {
	fmt.Print(`almen-int - shot-peening Almen arc height accounting

Usage:
  almen-int height <case.json>   evaluate arc height from a case file
  almen-int help                 show this message

The case file is a JSON document with four nested objects:
  shot       velocity (m/s), diameter (mm), density (kg/m^3)
  strip      thickness (mm), width (mm), length (mm), modulus (GPa)
  residual   stress (MPa at the reference velocity), layer_depth (mm)
  process    coverage (0, 2], rate_constant (1/min), gain_coefficient
  reference  velocity (m/s) that anchors the velocity power law

The arc height follows the thin-strip bending model h = M*L^2/(8*E*I) with
I = w*t^3/12 and M = sigma*w*d*(t-d)/2; the residual stress scales with the
square of the velocity. Coverage follows C(t) = 1 - exp(-lambda*t), and the
process is saturated when doubling the peening time raises the arc height by
no more than ten percent. A grade letter is only reported once saturated.

Example:
  almen-int height example/a2-steel.json
`)
}
