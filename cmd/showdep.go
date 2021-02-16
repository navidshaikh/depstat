package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// showdepCmd represents the showdep command
var showdepCmd = &cobra.Command{
	Use:   "showdep",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("showdep called")

		url, _ := cmd.Flags().GetString("url")
		file := "output.txt"

		// get contents from url
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()

		// write contents to "output.txt"
		out, _ := os.Create(file)
		defer out.Close()
		src := &PassThru{Reader: resp.Body, total: float64(resp.ContentLength)}
		size, err := io.Copy(out, src)
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Printf("\nFile Transferred. (%.1f MB)\n", float64(size)/bytesToMegaBytes)

		// search file for key
		res, err := searchFile("output.txt", "b/LICENSES", false)
		if err != nil {
			log.Fatal(err)
		}

		// replace file contents with search result
		err = os.Remove("output.txt")
		if err != nil {
			log.Fatal(err)
		}
		f, err := os.Create("output.txt")

		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()

		_, err2 := f.WriteString(res)

		if err2 != nil {
			log.Fatal(err2)
		}

		fmt.Println("re-created output.txt")

		// search updated file for key
		res, err = searchFile("output.txt", "diff", true)
		if err != nil {
			log.Fatal(err)
		}

		// do some modifications

		// show result and clean up
		fmt.Println(res)
	},
}

func init() {
	rootCmd.AddCommand(showdepCmd)
	showdepCmd.Flags().StringP("url", "u", "", "URL | URL to github patch")
	//showdepCmd.Flags().StringP("file", "f", "", "Filename | Name of txt file")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showdepCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showdepCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
