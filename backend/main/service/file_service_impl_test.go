package service_test

// func TestGetLocalFileContent(t *testing.T) {
// 	mockRepo := mocks.NewMockFileRepository(t)
// 	service := NewFileService(mockRepo)

// 	expectedFile := entities.NewFile(
// 		"",
// 		"test.js",
// 		[]byte("function main() {\n\tconsole.log('Hello, World!');\n}"),
// 	)
// 	mockRepo.EXPECT().ReadLocalFile("test.js").Return(*expectedFile, nil)
// 	file, err := service.GetLocalFileContent("test.js")

// 	assert.NoError(t, err)
// 	assert.Equal(t, *expectedFile, file)
// 	mockRepo.AssertExpectations(t)
// }

// func TestGetLocalFileContent_Error(t *testing.T) {
// 	mockRepo := mocks.NewMockFileRepository(t)
// 	service := NewFileService(mockRepo)

// 	mockRepo.EXPECT().ReadLocalFile("test.js").Return(entities.File{}, errors.New("file not found"))

// 	file, err := service.GetLocalFileContent("test.js")

// 	assert.Error(t, err)
// 	assert.Equal(t, entities.File{}, file)
// 	mockRepo.AssertExpectations(t)
// }

// func TestGetRemoteFileContent(t *testing.T) {
// 	mockRepo := mocks.NewMockFileRepository(t)
// 	service := NewFileService(mockRepo)

// 	expectedFile := entities.NewFile(
// 		"",
// 		"test.js",
// 		[]byte("function main() {\n\tconsole.log('Hello, World!');\n}"),
// 	)

// 	mockRepo.EXPECT().ReadRemoteFileByMD5("remote.js", "test-md5").Return(*expectedFile, nil)

// 	file, err := service.GetRemoteFileContent("remote.js", "test-md5")

// 	assert.NoError(t, err)
// 	assert.Equal(t, *expectedFile, file)
// 	mockRepo.AssertExpectations(t)
// }

// func TestGetRemoteFileContent_Error(t *testing.T) {
// 	mockRepo := mocks.NewMockFileRepository(t)
// 	service := NewFileService(mockRepo)

// 	mockRepo.EXPECT().ReadRemoteFileByMD5("remote.js", "test-md5").Return(entities.File{}, errors.New("file not found"))

// 	file, err := service.GetRemoteFileContent("remote.js", "test-md5")

// 	assert.Error(t, err)
// 	assert.Equal(t, entities.File{}, file)
// 	mockRepo.AssertExpectations(t)
// }
