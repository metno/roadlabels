#!/usr/bin/env python3




import datetime
from datetime import timedelta
import sys
import os
import calendar
import _thread

from threading import Thread
from multiprocessing import Process
import pyproj
import netCDF4
import numpy as np

def date_to_unixtime(date):
    """ Convert YYYYMMDDHH to unixtime

    Arguments:
       date (int): YYYYMMDDHH

    Returns:
       int: unixtime
    """
    year = date / 1000000
    month = date / 10000 % 100
    day = date / 100 % 100
    hour = date % 100
    unix_time = calendar.timegm(datetime.datetime(int(year), int(month), int(day)).timetuple()) + 3600 * int(hour)
    return unix_time


analysis_dir = "/lustre/storeB/immutable/archive/projects/metproduction/yr_short"

def get_forecast_file(time):
   
    date = time / 100
    year = date / 10000
    month = date / 100 % 100
    day = date % 100
    hour = int((time % 100) / 1) * 1

    filedate_str = "%04d%02d%02dT%02dZ" % (year, month, day, hour)
    dir =  "%04d/%02d/%02d" % (year, month, day)
    
    filename = "%s/%s/met_analysis_1_0km_nordic_%s.nc" %(analysis_dir, dir, filedate_str)
    if not os.path.isfile(filename):
        return None
    return filename




def get_values(ds, lats, lons, time, variable):
    """
    lats (list): latitudes
    lons (list) longitudes:

    Returns list of variable ( ir 2 m temp, cloud_area_fractons .. )
    """
    N = len(lats)
    projection = ds.variables["projection_lcc"].proj4
    proj = pyproj.Proj(projection)
    x = ds.variables["x"][:]
    y = ds.variables["y"][:]
    dstime = ds.variables["time"][:]
    unixtime = date_to_unixtime(time)
    t = np.argmin(np.abs(dstime - unixtime))
    #sys.stderr.write("Fetching timestep %d\n" % t)
    values = list()
    for i in range(N):
        X, Y = proj(lons[i], lats[i])
        
        Ix = np.argmin(np.abs(x - X))
        Iy = np.argmin(np.abs(y - Y))
      
        values += [ds.variables[variable][t, Iy, Ix]]
    return values

def getvarval(forecast_file, lat, lon, reftime, variable): 
    ds = netCDF4.Dataset(forecast_file, 'r')
    #print("Using analysis %s" % forecast_file)
    values = get_values(ds, [lat], [lon], reftime, "air_temperature_2m")
    temp = values[0] - 273.15
    
    ds.close()
    return temp


def print_t2m(reftime_long, lat, lon ):
    analysis_file = get_forecast_file(reftime_long)
    if analysis_file is None:
        raise RuntimeError("No suitable analysis file found")

    temp = getvarval(analysis_file, lat, lon, reftime_long, "air_temperature_2m")
    print(temp)

def get_t2m_24hrs(anal_time, lat, lon): 
   
    currdate = datetime.datetime.now().utcnow()
    currdatehour = currdate
    currdate = currdate.replace(hour=0)
    currdate = currdate.replace(minute=0)
    currdate = currdate.replace(second=0)
    currdate = currdate.replace(microsecond=0)


    for _ in range(0,24 ):
        today_str = currdatehour.strftime('%Y%m%d')
        hour_str = currdatehour.strftime('%H')
        reftime = int("%s%s" % (today_str, hour_str)) # sic
        forecast_file = get_forecast_file(reftime)
        
        if forecast_file is None:
            filedate_str = "%04d%02d%02dT%02dZ" % (currdatehour.year, currdatehour.month, currdatehour.day, currdatehour.hour)
            sys.stderr.write("No suitable analysis found for %s \n" % filedate_str)
        else :
            temp = getvarval(forecast_file, lat, lon, anal_time, "air_temperature_2m")
            print("%s lat: %.2f, lon: %.2f, celsius: %.2f" % (forecast_file, lat, lon, temp))

        currdate = currdate - timedelta(hours=1)
        currdatehour = currdatehour - timedelta(hours=1)

if __name__ == '__main__':

    #  Aker sykehus, Oslo, Lat: 59.94, Lon: 10.80
    get_t2m_24hrs(2023042510, 59.94,  10.80)

  


