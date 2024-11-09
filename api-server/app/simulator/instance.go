package simulator

import "log"

// Global instance of the `Simulator`.
var S Simulator = func() Simulator {
	sim := SimulatorRealWorld{}
	log.SetOutput(sim.LogOutput())
	return &sim
}()
