package mobileqq

import (
	"log"
	"os"
)

func init() {
	for _, dir := range []string{baseDir, cacheDir, logDir} {
		_, err := os.Stat(dir)
		if os.IsNotExist(err) {
			err = os.Mkdir(dir, 0777)
		}
		if err != nil {
			log.Fatalf("failed to mkdir %s, error %s", dir, err.Error())
		}
	}
	// logFile, err := os.OpenFile(path.Join(
	// 	logDir,
	// 	fmt.Sprintf(
	// 		"goqq-%s.log",
	// 		time.Now().Local().Format("20060102150405"),
	// 	),
	// ), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)
	// if err != nil {
	// 	log.Fatalf(err.Error())
	// }
	// log.SetOutput(io.MultiWriter(os.Stdout, logFile))
}
