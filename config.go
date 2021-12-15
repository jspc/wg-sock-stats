package main

import (
	"github.com/BurntSushi/toml"
)

type Owner struct {
	Name  string `toml:"name"`
	Email string `toml:"email"`
}

type Config struct {
	// Longitude and Latitude point (respectively)
	// to the rough longitude and latitude
	// of this server. Leaving it blank is fine
	Longitude float64 `toml:"long"`
	Latitude  float64 `toml:"lat"`

	// CheckPTR governs whether we lookup a peer's IP
	// to provide some kind of human readable information
	CheckPTR bool `toml:"check_ptr"`

	// MapOwners governs whether we report the owner name/ email
	// address for a public key. This is useful for personal/ company
	// wireguard installations, where a public key is less useful than
	// a name.
	//
	// However, for a more privacy focused wireguard implementation, this
	// information is not helpful to share.
	//
	// The default value of this field is false, and must be explicitly
	// set, even if the Owners field is set
	MapOwners bool `toml:"map_owners"`

	// Owners maps a public key to name and email address. See MapOwners
	// for more information, including caveats
	Owners map[string]Owner `toml:"owners"`
}

func ParseConfig(fn string) (c Config, err error) {
	_, err = toml.DecodeFile(fn, &c)

	return
}
