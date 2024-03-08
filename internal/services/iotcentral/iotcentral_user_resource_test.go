// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iotcentral_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/iotcentral/2021-11-01-preview/apps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	iotcentral "github.com/hashicorp/terraform-provider-azurerm/internal/services/iotcentral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iotcentral/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

const (
	appAdminRoleId   = "ca310b8d-2f4a-44e0-a36e-957c202cd8d4"
	appBuilderRoleId = "344138e9-8de4-4497-8c54-5237e96d6aaf"
	orgAdminRoleId   = "c495eb57-eb18-489e-9802-62c474e5645c"
)

type IotCentralUserResource struct{}

func TestAccIoTCentralUser_basic_type_servicePrincipal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotcentral_user", "test")
	r := IotCentralUserResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic_type_servicePrincipal(data, appAdminRoleId),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("type").HasValue("ServicePrincipal"),
				check.That(data.ResourceName).Key("object_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("tenant_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("role.0.role_id").HasValue(appAdminRoleId),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIoTCentralUser_basic_type_group(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotcentral_user", "test")
	r := IotCentralUserResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic_type_group(data, appAdminRoleId),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("type").HasValue("Group"),
				check.That(data.ResourceName).Key("object_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("tenant_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("role.0.role_id").HasValue(appAdminRoleId),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIoTCentralUser_basic_type_email(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotcentral_user", "test")
	r := IotCentralUserResource{}

	email := fmt.Sprintf("basic%d@example.ex", data.RandomInteger)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic_type_email(data, email, appAdminRoleId),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("type").HasValue("Email"),
				check.That(data.ResourceName).Key("email").HasValue(email),
				check.That(data.ResourceName).Key("role.0.role_id").HasValue(appAdminRoleId),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIoTCentralUser_complete_type_servicePrincipal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotcentral_user", "test")
	r := IotCentralUserResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete_type_servicePrincipal(data, orgAdminRoleId),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("type").HasValue("ServicePrincipal"),
				check.That(data.ResourceName).Key("object_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("tenant_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("role.0.role_id").HasValue(orgAdminRoleId),
				check.That(data.ResourceName).Key("role.0.organization_id").IsNotEmpty(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIoTCentralUser_complete_type_group(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotcentral_user", "test")
	r := IotCentralUserResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete_type_group(data, orgAdminRoleId),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("type").HasValue("Group"),
				check.That(data.ResourceName).Key("object_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("tenant_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("role.0.role_id").HasValue(orgAdminRoleId),
				check.That(data.ResourceName).Key("role.0.organization_id").IsNotEmpty(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIoTCentralUser_complete_type_email(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotcentral_user", "test")
	r := IotCentralUserResource{}

	email := fmt.Sprintf("complete%d@example.ex", data.RandomInteger)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete_type_email(data, email, orgAdminRoleId),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("type").HasValue("Email"),
				check.That(data.ResourceName).Key("email").HasValue(email),
				check.That(data.ResourceName).Key("role.0.role_id").HasValue(orgAdminRoleId),
				check.That(data.ResourceName).Key("role.0.organization_id").IsNotEmpty(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIoTCentralUser_basicUpdateRole_type_servicePrincipal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotcentral_user", "test")
	r := IotCentralUserResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic_type_servicePrincipal(data, appAdminRoleId),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("type").HasValue("ServicePrincipal"),
				check.That(data.ResourceName).Key("role.0.role_id").HasValue(appAdminRoleId),
			),
		},
		{
			Config: r.basic_type_servicePrincipal(data, appBuilderRoleId),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("type").HasValue("ServicePrincipal"),
				check.That(data.ResourceName).Key("role.0.role_id").HasValue(appBuilderRoleId),
			),
		},
	})
}

func TestAccIoTCentralUser_basicUpdateRole_type_group(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotcentral_user", "test")
	r := IotCentralUserResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic_type_group(data, appAdminRoleId),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("type").HasValue("Group"),
				check.That(data.ResourceName).Key("role.0.role_id").HasValue(appAdminRoleId),
			),
		},
		{
			Config: r.basic_type_group(data, appBuilderRoleId),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("type").HasValue("Group"),
				check.That(data.ResourceName).Key("role.0.role_id").HasValue(appBuilderRoleId),
			),
		},
	})
}

func TestAccIoTCentralUser_basicUpdateRole_type_email(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotcentral_user", "test")
	r := IotCentralUserResource{}

	email := fmt.Sprintf("basic%d@example.ex", data.RandomInteger)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic_type_email(data, email, appAdminRoleId),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("type").HasValue("Email"),
				check.That(data.ResourceName).Key("role.0.role_id").HasValue(appAdminRoleId),
			),
		},
		{
			Config: r.basic_type_email(data, email, appBuilderRoleId),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("type").HasValue("Email"),
				check.That(data.ResourceName).Key("role.0.role_id").HasValue(appBuilderRoleId),
			),
		},
	})
}

