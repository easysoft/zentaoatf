package main

import (
	"fmt"
	"log"
	"regexp"
)

func cleanup() {
	fmt.Println("cleanup")
}

func main() {
	pass, _ := regexp.MatchString(`^[0-9]{4}-[0-9]{2}-[0-9]{2}$`, "2020-12-18")
	if pass {
		log.Print(pass)
	}

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

	//html := "<html><head><title>Website Title</title></head><body><div class=\"random-class\"><h1>I like pie</h1><p>It's true!</p></div></body></html>"
	//fmt.Println(gohtml.FormatWithLineNo(html))
}
