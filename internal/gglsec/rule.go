package gglsec

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/xanzy/go-gitlab"
)

const (
	STATUS_OK     = "OK"
	STATUS_FAILED = "FAILED"

	SEVERITY_INFO     = "INFO"
	SEVERITY_WARNING  = "WARNING"
	SEVERITY_HIGH     = "HIGH"
	SEVERITY_CRITICAL = "CRITICAL"
)

type Rule interface {
	Apply() *RuleResult
}

type RuleMeta struct {
	Name        string
	Description string
	Severity    string
}

type RuleMixin struct {
	EntityId     string
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
	const OVERALL_STATUS_MESSAGE = "\n***OVERALL STATUS***\n\nSUCCESS: %d\nFAILED: %d\nSUCCESS PERCENT: %.f\n"

	tableWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.TabIndent)

	for _, res := range rrs.ruleResults {
		res.Fprintln(tableWriter)
	}

	tableWriter.Flush()

	fmt.Printf(OVERALL_STATUS_MESSAGE, rrs.success, rrs.failed, rrs.SuccessOverall())
}

func (rrs *RuleResults) SuccessOverall() float32 {
	return float32(rrs.success) / float32(rrs.count) * 100
}
