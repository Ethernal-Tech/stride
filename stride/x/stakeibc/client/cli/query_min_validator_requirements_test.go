package cli_test

import (
	"fmt"
	"testing"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"google.golang.org/grpc/status"

	"github.com/Stride-labs/stride/testutil/network"
	"github.com/Stride-labs/stride/testutil/nullify"
	"github.com/Stride-labs/stride/x/stakeibc/client/cli"
	"github.com/Stride-labs/stride/x/stakeibc/types"
)

func networkWithMinValidatorRequirementsObjects(t *testing.T) (*network.Network, types.MinValidatorRequirements) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	minValidatorRequirements := &types.MinValidatorRequirements{}
	nullify.Fill(&minValidatorRequirements)
	state.MinValidatorRequirements = minValidatorRequirements
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), *state.MinValidatorRequirements
}

func TestShowMinValidatorRequirements(t *testing.T) {
	net, obj := networkWithMinValidatorRequirementsObjects(t)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	for _, tc := range []struct {
		desc string
		args []string
		err  error
		obj  types.MinValidatorRequirements
	}{
		{
			desc: "get",
			args: common,
			obj:  obj,
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			var args []string
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowMinValidatorRequirements(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryGetMinValidatorRequirementsResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.MinValidatorRequirements)
				require.Equal(t,
					nullify.Fill(&tc.obj),
					nullify.Fill(&resp.MinValidatorRequirements),
				)
			}
		})
	}
}