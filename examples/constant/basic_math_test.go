// Copyright 2025 Robert Snedegar
//
// Licensed under the Apache License, Version 2.0 (the License);
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an AS IS BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package constant

import (
	"testing"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"positive numbers", 5, 3, 8},
		{"negative numbers", -5, -3, -8},
		{"mixed signs", 5, -3, 2},
		{"zero with positive", 0, 5, 5},
		{"zero with negative", 0, -5, -5},
		{"both zero", 0, 0, 0},
		{"large numbers", 1000000, 2000000, 3000000},
		{"commutative property", 3, 5, 8}, // same as first test but reversed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Add(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Add(%d, %d) = %d, want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"positive result", 5, 3, 2},
		{"negative result", 3, 5, -2},
		{"zero result", 5, 5, 0},
		{"negative numbers", -5, -3, -2},
		{"mixed signs", 5, -3, 8},
		{"zero with positive", 0, 5, -5},
		{"positive with zero", 5, 0, 5},
		{"both zero", 0, 0, 0},
		{"large numbers", 2000000, 1000000, 1000000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Subtract(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Subtract(%d, %d) = %d, want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"positive numbers", 5, 3, 15},
		{"negative numbers", -5, -3, 15},
		{"mixed signs", 5, -3, -15},
		{"zero with positive", 0, 5, 0},
		{"positive with zero", 5, 0, 0},
		{"both zero", 0, 0, 0},
		{"identity element", 5, 1, 5},
		{"commutative property", 3, 5, 15}, // same as first test but reversed
		{"large numbers", 1000, 2000, 2000000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Multiply(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Multiply(%d, %d) = %d, want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"positive numbers", 15, 3, 5},
		{"negative numbers", -15, -3, 5},
		{"mixed signs", 15, -3, -5},
		{"zero dividend", 0, 5, 0},
		{"identity element", 5, 1, 5},
		{"large numbers", 2000000, 1000, 2000},
		{"division with remainder", 17, 5, 3}, // integer division truncates
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Divide(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Divide(%d, %d) = %d, want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestDivideByZero(t *testing.T) {
	// Test that division by zero panics
	defer func() {
		if r := recover(); r == nil {
			t.Error("Divide(5, 0) should panic")
		}
	}()

	_ = Divide(5, 0)
}

// Benchmark functions for basic math operations

func BenchmarkAdd(b *testing.B) {
	for b.Loop() {
		_ = Add(42, 27)
	}
}

func BenchmarkSubtract(b *testing.B) {
	for b.Loop() {
		_ = Subtract(42, 27)
	}
}

func BenchmarkMultiply(b *testing.B) {
	for b.Loop() {
		_ = Multiply(42, 27)
	}
}

func BenchmarkDivide(b *testing.B) {
	for b.Loop() {
		_ = Divide(42, 27)
	}
}
