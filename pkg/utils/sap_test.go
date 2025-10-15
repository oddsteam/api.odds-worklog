package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddBlank(t *testing.T) {
	t.Run("should add blank 5", func(t *testing.T) {
		actual := AddBlank("A", 5)
		assert.Equal(t, "A    ", actual)
	})

	t.Run("should not blank when length is over length of value", func(t *testing.T) {
		actual := AddBlank("AAAA", 2)
		assert.Equal(t, "AAAA", actual)
	})

	t.Run("should add valid blank when value is Thai", func(t *testing.T) {
		actual := AddBlank("บจก. โซโล่ เลเวลลิ่ง", 130)
		assert.Equal(t, "บจก. โซโล่ เลเวลลิ่ง                                                                                                              ", actual)
	})
}

func TestAmountStr(t *testing.T) {
	type args struct {
		amount float64
		n      int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "pads with leading zeros when shorter than n",
			args: args{amount: 123.456, n: 7},
			want: "0123.46",
		},
		{
			name: "exact length no padding",
			args: args{amount: 123.456, n: 6},
			want: "123.46",
		},
		{
			name: "truncate when longer than n",
			args: args{amount: 123.456, n: 5},
			want: "123.4",
		},
		{
			name: "negative amount with padding",
			args: args{amount: -1.234, n: 7},
			want: "00-1.23",
		},
		{
			name: "rounding carryover",
			args: args{amount: 9.995, n: 7},
			want: "0010.00",
		},
		{
			name: "zero with truncation",
			args: args{amount: 0, n: 2},
			want: "0.00",
		},
		{
			name: "small amount padding",
			args: args{amount: 0.5, n: 6},
			want: "000.50",
		},
		{
			name: "large amount padding",
			args: args{amount: 54704.00, n: 15},
			want: "000000054704.00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, AmountStr(tt.args.amount, tt.args.n), "AmountStr(%v, %v)", tt.args.amount, tt.args.n)
		})
	}
}

func TestLeftN(t *testing.T) {
	type args struct {
		s string
		n int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "show first 5 digit", args: args{s: "HelloWorld", n: 5}, want: "Hello"},
		{name: "n มากกว่าความยาว string → return ทั้งหมด", args: args{s: "Hello", n: 10}, want: "Hello"},
		{name: "n = 0 → ควรได้ string ว่าง", args: args{s: "Hello", n: 0}, want: ""},
		{name: "input ว่าง → return ว่าง ", args: args{s: "", n: 10}, want: ""},
		{name: "Unicode (multi-byte) กรณีนี้จะ fail เพราะ len() นับ byte", args: args{s: "こんにちは世界", n: 3}, want: "こんに"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, LeftN(tt.args.s, tt.args.n), "LeftN(%v, %v)", tt.args.s, tt.args.n)
		})
	}
}

func TestReceiveAcCode(t *testing.T) {
	type args struct {
		bnkCode string
		acCode  string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "valid bnkCode 011 with 10 digits",
			args:    args{bnkCode: "011", acCode: "1234567890"},
			want:    "01234567890",
			wantErr: assert.NoError,
		},
		{
			name:    "valid bnkCode 011 with 9 digits",
			args:    args{bnkCode: "011", acCode: "123456789"},
			want:    "123456789",
			wantErr: assert.NoError,
		},
		{
			name:    "valid bnkCode 011 with spaces and dashes",
			args:    args{bnkCode: "011", acCode: " 123-456-7890 "},
			want:    "01234567890",
			wantErr: assert.NoError,
		},
		{
			name:    "invalid bnkCode",
			args:    args{bnkCode: "002", acCode: "1234567890"},
			want:    "",
			wantErr: assert.Error,
		},
		{
			name:    "empty acCode",
			args:    args{bnkCode: "011", acCode: ""},
			want:    "",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReceiveAcCode(tt.args.bnkCode, tt.args.acCode)
			if !tt.wantErr(t, err, fmt.Sprintf("ReceiveAcCode(%v, %v)", tt.args.bnkCode, tt.args.acCode)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ReceiveAcCode(%v, %v)", tt.args.bnkCode, tt.args.acCode)
		})
	}
}

func TestReceiveBRCode(t *testing.T) {
	type args struct {
		bnkCode  string
		brchCode string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "case 002, brchCode >= 3",
			args: args{bnkCode: "002", brchCode: "12345"},
			want: "0123",
		},
		{
			name: "case 002, brchCode < 3",
			args: args{bnkCode: "002", brchCode: "12"},
			want: "012",
		},
		{
			name: "case 030, brchCode >= 4",
			args: args{bnkCode: "030", brchCode: "12345"},
			want: "1234",
		},
		{
			name: "case 030, brchCode < 4",
			args: args{bnkCode: "030", brchCode: "12"},
			want: "12",
		},
		{
			name: "case 034 always 0000",
			args: args{bnkCode: "034", brchCode: "999"},
			want: "0000",
		},
		{
			name: "case 045 always 0010",
			args: args{bnkCode: "045", brchCode: "999"},
			want: "0010",
		},
		{
			name: "default case",
			args: args{bnkCode: "999", brchCode: "999"},
			want: "0001",
		},
		{
			name: "empty brchCode",
			args: args{bnkCode: "002", brchCode: ""},
			want: "0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, ReceiveBRCode(tt.args.bnkCode, tt.args.brchCode), "ReceiveBRCode(%v, %v)", tt.args.bnkCode, tt.args.brchCode)
		})
	}
}
