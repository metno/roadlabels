//ALTER TABLE webcams RENAME TO TempWebcams;

CREATE TABLE IF NOT EXISTS webcams (
    id integer primary key autoincrement,
    foreign_id varchar,
    name varchar,
    url varchar,
    latitude real,
    longitude real ,
    location string,
    country_code varchar,
    address varchar,
    municipality varchar,
    county varchar,
    county_no integer,
    road_number varchar,
    status varchar default "",
    external_status varchar default "",
    is_public boolean default 0,
    created_date datetime default (datetime(current_timestamp)),
    provider varchar,
    timezone varchar, 	   	   
    unique(foreign_id, provider)
);


//ALTER TABLE road_labels RENAME TO tmp_road_labels;


CREATE TABLE IF NOT EXISTS road_labels (
    id integer primary key autoincrement,
    camera_id integer,
    label int,
    username varchar ,
    image_timestamp datetime ,           
    created_date datetime default (datetime(current_timestamp)),
    FOREIGN KEY(camera_id) REFERENCES webcams(id)
    UNIQUE (camera_id, image_timestamp) ON CONFLICT REPLACE
);



INSERT INTO road_labels (
    id,
    camera_id,
    label,
    username, 
    image_timestamp,
    created_date
) SELECT 
    id,
    camera_id,
    label,
    username, 
    image_timestamp,
    created_date
FROM 
    tmp_road_labels;


CREATE TABLE IF NOT EXISTS roadshoulder_labels (
    id integer primary key autoincrement,
    camera_id integer,
    label int,
    username varchar ,
    image_timestamp datetime ,
    created_date datetime default (datetime(current_timestamp)),
    FOREIGN KEY(camera_id) REFERENCES webcams(id)
    UNIQUE (camera_id, username, image_timestamp) ON CONFLICT REPLACE
);

