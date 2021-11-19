package mollie

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-querystring/query"
)

// Chargeback describes a forced transaction reversal initiated by the cardholder's bank
type Chargeback struct {
	Resource         string          `json:"resource,omitempty"`
	ID               string          `json:"id,omitempty"`
	Amount           *Amount         `json:"amount,omitempty"`
	SettlementAmount *Amount         `json:"settlementAmount,omitempty"`
	CreatedAt        *time.Time      `json:"createdAt,omitempty"`
	ReversedAt       *time.Time      `json:"reversedAt,omitempty"`
	PaymentID        string          `json:"paymentId,omitempty"`
	Links            ChargebackLinks `json:"_links,omitempty"`
}

// ChargebackLinks describes all the possible links to be returned with
// a chargeback object.
type ChargebackLinks struct {
	Self          *URL `json:"self,omitempty"`
	Payment       *URL `json:"payment,omitempty"`
	Settlement    *URL `json:"settlement,omitempty"`
	Documentation *URL `json:"documentation,omitempty"`
}

// ChargebackOptions describes chargeback endpoint valid query string parameters.
type ChargebackOptions struct {
	Include string `url:"include,omitempty"`
	Embed   string `url:"embed,omitempty"`
}

// ListChargebackOptions describes list chargebacks endpoint valid query string parameters.
type ListChargebackOptions struct {
	Include   string `url:"include,omitempty"`
	Embed     string `url:"embed,omitempty"`
	ProfileID string `url:"profileId,omitempty"`
}

// ChargebackList describes how a list of chargebacks will be retrieved by Mollie.
type ChargebackList struct {
	Count    int `json:"count,omitempty"`
	Embedded struct {
		Chargebacks []Chargeback
	} `json:"_embedded,omitempty"`
	Links PaginationLinks `json:"_links,omitempty"`
}

// ChargebacksService instance operates over chargeback resources
type ChargebacksService service

// Get retrieves a single chargeback by its ID.
// Note the original payment’s ID is needed as well.
//
// See: https://docs.mollie.com/reference/v2/chargebacks-api/get-chargeback
func (cs *ChargebacksService) Get(ctx context.Context, paymentID, chargebackID string, options *ChargebackOptions) (p Chargeback, err error) {
	u := fmt.Sprintf("v2/payments/%s/chargebacks/%s", paymentID, chargebackID)
	if options != nil {
		v, _ := query.Values(options)
		u = fmt.Sprintf("%s?%s", u, v.Encode())
	}
	req, err := cs.client.NewAPIRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return
	}
	res, err := cs.client.Do(req)
	if err != nil {
		return
	}
	if err = json.Unmarshal(res.content, &p); err != nil {
		return
	}
	return
}

// List retrieves a list of chargebacks associated with your account/organization.
//
// See: https://docs.mollie.com/reference/v2/chargebacks-api/list-chargebacks
func (cs *ChargebacksService) List(ctx context.Context, options *ListChargebackOptions) (cl *ChargebackList, err error) {
	u := "v2/chargebacks"
	if options != nil {
		v, _ := query.Values(options)
		u = fmt.Sprintf("%s?%s", u, v.Encode())
	}
	return cs.list(ctx, u)
}

// ListForPayment retrieves a list of chargebacks associated with a single payment.
//
// See: https://docs.mollie.com/reference/v2/chargebacks-api/list-chargebacks
func (cs *ChargebacksService) ListForPayment(ctx context.Context, paymentID string, options *ListChargebackOptions) (cl *ChargebackList, err error) {
	u := fmt.Sprintf("v2/payments/%s/chargebacks", paymentID)
	if options != nil {
		v, _ := query.Values(options)
		u = fmt.Sprintf("%s?%s", u, v.Encode())
	}
	return cs.list(ctx, u)
}

// encapsulates the shared list methods logic
func (cs *ChargebacksService) list(ctx context.Context, uri string) (cl *ChargebackList, err error) {
	req, err := cs.client.NewAPIRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return
	}
	res, err := cs.client.Do(req)
	if err != nil {
		return
	}
	if err = json.Unmarshal(res.content, &cl); err != nil {
		return
	}
	return
}
