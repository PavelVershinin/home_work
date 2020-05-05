package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	I   = interface{}
	In  = <-chan I
	Out = In
	Bi  = chan I
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	outCh := make(Bi)
	pipelineCh := chain(in, stages)

	go func() {
		defer close(outCh)
		for {
			select {
			case <-done:
				return
			case i, ok := <-pipelineCh:
				if !ok {
					return
				}
				outCh <- i
			}
		}
	}()

	return outCh
}

func chain(in In, stages []Stage) Out {
	switch len(stages) {
	case 0:
		return in
	case 1:
		return stages[0](in)
	}
	return chain(stages[0](in), stages[1:])
}
