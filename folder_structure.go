package main

const (
	folderNameAirTemperature     = "air_temperature"
	folderNameCloudType          = "cloud_type"
	folderNameCloudiness         = "cloudiness"
	folderNameDewPoint           = "dew_point"
	folderNameExtremeTemperature = "extreme_temperature"
	folderNameExtremeWind        = "extreme_wind"
	folderNameKl                 = "kl"
	folderNameMorePrecip         = "more_precip"
	folderNamePrecipitation      = "precipitation"
	folderNamePressure           = "pressure"
	folderNameSoilTemperature    = "soil_temperature"
	folderNameSolar              = "solar"
	folderNameSun                = "sun"
	folderNameVisibility         = "visibility"
	folderNameWaterEquiv         = "water_equiv"
	folderNameWeatherPhenomena   = "weather_phenomena"
	folderNameWind               = "wind"
	folderNameWindSynop          = "wind_synop"
	folderNameWindTest           = "wind_test"
)

var folderNamesDaily = []string{
	folderNameKl,
	folderNameMorePrecip,
	folderNameSoilTemperature,
	folderNameSolar,
	folderNameWaterEquiv,
	folderNameWeatherPhenomena,
}
var folderNamesHourly = []string{
	folderNameAirTemperature,
	folderNameCloudType,
	folderNameCloudiness,
	folderNameDewPoint,
	folderNamePrecipitation,
	folderNamePressure,
	folderNameSoilTemperature,
	folderNameSolar,
	folderNameSun,
	folderNameVisibility,
	folderNameWind,
	folderNameWindSynop,
}

var folderNames10Minutes = []string{
	folderNameAirTemperature,
	folderNameExtremeTemperature,
	folderNameExtremeWind,
	folderNamePrecipitation,
	folderNameSolar,
	folderNameWind,
	folderNameWindTest,
}

var folderNames = map[resolution][]string{
	resolutionDaily:     folderNamesDaily,
	resolutionHourly:    folderNamesHourly,
	resolution10Minutes: folderNames10Minutes,
}

var urlsDaily = map[string][]string{
	folderNameKl: {
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/daily/kl/historical/",
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/daily/kl/recent/",
	},
	folderNameMorePrecip: {
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/daily/more_precip/historical/",
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/daily/more_precip/recent/",
	},
	folderNameSoilTemperature: {
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/daily/soil_temperature/historical/",
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/daily/soil_temperature/recent/",
	},
	folderNameSolar: {
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/daily/solar/",
	},
	folderNameWaterEquiv: {
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/daily/water_equiv/historical/",
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/daily/water_equiv/recent/",
	},
	folderNameWeatherPhenomena: {
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/daily/weather_phenomena/historical/",
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/daily/weather_phenomena/recent/",
	},
}

var urlsHourly = map[string][]string{

	folderNameAirTemperature: {
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/hourly/air_temperature/historical/",
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/hourly/air_temperature/recent/",
	},

	folderNameCloudType: {
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/hourly/cloud_type/historical/",
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/hourly/cloud_type/recent/",
	},

	folderNameCloudiness: {
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/hourly/cloudiness/historical/",
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/hourly/cloudiness/recent/",
	},

	folderNameDewPoint: {
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/hourly/dew_point/historical/",
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/hourly/dew_point/recent/",
	},

	folderNamePrecipitation: {
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/hourly/precipitation/historical/",
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/hourly/precipitation/recent/",
	},

	folderNamePressure: {
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/hourly/pressure/historical/",
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/hourly/pressure/recent/",
	},

	folderNameSoilTemperature: {
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/hourly/soil_temperature/historical/",
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/hourly/soil_temperature/recent/",
	},

	folderNameSolar: {
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/hourly/solar/",
	},

	folderNameSun: {
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/hourly/sun/historical/",
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/hourly/sun/recent/",
	},

	folderNameVisibility: {
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/hourly/visibility/historical/",
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/hourly/visibility/recent/",
	},

	folderNameWind: {
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/hourly/wind/historical/",
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/hourly/wind/recent/",
	},

	folderNameWindSynop: {
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/hourly/wind_synop/historical/",
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/hourly/wind_synop/recent/",
	},
}

var urls10Minutes = map[string][]string{
	folderNameAirTemperature: {
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/10_minutes/air_temperature/historical/",
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/10_minutes/air_temperature/recent/",
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/10_minutes/air_temperature/now/",
	},
	folderNameExtremeTemperature: {
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/10_minutes/extreme_temperature/historical/",
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/10_minutes/extreme_temperature/recent/",
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/10_minutes/extreme_temperature/now/",
	},
	folderNameExtremeWind: {
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/10_minutes/extreme_wind/historical/",
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/10_minutes/extreme_wind/recent/",
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/10_minutes/extreme_wind/now/",
	},
	folderNamePrecipitation: {
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/10_minutes/precipitation/historical/",
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/10_minutes/precipitation/recent/",
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/10_minutes/precipitation/now/",
	},
	folderNameSolar: {
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/10_minutes/solar/historical/",
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/10_minutes/solar/recent/",
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/10_minutes/solar/now/",
	},
	folderNameWind: {
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/10_minutes/wind/",
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/10_minutes/wind/recent/",
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/10_minutes/wind/now/",
	},
	folderNameWindTest: {
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/10_minutes/wind_test/historical/",
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/10_minutes/wind_test/recent/",
		"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/10_minutes/wind_test/now/",
	},
}

