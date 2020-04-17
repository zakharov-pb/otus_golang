package hw06_pipeline_execution //nolint:golint,stylecheck

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	sleepPerStage = time.Millisecond * 100
	fault         = sleepPerStage / 2
)

type Tester struct {
	ID      int       `json:"id"`
	Name    string    `json:"name"`
	Info    string    `json:"info"`
	Created time.Time `json:"created"`
}

func CreateTester(id int, name string, info string) *Tester {
	return &Tester{
		ID:      id,
		Name:    name,
		Info:    info,
		Created: time.Date(2020, time.April, 1, 0, 0, 0, 0, time.UTC),
	}
}

func (t *Tester) String() string {
	data, err := json.Marshal(t)
	if err != nil {
		return fmt.Sprintf(`{"error": "%v"}`, err)
	}
	return string(data)
}

func TestPipeline(t *testing.T) {
	// Stage generator
	g := func(name string, f func(v I) I) Stage {
		return func(in In) Out {
			out := make(Bi)
			go func() {
				defer close(out)
				for v := range in {
					time.Sleep(sleepPerStage)
					out <- f(v)
				}
			}()
			return out
		}
	}

	stages := []Stage{
		g("Dummy", func(v I) I { return v }),
		g("Multiplier (* 2)", func(v I) I { return v.(int) * 2 }),
		g("Adder (+ 100)", func(v I) I { return v.(int) + 100 }),
		g("Stringifier", func(v I) I { return strconv.Itoa(v.(int)) }),
	}

	t.Run("simple case", func(t *testing.T) {
		in := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, nil, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Equal(t, result, []string{"102", "104", "106", "108", "110"})
		require.Less(t,
			int64(elapsed),
			// ~0.8s for processing 5 values in 4 stages (100ms every) concurrently
			int64(sleepPerStage)*int64(len(stages)+len(data)-1)+int64(fault))
	})

	t.Run("done case", func(t *testing.T) {
		in := make(Bi)
		done := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		// Abort after 200ms
		abortDur := sleepPerStage * 2
		go func() {
			<-time.After(abortDur)
			close(done)
		}()

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, done, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Len(t, result, 0)
		require.Less(t, int64(elapsed), int64(abortDur)+int64(fault))
	})
}

func TestAdditionalPipeline(t *testing.T) {
	// Unformatted information
	data := []interface{}{
		0.5,
		`Hello
							
			World!`,
		1000,
		CreateTester(
			1,
			"User     01",
			`Unformatted     TeXt`,
		),
	}
	reference := []string{"0.5", "HELLO WORLD!", "1000", "{'ID':1,'NAME':'USER 01','INFO':'UNFORMATTED TEXT','CREATED':'2020-04-01T00:00:00Z'}"}
	// String stage generator
	g := func(name string, f func(v string) string) Stage {
		return func(in In) Out {
			out := make(Bi)
			go func() {
				defer close(out)
				for v := range in {
					time.Sleep(sleepPerStage)
					str := fmt.Sprintf("%v", v)
					out <- f(str)
				}
			}()
			return out
		}
	}
	stages := []Stage{
		g("Trim space", func() func(string) string {
			var re = regexp.MustCompile(`[ \t\n]+`)
			return func(v string) string { return re.ReplaceAllString(v, " ") }
		}()),
		g("Replace", func(v string) string { return strings.ReplaceAll(v, "\"", "'") }),
		g("ToUpper", func(v string) string { return strings.ToUpper(v) }),
	}

	t.Run("string case", func(t *testing.T) {
		in := make(Bi)
		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, nil, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Equal(t, result, reference)
		require.Less(t,
			int64(elapsed),
			int64(sleepPerStage)*int64(len(stages)+len(data)-1)+int64(fault),
		)
	})

	t.Run("stages nill case", func(t *testing.T) {
		start := time.Now()
		require.Nil(t, ExecutePipeline(nil, nil))
		elapsed := time.Since(start)
		require.Less(t,
			int64(elapsed),
			int64(sleepPerStage)*int64(len(stages)+len(data)-1)+int64(fault))
	})

	t.Run("input nill case", func(t *testing.T) {
		result := make([]interface{}, 0)
		start := time.Now()
		require.Nil(t, ExecutePipeline(nil, nil, stages...))
		elapsed := time.Since(start)

		require.Equal(t, len(result), 0)
		require.Less(t,
			int64(elapsed),
			int64(sleepPerStage)*int64(len(stages)+len(data)-1)+int64(fault))
	})

	t.Run("call done", func(t *testing.T) {
		done := make(Bi)
		in := make(Bi)
		go func() {
			defer close(done)
			for index, v := range data {
				if index == 2 {
					break
				}
				in <- v
			}
			time.Sleep(time.Millisecond * 10)
		}()
		result := make([]string, 0, 10)
		for s := range ExecutePipeline(in, done, stages...) {
			result = append(result, s.(string))
		}
		require.Less(t, len(result), len(reference))
	})
}
