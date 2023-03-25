// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package proxmox

import (
	"crypto/tls"
	"log"
	"os"

	"github.com/Telmate/proxmox-api-go/proxmox"
)

func newProxmoxClient(config Config) (*proxmox.Client, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: config.SkipCertValidation,
	}

	client, err := proxmox.NewClient(config.proxmoxURL.String(), nil, config.ProxmoxHttpHeaders, tlsConfig, config.ProxmoxProxyServer, int(config.TaskTimeout.Seconds()))
	if err != nil {
		return nil, err
	}

	*proxmox.Debug = config.PackerDebug || os.Getenv("PACKER_LOG") == "1"
	log.Printf("Proxmox API debug %t", *proxmox.Debug)

	if config.Token != "" {
		// configure token auth
		log.Print("using token auth")
		client.SetAPIToken(config.Username, config.Token)
	} else {
		// fallback to login if not using tokens
		log.Print("using password auth")
		err = client.Login(config.Username, config.Password, "")
		if err != nil {
			return nil, err
		}
	}

	return client, nil
}
