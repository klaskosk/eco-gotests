// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// APIVipConnectivityRequest api vip connectivity request
//
// swagger:model api_vip_connectivity_request
type APIVipConnectivityRequest struct {

	// A CA certficate to be used when contacting the URL via https.
	CaCertificate *string `json:"ca_certificate,omitempty"`

	// A string which will be used as Authorization Bearer token to fetch the ignition from ignition_endpoint_url (DEPRECATED use request_headers to pass this token).
	IgnitionEndpointToken *string `json:"ignition_endpoint_token,omitempty"`

	// Additional request headers to include when fetching the ignition from ignition_endpoint_url.
	RequestHeaders []*APIVipConnectivityAdditionalRequestHeader `json:"request_headers,omitempty"`

	// URL address of the API.
	// Required: true
	URL *string `json:"url"`

	// Whether to verify if the API VIP belongs to one of the interfaces (DEPRECATED).
	VerifyCidr bool `json:"verify_cidr,omitempty"`
}

// Validate validates this api vip connectivity request
func (m *APIVipConnectivityRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateRequestHeaders(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateURL(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *APIVipConnectivityRequest) validateRequestHeaders(formats strfmt.Registry) error {
	if swag.IsZero(m.RequestHeaders) { // not required
		return nil
	}

	for i := 0; i < len(m.RequestHeaders); i++ {
		if swag.IsZero(m.RequestHeaders[i]) { // not required
			continue
		}

		if m.RequestHeaders[i] != nil {
			if err := m.RequestHeaders[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("request_headers" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("request_headers" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *APIVipConnectivityRequest) validateURL(formats strfmt.Registry) error {

	if err := validate.Required("url", "body", m.URL); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this api vip connectivity request based on the context it is used
func (m *APIVipConnectivityRequest) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateRequestHeaders(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *APIVipConnectivityRequest) contextValidateRequestHeaders(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.RequestHeaders); i++ {

		if m.RequestHeaders[i] != nil {
			if err := m.RequestHeaders[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("request_headers" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("request_headers" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *APIVipConnectivityRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *APIVipConnectivityRequest) UnmarshalBinary(b []byte) error {
	var res APIVipConnectivityRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}