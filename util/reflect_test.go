package util

import (
	"reflect"
	. "reflect"
	"testing"
)

func TestGetSliceElemKind(t *testing.T) {
	type args struct {
		st Type
	}
	tests := []struct {
		name string
		args args
		want Kind
	}{
		{
			name: "[]int",
			args: args{reflect.TypeOf([]int{})},
			want: reflect.Int,
		},
		{
			name: "[]string",
			args: args{reflect.TypeOf([]string{})},
			want: reflect.String,
		},
		{
			name: "[]int64",
			args: args{reflect.TypeOf([]int64{})},
			want: reflect.Int64,
		},
		{
			name: "[]struct",
			args: args{reflect.TypeOf([]args{})},
			want: reflect.Struct,
		},
		{
			name: "[]ptr",
			args: args{reflect.TypeOf([]*args{})},
			want: reflect.Pointer,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSliceElemKind(tt.args.st); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSliceElemKind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSliceElemType(t *testing.T) {
	type args struct {
		st Type
	}
	tests := []struct {
		name string
		args args
		want Type
	}{
		{
			name: "[]int",
			args: args{reflect.TypeOf([]int{})},
			want: reflect.TypeOf(1),
		},
		{
			name: "[]string",
			args: args{reflect.TypeOf([]string{})},
			want: reflect.TypeOf("1"),
		},
		{
			name: "[]struct",
			args: args{reflect.TypeOf([]args{})},
			want: reflect.TypeOf(args{}),
		},
		{
			name: "[]ptr",
			args: args{reflect.TypeOf([]*args{})},
			want: reflect.TypeOf(&args{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSliceElemType(tt.args.st); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSliceElemType() = %v, want %v", got, tt.want)
			}
		})
	}
}
