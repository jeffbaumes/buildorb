package common

var (
	generators map[string](func(*Planet, CellLoc) int)
	systems    map[string](func() []*PlanetState)
)

func init() {
	generators = make(map[string](func(*Planet, CellLoc) int))

	generators["sphere"] = func(p *Planet, loc CellLoc) int {
		if float64(loc.Alt)/float64(p.AltCells) < 0.5 {
			return Stone
		}
		return Air
	}
	generators["moon"] = func(p *Planet, loc CellLoc) int {
		if float64(loc.Alt)/float64(p.AltCells) < 0.5 {
			return Moon
		}
		return Air
	}
	generators["sun"] = func(p *Planet, loc CellLoc) int {
		if float64(loc.Alt)/float64(p.AltCells) < 0.5 {
			return Sun
		}
		return Air
	}

	generators["rings"] = func(p *Planet, loc CellLoc) int {
		scale := 1.0
		n := p.noise.Eval2(float64(loc.Alt)*scale, 0)
		fracHeight := float64(loc.Alt) / float64(p.AltCells)
		if fracHeight < 0.5 {
			return Grass
		}
		if fracHeight > 0.6 && int(loc.Lat) == p.LatCells/2 {
			if n > 0.1 {
				return YellowBlock
			}
			return RedBlock
		}
		return Air
	}

	generators["bumpy"] = func(p *Planet, loc CellLoc) int {
		pos := p.CellLocToCartesian(loc).Normalize().Mul(float32(p.AltCells / 2))
		scale := 0.1
		height := float64(p.AltCells)/2 + p.noise.Eval3(float64(pos[0])*scale, float64(pos[1])*scale, float64(pos[2])*scale)*8
		if float64(loc.Alt) <= height {
			if float64(loc.Alt) > float64(p.AltCells)/2+2 {
				return Dirt
			}
			return Grass
		}
		if float64(loc.Alt) < float64(p.AltCells)/2+1 {
			return BlueBlock
		}
		return Air
	}

	generators["caves"] = func(p *Planet, loc CellLoc) int {
		pos := p.CellLocToCartesian(loc)
		const scale = 0.05
		height := (p.noise.Eval3(float64(pos[0])*scale, float64(pos[1])*scale, float64(pos[2])*scale) + 1.0) * float64(p.AltCells) / 2.0
		if height > float64(p.AltCells)/2 {
			return Stone
		}
		return Air
	}

	generators["rocks"] = func(p *Planet, loc CellLoc) int {
		pos := p.CellLocToCartesian(loc)
		const scale = 0.05
		noise := p.noise.Eval3(float64(pos[0])*scale, float64(pos[1])*scale, float64(pos[2])*scale)
		if noise > 0.5 {
			return Stone
		}
		return Air
	}

	systems = make(map[string](func() []*PlanetState))

	systems["planet"] = func() []*PlanetState {
		return []*PlanetState{
			&PlanetState{
				ID:              0,
				Name:            "Spawn",
				GeneratorType:   "bumpy",
				Radius:          64.0,
				AltCells:        64,
				RotationSeconds: 10,
			},
		}
	}

	systems["moon"] = func() []*PlanetState {
		return []*PlanetState{
			&PlanetState{
				ID:              0,
				Name:            "Spawn",
				GeneratorType:   "bumpy",
				Radius:          64.0,
				AltCells:        64,
				RotationSeconds: 10,
			},
			&PlanetState{
				ID:              1,
				Name:            "Moon",
				GeneratorType:   "moon",
				Radius:          32.0,
				AltCells:        32,
				OrbitPlanet:     0,
				OrbitDistance:   100,
				OrbitSeconds:    5,
				RotationSeconds: 10,
			},
		}
	}

	systems["sun-moon"] = func() []*PlanetState {
		return []*PlanetState{
			&PlanetState{
				ID:              0,
				Name:            "Spawn",
				GeneratorType:   "bumpy",
				Radius:          64.0,
				AltCells:        64,
				OrbitPlanet:     2,
				OrbitDistance:   300,
				OrbitSeconds:    1095,
				RotationSeconds: 180,
			},
			&PlanetState{
				ID:              1,
				Name:            "Moon",
				GeneratorType:   "moon",
				Radius:          32.0,
				AltCells:        32,
				OrbitPlanet:     0,
				OrbitDistance:   100,
				OrbitSeconds:    90,
				RotationSeconds: -90,
			},
			&PlanetState{
				ID:              2,
				Name:            "Sun",
				GeneratorType:   "sun",
				Radius:          64.0,
				AltCells:        64,
				OrbitPlanet:     2,
				RotationSeconds: 1e10,
			},
		}
	}

	systems["many"] = func() []*PlanetState {
		planets := []*PlanetState{
			&PlanetState{
				ID:              0,
				Name:            "Sun",
				GeneratorType:   "sun",
				Radius:          64.0,
				AltCells:        64,
				OrbitPlanet:     0,
				RotationSeconds: 1e10,
			},
		}
		for i := 0; i < 100; i++ {
			planets = append(planets, &PlanetState{
				ID:              2*i + 1,
				Name:            "Spawn",
				GeneratorType:   "sphere",
				Radius:          32.0,
				AltCells:        32,
				OrbitPlanet:     0,
				OrbitDistance:   70 * float64(i+1),
				OrbitSeconds:    10 + float64(i),
				RotationSeconds: 1e10,
			})
			planets = append(planets, &PlanetState{
				ID:              2*i + 2,
				Name:            "Spawn",
				GeneratorType:   "sphere",
				Radius:          16.0,
				AltCells:        16,
				OrbitPlanet:     2*i + 1,
				OrbitDistance:   30,
				OrbitSeconds:    5,
				RotationSeconds: 1e10,
			})
		}
		return planets
	}

}
