package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/metno/frostclient-roadweather"
	"github.com/metno/objectstore-stuff/objectstore"
	"github.com/metno/roadlabels/pkg/db"
	"github.com/metno/roadlabels/pkg/exttools"
	"github.com/metno/roadlabels/pkg/imgutil"
	authandlers "github.com/myggen/wwwauth/pkg/handlers"
	"github.com/myggen/wwwauth/pkg/userrepo"
	"gocv.io/x/gocv"
)

var BuildTime = ""
var Version = ""
var AppRoot = "roadlabels"
var Templates map[string]*template.Template

var S3endpoint string
var H2s3accessKey string
var H2s3secretKey string
var DbFile string
var FrostObses []frostclient.ObsRoadweather

var UserDBPath string

// Frost data examples
var Class2FrostObses = make(map[string][]frostclient.ObsRoadweather)

type Vars struct {
	BuildTime string
	Version   string
	AppRoot   string
	Templates map[string]*template.Template
}

func NewHandlesHandler(
	BuildTime string,
	Version string,
	AppRoot string,
	Templates map[string]*template.Template) Vars {
	return Vars{BuildTime, Version, AppRoot, Templates}
}

func (t Vars) Print() {

}

type StoreHandler interface {
	ObjectExists(string, string) (bool, error)
	string
}

type S3Handler struct {
	Prefix string // S3 Bucket
}

func (es3 S3Handler) objectExists(name string) (bool, error) {
	s3client, err := objectstore.NewClientWithBucket(S3endpoint, H2s3accessKey, H2s3secretKey, es3.Prefix)
	if err != nil {
		return false, err
	}

	objExists, err := s3client.ObjectExists(name)
	if err != nil {
		return false, fmt.Errorf("s3client.ObjectExists %s: %v", name, err)

	}
	return objExists, nil
}

func (es3 S3Handler) getBytes(name string) ([]byte, error) {
	s3client, err := objectstore.NewClientWithBucket(S3endpoint, H2s3accessKey, H2s3secretKey, es3.Prefix)
	if err != nil {
		return []byte{}, err
	}

	bytebuf, err := s3client.GetS3ObjectBytes(name)
	if err != nil {
		return []byte{}, fmt.Errorf("getBytes(%s): %v", name, err)

	}
	return bytebuf, nil
}

func getCaminfo(camera db.Camera) string {
	if camera.Location != "" {
		return fmt.Sprintf("Id: %d, SVVID: %s, Road: %s, Location: %s, Lat: %.2f, Lon: %.2f\n",
			camera.ID, camera.ForeignID, camera.RoadNumber, camera.Location, camera.Latitude, camera.Longitude)
	}
	return fmt.Sprintf("Id: %d, SVVID: %s, Name: %s, Lat: %.2f, Lon: %2.f\n",
		camera.ID, camera.ForeignID, camera.Name, camera.Latitude, camera.Longitude)
}

func getPathNext(date time.Time, path string) string {

	nextd := date.Add(6 * time.Hour)

	pathNext := strings.Replace(path, date.Format("20060102T1504Z"), nextd.Format("20060102T1504Z"), -1)
	pathNext = strings.Replace(pathNext, date.Format("2006/01/02"), nextd.Format("2006/01/02"), -1)
	return pathNext
}

func getPathPrev(date time.Time, path string) string {
	pathPrev := ""
	prevd := date.Add(-6 * time.Hour)

	pathPrev = strings.Replace(path, date.Format("20060102T1504Z"), prevd.Format("20060102T1504Z"), -1)
	pathPrev = strings.Replace(pathPrev, date.Format("2006/01/02"), prevd.Format("2006/01/02"), -1)

	return pathPrev
}

