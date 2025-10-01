package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/CodeAndCraft-Online/cortex-api/internal/swat"
	"github.com/CodeAndCraft-Online/cortex-api/swat/tests"
)

type TestResult = swat.TestResult

type TestGroup struct {
	Category string            `json:"category"`
	Results  []swat.TestResult `json:"results"`
	Passed   int               `json:"passed"`
	Failed   int               `json:"failed"`
}

type SWATReport struct {
	Title     string        `json:"title"`
	Version   string        `json:"version"`
	BaseURL   string        `json:"base_url"`
	Timestamp time.Time     `json:"timestamp"`
	Groups    []TestGroup   `json:"groups"`
	Total     int           `json:"total_tests"`
	Passed    int           `json:"passed"`
	Failed    int           `json:"failed"`
	Coverage  float64       `json:"coverage_percent"`
	Runtime   time.Duration `json:"runtime"`
}

const SWAT_VERSION = "1.0.0"

func main() {
	// Command line flags
	baseURL := flag.String("base-url", "http://codeandcraft.online:4321/api", "Base URL of the API to test")
	runTests := flag.String("run", "all", "Comma-separated list of test categories to run (auth,posts,comments,subs,votes,users,health)")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	reportFormat := flag.String("report", "console", "Report format: console, json")
	flag.Parse()

	fmt.Printf("🧠 CORTEX API SWAT (Software Assurance Testing) Suite v%s\n", SWAT_VERSION)
	fmt.Printf("🔗 Testing API at: %s\n\n", *baseURL)

	// Initialize SWAT client
	client := tests.NewSWATClient(*baseURL)

	// Determine which test groups to run
	testCategories := parseTestCategories(*runTests)

	report := SWATReport{
		Title:     "Cortex API SWAT Report",
		Version:   SWAT_VERSION,
		BaseURL:   *baseURL,
		Timestamp: time.Now(),
		Groups:    []TestGroup{},
	}

	startTime := time.Now()

	// Run test groups
	for _, category := range testCategories {
		group := runTestGroup(client, category, *verbose)
		report.Groups = append(report.Groups, group)
		report.Total += len(group.Results)
		report.Passed += group.Passed
		report.Failed += group.Failed
	}

	report.Runtime = time.Since(startTime)
	if report.Total > 0 {
		report.Coverage = float64(report.Passed) / float64(report.Total) * 100
	}

	// Clean up test data
	client.Cleanup()

	// Output report
	outputReport(report, *reportFormat, *verbose)

	// Exit with appropriate code
	if report.Failed > 0 {
		os.Exit(1)
	}
	os.Exit(0)
}

func parseTestCategories(runTests string) []string {
	if runTests == "all" {
		return []string{"health", "auth", "posts", "comments", "subs", "votes", "users", "security"}
	}
	return strings.Split(runTests, ",")
}

func runTestGroup(client *tests.SWATClient, category string, verbose bool) TestGroup {
	group := TestGroup{
		Category: strings.Title(category),
		Results:  []swat.TestResult{},
	}

	fmt.Printf("🔍 Running %s Tests\n", strings.Title(category))

	var results []swat.TestResult
	var err error

	switch category {
	case "health":
		results, err = tests.RunHealthTests(client, verbose)
	case "auth":
		results, err = tests.RunAuthTests(client, verbose)
	case "posts":
		results, err = tests.RunPostsTests(client, verbose)
	case "comments":
		results, err = tests.RunCommentsTests(client, verbose)
	case "subs":
		results, err = tests.RunSubsTests(client, verbose)
	case "votes":
		results, err = tests.RunVotesTests(client, verbose)
	case "users":
		results, err = tests.RunUsersTests(client, verbose)
	case "security":
		results, err = tests.RunSecurityTests(client, verbose)
	default:
		fmt.Printf("   ⚠️  Unknown test category: %s\n", category)
		return group
	}

	if err != nil {
		fmt.Printf("   ❌ Error running %s tests: %v\n", category, err)
		return group
	}

	group.Results = results
	for _, result := range results {
		if result.Status == "PASS" {
			group.Passed++
			printResult(result, "✅", verbose)
		} else {
			group.Failed++
			printResult(result, "❌", verbose)
		}
	}

	fmt.Printf("   📊 %s: %d passed, %d failed\n\n", strings.Title(category), group.Passed, group.Failed)
	return group
}

func printResult(result TestResult, icon string, verbose bool) {
	fmt.Printf("   %s %s", icon, result.Name)
	if result.Error != "" && verbose {
		fmt.Printf(" - %s", result.Error)
	}
	fmt.Printf(" (%dms)\n", result.Duration.Milliseconds())
}

func outputReport(report SWATReport, format string, verbose bool) {
	fmt.Println("==================")
	fmt.Println("SWAT SUMMARY")
	fmt.Printf("Total Tests: %d\n", report.Total)
	fmt.Printf("Passed: %d\n", report.Passed)
	fmt.Printf("Failed: %d\n", report.Failed)
	fmt.Printf("Coverage: %.1f%%\n", report.Coverage)
	fmt.Printf("Total Runtime: %v\n", report.Runtime)
	fmt.Println("==================")

	switch format {
	case "json":
		outputJSON(report)
	default:
		// console is default, already printed above
	}
}

func outputJSON(report SWATReport) {
	jsonData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		log.Fatal("Failed to marshal JSON report:", err)
	}
	fmt.Println("\n--- JSON REPORT ---")
	fmt.Println(string(jsonData))
}
