package main

import (
	"encoding/json"
	"fmt"
)

type SampleBody interface {
	GetValue() string
}

type SampleType string

const (
	SampleType1 = "one"
	SampleType2 = "two"
	SampleType3 = "three"
)

type Sample struct {
	Type SampleType `json:"type"`
	Body SampleBody `json:"body"`
}

func (s *Sample) UnmarshalJSON(d []byte) error {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(d, &m); err != nil {
		return err
	}

	rawBody, ok := m["body"]
	if !ok {
		return fmt.Errorf("no body: %s", string(d))
	}
	delete(m, "body")

	// Convert m (a map of raw JSON fields) into a Sample instance.
	// We're using json.Marshal + Unmarshal here for simplicity.
	// This avoids manual field-by-field decoding or using reflection.
	//
	// Note: This approach is less efficient than decoding each field directly,
	// but it's acceptable for example or non-performance-critical use cases.
	//
	// The block scope ({}) is used here to group related code together in this example.
	{
		dataSample, err := json.Marshal(m)
		if err != nil {
			return err
		}

		// We define a type alias (sampleAlias) without the UnmarshalJSON method
		// to prevent infinite recursion. If we unmarshal directly into Sample,
		// it would call Sample.UnmarshalJSON again, causing a stack overflow.
		// Using sampleAlias allows safe decoding of base fields only.
		type sampleAlias Sample
		var newS sampleAlias
		if err := json.Unmarshal(dataSample, &newS); err != nil {
			return err
		}
		*s = Sample(newS)
	}

	unmarshalBody := func() (SampleBody, error) {
		switch s.Type {
		case SampleType1:
			var dest SampleBody1
			return &dest, json.Unmarshal(rawBody, &dest)
		case SampleType2:
			var dest SampleBody2
			return &dest, json.Unmarshal(rawBody, &dest)
		case SampleType3:
			var dest SampleBody3
			return &dest, json.Unmarshal(rawBody, &dest)
		default:
			return nil, fmt.Errorf("invalid type %s", s.Type)
		}
	}

	sampleBody, err := unmarshalBody()
	if err != nil {
		return err
	}
	s.Body = sampleBody

	return nil
}

type SampleBody1 struct {
	Value string `json:"value"`
}

func (s *SampleBody1) GetValue() string {
	return s.Value
}

type SampleBody2 struct {
	Value string `json:"value"`
}

func (s SampleBody2) GetValue() string {
	return s.Value
}

type SampleBody3 struct {
	Value string `json:"value"`
}

func (s *SampleBody3) GetValue() string {
	return s.Value
}

type SampleGroup struct {
	ID      string            `json:"id"`
	Samples map[string]Sample `json:"samples"`
}

func main() {
	input := `
	{
		"id": "123",
		"samples": {
			"1": {
				"type": "one",
				"body": {
					"value": "this is a value 1"
				}
			},
			"2": {
				"type": "two",
				"body": {
					"value": "this is a value 2"
				}
			},
			"3": {
				"type": "three",
				"body": {
					"value": "this is a value 3"
				}
			}
		}
	}
	`

	var group SampleGroup
	if err := json.Unmarshal([]byte(input), &group); err != nil {
		panic(err)
	}

	fmt.Println(group.ID)

	for _, s := range group.Samples {
		fmt.Println("s.Body.GetValue():", s.Body.GetValue())
	}
}
