package iofmt

import (
	"fmt"
	"strings"

	"almen-int/internal/coverage"
	"almen-int/internal/grade"
	"almen-int/internal/model"
)

type RenderOptions struct {
	Verbose bool
}

func RenderResult(r grade.Result, mo model.ModelReport, co coverage.CoverageReport, opts RenderOptions) string {
	var b strings.Builder

	b.WriteString("almen arc height report\n")
	b.WriteString("------------------------\n")
	b.WriteString(Labeled("arc height", F4(r.ArcHeight)+" mm", width) + "\n")
	b.WriteString(Labeled("coverage", F4(r.Coverage)+" ("+Pct1(r.Coverage)+")", width) + "\n")
	b.WriteString(Labeled("saturated", r.SaturationStateText(), width) + "\n")
	b.WriteString(Labeled("recommended strip", recommendedText(r), width) + "\n")
	if r.HasGrade() {
		b.WriteString(Labeled("almen intensity", r.IntensityText(), width) + "\n")
	}

	if opts.Verbose {
		b.WriteString("\ncase\n")
		b.WriteString("----\n")
		b.WriteString(Labeled("geometry", DescribeParams(mo.Params), width) + "\n")
		b.WriteString(Labeled("shot", DescribeShot(mo.Params), width) + "\n")

		b.WriteString("\npeening and coverage\n")
		b.WriteString("--------------------\n")
		b.WriteString(Labeled("implied peening time", co.ImpliedTimeText(), width) + "\n")
		b.WriteString(Labeled("doubled time", co.DoubledTimeText(), width) + "\n")
		b.WriteString(Labeled("doubled coverage", F4(co.Saturation.DoubledCoverage), width) + "\n")
		b.WriteString(Labeled("arc-height growth", F2(co.GrowthPct)+"% on time doubling", width) + "\n")
		b.WriteString(Labeled("gain (coverage)", F6(co.Gain), width) + "\n")
		b.WriteString(Labeled("plateau arc height", F4(r.Plateau)+" mm", width) + "\n")
		b.WriteString(Labeled("coverage half-life", F2(coverage.HalfLife(co.Params.RateConstant))+" min", width) + "\n")
		b.WriteString(Labeled("remaining gain", F4(coverage.RemainingGain(co.Params.GainCoefficient, co.Params.Coverage)), width) + "\n")

		b.WriteString("\nbending model\n")
		b.WriteString("-------------\n")
		b.WriteString(Labeled("strip geometry", model.StripGeometryName(mo.Params), width) + "\n")
		b.WriteString(Labeled("shot kinetic energy", FG(mo.KineticEnergy)+" J", width) + "\n")
		b.WriteString(Labeled("residual stress", F2(mo.ResidualStress)+" MPa", width) + "\n")
		b.WriteString(Labeled("layer force", F3(mo.LayerForce)+" N", width) + "\n")
		b.WriteString(Labeled("lever arm", F4(mo.LeverArm)+" mm", width) + "\n")
		b.WriteString(Labeled("bending moment", F4(mo.BendingMoment)+" N*mm", width) + "\n")
		b.WriteString(Labeled("moment of inertia", F6(mo.MomentOfInertia)+" mm^4", width) + "\n")
		b.WriteString(Labeled("section modulus", F4(mo.SectionModulus)+" mm^3", width) + "\n")
		b.WriteString(Labeled("curvature", FG(mo.Curvature)+" 1/mm", width) + "\n")
		b.WriteString(Labeled("exact sagitta", F4(mo.ExactSagitta)+" mm", width) + "\n")
		b.WriteString(Labeled("small-deflection err", FG(mo.RelativeError)+" relative", width) + "\n")

		regime := model.EvaluateRegime(mo.Params)
		b.WriteString("\ndeflection regime\n")
		b.WriteString("-----------------\n")
		b.WriteString(Labeled("thickness/span", F4(regime.ThicknessOverSpan), width) + "\n")
		b.WriteString(Labeled("arc height/span", F4(regime.ArcHeightOverSpan), width) + "\n")
		b.WriteString(Labeled("radius/span", F4(regime.RadiusOverSpan), width) + "\n")
		b.WriteString(Labeled("surface strain", FG(regime.SurfaceStrain), width) + "\n")
		b.WriteString(Labeled("verdict", regime.RegimeSummary(), width) + "\n")
	}

	return b.String()
}

func recommendedText(r grade.Result) string {
	if !r.HasGrade() {
		return "(not saturated - no grade reported)"
	}
	return fmt.Sprintf("%s (%s), arc height %s mm",
		r.Recommend.Strip.Letter,
		r.Recommend.Strip.Note,
		F4(r.ArcHeight),
	)
}

const width = 18

func F3(v float64) string { return fmt.Sprintf("%.3f", v) }
