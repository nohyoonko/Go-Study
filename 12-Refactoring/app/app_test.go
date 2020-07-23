package app

import (
	"12-Refactoring/model"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTodos(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(MakeHandler())
	defer ts.Close() //함수가 끝나기 전에 서버 닫아줘야함

	/* add id1 */
	resp, err := http.PostForm(ts.URL+"/todos", url.Values{"name": {"Test todo"}}) //FormValue로 받기 때문에 PostForm
	assert.NoError(err)                                                            //에러가 없어야한다.
	assert.Equal(http.StatusCreated, resp.StatusCode)                              //우리가 만든거니까 StatusCreated인지 확인
	var todo model.Todo
	err = json.NewDecoder(resp.Body).Decode(&todo)
	assert.NoError(err)
	assert.Equal(todo.Name, "Test todo")
	id1 := todo.ID

	/* add id2 */
	resp, err = http.PostForm(ts.URL+"/todos", url.Values{"name": {"Test todo2"}})
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)
	err = json.NewDecoder(resp.Body).Decode(&todo)
	assert.NoError(err)
	assert.Equal(todo.Name, "Test todo2")
	id2 := todo.ID

	/* check id1 and id2 */
	resp, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	todos := []*model.Todo{}
	err = json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 2)
	for _, t := range todos {
		if t.ID == id1 {
			assert.Equal("Test todo", t.Name)
		} else if t.ID == id2 {
			assert.Equal("Test todo2", t.Name)
		} else {
			assert.Error(fmt.Errorf("testID should be id1 or id2"))
		}
	}

	//?complete가 FormValue로 받는 부분
	resp, err = http.Get(ts.URL + "/complete-todo/" + strconv.Itoa(id1) + "?complete=true")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	resp, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	todos = []*model.Todo{}
	err = json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 2)
	for _, t := range todos {
		if t.ID == id1 {
			assert.True(t.Completed)
		}
	}

	//GoLang에서 DELETE는 제공하지 않아서 따로 request 만들어줘야함
	req, _ := http.NewRequest("DELETE", ts.URL+"/todos/"+strconv.Itoa(id1), nil)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	resp, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	todos = []*model.Todo{}
	err = json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 1)
	for _, t := range todos {
		assert.Equal(t.ID, id2)
	}
}
