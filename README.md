# almen-int

`almen-int` is a small command-line tool that accounts for the arc height of a
shot-peened Almen strip. Given a shot velocity, a strip geometry and an elastic
modulus, it evaluates the bending model, applies an exponential coverage law,
and reports the arc height, the coverage, the saturation state and the
recommended Almen strip.

The core idea is the classic thin-strip bending picture: a shot-peened strip
carries a thin residual compressive layer on the peened face, the layer applies
a bending moment to the strip, and the strip bends into a shallow circular arc.
The sagitta of that arc is the Almen arc height.

## What it does

- Reads a JSON case file describing the shot stream, the Almen strip, the
  residual layer and the process coverage.
- Computes the bending moment from the residual layer and the rectangular
  moment of inertia of the strip, then the arc height of the deformed strip.
- Applies the exponential coverage law `C(t) = 1 - exp(-lambda*t)` and decides
  whether the process is saturated.
- Recommends an Almen strip (N, A or C) from the arc height once the process
  has saturated; no grade letter is reported for an unsaturated process.
- Validates every input and prints a clear diagnostic on stderr with a non-zero
  exit code for invalid cases.

It is intentionally narrow: it models Almen arc height accounting, not surface
treatment scheduling, not plating orders, not Hertz contact pressure and not
fatigue damage accumulation.

## Model

The pinned bending model, in the units used by the case files:

- Residual stress follows a velocity power law with exponent 2:

      sigma(v) = sigma_ref * (v / v_ref)^2

- Bending moment of the residual layer:

      M = sigma * w * d * (t - d) / 2

- Second moment of area of the rectangular strip:

      I = w * t^3 / 12

- Curvature and small-deflection sagitta:

      kappa = M / (E * I)
      h = kappa * L^2 / 8

So the arc height is proportional to `M*L^2/(E*I)`. The thickness enters the
inertia at the third power, which is why a thicker strip bends less for the
same peening.

Coverage accumulates exponentially with the peening time:

      C(t) = 1 - exp(-lambda * t)

and the arc height at a partial coverage is the plateau value scaled by the
gain `g(C) = 1 - exp(-kappa*C)`. The coverage never enters the bending moment
itself. Saturation is pinned to one rule: doubling the peening time must not
raise the arc height by more than ten percent. Because the gain curve is
concave, raising the coverage from 0.5 to 1.0 always adds less arc height than
raising it from 0 to 0.5.

## Usage

Build the binary:

    go build ./...

Evaluate a case:

    almen-int height example/a2-steel.json

The bundled `example/a2-steel.json` is a standard A-strip case on A2 tool
steel: the computed arc height lands inside the 0.1-0.6 mm band where the A
strip is the usual choice, and the coverage is high enough that the process is
saturated.

Sample output:

    almen arc height report
    ------------------------
    arc height        : 0.4407 mm
    coverage          : 0.9800 (98.0%)
    saturated         : yes
    recommended strip : A (standard strip, medium intensity), arc height 0.4407 mm
    almen intensity   : 0.4407 mmA

Add `-v` for the intermediate quantities (bending moment, inertia, curvature,
implied peening time, half-life, saturation ratio, regime checks) and `-json`
for machine-readable output:

    almen-int height -v example/a2-steel.json
    almen-int height -json example/a2-steel.json

An invalid case fails loudly:

    $ almen-int height /nonexistent.json
    almen-int: open case file: open /nonexistent.json: no such file or directory
    $ echo 1

## Case file format

A case file is a single JSON document:

| section    | field            | meaning                                            |
|------------|------------------|----------------------------------------------------|
| `shot`     | `velocity`       | shot velocity, m/s                                 |
| `shot`     | `diameter`       | shot diameter, mm                                  |
| `shot`     | `density`        | shot material density, kg/m^3                      |
| `strip`    | `thickness`      | strip thickness, mm                                |
| `strip`    | `width`          | strip width, mm                                    |
| `strip`    | `length`         | strip span length, mm                              |
| `strip`    | `modulus`        | strip elastic modulus, GPa                         |
| `residual` | `stress`         | residual stress at the reference velocity, MPa     |
| `residual` | `layer_depth`    | depth of the compressive layer, mm                 |
| `process`  | `coverage`       | target coverage, must be in (0, 2]                 |
| `process`  | `rate_constant`  | coverage rate constant lambda, 1/min               |
| `process`  | `gain_coefficient` | gain curve exponent kappa                         |
| `reference`| `velocity`       | velocity the residual stress is quoted at, m/s     |

Every numeric field is required. The validator rejects a zero or negative
velocity, a zero or negative strip size, a non-positive modulus, a residual
layer deeper than half the strip, and a coverage outside `(0, 2]`. The decoder
also refuses unknown JSON keys so a typo is reported instead of ignored.

## Saturation rule

Given the target coverage `C`, the implied peening time is
`t = -ln(1-C)/lambda`. Doubling that time reaches the coverage
`C2 = 1-(1-C)^2`, and the arc height grows by the factor
`g(C2)/g(C)`. The process is saturated when that factor is at most 1.10. The
threshold coverage sits around 0.65 for the reference gain coefficient:
below it the process is not saturated and no strength grade is reported.

## Build and test

    go build ./...
    go test ./...

The project requires Go 1.21 or later and uses only the standard library.

## License

MIT, see `LICENSE`.