var urls = map[resolution]map[string][]string{
	resolutionDaily:     urlsDaily,
	resolutionHourly:    urlsHourly,
	resolution10Minutes: urls10Minutes,
}

const colNameStationID = "STATIONS_ID"
const colNameDateTime = "MESS_DATUM"
const colNameEndOfRow = "eor"

var columns = map[resolution]map[string][]string{
	resolutionDaily: {
		folderNameKl:               {colNameStationID, colNameDateTime, "QN_3", "FX", "FM", "QN_4", "RSK", "RSKF", "SDK", "SHK_TAG", "NM", "VPM", "PM", "TMK", "UPM", "TXK", "TNK", "TGK", colNameEndOfRow},
		folderNameMorePrecip:       {colNameStationID, colNameDateTime, "QN_6", "RS", "RSF", "SH_TAG", "NSH_TAG", colNameEndOfRow},
		folderNameSoilTemperature:  {colNameStationID, colNameDateTime, "QN_2", "V_TE002M", "V_TE005M", "V_TE010M", "V_TE020M", "V_TE050M", colNameEndOfRow},
		folderNameSolar:            {colNameStationID, colNameDateTime, "QN_592", "ATMO_STRAHL", "FD_STRAHL", "FG_STRAHL", "SD_STRAHL", colNameEndOfRow},
		folderNameWaterEquiv:       {colNameStationID, colNameDateTime, "QN_6", "ASH_6", "SH_TAG", "WASH_6", "WAAS_6", colNameEndOfRow},
		folderNameWeatherPhenomena: {colNameStationID, colNameDateTime, "QN_4", "NEBEL", "GEWITTER", "STURM_6", "STURM_8", "TAU", "GLATTEIS", "REIF", "GRAUPEL", "HAGEL", colNameEndOfRow},
	},
	resolutionHourly: {
		folderNameAirTemperature:  {colNameStationID, colNameDateTime, "QN_9", "TT_TU", "RF_TU", colNameEndOfRow},
		folderNameCloudType:       {colNameStationID, colNameDateTime, "QN_8", "V_N", "V_N_I", "V_S1_CS", "V_S1_CSA", "V_S1_HHS", "V_S1_NS", "V_S2_CS", "V_S2_CSA", "V_S2_HHS", "V_S2_NS", "V_S3_CS", "V_S3_CSA", "V_S3_HHS", "V_S3_NS", "V_S4_CS", "V_S4_CSA", "V_S4_HHS", "V_S4_NS", colNameEndOfRow},
		folderNameCloudiness:      {colNameStationID, colNameDateTime, "QN_8", "V_N_I", "V_N", colNameEndOfRow},
		folderNameDewPoint:        {colNameStationID, colNameDateTime, "QN_8", "TT", "TD", colNameEndOfRow},
		folderNamePrecipitation:   {colNameStationID, colNameDateTime, "QN_8", "R1", "RS_IND", "WRTR", colNameEndOfRow},
		folderNamePressure:        {colNameStationID, colNameDateTime, "QN_8", "P", "P0", colNameEndOfRow},
		folderNameSoilTemperature: {colNameStationID, colNameDateTime, "QN_2", "V_TE002", "V_TE005", "V_TE010", "V_TE020", "V_TE050", "V_TE100", colNameEndOfRow},
		folderNameSolar:           {colNameStationID, colNameDateTime, "QN_592", "ATMO_LBERG", "FD_LBERG", "FG_LBERG", "SD_LBERG", "ZENIT", "MESS_DATUM_WOZ", colNameEndOfRow},
		folderNameSun:             {colNameStationID, colNameDateTime, "QN_7", "SD_SO", colNameEndOfRow},
		folderNameVisibility:      {colNameStationID, colNameDateTime, "QN_8", "V_VV_I", "V_VV", colNameEndOfRow},
		folderNameWind:            {colNameStationID, colNameDateTime, "QN_3", "F", "D", colNameEndOfRow},
		folderNameWindSynop:       {colNameStationID, colNameDateTime, "QN_8", "FF", "DD", colNameEndOfRow},
	},
	resolution10Minutes: {
		folderNameAirTemperature:     {colNameStationID, colNameDateTime, "QN", "PP_10", "TT_10", "TM5_10", "RF_10", "TD_10", colNameEndOfRow},
		folderNameExtremeTemperature: {colNameStationID, colNameDateTime, "QN", "TX_10", "TX5_10", "TN_10", "TN5_10", colNameEndOfRow},
		folderNameExtremeWind:        {colNameStationID, colNameDateTime, "QN", "FX_10", "FNX_10", "FMX_10", "DX_10", colNameEndOfRow},
		folderNamePrecipitation:      {colNameStationID, colNameDateTime, "QN", "RWS_DAU_10", "RWS_10", "RWS_IND_10", colNameEndOfRow},
		folderNameSolar:              {colNameStationID, colNameDateTime, "QN", "DS_10", "GS_10", "SD_10", "LS_10", colNameEndOfRow},
		folderNameWind:               {colNameStationID, colNameDateTime, "QN", "FF_10", "DD_10", colNameEndOfRow},
		folderNameWindTest:           {colNameStationID, colNameDateTime, "QN", "SLA_10", "SLO_10", "FF_10", "DD_10", colNameEndOfRow},
	},
}
