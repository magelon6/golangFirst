package data

import "testing"

func TestCheckValidation(t *testing.T) {
	p := &Product{
		Name: "keks",
		Price: 1.12,
		SKU: "abc-abc-rfc",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}