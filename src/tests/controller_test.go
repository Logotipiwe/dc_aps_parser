package tests

/*
func TestParserNew(t *testing.T) {
	_, app := initAppWithMocks([]*domain.ParseResult{{ID: 0, ApsNum: 0}})
	defer app.StopAllParsersSync()
	router := gin.Default()
	_ = input.NewParserController(router, app.ParserService)

	r, err := http.NewRequest("POST", "/parser/new", nil)
	w := httptest.NewRecorder()
	assert.NoError(t, err)

	router.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"ID\":0}", w.Body.String())
}

func TestParsersManyNew(t *testing.T) {
	_, app := initAppWithMocks([]*domain.ParseResult{{ID: 0, ApsNum: 0}})
	defer app.StopAllParsersSync()
	router := gin.Default()
	_ = input.NewParserController(router, app.ParserService)

	r, err := http.NewRequest("POST", "/parser/new", nil)
	w := httptest.NewRecorder()
	assert.NoError(t, err)
	router.ServeHTTP(w, r)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"ID\":0}", w.Body.String())

	r, err = http.NewRequest("POST", "/parser/new", nil)
	w = httptest.NewRecorder()
	assert.NoError(t, err)
	router.ServeHTTP(w, r)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"ID\":1}", w.Body.String())

	r, err = http.NewRequest("POST", "/parser/new", nil)
	w = httptest.NewRecorder()
	assert.NoError(t, err)
	router.ServeHTTP(w, r)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"ID\":2}", w.Body.String())
}*/
