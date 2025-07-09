package main

import (
	"encoding/json"
	"fmt"
	"log"
	"slices"
	"strings"

	"git.sxxfuture.net/taojiayi/super-cp/core"
	"git.sxxfuture.net/taojiayi/super-cp/utils"

	_ "git.sxxfuture.net/taojiayi/super-cp/uploaders/s3"
	"github.com/dustin/go-humanize"
	"github.com/joho/godotenv"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Short: "Copy files to object storage as a static website",
	Use:   "super-cp [flags] [job-names...]",
	Run: func(cmd *cobra.Command, jobNames []string) {
		dryRun := viper.GetBool("dry-run")
		configPath := viper.GetString("config")

		utils.Verbose("")
		utils.Verbosef("DryRun: %v", dryRun)
		utils.Verbosef("Config: %s", configPath)
		utils.Verbosef("  Jobs: %v", strings.Join(jobNames, ", "))
		utils.Verbose("")

		if len(jobNames) == 0 {
			jobNames = []string{"default"}
		}

		conf := core.MustLoadConfig(configPath)

		// validate job names
		for _, jobName := range jobNames {
			if _, ok := conf.Jobs[jobName]; !ok {
				log.Panicf("job '%s' not found\n", jobName)
			}
		}

		for _, jobName := range jobNames {
			utils.Verbosef("job started: %s", jobName)

			job := conf.Jobs[jobName]

			files, err := job.Source.WalkMatch()
			if err != nil {
				log.Panicf("walk source failed: %v\n", err)
			}

			for _, file := range files {
				for _, rule := range job.Rules {
					rule.Apply(file)
				}
			}

			files = lo.Filter(files, func(file *core.SourceFile, _ int) bool {
				if file.Excluded {
					utils.Verbosef("Excluded: %s", file.LocalPath)
				}
				return !file.Excluded
			})

			slices.SortFunc(files, func(a, b *core.SourceFile) int {
				return a.Index - b.Index
			})

			if viper.GetBool("verbose") {
				for _, file := range files {
					metadata, _ := json.MarshalIndent(file.Metadata, "          ", "    ")
					utils.Verbose("--------------------------------------------------")
					utils.Verbosef("   Local: %s", file.LocalPath)
					utils.Verbosef("  Remote: %s", file.RemotePath)
					utils.Verbosef("   Index: %d", file.Index)
					utils.Verbosef("    Size: %s", humanize.Bytes(uint64(file.Info.Size())))
					utils.Verbosef("Metadata: %s", string(metadata))
				}
			}

			if viper.GetBool("analyze") {
				totalSize := int64(0)
				for _, file := range files {
					totalSize += file.Info.Size()
				}

				fmt.Printf("\n╔═══ Job: %s\n", jobName)
				fmt.Printf("║     Source: %s\n", job.Source.Pattern)
				fmt.Printf("║     Target: %v\n", job.Dist.Uploader)
				fmt.Printf("║      Files: %d\n", len(files))
				fmt.Printf("║ Total Size: %s\n", humanize.Bytes(uint64(totalSize)))
				fmt.Printf("╚═%s\n", strings.Repeat("═", 50-2))
			}

			indexes := lo.Union(lo.Map(files, func(file *core.SourceFile, _ int) int {
				return file.Index
			}))

			slices.Sort(indexes)

			utils.Verbosef("Indexes: %v", indexes)

			for _, index := range indexes {
				utils.Verbosef("Upload index level: %d", index)

				files := lo.Filter(files, func(file *core.SourceFile, _ int) bool {
					return file.Index == index
				})

				if err := job.Dist.Uploader.Upload(files); err != nil {
					log.Panicf("upload files failed: %v\n", err)
				}
			}

			utils.Verbosef("Job '%s' done", jobName)
		}
	},
}

func init() {
	godotenv.Load()

	rootCmd.Flags().Bool("dry-run", false, "print all actions without actually do")
	rootCmd.Flags().BoolP("verbose", "v", false, "print verbose output")
	rootCmd.Flags().BoolP("analyze", "a", true, "print analyze output")
	rootCmd.Flags().IntP("concurrency", "j", 8, "concurrency level")

	rootCmd.Flags().StringP("config", "c", ".super-cp.yml", "config file path")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.BindPFlags(rootCmd.Flags())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Panic(err)
	}
}
