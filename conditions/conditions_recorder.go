package conditions

import (
	"database/sql"
	"io"
	"log"
	"net"
	"os"
	"sync/atomic"
	"time"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

type ConditionsRecorder struct{}

var connStr = "user=postgres dbname=weather sslmode=disable password=" + os.Getenv("PGPASSWORD") + " host=" + os.Getenv("PGHOST")

func StartServer() error {
	l, err := net.Listen("tcp", "127.0.0.1:10000")
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	RegisterConditionsServer(s, &ConditionsRecorder{})

	log.Print("running server")
	return s.Serve(l)
}

func (c *ConditionsRecorder) Report(stream Conditions_ReportServer) error {
	var reportCount uint64
	var startedAt time.Time

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("opening client stream")

	summary := ReportSummary{}

	db.SetMaxOpenConns(200)

	startedAt = time.Now()

	go func() {
		for {
			time.Sleep(time.Second)
			duration := time.Since(startedAt).Seconds()

			log.Println("inserted", reportCount, "reports in ", duration, "seconds")
			log.Println((float64(reportCount) / duration), "records per second")
		}
	}()

	for {
		report, err := stream.Recv()

		if err == io.EOF {
			log.Println("no further data from client, closing stream")
			return stream.SendAndClose(&summary)
		}

		go func(report *Condition) {
			_, err := db.Exec(
				"INSERT INTO conditions(time, location, temperature, humidity) VALUES ($1, $2, $3, $4)",
				time.Unix(report.GetTime().GetSeconds(), int64(report.GetTime().GetNanos())),
				report.GetLocation(),
				report.GetTemperature(),
				report.GetHumidity(),
			)

			if err != nil {
				log.Fatal(err)
			}
			atomic.AddUint64(&reportCount, 1)
		}(report)
	}
}
