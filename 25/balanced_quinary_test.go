package main

import (
	"fmt"
	"testing"
)

// test conversion to and from the following cases:
var cases = []struct {
	decimal int
	bq      string
}{
	{1, "1"},
	{2, "2"},
	{3, "1="},
	{4, "1-"},
	{5, "10"},
	{6, "11"},
	{7, "12"},
	{8, "2="},
	{9, "2-"},
	{10, "20"},
	{15, "1=0"},
	{20, "1-0"},
	{2022, "1=11-2"},
	{12345, "1-0---0"},
	{314159265, "1121-1110-1=0"},
}

func TestBalancedQuinaryToInt(t *testing.T) {
	for _, c := range cases {
		t.Run(fmt.Sprintf("%s to int", c.bq), func(t *testing.T) {
			bq := BalancedQuinary(c.bq)
			if bq.Int() != c.decimal {
				t.Errorf("Expected %d, got %d", c.decimal, bq.Int())
			}
		})
	}
}

func TestBalancedQuinaryFromInt(t *testing.T) {
	for _, c := range cases {
		t.Run(fmt.Sprintf("%d to BalancedQuinary", c.decimal), func(t *testing.T) {
			bq := BalancedQuinaryFromInt(c.decimal)
			if bq.String() != c.bq {
				t.Errorf("Expected %s, got %s", c.bq, bq.String())
			}
		})
	}
}
