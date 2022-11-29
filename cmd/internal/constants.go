package internal

const (
	CopyBufferSize         = 1024 * 1024      // 1 MB
	InitialStdinBufferSize = 10 * 1024 * 1024 // 10 MB
)

const (
	Bannor = `==================================================================
 PipeCheck: The following was read in and will be passed through:
==================================================================`

	LongDescription = `PipeCheck will output data recieved from STDIN to STDERR. It will then block on user confirmation in order to proceed with the chained (piped) commands.

If the user decides to proceed, the data provided to STDIN will be directed to STDOUT and the process will terminate successfully. If the user decides **not** to proceed, the process will not send any data to STDOUT and it will terminate with a 1.

Warnings:
 - When the process terminates unsuccessfully, subsequent pipes may still run depending on the environment settings.
 - PipeCheck is **not** designed to function with updating input (ex: tail -f).
`

	SkipValidationMessage = `===================================================
 Skipping Validation, will proceed to next command
===================================================`

	ValidationMessage = `===============================
 Proceed with this data (y/N):`

	ValidationMessageBannor = `===============================`
)
