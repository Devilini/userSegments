package main

import (
	"github.com/sirupsen/logrus"
	"userSegments/interanal/app"
	"userSegments/interanal/config"
)

func main() {
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	//
	//logging.L(ctx).Info("config initializing")
	cfg := config.GetConfig()

	//ctx = logging.ContextWithLogger(ctx, logging.NewLogger())

	a, err := app.NewApp(cfg)
	if err != nil {
		logrus.WithError(err).Fatal("app.NewApp")
	}

	//logging.L(ctx).Info("Running Application")
	a.Run()
	//if err != nil {
	//	logging.WithError(ctx, err).Fatal("app.Run")
	//	return
	//}
}

//package main
//
//import (
//	"encoding/json"
//	"log"
//	"net/http"
//)
//
//type test_struct struct {
//	Add string
//}
//
//func process(w http.ResponseWriter, r *http.Request) {
//	decoder := json.NewDecoder(r.Body)
//	var t test_struct
//	err := decoder.Decode(&t)
//	if err != nil {
//		log.Println(err)
//	}
//	log.Println(t.Add)
//
//	//buf, err := io.ReadAll(r.Body)
//	//if err != nil {
//	//	panic(err)
//	//}
//	//fmt.Printf("Body: %s\n", buf)
//}
//
//func main() {
//	server := http.Server{
//		Addr: ":8000",
//	}
//	http.HandleFunc("/api/users/5/segments", process)
//	server.ListenAndServe()
//}
