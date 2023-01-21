package processor

//
//Copyright 2018 Telenor Digital AS
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
import (
	"github.com/lab5e/lospan/pkg/lg"
	"github.com/lab5e/lospan/pkg/protocol"
	"github.com/lab5e/lospan/pkg/server"
)

// Decoder is the process that decodes the bytes received from the gateway interface into go structs.
type Decoder struct {
	input   <-chan server.GatewayPacket
	output  chan server.LoRaMessage
	context *server.Context
}

// Start launches the decoder. It will terminate when the input channel
// is closed. On exit the output channel will be closed.
func (d *Decoder) Start() {
	for p := range d.input {
		go func(raw server.GatewayPacket) {
			// The initial message type isn't important
			decoded := protocol.NewPHYPayload(protocol.Proprietary)
			if err := decoded.UnmarshalBinary(raw.RawMessage); err != nil {
				lg.Info("Error unmarshalling payload: %v", err)
				return
			}
			context := server.FrameContext{
				GatewayContext: raw,
			}
			msg := server.LoRaMessage{
				Payload:      decoded,
				FrameContext: context,
			}
			d.output <- msg
		}(p)
	}
	lg.Debug("Input channel for Decoder closed. Terminating")
	close(d.output)
}

// Output returns the output channel from the decoder. This channel will receive
// a message every time a message is successfully decoded.
func (d *Decoder) Output() <-chan server.LoRaMessage {
	return d.output
}

// NewDecoder creates a new decoder.
func NewDecoder(context *server.Context, input <-chan server.GatewayPacket) *Decoder {
	return &Decoder{
		input:   input,
		output:  make(chan server.LoRaMessage),
		context: context,
	}
}
