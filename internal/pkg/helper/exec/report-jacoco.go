package execHelper

import (
	"encoding/xml"
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	dateUtils "github.com/easysoft/zentaoatf/pkg/lib/date"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	"github.com/fatih/color"
	"time"
)

const (
	Tmpl = `Counter              Missed     Covered
------------------------------------------
%11s          %s          %s     
%11s          %s          %s     
%11s          %s          %s     
%11s          %s          %s     
%11s          %s          %s     
%11s          %s          %s     
`
)

func GenJacocoCovReport() (report *commDomain.JacocoResult) {
	report = &commDomain.JacocoResult{}

	content := fileUtils.ReadFileBuf(commConsts.JacocoReport)

	xml.Unmarshal(content, &report)

	var params []interface{}
	for _, counter := range report.Counter {
		params = append(params, string(counter.Type))
		params = append(params, counter.Missed)
		params = append(params, counter.Covered)
	}

	title := i118Utils.Sprintf("jacoco_report")

	msg := dateUtils.DateTimeStr(time.Now()) + " " + title + " \n" + i118Utils.Sprintf(Tmpl, params...)

	logUtils.ExecConsole(color.FgCyan, msg)
	logUtils.ExecResult(msg)

	return
}
