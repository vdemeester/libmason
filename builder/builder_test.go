package builder

import (
	"fmt"
	"reflect"
	"testing"

	"golang.org/x/net/context"

	"github.com/vdemeester/libmason"
)

const Value = "value"

type UpdateConfigStep struct{}

func (s *UpdateConfigStep) Execute(ctx context.Context, helper libmason.Helper, config *Config) (*Config, error) {
	value, ok := config.Get(Value)
	if !ok {
		config.Put(Value, 1)
		return config, nil
	}
	config.Put(Value, value.(int)+1)
	return config, nil
}

func (s *UpdateConfigStep) String() string {
	return "UpdateConfigStep"
}

func TestBuilderRunNoSteps(t *testing.T) {
	helper := &NopHelper{}
	builder := NewBuilder(helper)
	imageID, err := builder.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if imageID != "" {
		t.Fatalf("expected empty ImageID, got %s", imageID)
	}
}

func TestBuilderRunFailAtFirstError(t *testing.T) {
	steps := []Step{
		&UpdateConfigStep{},
		&UpdateConfigStep{},
		&ErrorStep{},
		&UpdateConfigStep{},
	}
	helper := &NopHelper{}
	builder := WithSteps(NewBuilder(helper), steps)
	_, err := builder.Run(context.Background())
	if err == nil || err != errStep {
		t.Fatalf("expected an errStep error, got %v", err)
	}
	actual, _ := builder.currentConfig.Get(Value)
	if actual.(int) != 2 {
		t.Fatalf("expected current config value to be 2, got %d", actual.(int))
	}
}

type logAggregator struct {
	logs []string
}

func (l *logAggregator) logF(format string, values ...interface{}) {
	l.logs = append(l.logs, fmt.Sprintf(format, values...))
}

func TestBuilderWithCustomLogFunc(t *testing.T) {
	m := logAggregator{}
	steps := []Step{
		&UpdateConfigStep{},
		&UpdateConfigStep{},
		&UpdateConfigStep{},
	}
	helper := &NopHelper{}
	builder := WithLogFunc(WithSteps(NewBuilder(helper), steps), m.logF)
	_, err := builder.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(m.logs) != 3 {
		t.Fatalf("expected 3 lines of log, got %v", m.logs)
	}
	expectedLogs := []string{
		"Step 0: UpdateConfigStep",
		"Step 1: UpdateConfigStep",
		"Step 2: UpdateConfigStep",
	}
	if !reflect.DeepEqual(m.logs, expectedLogs) {
		t.Fatalf("expected logs %v, got %v", expectedLogs, m.logs)
	}
}
