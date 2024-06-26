package processor

import (
	"github.com/lab5e/lospan/pkg/lg"
	"github.com/lab5e/lospan/pkg/server"
)

// Pipeline is the main processing pipeline for the server. Each step in
// the pipeline is handled by one or more goroutines. Channels are used to
// forward messages between the steps in the pipeline the channels are
// unbuffered at the moment and each step runs as a single goroutine. If one
// of the steps end up being a bottleneck we can increase the number of outputs
// and at the same time buffer the channels.
//
// The pipeline is roughly built like this:
//
//	GW Forwarder -> Decoder -> Decrypter -> MAC Processor
//	      => Scheduler => Encoder -> GW Forwarder
type Pipeline struct {
	Decoder      *Decoder
	Decrypter    *Decrypter
	MACProcessor *MACProcessor
	Scheduler    *Scheduler
	Encoder      *Encoder
}

// Start launches the pipeline
func (p *Pipeline) Start() {
	go p.Decoder.Start()
	go p.Decrypter.Start()
	go p.MACProcessor.Start()
	go p.Scheduler.Start()
	go p.Encoder.Start()
}

// NewPipeline creates a new pipeline. The pipeline will stop automatically
// when the forwarder is terminated
func NewPipeline(context *server.Context, forwarder GwForwarder) *Pipeline {
	ret := Pipeline{}

	lg.Debug("Creating decoder...")
	ret.Decoder = NewDecoder(context, forwarder.Output())

	lg.Debug("Creating decrypter...")
	ret.Decrypter = NewDecrypter(context, ret.Decoder.Output())

	lg.Debug("Creating MAC processor...")
	ret.MACProcessor = NewMACProcessor(context, ret.Decrypter.Output())

	lg.Debug("Creating scheduler...")
	ret.Scheduler = NewScheduler(context, ret.MACProcessor.CommandNotifier())

	lg.Debug("Creating encoder...")
	ret.Encoder = NewEncoder(context, ret.Scheduler.Output(), forwarder.Input())

	return &ret
}
