// main.go - Katzenpost ping tool
// Copyright (C) 2018, 2019  David Stainton
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"flag"
	"fmt"

	client "github.com/hashcloak/Meson-client"
	"github.com/hashcloak/Meson-client/config"
	"github.com/katzenpost/core/crypto/ecdh"
)

func register(configFile string) (*config.Config, *ecdh.PrivateKey) {
	cfg, err := config.LoadFile(configFile)
	if err != nil {
		panic(err)
	}
	_, linkKey := client.AutoRegisterRandomClient(cfg)
	return cfg, linkKey
}

func main() {
	var configFile string
	var service string
	var count int
	flag.StringVar(&configFile, "c", "", "configuration file")
	flag.StringVar(&service, "s", "", "service name")
	flag.IntVar(&count, "n", 5, "count")
	flag.Parse()

	if service == "" {
		panic("must specify service name with -s")
	}

	cfg, linkKey := register(configFile)

	// create a client and connect to the mixnet Provider
	c, err := client.NewFromConfig(cfg, service)
	if err != nil {
		panic(err)
	}

	s, err := c.NewSession(linkKey)
	if err != nil {
		panic(err)
	}

	serviceDesc, err := s.GetService(service)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Sending %d Sphinx packet payloads to: %s@%s\n", count, serviceDesc.Name, serviceDesc.Provider)
	passed := 0
	failed := 0
	for i := 0; i < count; i++ {
		_, err := s.BlockingSendUnreliableMessage(serviceDesc.Name, serviceDesc.Provider, []byte(`Data encryption is used widely to protect the content of Internet
communications and enables the myriad of activities that are popular today,
from online banking to chatting with loved ones. However, encryption is not
sufficient to protect the meta-data associated with the communications.

Modern encrypted communication networks are vulnerable to traffic analysis and
can leak such meta-data as the social graph of users, their geographical
location, the timing of messages and their order, message size, and many other
kinds of meta-data.

Since 1979, there has been active academic research into communication
meta-data protection, also called anonymous communication networking, that has
produced various designs. Of these, mix networks are among the most practical
and can readily scale to millions of users.

The Mix Network workshop will focus on bringing together experts from
the research and practitioner communities to give technical lectures on key
Mix networking topics in relation to attacks, defences, and practical
applications and usage considerations.`))
		if err != nil {
			failed++
			fmt.Printf("Send error: %+v", err)
			continue
		}
		fmt.Println("Sended....")
		passed++
	}
	fmt.Printf("\n")

	percent := (float64(passed) * float64(100)) / float64(count)
	fmt.Printf("Success rate is %f percent %d/%d)\n", percent, passed, count)

	c.Shutdown()
}
