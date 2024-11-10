package simulator

import "log"

// Global instance of the `Simulator`.
var S Simulator = func() Simulator {
	sim := FromBase(&SimulatorRealWorld{})
	// sim.AlwaysPassGoogleOauth()
	// sim.WithDb("test.db")
	log.SetOutput(sim.LogOutput())
	return sim
}()
