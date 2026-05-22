package formula

import "testing"

type seqRand struct {
	vals []int
	i    int
}

func (r *seqRand) Intn(n int) int {
	if n <= 0 {
		return 0
	}
	v := r.vals[r.i%len(r.vals)]
	r.i++
	return v % n
}

func TestEvalBooleanConvention(t *testing.T) {
	f, err := Parse("1 == 2")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Eval(f.Expr, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if v != 0 {
		t.Fatalf("expected 0, got %v", v)
	}
	f, err = Parse("3 > 1")
	if err != nil {
		t.Fatal(err)
	}
	v, err = Eval(f.Expr, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if v != 1 {
		t.Fatalf("expected 1, got %v", v)
	}
}

func TestEvalIf(t *testing.T) {
	f, err := Parse("if 1 then 10 else 5")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Eval(f.Expr, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if v != 10 {
		t.Fatalf("expected 10, got %v", v)
	}
}

func TestEvalVariables(t *testing.T) {
	f, err := Parse("{strength} + 2")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Eval(f.Expr, map[string]float64{"strength": 3}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if v != 5 {
		t.Fatalf("expected 5, got %v", v)
	}
}

func TestEvalDiceRequiresRand(t *testing.T) {
	f, err := Parse("2d6")
	if err != nil {
		t.Fatal(err)
	}
	_, err = Eval(f.Expr, nil, nil)
	if err == nil {
		t.Fatal("expected error for nil RandSource with dice")
	}
}

func TestEvalDiceDeterministic(t *testing.T) {
	f, err := Parse("2d6")
	if err != nil {
		t.Fatal(err)
	}
	rng := &seqRand{vals: []int{5, 1}}
	v, err := Eval(f.Expr, nil, rng)
	if err != nil {
		t.Fatal(err)
	}
	if v != 8 {
		t.Fatalf("expected 8, got %v", v)
	}
}

func TestEvalClampImath(t *testing.T) {
	f, err := Parse("clamp(5, 0, 10)")
	if err != nil {
		t.Fatal(err)
	}
	v, err := Eval(f.Expr, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if v != 5 {
		t.Fatalf("expected 5, got %v", v)
	}
}
