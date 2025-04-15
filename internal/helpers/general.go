package helpers

import (
  "strconv"
  "time"
)

func GetTimeStamp()string {
  version := strconv.FormatInt(time.Now().Unix(), 10)
  return version
}
