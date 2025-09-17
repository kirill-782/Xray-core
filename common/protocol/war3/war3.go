package war3

import (
	"encoding/binary"
	"errors"

	"github.com/xtls/xray-core/common"
)

// Общая шапка сниффера для Warcraft III протоколов (W3GS/GPS).
type SniffHeader struct {
	Kind string // "w3gs" или "gps"
}

func (h *SniffHeader) Protocol() string { return "warcraft3" }
func (h *SniffHeader) Domain() string   { return "" }

var errNotWar3 = errors.New("not warcraft3 (w3gs/gps) header")

// SniffWar3 определяет, принадлежит ли первый пакет к W3GS или GPS.
// Формат первых сообщений:
//
//	W3GS: UINT8 0xF7 (247), UINT8 0x1E, UINT16 len(le), body[len-4]
//	GPS : UINT8 0xF8 (248), UINT8 0x02 или 0x0A, UINT16 len(le), body[len-4]
func SniffWar3(b []byte) (*SniffHeader, error) {
	if len(b) < 4 {
		return nil, common.ErrNoClue
	}

	// Проверяем маркер/тип.
	switch b0, b1 := b[0], b[1]; b0 {
	case 0xF7: // W3GS
		if b1 != 0x1E {
			return nil, errNotWar3
		}
		// длина хранится в little-endian
		total := int(binary.LittleEndian.Uint16(b[2:4]))
		if total < 4 {
			return nil, errNotWar3
		}
		if total > len(b) {
			// Пакет ещё не полностью пришёл — нет уверенности.
			return nil, common.ErrNoClue
		}
		return &SniffHeader{Kind: "w3gs"}, nil

	case 0xF8: // GPS
		if b1 != 0x02 && b1 != 0x0A {
			return nil, errNotWar3
		}
		total := int(binary.LittleEndian.Uint16(b[2:4]))
		if total < 4 {
			return nil, errNotWar3
		}
		if total > len(b) {
			return nil, common.ErrNoClue
		}
		return &SniffHeader{Kind: "gps"}, nil

	default:
		return nil, errNotWar3
	}
}
