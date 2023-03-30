package main

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/taramakage/gon-verifier/internal/scorecard"
	"github.com/taramakage/gon-verifier/internal/verifier"
	"os"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "gon-verify",
		Short: "GoN evidence verify tools",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("invalid argument")
			}
			filePath := args[0]

			{
				task := []string{"A1", "A2", "A3", "A4", "A5", "A6"}
				opt := &verifier.Options{
					task,
					scorecard.DefaultStageOneTaskPoint,
					verifier.VerifyRegistryStageOne,
				}
				gv := verifier.NewGonVerifier("", opt)
				err := gv.Verify(filePath)
				if err != nil {
					return err
				}
			}

			{
				task := []string{"A7", "A8", "A9", "A10", "A11", "A12", "A13", "A14", "A15", "A16", "A17", "A18", "A19", "A20"}
				opt := &verifier.Options{
					task,
					scorecard.DefaultStageTwoTaskPoint,
					verifier.VerifyRegistryStageTwo,
				}
				gv := verifier.NewGonVerifier("", opt)
				err := gv.Verify(filePath)
				if err != nil {
					return err
				}
			}

			{
				task := []string{"A7", "A9", "A11", "A16", "A17", "A19"}
				opt := &verifier.Options{
					task,
					scorecard.DefaultStageTwoBTaskPoint,
					verifier.VerifyRegistryStageTwoShadow,
				}
				gv := verifier.NewGonVerifier("", opt)
				err := gv.Verify(filePath)
				if err != nil {
					return err
				}
			}

			{
				task := []string{"B1", "B2", "B5", "B6", "B7"}
				opt := &verifier.Options{
					task,
					scorecard.DefaultStageThreeTaskPoint,
					verifier.VerifyRegistryStageThree,
				}
				gv := verifier.NewGonVerifier("", opt)
				err := gv.Verify(filePath)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
