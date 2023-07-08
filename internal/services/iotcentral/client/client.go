// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-sdk/resource-manager/iotcentral/2021-11-01-preview/apps"
	authWrapper "github.com/hashicorp/go-azure-sdk/sdk/auth/autorest"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	iotcentralDataplane "github.com/tombuildsstuff/kermit/sdk/iotcentral/2022-10-31-preview/iotcentral"
)

type Client struct {
	AppsClient          *apps.AppsClient
	authorizerFunc      common.ApiAuthorizerFunc
	configureClientFunc func(c *autorest.Client, authorizer autorest.Authorizer)
	endpoint            string
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	appsClient, err := apps.NewAppsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Apps Client: %+v", err)
	}
	o.Configure(appsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		AppsClient:          appsClient,
		authorizerFunc:      o.Authorizers.AuthorizerFunc,
		configureClientFunc: o.ConfigureClient,
		endpoint:            "https://apps.azureiotcentral.com",
	}, nil
}

func (c *Client) OrganizationsClient(ctx context.Context, subdomain string) (*iotcentralDataplane.OrganizationsClient, error) {
	api := environments.NewApiEndpoint("IotCentral", c.endpoint, nil)
	iotCentralAuth, err := c.authorizerFunc(api)
	if err != nil {
		return nil, fmt.Errorf("obtaining auth token for %q: %+v", c.endpoint, err)
	}

	client := iotcentralDataplane.NewOrganizationsClient(subdomain)
	c.configureClientFunc(&client.Client, authWrapper.AutorestAuthorizer(iotCentralAuth))

	return &client, nil
}

func (c *Client) RolesClient(ctx context.Context, subdomain string) (*iotcentralDataplane.RolesClient, error) {
	api := environments.NewApiEndpoint("IotCentral", c.endpoint, nil)
	iotCentralAuth, err := c.authorizerFunc(api)
	if err != nil {
		return nil, fmt.Errorf("obtaining auth token for %q: %+v", c.endpoint, err)
	}

	client := iotcentralDataplane.NewRolesClient(subdomain)
	c.configureClientFunc(&client.Client, authWrapper.AutorestAuthorizer(iotCentralAuth))

	return &client, nil
}
