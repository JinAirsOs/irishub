package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire" // XXX fix
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/irisnet/irishub/client/context"
)

// GetCmdQuerySigningInfo implements the command to query signing info.
func GetCmdQuerySigningInfo(storeName string, cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "signing-info [validator-pubkey]",
		Short: "Query a validator's signing information",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pk, err := sdk.GetValPubKeyBech32(args[0])
			if err != nil {
				return err
			}

			key := slashing.GetValidatorSigningInfoKey(sdk.ValAddress(pk.Address()))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryStore(key, storeName)
			if err != nil {
				return err
			}
			if len(res) ==0 {
				return fmt.Errorf("the signing information of this validator %s is empty, please make sure its existence", args[0])
			}

			signingInfo := new(slashing.ValidatorSigningInfo)
			cdc.MustUnmarshalBinary(res, signingInfo)

			switch viper.Get(cli.OutputFlag) {

			case "text":
				human := signingInfo.HumanReadableString()
				fmt.Println(human)

			case "json":
				// parse out the signing info
				output, err := wire.MarshalJSONIndent(cdc, signingInfo)
				if err != nil {
					return err
				}
				fmt.Println(string(output))
			}

			return nil
		},
	}

	return cmd
}
