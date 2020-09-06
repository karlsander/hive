package hive

import (
	"log"
	"testing"
	"time"

	"github.com/pkg/errors"
)

func TestHiveJobWithPool(t *testing.T) {
	h := New()

	doGeneric := h.Handle("generic", generic{}, PoolSize(3))

	grp := NewGroup()
	grp.Add(doGeneric("first"))
	grp.Add(doGeneric("first"))
	grp.Add(doGeneric("first"))

	if err := grp.Wait(); err != nil {
		log.Fatal(err)
	}
}

type badRunner struct{}

// Run runs a badRunner job
func (g badRunner) Run(job Job, do DoFunc) (interface{}, error) {
	return job.String(), nil
}

func (g badRunner) OnStart() error {
	return errors.New("fail")
}

func TestRunnerWithError(t *testing.T) {
	h := New()

	doBad := h.Handle("badRunner", badRunner{})

	_, err := doBad(nil).Then()
	if err == nil {
		t.Error("expected error, did not get one")
	}
}

func TestRunnerWithOptionsAndError(t *testing.T) {
	h := New()

	doBad := h.Handle("badRunner", badRunner{}, RetrySeconds(1), MaxRetries(1))

	_, err := doBad(nil).Then()
	if err == nil {
		t.Error("expected error, did not get one")
	}
}

type timeoutRunner struct{}

// Run runs a timeoutRunner job
func (g timeoutRunner) Run(job Job, do DoFunc) (interface{}, error) {
	time.Sleep(time.Duration(time.Second * 3))

	return nil, nil
}

func (g timeoutRunner) OnStart() error {
	return nil
}

func TestRunnerWithJobTimeout(t *testing.T) {
	h := New()

	doTimeout := h.Handle("timeout", timeoutRunner{}, TimeoutSeconds(1))

	if _, err := doTimeout("hello").Then(); err != ErrJobTimeout {
		t.Error("job should have timed out, but did not")
	}
}
