// Copyright 2024 Testkube.
//
// Licensed as a Testkube Pro file under the Testkube Community
// License (the "License"); you may not use this file except in compliance with
// the License. You may obtain a copy of the License at
//
//	https://github.com/kubeshop/testkube/blob/main/licenses/TCL.txt

package devutils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"

	common2 "github.com/kubeshop/testkube/cmd/kubectl-testkube/commands/common"
	"github.com/kubeshop/testkube/cmd/kubectl-testkube/config"
	client2 "github.com/kubeshop/testkube/pkg/api/v1/client"
	"github.com/kubeshop/testkube/pkg/cloud/client"
)

type cloudObj struct {
	cfg       config.CloudContext
	envClient *client.EnvironmentsClient
	list      []client.Environment

	clientMu sync.Mutex
	client   client2.Client
	clientTs time.Time

	cmd *cobra.Command
}

func NewCloud(cfg config.CloudContext, cmd *cobra.Command) (*cloudObj, error) {
	if cfg.ApiKey == "" || cfg.OrganizationId == "" || cfg.OrganizationName == "" {
		return nil, errors.New("login to the organization first")
	}
	if strings.HasPrefix(cfg.AgentUri, "https://") {
		cfg.AgentUri = strings.TrimPrefix(cfg.AgentUri, "https://")
		if !regexp.MustCompile(`:\d+$`).MatchString(cfg.AgentUri) {
			cfg.AgentUri += ":443"
		}
	} else if strings.HasPrefix(cfg.AgentUri, "http://") {
		cfg.AgentUri = strings.TrimPrefix(cfg.AgentUri, "http://")
		if !regexp.MustCompile(`:\d+$`).MatchString(cfg.AgentUri) {
			cfg.AgentUri += ":80"
		}
	}
	// TODO: FIX THAT
	if strings.HasPrefix(cfg.AgentUri, "api.") {
		cfg.AgentUri = "agent." + strings.TrimPrefix(cfg.AgentUri, "api.")
	}
	envClient := client.NewEnvironmentsClient(cfg.ApiUri, cfg.ApiKey, cfg.OrganizationId)
	obj := &cloudObj{
		cfg:       cfg,
		envClient: envClient,
		cmd:       cmd,
	}

	err := obj.UpdateList()
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (c *cloudObj) List() []client.Environment {
	return c.list
}

func (c *cloudObj) ListObsolete() []client.Environment {
	obsolete := make([]client.Environment, 0)
	for _, env := range c.list {
		if !env.Connected {
			obsolete = append(obsolete, env)
		}
	}
	return obsolete
}

func (c *cloudObj) UpdateList() error {
	list, err := c.envClient.List()
	if err != nil {
		return err
	}
	result := make([]client.Environment, 0)
	for i := range list {
		if strings.HasPrefix(list[i].Name, "devbox-") {
			result = append(result, list[i])
		}
	}
	c.list = result
	return nil
}

func (c *cloudObj) Client(environmentId string) (client2.Client, error) {
	c.clientMu.Lock()
	defer c.clientMu.Unlock()

	if c.client == nil || c.clientTs.Add(5*time.Minute).Before(time.Now()) {
		common2.GetClient(c.cmd) // refresh token
		var err error
		c.client, err = client2.GetClient(client2.ClientCloud, client2.Options{
			Insecure:           c.AgentInsecure(),
			ApiUri:             c.ApiURI(),
			CloudApiKey:        c.ApiKey(),
			CloudOrganization:  c.cfg.OrganizationId,
			CloudEnvironment:   environmentId,
			CloudApiPathPrefix: fmt.Sprintf("/organizations/%s/environments/%s/agent", c.cfg.OrganizationId, environmentId),
		})
		if err != nil {
			return nil, err
		}
		c.clientTs = time.Now()
	}
	return c.client, nil
}

func (c *cloudObj) AgentURI() string {
	return c.cfg.AgentUri
}

func (c *cloudObj) AgentInsecure() bool {
	return strings.HasPrefix(c.cfg.ApiUri, "http://")
}

func (c *cloudObj) ApiURI() string {
	return c.cfg.ApiUri
}

func (c *cloudObj) ApiKey() string {
	return c.cfg.ApiKey
}

func (c *cloudObj) ApiInsecure() bool {
	return strings.HasPrefix(c.cfg.ApiUri, "http://")
}

func (c *cloudObj) DashboardUrl(id, path string) string {
	return strings.TrimSuffix(fmt.Sprintf("%s/organization/%s/environment/%s/", c.cfg.UiUri, c.cfg.OrganizationId, id)+strings.TrimPrefix(path, "/"), "/")
}

func (c *cloudObj) CreateEnvironment(name string) (*client.Environment, error) {
	env, err := c.envClient.Create(client.Environment{
		Name:           name,
		Owner:          c.cfg.OrganizationId,
		OrganizationId: c.cfg.OrganizationId,
	})
	if err != nil {
		return nil, err
	}
	// TODO: POST request is not returning slug - if it will, delete the fallback path
	if env.Slug != "" {
		c.list = append(c.list, env)
	} else {
		err = c.UpdateList()
		if err != nil {
			return nil, err
		}
		for i := range c.list {
			if c.list[i].Id == env.Id {
				env = c.list[i]
				break
			}
		}
	}
	// Hack to build proper URLs even when slug is missing
	if env.Slug == "" {
		env.Slug = env.Id
	}
	return &env, nil
}

func (c *cloudObj) DeleteEnvironment(id string) error {
	return c.envClient.Delete(id)
}
