package capacityreservationgroups

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListBySubscriptionOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]CapacityReservationGroup
}

type ListBySubscriptionCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []CapacityReservationGroup
}

type ListBySubscriptionOperationOptions struct {
	Expand *ExpandTypesForGetCapacityReservationGroups
}

func DefaultListBySubscriptionOperationOptions() ListBySubscriptionOperationOptions {
	return ListBySubscriptionOperationOptions{}
}

func (o ListBySubscriptionOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListBySubscriptionOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ListBySubscriptionOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Expand != nil {
		out.Append("$expand", fmt.Sprintf("%v", *o.Expand))
	}
	return &out
}

// ListBySubscription ...
func (c CapacityReservationGroupsClient) ListBySubscription(ctx context.Context, id commonids.SubscriptionId, options ListBySubscriptionOperationOptions) (result ListBySubscriptionOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/providers/Microsoft.Compute/capacityReservationGroups", id.ID()),
		OptionsObject: options,
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.ExecutePaged(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	var values struct {
		Values *[]CapacityReservationGroup `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListBySubscriptionComplete retrieves all the results into a single object
func (c CapacityReservationGroupsClient) ListBySubscriptionComplete(ctx context.Context, id commonids.SubscriptionId, options ListBySubscriptionOperationOptions) (ListBySubscriptionCompleteResult, error) {
	return c.ListBySubscriptionCompleteMatchingPredicate(ctx, id, options, CapacityReservationGroupOperationPredicate{})
}

// ListBySubscriptionCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c CapacityReservationGroupsClient) ListBySubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, options ListBySubscriptionOperationOptions, predicate CapacityReservationGroupOperationPredicate) (result ListBySubscriptionCompleteResult, err error) {
	items := make([]CapacityReservationGroup, 0)

	resp, err := c.ListBySubscription(ctx, id, options)
	if err != nil {
		err = fmt.Errorf("loading results: %+v", err)
		return
	}
	if resp.Model != nil {
		for _, v := range *resp.Model {
			if predicate.Matches(v) {
				items = append(items, v)
			}
		}
	}

	result = ListBySubscriptionCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
