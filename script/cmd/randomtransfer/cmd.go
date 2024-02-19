package randomtransfer

import (
	"github.com/spf13/cobra"
)

var (
	privateKey             string
	rpc                    string
	receiverNumber         int64
	maxValuePerTransaction float64
	token                  string
	timeCrossTransfer      int64
	timeBetweenTrans       int64
	fileName               string
	Transfer               = &cobra.Command{
		Use:   "transfer",
		Short: "transfer",
		//Args:  cobra.MinimumNArgs(2),
		// Run: func(cmd *cobra.Command, args []string) {
		// 	transferRandom()
		// },
	}
	RandomTransfer = &cobra.Command{
		Use:   "random-transfer",
		Short: "Random transfer",
		//Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			transferRandom()
		},
	}
	GenKey = &cobra.Command{
		Use:   "gen-key",
		Short: "Gen key to file",

		//Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			genKeyToFile()

		},
	}
	ImportFileAndTransfer = &cobra.Command{
		Use:   "import-file-and-transfer",
		Short: "Import file and transfer",

		//Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			importFileAndTransfer()

		},
	}
)

func init() {
	Transfer.AddCommand(RandomTransfer, GenKey, ImportFileAndTransfer)

	Transfer.PersistentFlags().StringVarP(&privateKey, "prvkey", "p", "", "private key")
	Transfer.PersistentFlags().StringVarP(&rpc, "rpc", "r", "", "rpc")
	Transfer.PersistentFlags().Int64Var(&receiverNumber, "receiver-number", 1000, "receiver-number")
	Transfer.PersistentFlags().Int64Var(&timeCrossTransfer, "time-cross-transfer", 30, "time-cross-transfer")
	Transfer.PersistentFlags().Float64Var(&maxValuePerTransaction, "max-value-per-transaction", 0.01, "max-value-per-transaction")
	Transfer.PersistentFlags().StringVar(&token, "token", "", "token")
	Transfer.PersistentFlags().Int64Var(&timeBetweenTrans, "time-between-trans", 500, "time-between-trans")
	Transfer.PersistentFlags().StringVarP(&fileName, "file-name", "f", "", "file-name")
	Transfer.MarkFlagsRequiredTogether("prvkey", "rpc")

}
