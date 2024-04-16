package internal

import "bytes"

// setDecoder handles both parameters and set decoding [] ().
// There was no reason to separate them into two different methods.
type setDecoder struct {
	*scanner
	activeBuffer *bytes.Buffer
	tokens       []string
}

func (d setDecoder) decode() []string {
	for {
		c, stop := d.nextChar()
		if stop {
			switch c {
			case 0, ']', ')':
				if d.activeBuffer.Len() > 0 {
					d.tokens = append(d.tokens, d.activeBuffer.String())
					d.activeBuffer.Reset()
				}
				return d.tokens
			case ',':
				if d.activeBuffer.Len() > 0 {
					d.tokens = append(d.tokens, d.activeBuffer.String())
					d.activeBuffer.Reset()
				}
				continue
			case '"':
				d.tokens = append(d.tokens, stringDecoder{scanner: d.scanner, activeBuffer: &bytes.Buffer{}}.decode())
			case '~':
				d.activeBuffer.WriteByte(c)
			default:
				// Default skips through the character
				// Probably want to return an error in a future refactoring
			}
		} else {
			d.activeBuffer.WriteByte(c)
		}
	}
}

// stringDecoder decodes quoted strings, tokens are still decoded using the tokenDecoder.
type stringDecoder struct {
	activeBuffer *bytes.Buffer
	*scanner
}

func (d stringDecoder) decode() string {
	var prevByte byte
	var prevByte2 byte
	// Write a " to force it to a quote, tokens that don't use string decoder can't use "'s as a legal character.
	d.activeBuffer.WriteByte('"')
	for {
		c, stop := d.nextChar()
		if stop {
			if (c == '"' && prevByte != '\\') || c == 0 {
				// Flush out the string if we hit the end of the string and an escape wasn't provided
				return d.activeBuffer.String()
			}
		}
		if c != '\\' || (c == '\\' && prevByte == '\\') {
			// Strip escaped characters
			d.activeBuffer.WriteByte(c)
		}
		if c == '\\' {
			prevByte2 = prevByte
		}
		if prevByte == '\\' {
			prevByte = prevByte2
		} else {
			prevByte = c
		}
	}
}

// tokenDecoder handles decoding a token (any non-quoted string)
type tokenDecoder struct {
	*scanner
	activeBuffer *bytes.Buffer
}

func (d tokenDecoder) decode() string {
	for {
		c, stop := d.nextChar()
		if stop {
			if c != 0 {
				d.scanner.pos--
			}
			return d.activeBuffer.String()
		}
		d.activeBuffer.WriteByte(c)
	}
}