func (IotCentralUserResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.UserID(state.ID)
	if err != nil {
		return nil, err
	}

	appId, err := apps.ParseIotAppID(state.Attributes["iotcentral_application_id"])
	if err != nil {
		return nil, err
	}

	app, err := clients.IoTCentral.AppsClient.Get(ctx, *appId)
	if err != nil || app.Model == nil {
		return nil, fmt.Errorf("checking for the presence of existing %q: %+v", appId, err)
	}

	userClient, err := clients.IoTCentral.UsersClient(ctx, *app.Model.Properties.Subdomain)
	if err != nil {
		return nil, fmt.Errorf("creating user client: %+v", err)
	}

	resp, err := userClient.Get(ctx, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	_, _, isValid := iotcentral.TryValidateUserExistence(resp.Value)

	return &isValid, nil
}

func (r IotCentralUserResource) basic_type_servicePrincipal(data acceptance.TestData, roleId string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {}

data "azurerm_client_config" "current" {}

resource "azuread_application" "test" {
  display_name = "acctest-iotcentralsp-%d"
}

resource "azuread_service_principal" "test" {
  application_id = azuread_application.test.application_id
}

	%s

resource "azurerm_iotcentral_user" "test" {
  iotcentral_application_id = azurerm_iotcentral_application.test.id
  user_id                   = "test-user-id"
  object_id                 = azuread_service_principal.test.object_id
  tenant_id                 = data.azurerm_client_config.current.tenant_id

  type = "ServicePrincipal"

  role {
    role_id = "%s"
  }
}
	`, data.RandomInteger, r.templateBasic(data), roleId)
}

func (r IotCentralUserResource) basic_type_group(data acceptance.TestData, roleId string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {}

data "azurerm_client_config" "current" {}

resource "azuread_group" "test" {
  display_name     = "acctest-iotcentraladgroup-%d"
  security_enabled = true
}

	%s

resource "azurerm_iotcentral_user" "test" {
  iotcentral_application_id = azurerm_iotcentral_application.test.id
  user_id                   = "test-user-id"
  object_id                 = azuread_group.test.object_id
  tenant_id                 = data.azurerm_client_config.current.tenant_id

  type = "Group"

  role {
    role_id = "%s"
  }
}
	`, data.RandomInteger, r.templateBasic(data), roleId)
}

func (r IotCentralUserResource) basic_type_email(data acceptance.TestData, email string, roleId string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

	%s

resource "azurerm_iotcentral_user" "test" {
  iotcentral_application_id = azurerm_iotcentral_application.test.id
  user_id                   = "test-user-id"
  email                     = "%s"

  type = "Email"

  role {
    role_id = "%s"
  }
}
	`, r.templateBasic(data), email, roleId)
}

func (r IotCentralUserResource) complete_type_servicePrincipal(data acceptance.TestData, roleId string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {}

data "azurerm_client_config" "current" {}

resource "azuread_application" "test" {
  display_name = "acctest-iotcentralsp-%d"
}

resource "azuread_service_principal" "test" {
  application_id = azuread_application.test.application_id
}

%s

resource "azurerm_iotcentral_user" "test" {
  iotcentral_application_id = azurerm_iotcentral_application.test.id
  user_id                   = "test-user-id"
  object_id                 = azuread_service_principal.test.object_id
  tenant_id                 = data.azurerm_client_config.current.tenant_id

  type = "ServicePrincipal"

  role {
    role_id         = "%s"
    organization_id = azurerm_iotcentral_organization.test.organization_id
  }
}
`, data.RandomInteger, r.templateComplete(data), roleId)
}

func (r IotCentralUserResource) complete_type_group(data acceptance.TestData, roleId string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {}

data "azurerm_client_config" "current" {}

resource "azuread_group" "test" {
  display_name     = "acctest-iotcentraladgroup-%d"
  security_enabled = true
}

%s

resource "azurerm_iotcentral_user" "test" {
  iotcentral_application_id = azurerm_iotcentral_application.test.id
  user_id                   = "test-user-id"
  object_id                 = azuread_group.test.object_id
  tenant_id                 = data.azurerm_client_config.current.tenant_id

  type = "Group"

  role {
    role_id         = "%s"
    organization_id = azurerm_iotcentral_organization.test.organization_id
  }
}
`, data.RandomInteger, r.templateComplete(data), roleId)
}

func (r IotCentralUserResource) complete_type_email(data acceptance.TestData, email string, roleId string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_iotcentral_user" "test" {
  iotcentral_application_id = azurerm_iotcentral_application.test.id
  user_id                   = "test-user-id"
  email                     = "%s"

  type = "Email"

  role {
    role_id         = "%s"
    organization_id = azurerm_iotcentral_organization.test.organization_id
  }
}
`, r.templateComplete(data), email, roleId)
}

func (IotCentralUserResource) templateBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_iotcentral_application" "test" {
  name                = "acctest-iotcentralapp-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sub_domain          = "subdomain-%[1]d"
  display_name        = "some-display-name"
  sku                 = "ST0"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (IotCentralUserResource) templateComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_iotcentral_application" "test" {
  name                = "acctest-iotcentralapp-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sub_domain          = "subdomain-%[1]d"
  display_name        = "some-display-name"
  sku                 = "ST0"
}

resource "azurerm_iotcentral_organization" "test" {
  iotcentral_application_id = azurerm_iotcentral_application.test.id
  organization_id           = "org-test-id"
  display_name              = "Org"
}
`, data.RandomInteger, data.Locations.Primary)
}
