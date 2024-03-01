package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/metno/roadlabels/pkg/db"
	authandlers "github.com/myggen/wwwauth/pkg/handlers"
	"github.com/myggen/wwwauth/pkg/userrepo"
)

var BuildTime = ""
var Version = ""
var AppRoot = "roadlabels"
var Templates map[string]*template.Template

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
<<<<<<< Updated upstream
		<title>Cams list </title>
=======
		<title>Cams List</title>
>>>>>>> Stashed changes
		
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

func init() {
	fmt.Printf("handlers.init()\n")
}
