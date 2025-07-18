package webrtc

import (
	"bytes"
	"encoding/binary"
	"strconv"

	"github.com/m1k1o/neko/server/pkg/types"

	"github.com/rs/zerolog"
)

const (
	OP_MOVE     = 0x01
	OP_SCROLL   = 0x02
	OP_KEY_DOWN = 0x03
	OP_KEY_UP   = 0x04
	OP_KEY_CLK  = 0x05
)

type PayloadHeader struct {
	Event  uint8
	Length uint16
}

type PayloadMove struct {
	PayloadHeader
	X uint16
	Y uint16
}

type PayloadScroll struct {
	PayloadHeader
	X int16
	Y int16
}

type PayloadKey struct {
	PayloadHeader
	Key uint64 // TODO: uint32
}

func (manager *WebRTCManagerCtx) handleLegacy(
	logger zerolog.Logger, data []byte,
	session types.Session,
) error {
	// continue only if session is host
	if !session.LegacyIsHost() {
		return nil
	}

	buffer := bytes.NewBuffer(data)
	header := &PayloadHeader{}
	hbytes := make([]byte, 3)

	if _, err := buffer.Read(hbytes); err != nil {
		return err
	}

	if err := binary.Read(bytes.NewBuffer(hbytes), binary.LittleEndian, header); err != nil {
		return err
	}

	buffer = bytes.NewBuffer(data)

	switch header.Event {
	case OP_MOVE:
		payload := &PayloadMove{}
		if err := binary.Read(buffer, binary.LittleEndian, payload); err != nil {
			return err
		}

		manager.desktop.Move(int(payload.X), int(payload.Y))
	case OP_SCROLL:
		payload := &PayloadScroll{}
		if err := binary.Read(buffer, binary.LittleEndian, payload); err != nil {
			return err
		}

		logger.
			Debug().
			Str("x", strconv.Itoa(int(payload.X))).
			Str("y", strconv.Itoa(int(payload.Y))).
			Msg("scroll")

		manager.desktop.Scroll(int(payload.X), int(payload.Y), false)
	case OP_KEY_DOWN:
		payload := &PayloadKey{}
		if err := binary.Read(buffer, binary.LittleEndian, payload); err != nil {
			return err
		}

		if payload.Key < 8 {
			err := manager.desktop.ButtonDown(uint32(payload.Key))
			if err != nil {
				logger.Warn().Err(err).Msg("button down failed")
				return nil
			}

			logger.Debug().Msgf("button down %d", payload.Key)
		} else {
			err := manager.desktop.KeyDown(uint32(payload.Key))
			if err != nil {
				logger.Warn().Err(err).Msg("key down failed")
				return nil
			}

			logger.Debug().Msgf("key down %d", payload.Key)
		}
	case OP_KEY_UP:
		payload := &PayloadKey{}
		err := binary.Read(buffer, binary.LittleEndian, payload)
		if err != nil {
			return err
		}

		if payload.Key < 8 {
			err := manager.desktop.ButtonUp(uint32(payload.Key))
			if err != nil {
				logger.Warn().Err(err).Msg("button up failed")
				return nil
			}

			logger.Debug().Msgf("button up %d", payload.Key)
		} else {
			err := manager.desktop.KeyUp(uint32(payload.Key))
			if err != nil {
				logger.Warn().Err(err).Msg("key up failed")
				return nil
			}

			logger.Debug().Msgf("key up %d", payload.Key)
		}
	case OP_KEY_CLK:
		// unused
		break
	}

	return nil
}
