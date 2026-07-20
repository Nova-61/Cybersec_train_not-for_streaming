package main

import "testing"

func TestAdd(t *testing.T) {
	result := Add(2, 3)
	expected := 5
	if result != expected {
		t.Errorf("Add(2,3) = %d; want %d", result, expected)
	}
}

func TestAddNegative(t *testing.T) {
	result := Add(-2, -3)
	expected := -5
	if result != expected {
		t.Errorf("Add(-2,-3) = %d; want %d", result, expected)
	}
}

func TestAddZero(t *testing.T) {
	result := Add(5, 0)
	expected := 5
	if result != expected {
		t.Errorf("Add(5,0) = %d; want %d", result, expected)
	}
}

func TestMultiply(t *testing.T) {
	result := Multiply(2, 3)
	expected := 6
	if result != expected {
		t.Errorf("Multiply(2,3) = %d; want %d", result, expected)
	}
}

func TestMultiplyNegative(t *testing.T) {
	result := Multiply(-2, 3)
	expected := -6
	if result != expected {
		t.Errorf("Multiply(-2,3) = %d; want %d", result, expected)
	}
}

func TestMultiplyZero(t *testing.T) {
	result := Multiply(5, 0)
	expected := 0
	if result != expected {
		t.Errorf("Multiply(5,0) = %d; want %d", result, expected)
	}
}

func TestDivide(t *testing.T) {
	result, err := Divide(6, 3)
	expected := 2
	if err != nil {
		t.Errorf("Divide(6,3) error: %v", err)
	}
	if result != expected {
		t.Errorf("Divide(6,3) = %d; want %d", result, expected)
	}
}

func TestDivideNegative(t *testing.T) {
	result, err := Divide(-6, 3)
	expected := -2
	if err != nil {
		t.Errorf("Divide(-6,3) error: %v", err)
	}
	if result != expected {
		t.Errorf("Divide(-6,3) = %d; want %d", result, expected)
	}
}

func TestDivideByZero(t *testing.T) {
	result, err := Divide(10, 0)
	if err == nil {
		t.Error("Divide(10,0) expected error, got nil")
	}
	if result != 0 {
		t.Errorf("Divide(10,0) = %d; want 0", result)
	}
}

func TestDivideByOne(t *testing.T) {
	result, err := Divide(7, 1)
	expected := 7
	if err != nil {
		t.Errorf("Divide(7,1) error: %v", err)
	}
	if result != expected {
		t.Errorf("Divide(7,1) = %d; want %d", result, expected)
	}
}
