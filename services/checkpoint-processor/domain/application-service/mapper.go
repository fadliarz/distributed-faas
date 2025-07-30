package application

func NewCheckpointProcessorDataMapper() CheckpointProcessorDataMapper {
	return &CheckpointProcessorDataMapperImpl{}
}

type CheckpointProcessorDataMapperImpl struct{}
