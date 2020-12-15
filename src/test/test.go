package main

import (
	"fmt"
	"github.com/yosssi/gohtml"
)

func cleanup() {
	fmt.Println("cleanup")
}

func main() {
	//c := make(chan os.Signal)
	//signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	//go func() {
	//	<-c
	//	cleanup()
	//	os.Exit(0)
	//}()
	//
	//for {
	//	fmt.Println("sleeping...")
	//	time.Sleep(10 * time.Second) // or runtime.Gosched() or similar per @misterbee
	//}

	//content := fileUtils.ReadFile("log/pytest-result.xml")
	//
	//pyTestSuite := model.PyTestSuites{}
	//err := xml.Unmarshal([]byte(content), &pyTestSuite)
	//if err == nil {
	//	testSuite := testingService.ConvertPyTestResult(pyTestSuite)
	//	log.Println(fmt.Sprintf("%v", testSuite))
	//}

	html := "<html><head><title>Website Title</title></head><body><div class=\"random-class\"><h1>I like pie</h1><p>It's true!</p></div></body></html>"
	fmt.Println(gohtml.FormatWithLineNo(html))
}
