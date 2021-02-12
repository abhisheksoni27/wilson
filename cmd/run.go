package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/abhisheksoni27/wilson/test_case"
	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run all tests in config directory",
	Long:  `run all tests in config directory`,
	Run:   Run,
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.PersistentFlags().StringP("config", "c", "", "config directory (default is $HOME/.wilson/)")

	runCmd.PersistentFlags().Int16P("max-tests-at-a-time", "m", 4, "max-tests-at-a-time number of tests to run in parallel")
}

func Run(cmd *cobra.Command, args []string) {
	fmt.Printf("hello there ✨\n\n")

	success := color.New(color.FgGreen).SprintFunc()
	primary := color.New(color.FgBlue).SprintFunc()
	failure := color.New(color.FgRed).SprintFunc()

	configDir := cmd.Flag("config").Value.String()

	if configDir == "" {
		usr, _ := user.Current()
		dir := usr.HomeDir
		configDir = dir + "/.wilson"
	}

	maxTestsAtATime, _ := strconv.Atoi(cmd.Flag("max-tests-at-a-time").Value.String())
	if maxTestsAtATime <= 0 {
		maxTestsAtATime = 4
	} else if maxTestsAtATime > 64 {
		maxTestsAtATime = 64
	}

	fmt.Printf("using %v as the config directory \n\n", success(configDir))
	fmt.Printf("will run %v tests in parallel \n\n", success(maxTestsAtATime))

	allTests := readTests(configDir)

	if len(allTests) == 0 {
		log.Fatal("no test cases found")
	}

	httpClient := resty.New()

	semaphore := make(chan struct{}, maxTestsAtATime)
	waitGroup := sync.WaitGroup{}

	for _, testCase := range allTests {
		semaphore <- struct{}{}
		waitGroup.Add(1)
		go func(testCase test_case.TestCase) {
			fmt.Printf("%v %v\n\n", primary("[RUNNING]"), testCase.URL)
			defer func() {
				<-semaphore
			}()

			defer waitGroup.Done()

			err := testCase.Run(httpClient)
			if err == nil {
				fmt.Printf("%v URL = %v\n\n", success("[PASSED]"), testCase.URL)
				return
			}

			fmt.Printf("%v URL = %v\n\tErr = %s\n\n", failure("[FAILED]"), testCase.URL, failure(err.Error()))
		}(testCase)
	}

	waitGroup.Wait()
}

func readTests(configDir string) []test_case.TestCase {
	allTests := make([]test_case.TestCase, 0)

	err := filepath.Walk(configDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && strings.HasSuffix(path, ".json") {
				fileData, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}

				var testCases []test_case.TestCase
				err = json.Unmarshal(fileData, &testCases)
				allTests = append(allTests, testCases...)
			}
			return nil
		})

	if err != nil {
		log.Fatal(err)
	}

	return allTests
}
