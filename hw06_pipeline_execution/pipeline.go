package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	I   = interface{}
	In  = <-chan I
	Out = In
	Bi  = chan I
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(Bi)

	switch len(stages) {
	case 0:
		defer close(out)
		return out
	case 1:
		return stages[0](in)
	}

	go func() {
		defer close(out)

		stream := in
		for _, stage := range stages {
			stream = stage(stream)
		}

		for {
			select {
			case <-done:
				return
			case i, ok := <-stream:
				if !ok {
					return
				}
				out <- i
			}
		}
	}()

	return out
}
