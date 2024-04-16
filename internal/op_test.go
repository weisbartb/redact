package internal

import "testing"

func Test_opCode_String(t *testing.T) {
	tests := []struct {
		name string
		o    opCode
		want string
	}{
		{
			name: "None",
			o:    opCodeNone,
			want: "nil",
		},
		{
			name: "Set",
			o:    opCodeSet,
			want: "set",
		},
		{
			name: "Params",
			o:    opCodeParams,
			want: "params",
		},
		{
			name: "String",
			o:    opCodeString,
			want: "string",
		},
		{
			name: "Float",
			o:    opCodeFloat,
			want: "float",
		},
		{
			name: "Int",
			o:    opCodeInt,
			want: "int",
		},
		{
			name: "Bool",
			o:    opCodeBool,
			want: "bool",
		},
		{
			name: "Run",
			o:    opCodeRun,
			want: "run",
		},
		{
			name: "Nil",
			o:    opCodeNil,
			want: "nil",
		},
		{
			name: "Chain",
			o:    opCodeChain,
			want: "chain",
		},
		{
			name: "Unknown",
			o:    -1,
			want: "unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
