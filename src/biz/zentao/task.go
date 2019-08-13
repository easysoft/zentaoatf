package zentao

import (
	"encoding/json"
	"github.com/easysoft/zentaoatf/src/client"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/utils"
)

func GetTaskInfo(baseUrl string, taskId string) model.TestTask {
	params := [][]string{{"taskID", taskId}}

	myurl := baseUrl + utils.GenSuperApiUri("testtask", "getById", params)
	dataStr, ok := client.Get(myurl, nil)

	if ok {
		var task model.TestTask
		json.Unmarshal([]byte(dataStr), &task)

		return task
	}

	return model.TestTask{}
}
