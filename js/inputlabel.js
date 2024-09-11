
 
function reloadPage(camid){
    // If I had remebered how fucked upp js date / - formating is I had chosen a different format to match .. 
    datetime = document.getElementById('pickadate').value;
    dateq = datetime.replaceAll('-', '');
    dateiso = datetime;
    datetime = datetime.replaceAll('-', '/');
    
    console.log("datetime: "+ datetime);
    url = "/roadlabels/inputlabel?q=roadcams/" + datetime + "/" + camid + "/" + camid + "_" + dateq + "T1200Z" + ".jpg";
    // /roadlabels/inputlabel?q=roadcams/2023/02/04/3/3_20230204T0600Z.jpg
    window.location.href = url;
}

var obstype = "none";
document.addEventListener('keydown', function(event) {
    
    obstype = document.activeElement.id;
    console.log("obsobstype: " + obstype);
    

    var char = event.which || event.keyCode;

    //alert( char + 'was pressed');

    if ( char == 39) { // Right arrow
        document.getElementById("next").click(); // Click on the checkbox
        return;
    }

    if ( char == 37){ // Left arrow
        document.getElementById("prev").click(); // Click on the checkbox
        return;
    }

    
    var map = {}; // Map 0-9 keys 
    for (i=0; i<=9; i++) {
        map[48 + i] = i + 1;
    }
    

    var mapnumlock = {};
    for (i=0; i<=9; i++) {
        map[96 + i] = i + 1;
    }

    var radios = document.forms[0].elements["cc"];
    var radios2 = document.forms[0].elements["obs2"];

    if ( char == 78 ) { // 'n'
        if (obstype != "obs2")	// Road shoulders. Removed for now
            radios[0].click();
        else 
            radios2[0].click();
            
    }

    if (char in map) {
        if (obstype != "obs2")
            radios[map[char]].click();
        else 
            radios2[map[char]].click();
    }

    if (char in mapnumlock) {
        if (obstype != "obs2")
            radios[mapnumlock[char]].click();
        else 
            radios2[mapnumlock[char]].click();
    }
    
});

function getParameterByName(name) {
        var match = RegExp('[?&]' + name + '=([^&]*)').exec(window.location.search);
    return match && decodeURIComponent(match[1].replace(/\+/g, ' '));
}

function radioClicked() {
    var q = getParameterByName("q");
    // Get the non query part of the url
    u = window.location.href.split('?')[0];
    ccVal = document.querySelector('input[name="cc"]:checked').value;
    //obs2Val = document.querySelector('input[name="obs2"]:checked').value;
    //window.location.href = u + "?q=" + q + "&cc=" + ccVal + "&obs2=" + obs2Val;
    window.location.href = u + "?q=" + q + "&cc=" + ccVal + "&obs2=" + "-1";
}

window.onload=function() {
    var radios = document.forms[0].elements["cc"];
    for (var i = 0; i < radios.length; i++) {
        radios[i].onclick=radioClicked;
    }
}
function next(next_image, saveID, saveCC, label2, saveStamp, saveStampJS, temp) {
        
        
        date  =new Date(saveStampJS);
        dt = new Date()
        dt.setHours(dt.getHours());

        if (date >= dt) {
            alert("No pictures from the future yet (" + saveStampJS + ")");
            return;
        }

        u = window.location.href.split('?')[0];
        window.location.href = u + "?q=" + next_image + "&saveID=" + saveID + "&saveCC=" + saveCC + "&saveStamp=" + saveStamp  + "&label2="  + label2 + "&temp=" + temp ;
        
}

function prev(prev_image, saveID, saveCC, label2, saveStamp, temp) {
        u = window.location.href.split('?')[0];
        
        window.location.href = u + "?q=" + prev_image + "&saveID=" + saveID + "&saveCC=" + saveCC + "&saveStamp=" + saveStamp + "&label2="  + label2 + "&temp=" + temp ;
}

