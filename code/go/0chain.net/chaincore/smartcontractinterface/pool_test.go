package smartcontractinterface_test

import (
	"0chain.net/chaincore/smartcontractinterface"
	"0chain.net/chaincore/state"
	"0chain.net/chaincore/tokenpool"
	"0chain.net/smartcontract/minersc"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func makeTestPoolStats() *smartcontractinterface.PoolStats {
	return &smartcontractinterface.PoolStats{
		DelegateID:   "id",
		High:         2,
		Low:          1,
		InterestPaid: 1,
		RewardPaid:   1,
		NumRounds:    2,
		Status:       "status",
	}
}

func TestPoolStats_AddInterests(t *testing.T) {
	type fields struct {
		DelegateID   string
		High         state.Balance
		Low          state.Balance
		InterestPaid state.Balance
		RewardPaid   state.Balance
		NumRounds    int64
		Status       string
	}
	type args struct {
		value state.Balance
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *smartcontractinterface.PoolStats
	}{
		{
			name: "OK",
			fields: fields{
				InterestPaid: 1,
				High:         2,
				Low:          4,
			},
			args: args{value: 3},
			want: &smartcontractinterface.PoolStats{
				InterestPaid: 4,
				High:         3,
				Low:          3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &smartcontractinterface.PoolStats{
				DelegateID:   tt.fields.DelegateID,
				High:         tt.fields.High,
				Low:          tt.fields.Low,
				InterestPaid: tt.fields.InterestPaid,
				RewardPaid:   tt.fields.RewardPaid,
				NumRounds:    tt.fields.NumRounds,
				Status:       tt.fields.Status,
			}

			ps.AddInterests(tt.args.value)
			assert.Equal(t, tt.want, ps)
		})
	}
}

func TestPoolStats_AddRewards(t *testing.T) {
	type fields struct {
		DelegateID   string
		High         state.Balance
		Low          state.Balance
		InterestPaid state.Balance
		RewardPaid   state.Balance
		NumRounds    int64
		Status       string
	}
	type args struct {
		value state.Balance
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *smartcontractinterface.PoolStats
	}{
		{
			name: "OK",
			fields: fields{
				RewardPaid: 1,
				High:       2,
				Low:        4,
			},
			args: args{value: 3},
			want: &smartcontractinterface.PoolStats{
				RewardPaid: 4,
				High:       3,
				Low:        3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &smartcontractinterface.PoolStats{
				DelegateID:   tt.fields.DelegateID,
				High:         tt.fields.High,
				Low:          tt.fields.Low,
				InterestPaid: tt.fields.InterestPaid,
				RewardPaid:   tt.fields.RewardPaid,
				NumRounds:    tt.fields.NumRounds,
				Status:       tt.fields.Status,
			}

			ps.AddRewards(tt.args.value)
			assert.Equal(t, tt.want, ps)
		})
	}
}

func TestPoolStats_Encode(t *testing.T) {
	ps := makeTestPoolStats()
	blob, err := json.Marshal(ps)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		DelegateID   string
		High         state.Balance
		Low          state.Balance
		InterestPaid state.Balance
		RewardPaid   state.Balance
		NumRounds    int64
		Status       string
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name:   "OK",
			fields: fields(*ps),
			want:   blob,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &smartcontractinterface.PoolStats{
				DelegateID:   tt.fields.DelegateID,
				High:         tt.fields.High,
				Low:          tt.fields.Low,
				InterestPaid: tt.fields.InterestPaid,
				RewardPaid:   tt.fields.RewardPaid,
				NumRounds:    tt.fields.NumRounds,
				Status:       tt.fields.Status,
			}
			if got := ps.Encode(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPoolStats_Decode(t *testing.T) {
	ps := makeTestPoolStats()
	blob, err := json.Marshal(ps)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		DelegateID   string
		High         state.Balance
		Low          state.Balance
		InterestPaid state.Balance
		RewardPaid   state.Balance
		NumRounds    int64
		Status       string
	}
	type args struct {
		input []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    *smartcontractinterface.PoolStats
	}{
		{
			name:    "OK",
			args:    args{input: blob},
			wantErr: false,
			want:    ps,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &smartcontractinterface.PoolStats{
				DelegateID:   tt.fields.DelegateID,
				High:         tt.fields.High,
				Low:          tt.fields.Low,
				InterestPaid: tt.fields.InterestPaid,
				RewardPaid:   tt.fields.RewardPaid,
				NumRounds:    tt.fields.NumRounds,
				Status:       tt.fields.Status,
			}
			if err := ps.Decode(tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, ps)
		})
	}
}

func TestNewDelegatePool(t *testing.T) {
	tests := []struct {
		name string
		want *smartcontractinterface.DelegatePool
	}{
		{
			name: "OK",
			want: &smartcontractinterface.DelegatePool{
				ZcnLockingPool: &tokenpool.ZcnLockingPool{},
				PoolStats:      &smartcontractinterface.PoolStats{Low: -1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := smartcontractinterface.NewDelegatePool(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDelegatePool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func makeTestDelegatePool() *smartcontractinterface.DelegatePool {
	return &smartcontractinterface.DelegatePool{
		PoolStats: makeTestPoolStats(),
		ZcnLockingPool: &tokenpool.ZcnLockingPool{
			ZcnPool: tokenpool.ZcnPool{},
			TokenLockInterface: &minersc.ViewChangeLock{
				Owner: "owner",
			},
		},
	}
}

func TestDelegatePool_Encode(t *testing.T) {
	dp := makeTestDelegatePool()
	blob, err := json.Marshal(dp)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		PoolStats      *smartcontractinterface.PoolStats
		ZcnLockingPool *tokenpool.ZcnLockingPool
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "OK",
			fields: fields{
				PoolStats:      dp.PoolStats,
				ZcnLockingPool: dp.ZcnLockingPool,
			},
			want: blob,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dp := &smartcontractinterface.DelegatePool{
				PoolStats:      tt.fields.PoolStats,
				ZcnLockingPool: tt.fields.ZcnLockingPool,
			}
			if got := dp.Encode(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDelegatePool_Decode(t *testing.T) {
	dp := makeTestDelegatePool()
	blob, err := json.Marshal(dp)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		PoolStats      *smartcontractinterface.PoolStats
		ZcnLockingPool *tokenpool.ZcnLockingPool
	}
	type args struct {
		input     []byte
		tokenlock tokenpool.TokenLockInterface
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    *smartcontractinterface.DelegatePool
	}{
		{
			name:    "Unmarshal_Input_ERR",
			args:    args{input: []byte("}{")},
			wantErr: true,
		},
		{
			name:    "Unmarshal_Stats_ERR",
			args:    args{input: []byte("}{")},
			wantErr: true,
		},
		{
			name: "Decode_ERR",
			fields: fields{
				ZcnLockingPool: &tokenpool.ZcnLockingPool{
					TokenLockInterface: &minersc.ViewChangeLock{},
				},
			},
			args: args{
				input: blob,
			},
			wantErr: true,
		},
		{
			name: "OK",
			fields: fields{
				ZcnLockingPool: &tokenpool.ZcnLockingPool{
					TokenLockInterface: &minersc.ViewChangeLock{},
				},
			},
			args: args{
				input:     blob,
				tokenlock: &minersc.ViewChangeLock{},
			},
			want: dp,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dp := &smartcontractinterface.DelegatePool{
				PoolStats:      tt.fields.PoolStats,
				ZcnLockingPool: tt.fields.ZcnLockingPool,
			}
			if err := dp.Decode(tt.args.input, tt.args.tokenlock); (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				assert.Equal(t, tt.want, dp)
			}
		})
	}
}
