package w3w

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	endpoint string = "https://api.what3words.com"
)

// Default error codes
var (
	ErrNoAPIKey = errors.New("No API Key specified")
)

var (
	defs = Options{"en", false}
)

// ----------------------------------------------------------------------------
// What3Words type
// ----------------------------------------------------------------------------

// What3Words holds the 3 word position returned or passed to the server
type What3Words [3]string

// ----------------------------------------------------------------------------
// LatLng type
// ----------------------------------------------------------------------------

// LatLng type
type LatLng [2]float64

// Lat extracts the latitude from the latLng type
func (ll *LatLng) Lat() float64 {
	return ll[0]
}

// Lng extracts the longitude from the latLng type
func (ll *LatLng) Lng() float64 {
	return ll[1]
}

// ----------------------------------------------------------------------------
// BBox type
// ----------------------------------------------------------------------------

// BBox represents the corners of the W3W square
type BBox [2]*LatLng

// SW returns the southwest coordinates of the square
func (b *BBox) SW() *LatLng {
	return b[0]
}

// NE returns the northeast coordinates of the square
func (b *BBox) NE() *LatLng {
	return b[1]
}

// ----------------------------------------------------------------------------
// Position struct
// ----------------------------------------------------------------------------

// Position holds the server response from a query
type Position struct {
	Type     string     `json:"type"`
	Words    What3Words `json:"words"`
	Position *LatLng    `json:"position"`
	Corners  *BBox      `json:"corners"`
	Language string     `json:"language"`
}

// ----------------------------------------------------------------------------
// Language struct
// ----------------------------------------------------------------------------

// Language holds a single language type from the languages call
type Language struct {
	Code string `json:"code"`
	Name string `json:"name_display"`
}

// ----------------------------------------------------------------------------
// Languages struct
// ----------------------------------------------------------------------------

// Languages holds the full languages response from the server
type Languages struct {
	Languages []Language `json:"languages"`
}

// ----------------------------------------------------------------------------
// Options struct
// ----------------------------------------------------------------------------

// Options holds the various optional query params for all W3W calls.
type Options struct {
	Lang    string
	Corners bool
}

func (o *Options) add(v *url.Values) {
	if o.Lang == "" {
		v.Set("lang", "en")
	} else {
		v.Set("lang", o.Lang)
	}

	if o.Corners {
		v.Set("corners", "true")
	}

}

// ----------------------------------------------------------------------------
// W3W struct
// ----------------------------------------------------------------------------

// W3W holds details about the connection to the W3W service
type W3W struct {
	apikey   string
	client   *http.Client
	defaults *Options
}

// New returns a W3W with the given API key. The options defaults allows for sensible defaults to be
// associated with each W3W call.
//
// if the key is missing or empty, the returned error is `ErrNoAPIKey`
func New(apikey string, defaults *Options) (*W3W, error) {
	if apikey == "" || strings.TrimSpace(apikey) == "" {
		return nil, ErrNoAPIKey
	}

	if defaults == nil {
		return &W3W{apikey, &http.Client{}, &defs}, nil
	}

	return &W3W{apikey, &http.Client{}, defaults}, nil

}

// Words converts a 3 word string to LatLng position
func (w *W3W) Words(words What3Words, opts *Options) (*Position, error) {
	vals := url.Values{}

	vals.Set("key", w.apikey)
	vals.Set("string", strings.Join(words[:], "."))

	pos, err := w.exec(endpoint+"/w3w", &vals, opts, &Position{})

	return pos.(*Position), err
}

// Position converts a 3 word string to LatLng position
func (w *W3W) Position(ll LatLng, opts *Options) (*Position, error) {
	vals := url.Values{}

	vals.Set("key", w.apikey)
	vals.Set("position", fmt.Sprintf("%.15f,%.15f", ll[0], ll[1]))

	pos, err := w.exec(endpoint+"/position", &vals, opts, &Position{})
	return pos.(*Position), err
}

// LangsW3W obtains the list of available 3 word lanagues for a given W3W position
func (w *W3W) LangsW3W(words What3Words, opts *Options) (*Languages, error) {
	vals := url.Values{}
	vals.Set("key", w.apikey)
	vals.Set("string", strings.Join(words[:], "."))

	langs, err := w.exec(endpoint+"/get-languages", &vals, opts, &Languages{[]Language{}})

	return langs.(*Languages), err
}

// LangsPos obtains the list of available 3 word lanagues for a given LatLng position
func (w *W3W) LangsPos(ll LatLng, opts *Options) (*Languages, error) {
	vals := url.Values{}
	vals.Set("key", w.apikey)
	vals.Set("position", fmt.Sprintf("%.15f,%.15f", ll[0], ll[1]))

	langs, err := w.exec(endpoint+"/get-languages", &vals, opts, &Languages{[]Language{}})

	return langs.(*Languages), err
}

func (w *W3W) exec(url string, vals *url.Values, opts *Options, in interface{}) (interface{}, error) {

	if opts != nil {
		opts.add(vals)
	} else {
		w.defaults.add(vals)
	}

	req, err := http.NewRequest("GET", url+"?"+vals.Encode(), nil)
	req.Header.Add("Accept", "application/json")

	resp, err := w.client.Do(req)

	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(resp.Body).Decode(in)

	if err != nil {
		return nil, err
	}

	return in, nil
}
