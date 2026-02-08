# go-lab

**Experimental Laboratory for Go Internals**

This repository contains experimental code to verify the physical properties of the Go language (Runtime, GC, Memory, Scheduler).
It focuses on **measurable results** and **internal mechanisms**, not application logic.

## Experiment Protocol

1. **Hypothesis Driven**

- Every experiment must start with a hypothesis based on Go internals.
- Use the [Experiment Design Template](.github/ISSUE_TEMPLATE/experiment.md).

1. **Zero Dependencies**

- Use standard libraries (`unsafe`, `reflect`, `runtime`) whenever possible.
- Avoid black-box libraries to ensure clarity of the mechanism.

1. **Measurable**

- Results must be quantified using `testing.B` (Benchmark) or `runtime.ReadMemStats`.
- Focus on `allocs/op` and `ns/op`.

## Directory Structure

Managed by `go.work` (Go Workspace).

```text
go-lab/
├── go.work                # Workspace configuration
├── .agent.md              # AI Researcher guidelines
├── pkg/                   # Shared measurement tools (Local module)
└── experiments/           # Experiment logs
    ├── 01-struct-padding/
    ├── 02-goroutine-leak/
    └── ...
```

## Labeling Strategy

### Topics (Target)

- `topic:memory`: GC, Heap, Stack, Escape Analysis
- `topic:cpu`: Scheduler, Goroutines, Syscall
- `topic:data-structure`: Slice, Map, Channel internals

### Results (Outcome)

- `result:verified`: Hypothesis confirmed.
- `result:unexpected`: **Significant insight found**. (Behavior differed from hypothesis)
- `result:inconclusive`: No significant difference observed.

## Workflow

1. **Issue**: Create an issue to define the hypothesis.
2. **Code**: Implement the experiment in `experiments/<topic>`.
3. **Benchmark**: Run `go test -bench . -benchmem`.
4. **Report**: Post the results in the Issue and close it with a Result label.
