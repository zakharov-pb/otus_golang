package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	// I define interface
	I = interface{}
	// In channel of read data
	In = <-chan I
	// Out channel of read data
	Out = In
	// Bi input-output channel
	Bi = chan I
)

// Stage worker
type Stage func(in In) (out Out)

// ExecutePipeline function create pipeline
func ExecutePipeline(in In, done In, stages ...Stage) Out {
	countStages := len(stages)
	if countStages == 0 || in == nil {
		return nil
	}
	out := make(Bi)
	go func() {
		defer close(out)
		for {
			select {
			case val, ok := <-in:
				if !ok {
					return
				}
				select {
				case out <- val:
				case <-done:
					return
				}
			case <-done:
				return
			}
		}
	}()
	o := stages[0](out)
	if countStages > 1 { //nolint:gomnd
		return ExecutePipeline(o, done, stages[1:]...)
	}
	return o
}
