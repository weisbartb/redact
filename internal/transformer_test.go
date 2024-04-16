package internal

import "testing"

func TestArg_Float(t *testing.T) {
	type fields struct {
		OpCode opCode
		Value  any
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name: "String",
			fields: fields{
				OpCode: opCodeString,
				Value:  "6.43",
			},
			want: 6.43,
		},
		{
			name: "Int",
			fields: fields{
				OpCode: opCodeInt,
				Value:  6,
			},
			want: 6,
		},
		{
			name: "Bool",
			fields: fields{
				OpCode: opCodeBool,
				Value:  true,
			},
			want: 1,
		},
		{
			name: "Float",
			fields: fields{
				OpCode: opCodeFloat,
				Value:  float64(6.44),
			},
			want: 6.44,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Arg{
				OpCode: tt.fields.OpCode,
				Value:  tt.fields.Value,
			}
			if got := a.Float(); got != tt.want {
				t.Errorf("Float() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArg_Int(t *testing.T) {
	type fields struct {
		OpCode opCode
		Value  any
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "String",
			fields: fields{
				OpCode: opCodeString,
				Value:  "6",
			},
			want: 6,
		},
		{
			name: "Int",
			fields: fields{
				OpCode: opCodeInt,
				Value:  6,
			},
			want: 6,
		},
		{
			name: "Bool",
			fields: fields{
				OpCode: opCodeBool,
				Value:  true,
			},
			want: 1,
		},
		{
			name: "Float",
			fields: fields{
				OpCode: opCodeFloat,
				Value:  float64(6.44),
			},
			want: 6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Arg{
				OpCode: tt.fields.OpCode,
				Value:  tt.fields.Value,
			}
			if got := a.Int(); got != tt.want {
				t.Errorf("Int() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArg_String(t *testing.T) {
	type fields struct {
		OpCode opCode
		Value  any
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "String",
			fields: fields{
				OpCode: opCodeString,
				Value:  "6",
			},
			want: "6",
		},
		{
			name: "Int",
			fields: fields{
				OpCode: opCodeInt,
				Value:  6,
			},
			want: "6",
		},
		{
			name: "Bool",
			fields: fields{
				OpCode: opCodeBool,
				Value:  true,
			},
			want: "true",
		},
		{
			name: "Float",
			fields: fields{
				OpCode: opCodeFloat,
				Value:  float64(6.44),
			},
			want: "6.44",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Arg{
				OpCode: tt.fields.OpCode,
				Value:  tt.fields.Value,
			}
			if got := a.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArg_Bool(t *testing.T) {
	type fields struct {
		OpCode opCode
		Value  any
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "String",
			fields: fields{
				OpCode: opCodeString,
				Value:  "1",
			},
			want: true,
		},
		{
			name: "Int",
			fields: fields{
				OpCode: opCodeInt,
				Value:  1,
			},
			want: true,
		},
		{
			name: "Int",
			fields: fields{
				OpCode: opCodeInt,
				Value:  0,
			},
			want: false,
		},
		{
			name: "Bool",
			fields: fields{
				OpCode: opCodeBool,
				Value:  true,
			},
			want: true,
		},
		{
			name: "Float",
			fields: fields{
				OpCode: opCodeFloat,
				Value:  float64(1),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Arg{
				OpCode: tt.fields.OpCode,
				Value:  tt.fields.Value,
			}
			if got := a.Bool(); got != tt.want {
				t.Errorf("Bool() = %v, want %v", got, tt.want)
			}
		})
	}
}
