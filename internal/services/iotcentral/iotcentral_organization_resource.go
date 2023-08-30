// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iotcentral

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iotcentral/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	dataplane "github.com/tombuildsstuff/kermit/sdk/iotcentral/2022-10-31-preview/iotcentral"
)

type IotCentralOrganizationResource struct{}

var (
	_ sdk.ResourceWithUpdate = IotCentralOrganizationResource{}
)

type IotCentralOrganizationModel struct {
	SubDomain      string `tfschema:"sub_domain"`
	OrganizationId string `tfschema:"organization_id"`
	DisplayName    string `tfschema:"display_name"`
	Parent         string `tfschema:"parent"`
}

func (r IotCentralOrganizationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"sub_domain": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"organization_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"display_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"parent": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
	}
}

func (r IotCentralOrganizationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r IotCentralOrganizationResource) ResourceType() string {
	return "azurerm_iotcentral_organization"
}

func (r IotCentralOrganizationResource) ModelObject() interface{} {
	return &IotCentralOrganizationModel{}
}

func (r IotCentralOrganizationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.NestedItemId
}

func (r IotCentralOrganizationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTCentral
			var state IotCentralOrganizationModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			orgClient, err := client.OrganizationsClient(ctx, state.SubDomain)
			if err != nil {
				return fmt.Errorf("creating organization client: %+v", err)
			}

			model := dataplane.Organization{
				DisplayName: &state.DisplayName,
			}

			if state.Parent != "" {
				model.Parent = &state.Parent
			}

			org, err := orgClient.Create(ctx, state.OrganizationId, model)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", state.OrganizationId, err)
			}

			baseUrl, err := parse.ParseBaseUrl(state.SubDomain)
			if err != nil {
				return err
			}

			orgId, err := parse.NewNestedItemID(baseUrl, parse.NestedItemTypeOrganization, *org.ID)
			if err != nil {
				return err
			}

			metadata.SetID(orgId)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r IotCentralOrganizationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTCentral
			id, err := parse.ParseNestedItemID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			orgClient, err := client.OrganizationsClient(ctx, id.SubDomain)
			if err != nil {
				return fmt.Errorf("creating organization client: %+v", err)
			}

			org, err := orgClient.Get(ctx, id.Id)
			if err != nil {
				if org.ID == nil || *org.ID == "" {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := IotCentralOrganizationModel{
				SubDomain:      id.SubDomain,
				OrganizationId: *org.ID,
				DisplayName:    *org.DisplayName,
			}

			if org.Parent != nil {
				state.Parent = *org.Parent
			}

			return metadata.Encode(&state)
		},
		Timeout: 5 * time.Minute,
	}
}

func (r IotCentralOrganizationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTCentral
			id, err := parse.ParseNestedItemID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state IotCentralOrganizationModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			orgClient, err := client.OrganizationsClient(ctx, id.SubDomain)
			if err != nil {
				return fmt.Errorf("creating organization client: %+v", err)
			}

			existing, err := orgClient.Get(ctx, id.Id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if metadata.ResourceData.HasChange("display_name") {
				existing.DisplayName = &state.DisplayName
			}

			if metadata.ResourceData.HasChange("parent") {
				if state.Parent != "" {
					// A bug in the Client, currently doesn't allow you to unset the parent
					// autorest/azure: Service returned an error. Status=422 Code="InvalidBody" Message="ID exceeds maximum character limit of 48 ..."
					existing.Parent = &state.Parent
				}
			}

			_, err = orgClient.Update(ctx, *existing.ID, existing, "*")
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r IotCentralOrganizationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTCentral
			id, err := parse.ParseNestedItemID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			orgClient, err := client.OrganizationsClient(ctx, id.SubDomain)
			if err != nil {
				return fmt.Errorf("creating organization client: %+v", err)
			}

			_, err = orgClient.Remove(ctx, id.Id)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
