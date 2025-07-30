package application

// package service

// import (
// 	"context"
// 	"errors"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// 	"github.com/stretchr/testify/suite"

// 	"github.com/fadliarz/distributed-faas/services/invocation-service/domain/application-service/features/command"
// 	"github.com/fadliarz/distributed-faas/services/invocation-service/domain/domain-core"
// 	"github.com/fadliarz/distributed-faas/services/invocation-service/mocks"
// )

// // Test Suite with auto-generated mocks
// type InvocationApplicationServiceTestSuite struct {
// 	suite.Suite

// 	ctx               context.Context
// 	sut               *InvocationApplicationService
// 	dependencyManager *DependencyManager
// }

// type DependencyManager struct {
// 	Mapper         *mocks.MockMapper
// 	Service        *mocks.MockInvocationDomainService
// 	InvocationRepo *mocks.MockInvocationRepository
// 	FunctionRepo   *mocks.MockFunctionRepository
// }

// func (suite *InvocationApplicationServiceTestSuite) SetupTest() {
// 	suite.ctx = context.Background()

// 	suite.dependencyManager = &DependencyManager{
// 		Mapper:         mocks.NewMockMapper(suite.T()),
// 		Service:        mocks.NewMockInvocationDomainService(suite.T()),
// 		InvocationRepo: mocks.NewMockInvocationRepository(suite.T()),
// 		FunctionRepo:   mocks.NewMockFunctionRepository(suite.T()),
// 	}

// 	suite.sut = NewInvocationApplicationService(
// 		suite.dependencyManager.Mapper,
// 		suite.dependencyManager.Service,
// 		suite.dependencyManager.InvocationRepo,
// 		suite.dependencyManager.FunctionRepo,
// 	)
// }

// func (suite *InvocationApplicationServiceTestSuite) TestPersistInvocation_Success() {
// 	// Arrange
// 	cmd := createTestCommand()
// 	invocation := createTestInvocation()
// 	function := createTestFunction()

// 	// Set up expectations
// 	suite.dependencyManager.Mapper.EXPECT().CreateInvocationCommandToInvocation(cmd).Return(invocation, nil)
// 	suite.dependencyManager.FunctionRepo.EXPECT().FindByUserIDAndFunctionID(suite.ctx, domain.UserID(cmd.UserID), domain.FunctionID(cmd.FunctionID)).Return(function, nil)
// 	suite.dependencyManager.InvocationRepo.EXPECT().Save(suite.ctx, invocation).Return(invocation.InvocationID, nil)

// 	// Act
// 	result, err := suite.sut.PersistInvocation(suite.ctx, cmd)

// 	// Assert
// 	require.NoError(suite.T(), err, "")
// 	require.NotEmpty(suite.T(), result, "")
// 	assert.Equal(suite.T(), invocation.InvocationID, result)
// }

// func (suite *InvocationApplicationServiceTestSuite) TestPersistFunction_MappingError() {
// 	// Arrange
// 	cmd := createTestCommand()
// 	mappingError := errors.New("failed mapping")

// 	// Mock
// 	suite.dependencyManager.Mapper.EXPECT().CreateInvocationCommandToInvocation(cmd).Return(nil, mappingError)

// 	// Act
// 	result, err := suite.sut.PersistInvocation(suite.ctx, cmd)

// 	// Assert
// 	require.Error(suite.T(), err, "")
// 	require.Empty(suite.T(), result, "")
// 	require.ErrorContains(suite.T(), err, mappingError.Error())
// }

// func (suite *InvocationApplicationServiceTestSuite) TestPersistFunction_FunctionNotFoundError() {
// 	// Arrange
// 	cmd := createTestCommand()
// 	invocation := createTestInvocation()
// 	notFoundError := errors.New("")

// 	// Mock
// 	suite.dependencyManager.Mapper.EXPECT().CreateInvocationCommandToInvocation(cmd).Return(invocation, nil)
// 	suite.dependencyManager.FunctionRepo.EXPECT().FindByUserIDAndFunctionID(suite.ctx, domain.UserID(cmd.UserID), domain.FunctionID(cmd.FunctionID)).Return(nil, notFoundError)

// 	// Act
// 	result, err := suite.sut.PersistInvocation(suite.ctx, cmd)

// 	// Assert
// 	require.Error(suite.T(), err, "")
// 	require.Empty(suite.T(), result, "")
// 	require.ErrorContains(suite.T(), err, notFoundError.Error())
// }

// func TestInvocationApplicationServiceTestSuite(t *testing.T) {
// 	suite.Run(t, new(InvocationApplicationServiceTestSuite))
// }

// func createTestCommand() *command.CreateInvocationCommand {
// 	return &command.CreateInvocationCommand{
// 		UserID:     "user-123",
// 		FunctionID: "function-456",
// 	}
// }

// func createTestInvocation() *domain.Invocation {
// 	return &domain.Invocation{
// 		InvocationID: domain.InvocationID("invocation-789"),
// 		FunctionID:   domain.FunctionID("function-456"),
// 	}
// }

// func createTestFunction() *domain.Function {
// 	return &domain.Function{
// 		UserID:        domain.UserID("user-123"),
// 		FunctionID:    domain.FunctionID("function-456"),
// 		SourceCodeURL: "https://google.com",
// 	}
// }
