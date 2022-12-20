package gglsec

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/xanzy/go-gitlab"
)

const (
	STATUS_OK     string = "OK"
	STATUS_FAILED string = "FAILED"

	SEVERITY_INFO     Severity = "INFO"
	SEVERITY_WARNING  Severity = "WARNING"
	SEVERITY_HIGH     Severity = "HIGH"
	SEVERITY_CRITICAL Severity = "CRITICAL"

	ENTITY_TYPE_GROUP    EnityType = "GROUP"
	ENTITY_TYPE_PROJECTS EnityType = "PROJECT"
)

type Rule interface {
	Apply() *RuleResult
}

type Status string

type Severity string

type EnityType string

type RuleMeta struct {
	Name        string
	Description string
	Severity    Severity
	EntityId    string
	EntityType  EnityType
}

type RuleMixin struct {
	Meta         *RuleMeta
	GitlabClient *gitlab.Client
}

type resultStatus bool

func (rs resultStatus) String() string {
	if rs {
		return STATUS_OK
	}
	return STATUS_FAILED
}

type RuleResult struct {
	Meta    *RuleMeta
	Status  resultStatus
	Message string
}

func NewRuleResult(meta *RuleMeta) *RuleResult {
	return &RuleResult{
		Meta:    meta,
		Status:  false,
		Message: "No message",
	}
}

func (rr *RuleResult) Sprint() string {
	return fmt.Sprint(rr.Meta.Name, "\t", rr.Meta.Severity, "\t", rr.Status, "\t", rr.Message)
}

func (rr *RuleResult) Println() {
	fmt.Println(rr.Sprint())
}

func (rr *RuleResult) Fprintln(w io.Writer) {
	fmt.Fprintln(w, rr.Sprint())
}

type RuleResults struct {
	ruleResults []*RuleResult
	count       int
	success     int
	failed      int
}

func NewRuleResults(results ...*RuleResult) *RuleResults {
	var (
		count   int
		success int
		failed  int
	)

	for _, result := range results {
		count++
		if result.Status {
			success++
		} else {
			failed++
		}
	}

	return &RuleResults{
		ruleResults: results,
		count:       count,
		success:     success,
		failed:      failed,
	}
}

func (rrs *RuleResults) Append(result *RuleResult) {
	rrs.ruleResults = append(rrs.ruleResults, result)
	rrs.count++
	if result.Status {
		rrs.success++
	} else {
		rrs.failed++
	}
}

func (rrs *RuleResults) PrintReport() {
	const OVERALL_STATUS_MESSAGE = "***OVERALL STATUS***\n\nSUCCESS: %d\nFAILED: %d\nSUCCESS PERCENT: %.f\n"

	tableWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.TabIndent)

	type sortedResultsKey struct {
		eid   string
		etype EnityType
	}

	sortedResults := make(map[sortedResultsKey][]*RuleResult)
	for _, res := range rrs.ruleResults {
		key := sortedResultsKey{
			res.Meta.EntityId,
			res.Meta.EntityType,
		}
		sortedResults[key] = append(sortedResults[key], res)
	}

	for key, results := range sortedResults {
		fmt.Printf("Scan for %s with ID %s\n", strings.ToLower(string(key.etype)), key.eid)
		for _, res := range results {
			res.Fprintln(tableWriter)
		}
		tableWriter.Flush()
		fmt.Println()
	}

	fmt.Printf(OVERALL_STATUS_MESSAGE, rrs.success, rrs.failed, rrs.SuccessOverall())
}

func (rrs *RuleResults) SuccessOverall() float32 {
	return float32(rrs.success) / float32(rrs.count) * 100
}
