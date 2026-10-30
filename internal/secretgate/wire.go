package secretgate

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/godbus/dbus/v5"
)

// maxMessage is the D-Bus limit on a whole message (2^27).
const maxMessage = 1 << 27

// readFrame reads one whole D-Bus message as it was sent: the relay forwards bytes, never a
// re-encoding, so what reaches the keyring is exactly what its caller wrote.
func readFrame(r io.Reader) ([]byte, error) {
	var hdr [16]byte
	if _, err := io.ReadFull(r, hdr[:]); err != nil {
		return nil, err
	}
	var order binary.ByteOrder
	switch hdr[0] {
	case 'l':
		order = binary.LittleEndian
	case 'B':
		order = binary.BigEndian
	default:
		return nil, fmt.Errorf("bad endianness byte %#x", hdr[0])
	}
	body := uint64(order.Uint32(hdr[4:8]))
	fields := uint64(order.Uint32(hdr[12:16]))
	total := 16 + (fields+7)&^7 + body
	if total > maxMessage {
		return nil, fmt.Errorf("message of %d bytes exceeds the D-Bus limit", total)
	}
	frame := make([]byte, total)
	copy(frame, hdr[:])
	if _, err := io.ReadFull(r, frame[16:]); err != nil {
		return nil, err
	}
	return frame, nil
}

func decodeFrame(frame []byte) (*dbus.Message, error) {
	return dbus.DecodeMessage(bytes.NewReader(frame))
}

func header[T any](m *dbus.Message, f dbus.HeaderField) T {
	var zero T
	v, ok := m.Headers[f]
	if !ok {
		return zero
	}
	t, _ := v.Value().(T)
	return t
}

// encodeWithSerial encodes m with serial: godbus keeps a message's serial unexported.
func encodeWithSerial(m *dbus.Message, serial uint32) ([]byte, error) {
	var buf bytes.Buffer
	if err := m.EncodeTo(&buf, binary.LittleEndian); err != nil {
		return nil, err
	}
	raw := buf.Bytes()
	binary.LittleEndian.PutUint32(raw[8:12], serial)
	return raw, nil
}

// serverAuth is the bus side of the SASL handshake with the provider: EXTERNAL only, as uid (already
// checked by SO_PEERCRED), no fd passing. It returns the reader positioned at the first message.
func serverAuth(rw io.ReadWriter, uid uint32) (*bufio.Reader, error) {
	br := bufio.NewReader(rw)
	if b, err := br.ReadByte(); err != nil || b != 0 {
		return nil, fmt.Errorf("sasl: missing initial NUL byte")
	}
	want := hex.EncodeToString([]byte(strconv.FormatUint(uint64(uid), 10)))
	guid := make([]byte, 16)
	if _, err := rand.Read(guid); err != nil {
		return nil, err
	}
	authed := false
	waitingData := false
	for range 16 {
		line, err := br.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("sasl: %w", err)
		}
		cmd, arg, _ := strings.Cut(strings.TrimRight(line, "\r\n"), " ")
		reply := "ERROR\r\n"
		switch {
		case cmd == "AUTH" && strings.HasPrefix(arg, "EXTERNAL"):
			id := strings.TrimSpace(strings.TrimPrefix(arg, "EXTERNAL"))
			switch {
			case id == "":
				waitingData, reply = true, "DATA\r\n"
			case id == want:
				authed, reply = true, "OK "+hex.EncodeToString(guid)+"\r\n"
			default:
				reply = "REJECTED EXTERNAL\r\n"
			}
		case cmd == "AUTH" || cmd == "CANCEL":
			authed, waitingData, reply = false, false, "REJECTED EXTERNAL\r\n"
		case cmd == "DATA" && waitingData:
			waitingData = false
			if arg == "" || arg == want {
				authed, reply = true, "OK "+hex.EncodeToString(guid)+"\r\n"
			} else {
				reply = "REJECTED EXTERNAL\r\n"
			}
		case cmd == "BEGIN" && authed:
			return br, nil
		}
		if _, err := io.WriteString(rw, reply); err != nil {
			return nil, err
		}
	}
	return nil, errors.New("sasl: too many commands")
}

// clientAuth authenticates to the user bus as the daemon's own uid (root), without fd passing.
func clientAuth(rw io.ReadWriter, uid uint32) (*bufio.Reader, error) {
	id := hex.EncodeToString([]byte(strconv.FormatUint(uint64(uid), 10)))
	if _, err := io.WriteString(rw, "\x00AUTH EXTERNAL "+id+"\r\n"); err != nil {
		return nil, err
	}
	br := bufio.NewReader(rw)
	line, err := br.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("sasl: %w", err)
	}
	if !strings.HasPrefix(line, "OK ") {
		return nil, fmt.Errorf("sasl: bus refused EXTERNAL: %q", strings.TrimSpace(line))
	}
	if _, err := io.WriteString(rw, "BEGIN\r\n"); err != nil {
		return nil, err
	}
	return br, nil
}
