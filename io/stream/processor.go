package stream

// Processor is a function which handles input stream and send result into output stream.
// WARN: inputStream / inputErr should be closed by the previous Processor, while outputStream / outputErr should be closed by the current one.
type Processor func(inputStream *IOStream, inputErr *ErrorPasser) (outputStream *IOStream, outputErr *ErrorPasser)

func (cur Processor) Next(next Processor) Processor {
	return func(inputStream *IOStream, inputErr *ErrorPasser) (*IOStream, *ErrorPasser) {
		outputStream, outputErr := cur(inputStream, inputErr)
		return next(outputStream, outputErr)
	}
}

func BuildProcChain(procs ...Processor) Processor {

	if len(procs) == 0 {
		return nil
	}

	proc := procs[0]
	for i := 1; i < len(procs); i++ {
		proc = proc.Next(procs[i])
	}

	return proc

}
