package db

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

// Camera camera struc
type Camera struct {
	ID           int
	ForeignID    string
	Name         string
	Latitude     float64
	Longitude    float64
	Location     string
	RoadNumber   string
	_location    sql.NullString
	Status       string
	Url          string
	Municipality string
}

// Camera camera struc
type RoadLabel struct {
	ID             int
	CameraID       int
	UserName       string
	ImageTimestamp string
	Label          int
	CreatedDate    string
}

// DBFILE Path to sqlite database
var DBFILE string

func SaveOrUpdateTemp(camid int, temp float64, imageTimestamp string) error {
	log.Printf("DB Saving temperature  for camid %d, %f, %s", camid, temp, imageTimestamp)
	db, err := sql.Open("sqlite3", DBFILE)

	if err != nil {
		return err
	}
	defer db.Close()

	// Insert if not exist
	stmt, err := db.Prepare("INSERT OR IGNORE INTO temp_labels (camera_id, degrees_celsius, image_timestamp) VALUES(?,?,?)")

	if err != nil {
		return fmt.Errorf("SaveOrUpdateTemp: %v", err)
	}
	defer stmt.Close()
	res, err := stmt.Exec(camid, temp, imageTimestamp)
	if err != nil {
		return fmt.Errorf("SaveOrUpdateTemp: %v", err)
	}
	id, _ := res.LastInsertId()
	if id != 0 { //Did not exist
		return nil
	} // else existed update:
	stmt, err = db.Prepare("UPDATE  temp_labels SET degrees_celsius=? WHERE camera_id=? AND image_timestamp=?")

	if err != nil {
		return fmt.Errorf("SaveOrUpdateTemp: %v", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(temp, camid, imageTimestamp)

	if err != nil {
		return fmt.Errorf("SaveOrUpdateTemp: %v", err)
	}

	return nil
}

// SaveOrUpdateRoadLabel - Save road-label, update cc label if exists.
func SaveOrUpdateRoadLabel(camid int, label int, imageTimestamp string, userName string) error {
	//log.Printf("DB Saving SaveOrUpdateRoadLabel camid %d, %d, %s", camid, label, imageTimestamp)
	db, err := sql.Open("sqlite3", DBFILE)

	if err != nil {
		return fmt.Errorf(fmt.Sprintf("error SaveOrUpdateRoadLabel with dbfile '%s' : %v", DBFILE, err))
	}
	defer db.Close()
	if label < 0 {
		stmt, err := db.Prepare("DELETE FROM road_labels WHERE camera_id=? AND image_timestamp=?")

		if err != nil {
			return fmt.Errorf("SaveOrUpdateRoadLabel: %v", err)
		}
		defer stmt.Close()
		_, err = stmt.Exec(camid, imageTimestamp)
		if err != nil {
			return fmt.Errorf("SaveOrUpdateRoadLabel: %v", err)
		}
	} else {

		// Insert if not exist
		stmt, err := db.Prepare("INSERT OR IGNORE INTO road_labels (camera_id, label, image_timestamp, username) VALUES(?,?,?,?)")

		if err != nil {
			return fmt.Errorf("SaveOrUpdateRoadLabel: %v", err)
		}
		defer stmt.Close()
		res, err := stmt.Exec(camid, label, imageTimestamp, userName)
		if err != nil {
			return fmt.Errorf("SaveOrUpdateRoadLabel: %v", err)
		}
		id, _ := res.LastInsertId()
		//fmt.Printf("Inserted label with ID: %d\n", id)
		if id != 0 { //Did not exist
			return nil
		} // else existed update:
		stmt, err = db.Prepare("UPDATE  road_labels SET label=? WHERE camera_id=? AND image_timestamp=? AND username=?")

		if err != nil {
			return fmt.Errorf("SaveOrUpdateRoadLabel.UPDATE: %v", err)
		}
		defer stmt.Close()
		_, err = stmt.Exec(label, camid, imageTimestamp, userName)

		if err != nil {
			return fmt.Errorf("SaveOrUpdateRoadLabel.UPDATEXEC: %v", err)
		}
	}
	return nil
}

// GetLabel get label1 value for a cameraID and Datetime
func GetRoadLabel(camid int, imageTimestamp string) (RoadLabel, error) {
	db, err := sql.Open("sqlite3", DBFILE)

	label := RoadLabel{}
	if err != nil {
		label.Label = -1
		return label, err
	}
	defer db.Close()

	//log.Printf("SELECT id, camera_id, label, username, image_timestamp, created_date FROM road_labels WHERE camera_id=%d AND image_timestamp=%s", camid, imageTimestamp)
	row := db.QueryRow("SELECT id, camera_id, label, username, cast(image_timestamp as text), created_date FROM road_labels WHERE camera_id=? AND image_timestamp=?", camid, imageTimestamp)

	switch err := row.Scan(&label.ID, &label.CameraID, &label.Label, &label.UserName, &label.ImageTimestamp, &label.CreatedDate); err {
	case sql.ErrNoRows:
		return label, nil
	case nil:
		//
	default:
		log.Printf("GetRoadLabel %v", err)
	}

	return label, err
}

// GetAllRoadLabels all road labels
func GetAllRoadLabels() ([]RoadLabel, error) {
	labels := []RoadLabel{}

	db, err := sql.Open("sqlite3", DBFILE)

	if err != nil {
		return labels, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, camera_id, label, cast(image_timestamp as varchar) , username FROM road_labels")
	if err != nil {
		return labels, err
	}
	defer rows.Close()

	for rows.Next() {
		label := RoadLabel{}
		err = rows.Scan(&label.ID, &label.CameraID, &label.Label, &label.ImageTimestamp, &label.UserName)
		if err != nil {
			return labels, err
		}
		labels = append(labels, label)
	}
	return labels, err
}

// GetLabelRoadsoulderLabel get roadshoulder value for a cameraID and Datetime
func GetLabelShoulderLabel(camid int, imageTimestamp string, username string) (int, error) {
	db, err := sql.Open("sqlite3", DBFILE)

	if err != nil {
		return -1, err
	}
	defer db.Close()
	cc := -1

	row := db.QueryRow("SELECT label FROM roadshoulder_labels WHERE camera_id=? AND image_timestamp=? AND username=?", camid, imageTimestamp, username)

	switch err := row.Scan(&cc); err {
	case sql.ErrNoRows:
		//
	case nil:
		//
	default:

	}
	return cc, err
}

// GetLabelCountForCamID - Get numberof labels for a camera
func GetRoadLabelCountForCamID(camid int) (int, error) {
	var count int
	db, err := sql.Open("sqlite3", DBFILE)

	if err != nil {
		return 0, err
	}
	defer db.Close()

	rows, err := db.Query("select count(*) as total from (select distinct camera_id, image_timestamp, label from road_labels where label!=9 and camera_id=" + strconv.Itoa(camid) + ")")

	if err != nil {
		return 0, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return 0, err
		}
		break
	}
	return count, nil
}

// GetCam - Get a single Camera
func GetCam(camid int) (Camera, error) {
	camera := Camera{}
	db, err := sql.Open("sqlite3", DBFILE)

	if err != nil {
		return camera, err
	}
	defer db.Close()
	rows, err := db.Query("SELECT id, foreign_id, name, latitude, longitude, road_number, location, status FROM webcams where id=?", camid)

	if err != nil {
		return camera, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&camera.ID, &camera.ForeignID, &camera.Name, &camera.Latitude, &camera.Longitude, &camera.RoadNumber, &camera._location, &camera.Status)
		if err != nil {
			return camera, err
		}
		if camera._location.Valid {
			camera.Location = camera._location.String
		}
		break
	}
	return camera, err
}

// GetCams returns a list of valid webcams ( ie status ="")
func GetCams() ([]Camera, error) {
	cams := []Camera{}

	db, err := sql.Open("sqlite3", DBFILE)

	if err != nil {
		return cams, err
	}
	defer db.Close()
	// 36 is the stupid fisheye cam at Blidnern. Show it and not train it until we decide what to do about it
	rows, err := db.Query("SELECT id, foreign_id, name, latitude, longitude, road_number, location, status FROM webcams where status not like \"%down%\" ")

	if err != nil {
		return cams, err
	}
	defer rows.Close()

	for rows.Next() {
		cam := Camera{}
		err := rows.Scan(&cam.ID, &cam.ForeignID, &cam.Name, &cam.Latitude, &cam.Longitude, &cam.RoadNumber, &cam._location, &cam.Status)
		if cam._location.Valid {
			cam.Location = cam._location.String
		}
		if err != nil {
			fmt.Printf("Database error: %v\n", err)
		}
		cams = append(cams, cam)
	}

	sort.Slice(cams, func(p, q int) bool {
		return cams[p].ID < cams[q].ID
	})
	return cams, err
}

/*
func main() {
	cams, err := GetCams()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", cams)
}
*/
