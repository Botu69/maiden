package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/monome/maiden/pkg/catalog"
	"github.com/monome/maiden/pkg/lines"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var catalogUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "update configured catalogs",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		catalogUpdateRun(args)
	},
}

func init() {
	catalogCmd.AddCommand(catalogUpdateCmd)
}

// FIXME: merge this with pkg/dust/project.go:downloadURL
func downloadURLToFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}

// FIXME: refactor this to remove duplicate logic
func catalogUpdateRun(args []string) {
	sources := viper.GetStringMap("sources")
	for key, v := range sources {
		source := v.(map[string]interface{})

		kind, ok := source["kind"]
		if !ok {
			// FIXME: should probably log that the configuration is malformed
			continue
		}

		fmt.Printf("Updating: %s; ", key)
		switch kind {
		case "lines":
			rawpath, ok := source["output"]
			if !ok {
				fmt.Println("missing 'output' value config for source: ", key)
				break
			}
			path := os.ExpandEnv(rawpath.(string))
			fmt.Printf("fetching topics from lines... ")
			catalog := catalog.New()
			err := lines.GatherProjects(catalog)
			if err != nil {
				log.Fatalf("failed while gathering project: %s", err)
			}
			destination, err := os.Create(path)
			if err != nil {
				log.Fatalf("%s", err)
			}
			defer destination.Close()
			err = catalog.Store(destination)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("wrote: %s\n", path)

		case "download":
			rawpath, ok := source["output"]
			if !ok {
				fmt.Println("missing 'output' value config for source: ", key)
				break
			}
			path := os.ExpandEnv(rawpath.(string))
			rawurl, ok := source["url"]
			if !ok {
				fmt.Println("missing 'url' value config for source: ", key)
				break
			}
			url := rawurl.(string)
			fmt.Printf("fetching %s... ", url)
			if err := downloadURLToFile(path, url); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("wrote: %s\n", path)

		default:
			fmt.Println("unrecognized catalog source:", kind)
		}
	}
}
