package state

import (
	"0chain.net/core/datastore"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestNewMint(t *testing.T) {
	var (
		minter             = "minter"
		toClientID         = "to client id"
		amount     Balance = 5
	)

	type args struct {
		minter   datastore.Key
		Receiver datastore.Key
		amount   Balance
	}
	tests := []struct {
		name string
		args args
		want *Mint
	}{
		{
			name: "OK",
			args: args{
				minter:   minter,
				Receiver: toClientID,
				amount:   amount,
			},
			want: &Mint{
				Minter:   minter,
				Receiver: toClientID,
				Amount:   amount,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMint(tt.args.minter, tt.args.Receiver, tt.args.amount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMint_Encode(t *testing.T) {
	m := NewMint("minter", "client id", 5)
	blob, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		Minter   datastore.Key
		Receiver datastore.Key
		Amount   Balance
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name:   "OK",
			fields: fields(*m),
			want:   blob,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Mint{
				Minter:   tt.fields.Minter,
				Receiver: tt.fields.Receiver,
				Amount:   tt.fields.Amount,
			}
			if got := m.Encode(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMint_Decode(t *testing.T) {
	m := NewMint("minter", "client id", 5)
	blob, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		Minter     datastore.Key
		ToClientID datastore.Key
		Amount     Balance
	}
	type args struct {
		input []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Mint
		wantErr bool
	}{
		{
			name:    "OK",
			args:    args{input: blob},
			want:    m,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Mint{
				Minter:   tt.fields.Minter,
				Receiver: tt.fields.ToClientID,
				Amount:   tt.fields.Amount,
			}
			if err := m.Decode(tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, m)
		})
	}
}
