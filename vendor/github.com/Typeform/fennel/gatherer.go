package fennel

type SimpleGatherer interface {
	Gather(Datapoint) error
}

type BatchGatherer interface {
	GatherBatch([]Datapoint) error
}

type Gatherer interface {
	SimpleGatherer
	BatchGatherer
}