func DbDownloadHandler(w http.ResponseWriter, r *http.Request) {
	redirectIfNotLoggedin(w, r)
	log.Printf("serving: %s", DbFile)
	fileBytes, err := os.ReadFile(DbFile)
	if err != nil {
		log.Printf("os.ReadFile: %v", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
}

func UserDBDownloadHandler(w http.ResponseWriter, r *http.Request) {
	redirectIfNotLoggedin(w, r)
	log.Printf("serving: %s", UserDBPath)
	fileBytes, err := os.ReadFile(UserDBPath)
	if err != nil {
		log.Printf("os.ReadFile: %v", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
}

func FrontPageHandler(w http.ResponseWriter, r *http.Request, title string) {

	pi := authandlers.NewPageInfo()
	pi.Title = "Road Labels"
	pi.H1 = "Road Condition Annotations"
	pi.Version = Version
	pi.BuildTime = BuildTime

	t, ok := Templates["home.html"]
	if !ok {
		log.Printf("template %s not found", "home.html")
		return
	}

	uu, err := authandlers.IsLoggedIn(r)
	if err != nil {
		log.Printf("authandlers.IsLoggedIn: %v", err)
		pi.Errors = append(pi.Errors, fmt.Sprintf("authandlers.IsLoggedIn: %v", err))
		t.Execute(w, pi)
		return
	}

	if (uu == userrepo.User{}) {
		log.Printf("frontPageHandler: Cookie anyany.xyz_session_token does not exist ")
		http.Redirect(w, r, AppRoot+"/login", http.StatusFound)
		return
	}

	pi.User = uu
	pi.Info = append(pi.Info, "Welcome")

	//var labelcounts map[string]int
	labelcounts := make(map[int]int)
	labels, err := db.GetAllRoadLabels()
	if err != nil {
		log.Fatalf("db.getCams(): %v", err)
	}
	tot := 0
	for l := 0; l < len(labels); l++ {
		labelcounts[labels[l].Label]++
		tot++
	}

	jsarr := "["
	for k := 0; k < 10; k++ {
		jsarr += fmt.Sprintf("%d,", labelcounts[k])
	}
	jsarr = strings.TrimRight(jsarr, ",")
	jsarr += "]"
	s := template.JS(jsarr)
	xleg := fmt.Sprintf("Tot. labels %d", tot)

	pi.Any = append(pi.Any, s)
	pi.Any = append(pi.Any, xleg)
	t.Execute(w, pi)

}

func ShowLabelsHandler(w http.ResponseWriter, r *http.Request, title string) {

	pi := authandlers.NewPageInfo()
	pi.Title = "Road Labels"
	pi.H1 = "Road Condition Annotations"
	pi.Version = Version
	pi.BuildTime = BuildTime

	t, ok := Templates["labeled_images.html"]
	if !ok {
		log.Printf("template %s not found", "labeled_images.html")
		return
	}

	uu, err := authandlers.IsLoggedIn(r)
	if err != nil {
		log.Printf("authandlers.IsLoggedIn: %v", err)
		pi.Errors = append(pi.Errors, fmt.Sprintf("authandlers.IsLoggedIn: %v", err))
		t.Execute(w, pi)
		return
	}

	if (uu == userrepo.User{}) {
		log.Printf("frontPageHandler: Cookie anyany.xyz_session_token does not exist ")
		http.Redirect(w, r, AppRoot+"/login", http.StatusFound)
		return
	}

	labelS := r.URL.Query().Get("label")
	htmlLinks := ""
	for l := 0; l < 10; l++ {
		if labelS == fmt.Sprintf("%d", l) {
			htmlLinks += fmt.Sprintf(`%d&nbsp; `, l)
		} else {
			htmlLinks += fmt.Sprintf(`<a href="/roadlabels/showlabels?label=%d">%d</a>&nbsp; `, l, l)
		}
	}
	if labelS != "all" && labelS != "" {
		htmlLinks += `<a href="/roadlabels/showlabels?label=all">all</a>&nbsp;`
	} else {

		htmlLinks += `all&nbsp;`
	}
	pi.InfoHtml = template.HTML(htmlLinks)
	t.Execute(w, pi)

}

func ImagePagingHandler(w http.ResponseWriter, r *http.Request, title string) {
	db.DBFILE = DbFile

	descShort := map[int]string{
		0: "Dry",
		1: "Patchy water",
		2: "Light water",
		3: "Heavy water",
		4: "Ice (Very difficult)",
		5: "Light slush",
		6: "Heavy slush",
		7: "Light or patchy snow",
		8: "Heavy snow",
		9: "Obstructed.",
	}
	_ = descShort

	type Element struct {
		PathThumb string
		PathBig   string
		Label     int
		Desc      string
	}

	type Data struct {
		Data  []Element
		Total int
	}

	labels, err := db.GetAllRoadLabels()
	if err != nil {
		log.Fatalf("db.getCams(): %v", err)
	}

	labelS := r.URL.Query().Get("label")
	var labels2 []db.RoadLabel
	if labelS != "" && labelS != "all" {
		label, _ := strconv.Atoi(labelS)
		for l := 0; l < len(labels); l++ {
			if labels[l].Label == label {
				labels2 = append(labels2, labels[l])
			}
		}
		labels = labels2
	}

	pageNoS := r.URL.Query().Get("page")
	limitS := r.URL.Query().Get("limit")

	pageNo, _ := strconv.Atoi(pageNoS)
	limit, _ := strconv.Atoi(limitS)

	log.Printf("page no: %s, limitS: %s label: %s from: %d to:%d", pageNoS, limitS, labelS, (pageNo-1)*limit, (pageNo-1)*limit+limit)

	var data = Data{}
	from := (pageNo - 1) * limit
	to := (pageNo-1)*limit + limit
	if to > len(labels) {
		to = len(labels)
	}
	for l := from; l < to; l++ {
		rx := regexp.MustCompile(`(\d{4})(\d{2})(\d{2})T(\d{2})(\d{2})Z`)
		matches := rx.FindStringSubmatch(labels[l].ImageTimestamp)
		if len(matches) != 6 {
			fmt.Printf("Could not parse input query %s %v\n", labels[l].ImageTimestamp, matches)
			return
		}

		year := matches[1]
		month := matches[2]
		day := matches[3]
		pathThumb := fmt.Sprintf("roadcams/%s/%s/%s/%d/thumbs/%d_%s.jpg", year, month, day, labels[l].CameraID, labels[l].CameraID, labels[l].ImageTimestamp)
		pathBig := fmt.Sprintf("roadcams/%s/%s/%s/%d/%d_%s.jpg", year, month, day, labels[l].CameraID, labels[l].CameraID, labels[l].ImageTimestamp)

		elm := Element{}
		elm.PathThumb = pathThumb
		elm.PathBig = pathBig
		elm.Label = labels[l].Label
		elm.Desc = descShort[labels[l].Label]
		data.Data = append(data.Data, elm)

	}
	data.Total = len(labels)

	bytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("ImagePagingHandler  json.Marshal: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(bytes))
}

// TODO: Use templates
func InputLabelHandler(w http.ResponseWriter, r *http.Request, title string) {
	user := redirectIfNotLoggedin(w, r)

	desc := map[int]string{
		0: "Dry (No visible water or patches on the surface. Surrounding environment also looks dryish)",
		1: "Patchy water (Patches of water or not dry looking, with moistness in the surface)",
		2: "Light water (Thinner film of water, light reflection, often wet or snowy surroundings)",
		3: "Heavy water (Clearly very wet)",
		4: "Ice (Very difficult)",
		5: "Light slush (Small amounts of snow on a wet surface)",
		6: "Heavy slush (Large amounts of snow on a wet surface)",
		7: "Light or patchy snow (Snow on surface where road is also visible but not wet)",
		8: "Heavy snow (Snow on surface not visible road)",
		9: "Obstructed. Cannot decide class",
	}

	/*	// Road shoulders
		desc2 := map[int]string{
			0: "No snow",
			1: "Patchy",
			2: "Covered. (More than ~80%)",
		}
	*/

	fmt.Fprintf(w, `<!DOCTYPE html>
<head>
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Label page</title>


<script src="/roadlabels/js/inputlabel.js"></script>
<link rel="stylesheet" href="/roadlabels/css/flatpickr.min.css">
<script src="/roadlabels/js/jquery-3.3.1.min.js"></script>
<script src="/roadlabels/js/flatpickr.js"></script>

	<style>

	section {
	float: left;
	margin: 0 1.5%%;
	width: 63%%;
	}
	aside {
	float: right;
	margin: 0 1.5%%;
	width: 30%%;
	}

	footer {
	clear: both;
	}

	/* Container holding the image and the text */
	.container {
		position: relative;
		text-align: center;
		color: white;
	}

	/* Top left text */
	.top-left {
		position: absolute;
		top: 8px;
		left: 16px;
	}

</style>
</head>
<body>
<div style='float: right;'>Logged in as %s <a href="/%s/logout"> logout </a> </div>
		`, user, AppRoot)

	path := r.URL.Query().Get("q")
	rx := regexp.MustCompile(`(\d+)_(\d{4})(\d{2})(\d{2})T(\d{2})(\d{2})Z`)
	matches := rx.FindStringSubmatch(path)
	if len(matches) != 7 {
		fmt.Printf("Could not parse input query %s %v\n", path, matches)
		return
	}
	camid, _ := strconv.Atoi(matches[1])
	year := matches[2]
	month := matches[3]
	day := matches[4]
	hour := matches[5]
	min := matches[6]

	dateString := fmt.Sprintf("%s-%s-%sT%s:%s:00.000Z", year, month, day, hour, min)
	date, err := time.Parse(time.RFC3339, dateString)
	if err != nil {
		panic(fmt.Errorf("error while parsing date :%v", err))
	}

	// Display caminfo
	camera, err := db.GetCam(camid)
	if err != nil {
		fmt.Fprintf(w, "Could not get camera: %v\n", err)
		return
	}
	camInfo := getCaminfo(camera)
	fmt.Fprintf(w, "%s\n", camInfo)

	fmt.Fprintf(w, "<br/><br/>\n")

	oneDay := time.Hour * 24 * 1
	tomorrow := date.Add(oneDay)
	yesterday := date.Add(-oneDay)
	prevq := fmt.Sprintf("%s%d", yesterday.Format("2006/01/02/"), camid)
	nextq := fmt.Sprintf("%s%d", tomorrow.Format("2006/01/02/"), camid)
	currentq := fmt.Sprintf("%s%d", date.Format("2006/01/02/"), camid)
	todayq := fmt.Sprintf("%s%d", time.Now().UTC().Format("2006/01/02/"), camid)
	fmt.Fprintf(w, `<a href="/roadlabels">Home</a>`)
	fmt.Fprintf(w, `&nbsp; <a href="/roadlabels/thumbs?q=%s">Previous day</a>`, prevq)
	fmt.Fprintf(w, `&nbsp; <a href="/roadlabels/thumbs?q=%s">Current day</a>`, currentq)
	fmt.Fprintf(w, `&nbsp; <a href="/roadlabels/thumbs?q=%s">Today (UTC)</a>`, todayq)
	fmt.Fprintf(w, `&nbsp; <a href="/roadlabels/thumbs?q=%s">Next day</a>`, nextq)
	fmt.Fprintf(w, `&nbsp; Date: <input type="text" value="%s" id="pickadate"  data-utc=true class="pickadate" onchange="reloadPage(%d)" /> <br/><br/>`,
		date.Format("2006-01-02T15:04:05Z"), camid)

	pathNext := getPathNext(date, path)

	pathPrev := getPathPrev(date, path)

	imageTimestamp := date.Format("20060102T1504Z")

	// I am not able to get Javascript to parse the above format :-(
	nextdjs := date.Add(6 * time.Hour)
	dateNextJS := nextdjs.Format("2006-01-02T15:04Z")

	inputLabel := -1
	labeledBy := ""

	if r.URL.Query().Get("cc") != "" {
		inputLabel, _ = strconv.Atoi(r.URL.Query().Get("cc"))
	} else {
		l, err := db.GetRoadLabel(camid, imageTimestamp)
		if err != nil {
			log.Printf(" db.GetRoadLabel: %v", err)
			inputLabel = -1
		}
		if (l == db.RoadLabel{}) {
			inputLabel = -1
		} else {
			inputLabel = l.Label
			labeledBy = "labeled by: " + l.UserName
			fmt.Printf("*** %+v\n", l)
		}
	}

	inputLabel2 := -1

	if r.URL.Query().Get("saveID") != "" && r.URL.Query().Get("saveCC") != "" && r.URL.Query().Get("saveStamp") != "" {
		saveID, _ := strconv.Atoi(r.URL.Query().Get("saveID"))
		saveCC, _ := strconv.Atoi(r.URL.Query().Get("saveCC"))
		saveStamp := r.URL.Query().Get("saveStamp")
		err := db.SaveOrUpdateRoadLabel(saveID, saveCC, saveStamp, user)
		if err != nil {
			log.Printf("db.SaveOrUpdateLabel: %v", err)
		}
	}

	if r.URL.Query().Get("saveID") != "" && r.URL.Query().Get("temp") != "" && r.URL.Query().Get("saveStamp") != "" {
		saveID, _ := strconv.Atoi(r.URL.Query().Get("saveID"))
		temp := r.URL.Query().Get("temp")
		saveStamp := r.URL.Query().Get("saveStamp")
		// Only save temp if image has been annotated
		if r.URL.Query().Get("label2") != "-1" || r.URL.Query().Get("saveCC") != "-1" {
			if t, err := strconv.ParseFloat(temp, 64); err == nil {
				if t < 273.15 {
					err := db.SaveOrUpdateTemp(saveID, t, saveStamp)
					if err != nil {
						log.Printf("db.SaveOrUpdateTemp: %v", err)
					}
				}
			}
		}
	}

	fmt.Fprintf(w, "<section>")

	var bucket = fmt.Sprintf("roadcams-bucket-%s%s", year, month)
	var storeHandler = S3Handler{Prefix: bucket}

	exists, err := storeHandler.objectExists(path)
	if err != nil {
		log.Printf("Could not determine whether %s exists or not: %v", path, err)
		path = "/roadlabels/nodata/no-data.png"
		fmt.Fprintf(w, `<img width=640 height=480 src="%s" border="0"></img>`, path)
	} else if !exists {
		log.Printf("Does not exist: %s", path)
		path = "/roadlabels/nodata/no-data.png"
		fmt.Fprintf(w, `<img width=640 height=480 src="%s" border="0"></img>`, path)

	} else {
		fmt.Fprintf(w, `<a href="/roadlabels/labeledimage?q=%s&cc=%d&obs2=%d"> <img  style="height:50vh"; src=/roadlabels/labeledimage?q=%s&cc=%d&obs2=%d border="0"></img></a>`, path, inputLabel, inputLabel2, path, inputLabel, inputLabel2)
	}

	temperature := -314.15
	var err1 error
	if date.Compare(time.Now().UTC().Add(-1*time.Hour)) < 0 { // Else temp not available yet .
		temperature, err1 = exttools.GetTemp(date, float32(camera.Latitude), float32(camera.Longitude))
	}
	if err1 != nil {
		log.Printf("Error exttools.GetTemp: %v", err)
	}

	fmt.Fprintf(w, `<br/>Termin: %s &nbsp; Temp: %.2f&#8451; &nbsp;%s<br/><br/>`,
		date.Format("2006-01-02T15:04:05Z"), temperature, labeledBy)
	fmt.Fprintf(w, "</section>")

	fmt.Fprintf(w, `
	<script>
		$("#pickadate").flatpickr({enableTime: false,dateFormat: "Y-m-d",   minDate: "2023-02-04", maxDate: "%s"})
	</script>`, time.Now().UTC().Format("2006-01-02T15:04:05Z"))

	fmt.Fprintf(w, `
	`)

	checked := ""
	if inputLabel < 0 {
		checked = "checked"
	}

	fmt.Fprintf(w, `<div style="float:left; clear:both;">`)

	fmt.Fprintf(w, `Road condition label`)
	fmt.Fprintf(w, `<form action=""  method="post" class="ccForm" id="ccForm" name="ccForm">`)

	fmt.Fprintf(w, `<input id="None" type="radio" name="cc" value="-1" %s>None<br>`, checked)

	// Road conditions radios
	checked = ""
	for x := 0; x <= len(desc)-1; x++ {
		if x == inputLabel {
			checked = "checked"
		} else {
			checked = ""
		}
		fmt.Fprintf(w, `<input type="radio" name="cc" id="%d" value="%d" %s> %d - %s<br>`, x, x, checked, x, desc[x])
		fmt.Fprintf(w, "\n")
	}

	fmt.Fprintf(w, "</form>")

	fmt.Fprintf(w, `</div>`)
	fmt.Fprintf(w, "<footer>")
	fmt.Fprintf(w, `<button id="prev" onclick="prev('%s', '%d', %d, %d, '%s', '%f')" >Save and Prev</button>`+"\n", pathPrev, camid, inputLabel, inputLabel2, imageTimestamp, temperature)
	fmt.Fprintf(w, `<button id="next" onclick="next('%s', '%d', %d, %d, '%s', '%s', '%f')" >Save and Next</button>`+"\n", pathNext, camid, inputLabel, inputLabel2, imageTimestamp, dateNextJS, temperature)
	fmt.Fprintf(w, `
		<p style="color:green;"> Tips:<br> Set road condition class from the keyboard by pressing 0-9. <br> Press 'n' for None (Or not click at all). Save using the left/right arrow keys</p>
		<p style="color:green;"> Tips:<br> We start browsing from "yesterday" so move backwards (Use Save and Prev)</p>
	`)
	fmt.Fprintf(w, "</footer>")

	fmt.Fprintf(w, "</body></html>")
}

func labelImage(w http.ResponseWriter, path string, label string, label2 string) {

	rx := regexp.MustCompile(`(\d+)_(\d{4})(\d{2})(\d{2})T(\d{2})(\d{2})Z`)
	matches := rx.FindStringSubmatch(path)
	if len(matches) != 7 {
		fmt.Printf("Could not parse input query %s %v\n", path, matches)
		return
	}

	year := matches[2]
	month := matches[3]

	var bucket = fmt.Sprintf("roadcams-bucket-%s%s", year, month)
	var storeHandler = S3Handler{Prefix: bucket}

	bytebuf, err := storeHandler.getBytes(path)
	if err != nil {
		log.Printf("getBytes %s: %v", path, err)
		return
	}
	img, err := imgutil.NewImageFromBytes(bytebuf)
	if err != nil {
		log.Printf("labelImage: unable to write image: %v\n", err)
		return
	}
	defer img.Close()
	text := ""

	icc, _ := strconv.Atoi(label)
	if icc >= 0 {
		text = fmt.Sprintf("Road: %s", label)
	}
	il2, _ := strconv.Atoi(label2)
	if il2 >= 0 {
		text += fmt.Sprintf(" Shoulder: %s", label2)
	}
	img.PutText(text)

	buf, err := gocv.IMEncode(".jpg", img.Img)
	if err != nil {
		log.Printf("unable to encode matrix: %v", err)
		panic(err)
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buf.GetBytes())))
	if _, err := w.Write(buf.GetBytes()); err != nil {
		log.Printf("unable to write image: %v\n", err)
	}

}

func LabelImageHandler(w http.ResponseWriter, r *http.Request, title string) {

	path := r.URL.Query().Get("q")
	cc := r.URL.Query().Get("cc")
	obs2 := r.URL.Query().Get("obs2")
	labelImage(w, path, cc, obs2)
}

func LabelThumbHandler(w http.ResponseWriter, r *http.Request, title string) {

	path := r.URL.Query().Get("q")
	cc := r.URL.Query().Get("cc")
	obs2 := r.URL.Query().Get("obs2")
	labelImage(w, path, cc, obs2)
}

// TODO: Use templates
func ThumbsPageHandler(w http.ResponseWriter, r *http.Request, title string) {
	log.Printf("thumbsPageHandler")
	username := redirectIfNotLoggedin(w, r)

	q := r.URL.Query().Get("q")
	rx := regexp.MustCompile(`(\d{4})/(\d{2})/(\d{2})/(\d+)`)
	matches := rx.FindStringSubmatch(q)
	if len(matches) != 5 {
		fmt.Printf("Could not parse input query %s %v\n", q, matches)
		return
	}
	year := matches[1]
	month := matches[2]
	day := matches[3]
	camid, err := strconv.Atoi(matches[4])
	if err != nil {
		panic("Invalid camid")
	}

	fmt.Printf("thumbsPageHandler start\n")
	fmt.Fprintf(w, `<!DOCTYPE html>
		<head>
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<link rel="stylesheet" href="/roadlabels/css/flatpickr.min.css">
		<script src="/roadlabels/js/jquery-3.3.1.min.js"></script>
	 	<script src="/roadlabels/js/flatpickr.js"></script>
 
	 <script>
 
	 function reloadAllcams(){
		 // If I had remebered how fucked upp js date / - formating is I had chosen a different format to match .. 
		 datetime = document.getElementById('pickadate').value;
		 datetime = datetime.replaceAll('-', '/');
		
		 
		 window.location.href = "/roadlabels/thumbs?q=" + datetime + "/" + %d + "&sortBy=Name";
	 }
	 </script>
		</head>
	
	<html><body>`, camid)

	fmt.Fprintf(w, `<div style='float: right;'>Logged in as %s <a href="/%s/logout"> logout </a> </div>`, username, AppRoot)

	dateString := fmt.Sprintf("%s-%s-%sT00:00:00.000Z", year, month, day)
	currentDay, err := time.Parse(time.RFC3339, dateString)
	if err != nil {
		fmt.Fprintf(w, "error while parsing date: %v\n", err)
		return
	}

	camera, err := db.GetCam(camid)
	if err != nil {
		fmt.Fprintf(w, "Could not get camera: %v\n", err)
		return
	}
	camInfo := getCaminfo(camera)
	fmt.Fprintf(w, "%s\n", camInfo)

	fmt.Fprintf(w, `Date <input type="text" value="%s" id="pickadate"  data-utc=true class="pickadate" onchange="reloadAllcams()" />`, currentDay.Format("2006-01-02")+" 12:00")
	todayStr := time.Now().UTC().Add(-24 * time.Hour).Format(("2006-01-02"))
	fmt.Fprintf(w, `
	<script>
		$("#pickadate").flatpickr({enableTime: false,dateFormat: "Y-m-d",   minDate: "2023-02-04", maxDate: "%s"})
	</script>`, todayStr)

	oneDay := time.Hour * 24 * 1
	nextDay := currentDay.Add(oneDay)
	prevDay := currentDay.Add(-oneDay)
	prevq := fmt.Sprintf("%s%d", prevDay.Format("2006/01/02/"), camid)
	nextq := fmt.Sprintf("%s%d", nextDay.Format("2006/01/02/"), camid)
	todayq := fmt.Sprintf("%s%d", time.Now().UTC().Format("2006/01/02/"), camid)

	fmt.Fprintf(w, "<br/><br/><a href=%q>Home</a>", "/roadlabels")
	fmt.Fprintf(w, "&nbsp; <a href=\"/roadlabels/thumbs?q=%s\">Previous day</a>\n", prevq)
	fmt.Fprintf(w, "&nbsp; <a href=\"/roadlabels/thumbs?q=%s\">Today</a>\n", todayq)
	fmt.Fprintf(w, "&nbsp; <a href=\"/roadlabels/thumbs?q=%s\">Next day</a><br/><br/>\n", nextq)

	fmt.Fprintf(w, "<table cellspacing=\"0\" cellpadding=\"2\" width=\"500\">")
	fmt.Fprintf(w, "<tr>")
	count := 0

	dateCurrString := fmt.Sprintf("%s-%s-%sT00:00:00.000Z", year, month, day)
	dayCurr, err := time.Parse(time.RFC3339, dateCurrString)
	if err != nil {
		panic(fmt.Errorf("error while parsing date :%v", err))
	}
	fname := fmt.Sprintf("%d_%04d%02d%02dT%02d%02dZ.jpg", camid, dayCurr.Year(),
		dayCurr.Month(), dayCurr.Day(), dayCurr.Hour(), dayCurr.Minute())
	fmt.Printf("FNAME: %s\n", fname)

	var haveData bool
	var bucket = fmt.Sprintf("roadcams-bucket-%s%s", year, month)
	var storeHandler = S3Handler{Prefix: bucket}
	for i := 0; i < 4; i++ {

		imageTimestamp := fmt.Sprintf("%0.4d%0.2d%0.2dT%0.2d%0.2dZ", dayCurr.Year(), dayCurr.Month(),
			dayCurr.Day(), dayCurr.Hour(), dayCurr.Minute())
		pathOrig := fmt.Sprintf("roadcams/%0.4d/%0.2d/%0.2d/%d/%d_%s.jpg", dayCurr.Year(), dayCurr.Month(),
			dayCurr.Day(), camid, camid, imageTimestamp)
		pathThumb := fmt.Sprintf("roadcams/%0.4d/%0.2d/%0.2d/%d/thumbs/%d_%s.jpg", dayCurr.Year(), dayCurr.Month(),
			dayCurr.Day(), camid, camid, imageTimestamp)

		haveData = true

		l, err := db.GetRoadLabel(camid, imageTimestamp)

		cc := -1
		if err != nil {
			log.Printf("thumbsPageHandler db.GetRoadLabel %d: %v", camid, err)
		}
		if (l != db.RoadLabel{}) {
			cc = l.Label
		}

		termin := fmt.Sprintf("%0.2d%0.2dZ", dayCurr.Hour(), dayCurr.Minute())

		bytebuf, err := storeHandler.getBytes(pathThumb)
		if err != nil {
			log.Printf("thumbsPageHandler.getBytess %s: %v", pathThumb, err)
			fmt.Fprintf(w, `<td valign="middle" ><a title="%s" href="/roadlabels/inputlabel?q=%s&cc=%d">`, imageTimestamp, pathOrig, cc)
			fmt.Fprintf(w, `<img src="/roadlabels/nodata/no-data.png"  style="height:25vh"; border="0" /> <br>%s`, termin)
			fmt.Fprintf(w, "</a></td>\n")

		} else {
			base64Encoding := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(bytebuf)

			fmt.Fprintf(w, `<td valign="middle" ><a title="%s" href="/roadlabels/inputlabel?q=%s&cc=%d">`, imageTimestamp, pathOrig, cc)
			fmt.Fprintf(w, `<img src="%s"  style="height:25vh"; border="0" /> <br>%s`, base64Encoding, termin)
			fmt.Fprintf(w, "</a></td>\n")
		}
		count++

		if count%2 == 0 {
			fmt.Fprintf(w, "</tr><tr>")
		}
		dayCurr = dayCurr.Add(6 * time.Hour)
	}
	fmt.Fprintf(w, "</tr>")
	fmt.Fprintf(w, "</table>")
	if !haveData {
		fmt.Fprintf(w, "No data available for date")
	}
	fmt.Fprintf(w, "</body></html>")
	fmt.Printf("thumbsPageHandler finished\n")
}

func AllCamsHandler(w http.ResponseWriter, r *http.Request, title string) {

	user := redirectIfNotLoggedin(w, r)

	//  2023-02-04: Firs reliable regular downloads
	now := time.Now().UTC().Add(-24 * time.Hour)
	timestamp := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, time.UTC)

	dateStr := r.URL.Query().Get("date")
	fmt.Println("date =>", dateStr)
	if dateStr != "" {
		date, err := time.Parse("20060102", dateStr)
		if err != nil {
			date, err = time.Parse("20060102T1504Z", dateStr)
			if err != nil {
				fmt.Fprintf(w, "error: %v", err)
				return
			}
		}
		timestamp = date

	}

	sortBy := r.URL.Query().Get("sortBy")
	if sortBy == "" {
		sortBy = "Name"
	}

	db.DBFILE = DbFile
	cams, err := db.GetCams()
	if err != nil {
		log.Fatalf("db.getCams(): %v", err)
	}

	sort.Slice(cams, func(p, q int) bool {
		if sortBy == "Name" {
			return cams[p].Name < cams[q].Name
		}
		if sortBy == "ForeignID" {
			return cams[p].ForeignID < cams[q].ForeignID
		}
		if sortBy == "ID" {
			return cams[p].ID < cams[q].ID
		}

		return cams[p].Name < cams[q].Name
	})

	datePrev := timestamp.Add(-24 * time.Hour)
	dateNext := timestamp.Add(24 * time.Hour)

	todayStr := time.Now().UTC().Add(-24 * time.Hour).Format(("2006-01-02"))
	fmt.Fprintf(w, `<!DOCTYPE html><html>
	<head>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>Roadcams Imagelist</title>

    <link rel="stylesheet" href="/roadlabels/css/flatpickr.min.css">
   	<script src="/roadlabels/js/jquery-3.3.1.min.js"></script>
	<script src="/roadlabels/js/flatpickr.js"></script>

	<style>
	img {
		max-width: 33%%;
	  }
	</style>

	<script>

	function reloadAllcams(){
		// If I had remebered how fucked upp js date / - formating is I had chosen a different format to match .. 
		
		datetime = document.getElementById('pickadate').value
		datetime += "T1200";
		console.log("Date picked: " + datetime)
		datetime = datetime.replaceAll('-', '');
		datetime = datetime.replaceAll(':', '');
		datetime = datetime.replaceAll(' ', 'T');
		datetime += "Z";
		
		window.location.href = "/roadlabels/allcams?date=" + datetime + "&sortBy=Name";
	}

	</script>

	</head> `)

	count := 0
	fmt.Fprintf(w, `<body>
	<div style='float: right;'>Logged in as %s <a href="/%s/logout"> logout </a> </div>
	<a href="/roadlabels">Home</a> Camcount: %d, Date: <input type="text" value="%s" id="pickadate"  data-utc=true class="pickadate" onchange="reloadAllcams()" />UTC +0`, user, AppRoot, len(cams), timestamp.Format("2006-01-02 15:04"))

	fmt.Fprintf(w, `<br/><a href="/roadlabels/allcams?date=%s&sortBy=%s"><< Prev Date </a> | <a href="/roadlabels/allcams?date=%s&sortBy=%s">Next Date >> </a>    
	
	
	<script>
	$("#pickadate").flatpickr({enableTime: true,dateFormat: "Y-m-d",   minDate: "2023-02-04", maxDate: "%s"});
	</script>

		<br/>`,
		datePrev.Format("20060102T1504Z"), sortBy, dateNext.Format("20060102T1504Z"), sortBy, todayStr)

	fmt.Fprintf(w, "<table><tr>")

	//now := time.Now().UTC()
	//now = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, time.UTC)
	//now = now.Add(time.Hour * -1)

	var bucket = fmt.Sprintf("roadcams-bucket-%.02d%.02d", timestamp.Year(), timestamp.Month())

	destdirThumb := fmt.Sprintf("%s/%s", "roadcams", timestamp.Format("2006/01/02"))
	s3client, err := objectstore.NewClientWithBucket(S3endpoint, H2s3accessKey, H2s3secretKey, bucket)
	if err != nil {
		log.Fatalln(err)
	}
	for c := 0; c < len(cams); c++ {
		imageTimestamp := fmt.Sprintf("%0.4d%0.2d%0.2dT%0.2d%0.2dZ", timestamp.Year(), timestamp.Month(),
			timestamp.Day(), timestamp.Hour(), timestamp.Minute())
		pathOrig := fmt.Sprintf("%s/%0.4d/%0.2d/%0.2d/%d/%d_%s.jpg", "roadcams", timestamp.Year(), timestamp.Month(),
			timestamp.Day(), cams[c].ID, cams[c].ID, imageTimestamp)

		destpathThumb := fmt.Sprintf("%s/%d/thumbs/%d_%s.jpg", destdirThumb, cams[c].ID, cams[c].ID, timestamp.Format("20060102T1504Z"))
		bytebuf, err := s3client.GetS3ObjectBytes(destpathThumb)
		if err != nil {
			log.Printf("s3client.GetS3ObjectBytes(%s/%s): %v", s3client.Bucket, destpathThumb, err)

		}

		base64Encoding := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(bytebuf)
		fmt.Fprintf(w, `<td><a href="/roadlabels/inputlabel?q=%s"" ><img src=%s /><br/>%s, </br>Road:, %s SVV ID: %s</a> </td>`+"\n",
			pathOrig, base64Encoding, cams[c].Name, cams[c].RoadNumber, cams[c].ForeignID)
		count++

		if count%3 == 0 {
			fmt.Fprintf(w, "</tr><tr>")
		}

	}
	fmt.Fprintf(w, "</tr></table>")
	fmt.Fprintf(w, "</body></html>")
	fmt.Printf("Count: %d\n", count)
}

func FrostObsHandler(w http.ResponseWriter, r *http.Request, title string) {

	pi := authandlers.NewPageInfo()
	pi.Title = "Road Labels"

	pi.Version = Version
	pi.BuildTime = BuildTime
	pi.User.UserName = "noauth"

	t, ok := Templates["example_images_based_on_frost.html"]
	if !ok {
		log.Printf("template %s not found", "example_images_based_on_frost.html")
		return
	}

	uu, err := authandlers.IsLoggedIn(r)
	if err != nil {
		log.Printf("authandlers.IsLoggedIn: %v", err)
		pi.Errors = append(pi.Errors, fmt.Sprintf("authandlers.IsLoggedIn: %v", err))
		t.Execute(w, pi)
		return
	}

	if (uu == userrepo.User{}) {
		log.Printf("frontPageHandler: Cookie anyany.xyz_session_token does not exist ")
		http.Redirect(w, r, AppRoot+"/login", http.StatusFound)
		return
	}

	pi.User.UserName = uu.UserName

	class := r.URL.Query().Get("class")
	pi.H1 = class + " in frost"
	log.Printf("Classe req: %s", class)
	FrostObses = Class2FrostObses[class]
	fmt.Printf("LEN frostObses: %d", len(FrostObses))

	for key, element := range Class2FrostObses {
		fmt.Println("Key:", key, "=>", "len Element:", len(element))
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(FrostObses), func(i, j int) { FrostObses[i], FrostObses[j] = FrostObses[j], FrostObses[i] })
	t.Execute(w, pi)

}

func IceImagePagingHandler(w http.ResponseWriter, r *http.Request, title string) {
	db.DBFILE = DbFile

	type Element struct {
		PathThumb string
		PathBig   string
		Value     float32
		Desc      string
		Label     int
	}

	type Data struct {
		Data  []Element
		Total int
	}

	pageNoS := r.URL.Query().Get("page")
	limitS := r.URL.Query().Get("limit")

	pageNo, _ := strconv.Atoi(pageNoS)
	limit, _ := strconv.Atoi(limitS)

	labels, err := db.GetAllRoadLabels()
	if err != nil {
		log.Fatalf("db.getCams(): %v", err)
	}

	lmap := make(map[int]db.RoadLabel)
	for lc := 0; lc < len(labels); lc++ {
		lmap[labels[lc].CameraID] = labels[lc]
	}

	//log.Printf("page no: %s, limitS: %s label: %s from: %d to:%d", pageNoS, limitS, labelS, (pageNo-1)*limit, (pageNo-1)*limit+limit)

	var data = Data{}
	from := (pageNo - 1) * limit
	to := (pageNo-1)*limit + limit
	if to > len(FrostObses) {
		to = len(FrostObses)
	}
	for l := from; l < to; l++ {
		iob := FrostObses[l]
		refTime := iob.RefTime
		imageTimestamp := refTime.Format("20060102T1504Z")
		label := -1
		val, ok := lmap[iob.CamID]
		if ok { // Have labels on cam
			if imageTimestamp == val.ImageTimestamp {
				label = val.Label
			}
		}

		pathThumb := fmt.Sprintf("roadcams/%0.4d/%0.2d/%0.2d/%d/thumbs/%d_%s.jpg", refTime.Year(), refTime.Month(), refTime.Day(), iob.CamID, iob.CamID, imageTimestamp)
		pathBig := fmt.Sprintf("roadcams/%0.4d/%0.2d/%0.2d/%d/%d_%s.jpg", refTime.Year(), refTime.Month(), refTime.Day(), iob.CamID, iob.CamID, imageTimestamp)

		elm := Element{}
		elm.PathThumb = pathThumb
		elm.PathBig = pathBig
		elm.Desc = "road_ice_thickness"
		elm.Label = label

		//elm.Value = iob.Value

		var bucket = fmt.Sprintf("roadcams-bucket-%.02d%.02d", refTime.Year(), refTime.Month())
		var storeHandler = S3Handler{Prefix: bucket}
		exists, err := storeHandler.objectExists(pathBig)
		if err != nil {
			log.Printf("Could not determine whether %s exists or not: %v", pathBig, err)
			continue
		} else if !exists {
			log.Printf("Does not exist: %s", pathBig)
			continue

		}
		data.Data = append(data.Data, elm)
	}
	data.Total = len(FrostObses)

	bytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("IceImagePagingHandler json.Marshal: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(bytes))
}

// TODO: Use templates
func CamlistHandler(w http.ResponseWriter, r *http.Request, title string) {
	user := redirectIfNotLoggedin(w, r)

	cams, err := getCamsInfo()
	if err != nil {
		log.Printf("getCamsInfo() : %v", err)
	}

	now := time.Now().UTC().Add(-24 * time.Hour)
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	//startDate := time.Date(2023, 2, 4, 0, 0, 0, 0, time.UTC)

	fmt.Fprintf(w, `<!DOCTYPE html>
	<html>
	<head>
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<title>Cams List</title>		
	</head>
	<body>
		<div style='float: right;'>Logged in as %s <a href="/%s/logout"> logout </a> </div>
		<a href="/roadlabels">Home</a> 
		<table border="1">
		<tr><td style="width:6em">Camid</td> <td>SVV Id</td> <td>Municipality</td> <td>Road No</td> <td>Name</td> <td>Lat.</td> <td>Lon.</td> <td># of lables</td> <td>Camera Status</td></tr>
	`, user, AppRoot)
	fmt.Fprintf(w, "\n")

	for i := 0; i < len(cams); i++ {

		q := startDate.Format("2006/01/02/") + strconv.Itoa(cams[i].cam.ID)
		fmt.Fprintf(w, "<tr>")

		fmt.Fprintf(w, `<td>Id: %d.</td>
						<td>%s</td>
						<td><a href="/roadlabels/thumbs?q=%s">%s</a></td>
						<td>%s</td>
						<td>%s</td>`, cams[i].cam.ID, cams[i].cam.ForeignID, q, cams[i].cam.Municipality, cams[i].cam.RoadNumber, cams[i].cam.Name)
		fmt.Fprintf(w, `<td>%.02f</td>
						<td>%.02f</td>
						<td align="center">%d</td>
						<td>%s</td>`+"\n",
			cams[i].cam.Latitude, cams[i].cam.Longitude, cams[i].labelCount, cams[i].cam.Status)
		fmt.Fprintf(w, "</tr>")

	}

	fmt.Fprintf(w, `
		</table>
	</body>
	</html>
	`)
}

func redirectIfNotLoggedin(w http.ResponseWriter, r *http.Request) string {
	pi := authandlers.NewPageInfo()
	pi.Title = "Road Labels"
	pi.H1 = "Road Labels"
	pi.Version = Version
	pi.BuildTime = BuildTime

	t, ok := Templates["home.html"]
	if !ok {
		log.Printf("template %s not found", "response.html")
		return ""
	}

	uu, err := authandlers.IsLoggedIn(r)
	if err != nil {
		log.Printf("authandlers.IsLoggedIn: %v", err)
		pi.Errors = append(pi.Errors, fmt.Sprintf("authandlers.IsLoggedIn: %v", err))
		t.Execute(w, pi)
		return ""
	}

	if (uu == userrepo.User{}) {
		log.Printf("HomeHandler: Cookie anyany.xyz_session_token does not exist ")
		http.Redirect(w, r, AppRoot+"/login", http.StatusFound)
		return ""
	}
	return uu.UserName
}

type camInfo struct {
	cam        db.Camera
	labelCount int
}

func getCamsInfo() ([]camInfo, error) {
	cams, err := db.GetCams()
	var camsinfo []camInfo
	if err != nil {
		return camsinfo, err
	}

	for i := 0; i < len(cams); i++ {
		labelcnt, err := db.GetRoadLabelCountForCamID(cams[i].ID)
		if err != nil {
			return []camInfo{}, err
		}
		arr := strings.Split(cams[i].Location, ", ")

		if len(arr) == 2 {
			// massage for display
			cams[i].Name = arr[0]
			cams[i].Municipality = arr[1]
		}
		camInfo := camInfo{cams[i], labelcnt}
		camsinfo = append(camsinfo, camInfo)
	}
	sort.Slice(camsinfo, func(i, j int) bool {
		return camsinfo[i].cam.Municipality < camsinfo[j].cam.Municipality
	})

	return camsinfo, nil
}
