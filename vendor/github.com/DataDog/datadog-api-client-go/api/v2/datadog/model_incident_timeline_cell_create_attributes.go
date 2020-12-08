/*
 * Unless explicitly stated otherwise all files in this repository are licensed under the Apache-2.0 License.
 * This product includes software developed at Datadog (https://www.datadoghq.com/).
 * Copyright 2019-Present Datadog, Inc.
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package datadog

import (
	"encoding/json"
	"fmt"
)

// IncidentTimelineCellCreateAttributes - The timeline cell's attributes for a create request.
type IncidentTimelineCellCreateAttributes struct {
	IncidentTimelineCellMarkdownCreateAttributes *IncidentTimelineCellMarkdownCreateAttributes
}

// IncidentTimelineCellMarkdownCreateAttributesAsIncidentTimelineCellCreateAttributes is a convenience function that returns IncidentTimelineCellMarkdownCreateAttributes wrapped in IncidentTimelineCellCreateAttributes
func IncidentTimelineCellMarkdownCreateAttributesAsIncidentTimelineCellCreateAttributes(v *IncidentTimelineCellMarkdownCreateAttributes) IncidentTimelineCellCreateAttributes {
	return IncidentTimelineCellCreateAttributes{IncidentTimelineCellMarkdownCreateAttributes: v}
}

// Unmarshal JSON data into one of the pointers in the struct
func (dst *IncidentTimelineCellCreateAttributes) UnmarshalJSON(data []byte) error {
	var err error
	match := 0
	// try to unmarshal data into IncidentTimelineCellMarkdownCreateAttributes
	err = json.Unmarshal(data, &dst.IncidentTimelineCellMarkdownCreateAttributes)
	if err == nil {
		jsonIncidentTimelineCellMarkdownCreateAttributes, _ := json.Marshal(dst.IncidentTimelineCellMarkdownCreateAttributes)
		if string(jsonIncidentTimelineCellMarkdownCreateAttributes) == "{}" { // empty struct
			dst.IncidentTimelineCellMarkdownCreateAttributes = nil
		} else {
			match++
		}
	} else {
		dst.IncidentTimelineCellMarkdownCreateAttributes = nil
	}

	if match > 1 { // more than 1 match
		// reset to nil
		dst.IncidentTimelineCellMarkdownCreateAttributes = nil

		return fmt.Errorf("Data matches more than one schema in oneOf(IncidentTimelineCellCreateAttributes)")
	} else if match == 1 {
		return nil // exactly one match
	} else { // no match
		return fmt.Errorf("Data failed to match schemas in oneOf(IncidentTimelineCellCreateAttributes)")
	}
}

// Marshal data from the first non-nil pointers in the struct to JSON
func (src IncidentTimelineCellCreateAttributes) MarshalJSON() ([]byte, error) {
	if src.IncidentTimelineCellMarkdownCreateAttributes != nil {
		return json.Marshal(&src.IncidentTimelineCellMarkdownCreateAttributes)
	}

	return nil, nil // no data in oneOf schemas
}

// Get the actual instance
func (obj *IncidentTimelineCellCreateAttributes) GetActualInstance() interface{} {
	if obj.IncidentTimelineCellMarkdownCreateAttributes != nil {
		return obj.IncidentTimelineCellMarkdownCreateAttributes
	}

	// all schemas are nil
	return nil
}

type NullableIncidentTimelineCellCreateAttributes struct {
	value *IncidentTimelineCellCreateAttributes
	isSet bool
}

func (v NullableIncidentTimelineCellCreateAttributes) Get() *IncidentTimelineCellCreateAttributes {
	return v.value
}

func (v *NullableIncidentTimelineCellCreateAttributes) Set(val *IncidentTimelineCellCreateAttributes) {
	v.value = val
	v.isSet = true
}

func (v NullableIncidentTimelineCellCreateAttributes) IsSet() bool {
	return v.isSet
}

func (v *NullableIncidentTimelineCellCreateAttributes) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableIncidentTimelineCellCreateAttributes(val *IncidentTimelineCellCreateAttributes) *NullableIncidentTimelineCellCreateAttributes {
	return &NullableIncidentTimelineCellCreateAttributes{value: val, isSet: true}
}

func (v NullableIncidentTimelineCellCreateAttributes) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableIncidentTimelineCellCreateAttributes) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}