// Package i3bar provides a Go library for i3bar JSON protocol support.
//
// structs and much of the inspiration for this project comes from:
//
// https://github.com/davidscholberg/go-i3barjson
//
package i3bar

import (
	"encoding/json"
	"io"
)

// Header represents the header of an i3bar message.
type Header struct {
	Version     int  `json:"version"`
	StopSignal  int  `json:"stop_signal,omitempty"`
	ContSignal  int  `json:"cont_signal,omitempty"`
	ClickEvents bool `json:"click_events,omitempty"`
}

// Block represents a single block of an i3bar message.
type Block struct {
	FullText            string `json:"full_text"`
	ShortText           string `json:"short_text,omitempty"`
	Color               string `json:"color,omitempty"`
	MinWidth            string `json:"min_width,omitempty"`
	Align               string `json:"align,omitempty"`
	Name                string `json:"name,omitempty"`
	Instance            string `json:"instance,omitempty"`
	Urgent              bool   `json:"urgent,omitempty"`
	Separator           bool   `json:"separator"`
	SeparatorBlockWidth int    `json:"separator_block_width,omitempty"`
	Markup              string `json:"markup,omitempty"`
}

// StatusLine represents a full i3bar status line.
type StatusLine []*Block

// Click represents an i3bar click event.
type Click struct {
	Name     string `json:"name"`
	Instance string `json:"instance"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Button   int    `json:"button"`
}

// Encode encodes the structs to the writer as they arrive through the channel. Close the channel to stop writing.
func Encode(w io.Writer, header *Header, slChan chan StatusLine) error {
	content, err := json.Marshal(header)
	if err != nil {
		return err
	}
	if _, err := w.Write(content); err != nil {
		return err
	}

	if _, err := w.Write([]byte("[")); err != nil {
		return err
	}

	first := true

	for sl := range slChan {
		if !first {
			if _, err := w.Write([]byte(",")); err != nil {
				return err
			}
		} else {
			first = false
		}

		if err := json.NewEncoder(w).Encode(sl); err != nil {
			return err
		}
	}

	if _, err := w.Write([]byte("]")); err != nil {
		return err
	}

	return nil
}
