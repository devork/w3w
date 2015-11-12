# W3W Go Client

This is a simple Go client for the [W3W API](http://developer.what3words.com/api).

## Usage

Pull the code you need into your project:

`go get github.com/devork/w3w`

## Samples

Import the package as:

```
import "github.com/devork/w3w"
```

### Create a new `w3w` struct:

```
w, err := w3w.New("APIKEY", &w3w.Options{"en", true})
```
The new call allows a default set of options to be included with each call to the API. These can be
overridden on each call. if not provided, the defaults of `lang=en` and `corners=true` are used.

### Fetch the position of a W3W

```
pos, err := w.Words(w3w.What3Words{"prom", "cape", "pump"}, nil)
```

### Fetch W3W for a LatLng

```
pos, err = w.Position(w3w.LatLng{51.484463, -0.195405}, nil)
```

### Fetch the available languages for a W3W position

```
langs, err := w.LangsW3W(w3w.What3Words{"index", "home", "raft"}, nil)
```

### Override the defaults

```
opts := &w3w.Options{"de", false}
pos, err = w.Position(w3w.LatLng{51.484463, -0.195405}, opts)
```
