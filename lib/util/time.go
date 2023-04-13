package util

import "github.com/golang-module/carbon"

func GetShortDateTime() string {
	return carbon.Now().Layout(carbon.ShortDateTimeLayout)
}

func GetDateTime() string {
	return carbon.Now().Layout(carbon.DateTimeLayout)
}

func GetSecondsFromTime(timeStart carbon.Carbon) int {
	return int(carbon.Now().DiffAbsInSeconds(timeStart))
}
