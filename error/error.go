package error

type ErrorLevel uint8
type WorkStage uint8

// error level
const (
	LEVEL_WARNING ErrorLevel = 1
	LEVEL_ERROR ErrorLevel = 2
)

// error stage, i.e. during compiling, linking, ...
const(
	STAGE_COMPILE WorkStage = 1
	STAGE_LINK WorkStage = 2
)

func handleException(stage WorkStage, level ErrorLevel){

}
