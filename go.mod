module github.com/metno/roadlabels

go 1.25.0

//replace github.com/metno/frostclient-roadweather => /home/espenm/space/projects/frostclient-roadweather
//replace github.com/metno/frostclient-roadweather => /home/espenm/projects/frostclient-roadweather
//replace github.com/metno/objectstore-stuff => /home/espenm/projects/objectstore-stuff

require (
	github.com/mattn/go-sqlite3 v1.14.17
	github.com/metno/frostclient-roadweather v0.0.9
	github.com/metno/objectstore-stuff v0.0.1
	github.com/myggen/wwwauth v0.0.5
	gocv.io/x/gocv v0.42.0
)

require (
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/klauspost/cpuid/v2 v2.1.0 // indirect
	github.com/minio/md5-simd v1.1.2 // indirect
	github.com/minio/minio-go/v7 v7.0.47 // indirect
	github.com/minio/sha256-simd v1.0.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rs/xid v1.4.0 // indirect
	github.com/sirupsen/logrus v1.9.1 // indirect
	golang.org/x/crypto v0.52.0 // indirect
	golang.org/x/exp v0.0.0-20230905200255-921286631fa9 // indirect
	golang.org/x/net v0.54.0 // indirect
	golang.org/x/sys v0.45.0 // indirect
	golang.org/x/text v0.37.0 // indirect
	gopkg.in/ini.v1 v1.66.6 // indirect
)
