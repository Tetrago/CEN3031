package course

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type course struct {
	Identifier  string `json:"code"`
	Description string `json:"title"`
}

var cache = make(map[string][]course)

func queryDepartmentCourses(dep string) ([]course, error) {
	if v, ok := cache[dep]; ok {
		return v, nil
	}

	url := fmt.Sprintf("https://catalog.ufl.edu/course-search/api/?page=fose&route=search&subject=%s", dep)
	request := fmt.Sprintf(`{"other":{"srcdb":""},"criteria":[{"field":"subject","value":"%s"}]}`, dep)

	res, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(request)))
	if err != nil {
		cache[dep] = nil
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		cache[dep] = nil
		return nil, err
	}

	var dest struct {
		Courses []course `json:"results"`
	}

	if err := json.Unmarshal(body, &dest); err != nil {
		cache[dep] = nil
		return nil, err
	}

	cache[dep] = dest.Courses
	return dest.Courses, nil
}
