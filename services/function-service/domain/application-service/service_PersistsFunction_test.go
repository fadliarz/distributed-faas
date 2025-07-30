package application

// import (
// 	"context"
// 	"errors"
// 	"testing"

// 	"github.com/fadliarz/distributed-faas/services/function-service/domain/application-service/features/command"
// 	"github.com/fadliarz/distributed-faas/services/function-service/domain/domain-core"
// 	"github.com/fadliarz/distributed-faas/services/function-service/mocks"
// 	"github.com/stretchr/testify/require"
// 	"github.com/stretchr/testify/suite"
// )

// type FunctionApplicationServiceTestSuite struct {
// 	suite.Suite
// 	ctx context.Context

// 	sut               FunctionApplicationService
// 	dependencyManager *DependencyManager
// }

// type DependencyManager struct {
// 	Mapper       *mocks.MockMapper
// 	Service      *mocks.MockFunctionDomainService
// 	FunctionRepo *mocks.MockFunctionRepository
// }

// func (suite *FunctionApplicationServiceTestSuite) SetupTest() {
// 	suite.ctx = context.Background()

// 	suite.dependencyManager = &DependencyManager{
// 		Mapper:       mocks.NewMockMapper(suite.T()),
// 		Service:      mocks.NewMockFunctionDomainService(suite.T()),
// 		FunctionRepo: mocks.NewMockFunctionRepository(suite.T()),
// 	}

// 	suite.sut = NewFunctionApplicationService(suite.dependencyManager.Mapper, suite.dependencyManager.Service, suite.dependencyManager.FunctionRepo)
// }

// func (suite *FunctionApplicationServiceTestSuite) TestFunctionApplicationService_PersistsFunction_Success() {
// 	// Arrange
// 	cmd := createTestCommand()
// 	function := createTestFunction()

// 	// Set up expectations
// 	suite.dependencyManager.Mapper.EXPECT().CreateFunctionCommandToFunction(cmd).Return(function, nil)
// 	suite.dependencyManager.FunctionRepo.EXPECT().Save(function).Return(function.FunctionID, nil)

// 	// Act
// 	result, err := suite.sut.PersistFunction(cmd)

// 	// Assert
// 	require.NoError(suite.T(), err, "Expected no error when persisting function")
// 	require.NotEmpty(suite.T(), result, "Expected a valid Function ID to be returned")
// 	require.Equal(suite.T(), function.FunctionID.String(), result.String(), "Function ID should match")
// }

// func (suite *FunctionApplicationServiceTestSuite) TestFunctionApplicationService_PersistFunction_MapperError() {
// 	// Arrange
// 	cmd := createTestCommand()
// 	mapperError := errors.New("mapper error")

// 	// Set up expectations
// 	suite.dependencyManager.Mapper.EXPECT().CreateFunctionCommandToFunction(cmd).Return(nil, mapperError)

// 	// Act
// 	result, err := suite.sut.PersistFunction(cmd)

// 	// Assert
// 	require.Error(suite.T(), err, "Expected an error when mapper fails")
// 	require.Empty(suite.T(), result, "Expected no Function ID to be returned when mapper fails")
// }

// func (suite *FunctionApplicationServiceTestSuite) TestFunctionApplicationService_PersistFunction_SaveError() {
// 	// Arrange
// 	cmd := createTestCommand()
// 	function := createTestFunction()
// 	saveError := errors.New("")

// 	// Set up expectations
// 	suite.dependencyManager.Mapper.EXPECT().CreateFunctionCommandToFunction(cmd).Return(function, nil)
// 	suite.dependencyManager.FunctionRepo.EXPECT().Save(function).Return("", saveError)

// 	// Act
// 	result, err := suite.sut.PersistFunction(cmd)

// 	// Assert
// 	require.Error(suite.T(), err, "Expected an error when saving function fails")
// 	require.Empty(suite.T(), result, "Expected no Function ID to be returned when saving function fails")
// }

// func createTestCommand() *command.CreateFunctionCommand {
// 	return &command.CreateFunctionCommand{
// 		UserID:        "userid-123",
// 		SourceCodeURL: "https://google.com",
// 	}
// }

// func createTestFunction() *domain.Function {
// 	return &domain.Function{
// 		UserID:        "userid-123",
// 		FunctionID:    "functionid-123",
// 		SourceCodeURL: "https://google.com",
// 	}
// }

// func TestFunctionApplicationServiceTestSuite(t *testing.T) {
// 	suite.Run(t, new(FunctionApplicationServiceTestSuite))
// }
