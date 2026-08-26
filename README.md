# almen-int

`almen-int` 核算喷丸 Almen 试片弧高。用户给出丸粒速度、试片几何与弹性模量，算出薄片弯矩弧高、指数覆盖率、饱和判定和 N/A/C 建议号。必须同时成立：速度加大弧高加大、覆盖率趋近平台后时间加倍弧高增幅不超过 10%、试片加厚同一丸流弧高变小。交付为无前端的 HTTP 服务与小型 CLI。

## Model

Residual stress follows a velocity power law with exponent 2:

    sigma(v) = sigma_ref * (v / v_ref)^2

Bending moment of the residual layer:

    M = sigma * w * d * (t - d) / 2

Second moment of area of the rectangular strip:

    I = w * t^3 / 12

Curvature and small-deflection sagitta:

    kappa = M / (E * I)
    h = kappa * L^2 / 8

Coverage accumulates exponentially with the peening time:

    C(t) = 1 - exp(-lambda * t)

and the arc height at a partial coverage is the plateau value scaled by the
gain `g(C) = 1 - exp(-kappa*C)`. Saturation is pinned to one rule: doubling
the peening time must not raise the arc height by more than ten percent.

N, A and C strips share the same residual layer. Equivalent height on another
thickness follows `(t - d) / t^3`, so converting A→C→A recovers the original
reading.

## HTTP API

With no arguments, or with the `serve` subcommand, the process listens on
`:8080`:

```
go run .
go run . serve
```

`GET /api/health` returns `{"status":"ok"}`.

`POST /api/height` evaluates a peening case:

```
{
  "shot": {"velocity": 48.0, "diameter": 0.60, "density": 7800.0},
  "strip": {"thickness": 1.29, "width": 18.5, "length": 76.0, "modulus": 205.0},
  "residual": {"stress": 850.0, "layer_depth": 0.05},
  "process": {"coverage": 0.98, "rate_constant": 0.085, "gain_coefficient": 2.6},
  "reference": {"velocity": 50.0}
}
```

`POST /api/saturation` returns the same case's 10% saturation verdict, the
implied peening time, and the doubled-coverage gain ratio.

Illegal domain input (non-positive velocity, coverage outside (0, 2], layer
deeper than half the strip) returns HTTP 422. Malformed JSON returns HTTP 400.

## CLI usage

```
go build -o almen-int .
./almen-int                              # HTTP on :8080
./almen-int height example/a2-steel.json
./almen-int height -v example/a2-steel.json
./almen-int height -json example/a2-steel.json
```

The bundled `example/a2-steel.json` is a standard A-strip case on A2 tool
steel: the computed arc height lands inside the 0.1–0.6 mm band where the A
strip is the usual choice, and the coverage is high enough that the process is
saturated.

Sample CLI output:

    almen arc height report
    ------------------------
    arc height        : 0.4407 mm
    coverage          : 0.9800 (98.0%)
    saturated         : yes
    recommended strip : A (standard strip, medium intensity), arc height 0.4407 mm
    almen intensity   : 0.4407 mmA

An invalid case fails on stderr with a non-zero exit code.

## Case file format

| section    | field              | meaning                                         |
|------------|--------------------|-------------------------------------------------|
| `shot`     | `velocity`         | shot velocity, m/s                              |
| `shot`     | `diameter`         | shot diameter, mm                               |
| `shot`     | `density`          | shot material density, kg/m^3                   |
| `strip`    | `thickness`        | strip thickness, mm                             |
| `strip`    | `width`            | strip width, mm                                 |
| `strip`    | `length`           | strip span length, mm                           |
| `strip`    | `modulus`          | strip elastic modulus, GPa                      |
| `residual` | `stress`           | residual stress at the reference velocity, MPa  |
| `residual` | `layer_depth`      | depth of the compressive layer, mm              |
| `process`  | `coverage`         | target coverage, must be in (0, 2]              |
| `process`  | `rate_constant`    | coverage rate constant lambda, 1/min            |
| `process`  | `gain_coefficient` | gain curve exponent kappa                       |
| `reference`| `velocity`         | velocity the residual stress is quoted at, m/s  |

Every numeric field is required. The decoder refuses unknown JSON keys.

## Build and test

```
go build ./...
go test ./...
```

The project requires Go 1.21 or later and uses only the standard library.

## License

MIT, see `LICENSE`.
