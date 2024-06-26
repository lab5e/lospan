package main

import (
	b64 "encoding/base64"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/lab5e/lospan/pkg/protocol"
)

const numNames = 3

var payloadTemplate string
var payloadArgs [numNames]string

func init() {
	rand.Seed(time.Now().UnixNano())

	payloadTemplate = "It seemed to me, said %s the Sane, that any civilization that had so far lost its head as to need to include a set of detailed instructions for use in a packet of toothpicks, was no longer a civilization in which I could live and stay sane."

	payloadArgs[0] = "Hans Jørgen"
	payloadArgs[1] = "Bjørn"
	payloadArgs[2] = "Ståle"
}

// MessageGenerator generates LoRaWAN messages of various kinds.
type MessageGenerator struct {
	CorruptMICError     *Randomizer
	FrameCounterError   *Randomizer
	CorruptPayloadError *Randomizer
	MaxPayloadSize      int
	UseNumericalPayload bool
	deviceTick          float64
	function            func(float64) uint16
}

// NewMessageGenerator creates a new message generator
func randomFunc() func(float64) uint16 {
	switch rand.Intn(3) {
	case 0:
		return func(tick float64) uint16 {
			return uint16(math.Abs(math.Sin(100/tick*10)) * 100)
		}
	case 1:
		return func(tick float64) uint16 {
			return uint16((10*math.Sin(tick)*math.Cos(50*tick) + 10) * 10)
		}
	default:
		return func(tick float64) uint16 {
			return uint16(50*(math.Sin(1.0/50.0-tick*10*math.Cos(2*tick+1)+50)) + 200)
		}
	}
}

// NewMessageGenerator generates random messages
func NewMessageGenerator(config EagleConfig) MessageGenerator {
	return MessageGenerator{
		CorruptMICError:     NewRandomizer(config.CorruptMIC),
		FrameCounterError:   NewRandomizer(config.FrameCounterErrors),
		CorruptPayloadError: NewRandomizer(config.CorruptedPayload),
		MaxPayloadSize:      config.MaxPayloadSize,
		deviceTick:          0,
		function:            randomFunc(),
	}
}

// RandomMessageType returns a random message type; confirmed and/or unconfirmed
func (m *MessageGenerator) randomMessageType() protocol.MType {
	switch rand.Intn(2) {
	case 0:
		return protocol.ConfirmedDataUp
	case 1:
		return protocol.UnconfirmedDataUp
	}
	return protocol.Proprietary
}

func (m *MessageGenerator) buildPayload() []byte {
	if m.UseNumericalPayload {
		m.deviceTick += 0.1
		value := m.function(m.deviceTick)
		data := make([]byte, 2)
		data[0] = byte(value & 0x00FF)
		data[1] = byte(value >> 8)
		return data
	}
	payload := fmt.Sprintf(payloadTemplate, payloadArgs[rand.Intn(numNames)])
	payloadSize := rand.Intn(m.MaxPayloadSize-1) + 1
	return []byte(payload)[:payloadSize]
}

// Generate creates a new message, encodes it and returns a base64-encoded string. The second
// returned value indicates a valid message or not.
func (m *MessageGenerator) Generate(keys *DeviceKeys, fCnt uint16) (string, protocol.MType) {
	mt := m.randomMessageType()

	if m.CorruptPayloadError.Now() {
		buf := make([]byte, rand.Intn(1024))
		rand.Read(buf)
		return b64.StdEncoding.EncodeToString(buf), protocol.Proprietary
	}

	message := protocol.NewPHYPayload(mt)
	message.MACPayload.FHDR.DevAddr = keys.DevAddr
	message.MACPayload.FPort = uint8(rand.Intn(223) + 1)
	message.MACPayload.FHDR.FCnt = fCnt

	m.FrameCounterError.Maybe(func() {
		message.MACPayload.FHDR.FCnt = 0
		mt = protocol.Proprietary
	})

	message.MACPayload.FRMPayload = m.buildPayload()

	buffer, err := message.EncodeMessage(keys.NwkSKey, keys.AppSKey)
	if err != nil {
		log.Println("Error marshalling payload: ", err)
	}

	m.CorruptMICError.Maybe(func() {
		buffer[len(buffer)-1]++
		mt = protocol.Proprietary
	})

	return b64.StdEncoding.EncodeToString(buffer), mt
}
