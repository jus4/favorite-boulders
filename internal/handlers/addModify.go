package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"github.com/jus4/favorite-boulders/internal/helpers"
	"github.com/jus4/favorite-boulders/internal/model"
	"github.com/jus4/favorite-boulders/internal/store/queries"
	"github.com/jus4/favorite-boulders/internal/templates/components"
	"github.com/jus4/favorite-boulders/internal/templates/pages"
	"github.com/jus4/favorite-boulders/internal/utils"
	"github.com/labstack/echo/v4"
  "encoding/json"
  "log"
)

func AddModify(c echo.Context) error {
  userId, ok := c.Get("user_email").(string)
  if !ok {
    return helpers.Render(c, http.StatusForbidden, pages.NotAllowed(model.Metadata{Title: "Ei pääsyä"}, nil))
  }
  dbpool := helpers.GetDbPool()

  favourites := model.NewFavouritesList(dbpool)
  data, err := favourites.GetFavouriteList(context.Background(), userId)
  if err != nil {
    fmt.Print(err)
  }
  defer dbpool.Close()

  content := components.FavouritesList(data)

  return helpers.Render(c, http.StatusOK, pages.AddModify(model.Metadata{Title: "Suosikki reitit"}, content))
}

func EditRoute( c echo.Context) error {
  _, ok := c.Get("user_email").(string)
  if !ok {
    return helpers.Render(c, http.StatusBadRequest, pages.FavouriteClimbs(model.Metadata{Title: "Suosikki reitit"}, nil))
  }
  formrouteId := c.FormValue("climb-id")
  routeId, err := strconv.ParseInt(formrouteId, 10, 64)
  if err != nil {
      fmt.Println("Conversion error:", err)
      return helpers.Render(c, http.StatusBadRequest, pages.FavouriteClimbs(model.Metadata{Title: "Suosikki reitit"}, nil))
  }
  dbpool := helpers.GetDbPool()

  route := model.NewEditRouteSectorCrag(dbpool)
  routeData, getRouteDataErr := route.GetRouteById(context.Background(), routeId)
  if getRouteDataErr != nil {
    return helpers.Render(c, http.StatusOK, pages.FavouriteClimbs(model.Metadata{Title: "Suosikki reitit"}, nil))
  }

  return helpers.Render(c, http.StatusOK, components.EditRoute(components.EditRouteProps{
    Name: routeData.Title.String, 
    Grade: routeData.Grade.String,
    RouteType: routeData.RouteType.String,
    ID: strconv.FormatInt(routeData.RouteID, 10),
    Image: fmt.Sprint(routeData.ImageMain),
    Description: routeData.Description.String,
  }))
}

func UpdateRoute( c echo.Context) error {
  _, ok := c.Get("user_email").(string) // is logged in with session
  idStr := c.Param("id") // get the id as a string
  routeName := c.FormValue("name")
  routeType := c.FormValue("route_type")
  grade := c.FormValue("grade")
  description := c.FormValue("description")
  image, formImageError:= c.FormFile("image")
  var imagesJSON []byte

  // Update image
  if formImageError == nil {

    imageOpen, imageOpenErr := image.Open()
    if imageOpenErr != nil {
      return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to open image"})
    }

    defer imageOpen.Close()

    // Get Conversion
    imgConversions, imgConversionErr := utils.GenerateRouteImages(imageOpen)
    if imgConversionErr != nil {
      return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to open image"})
    }

    // Upload to aws
    uploader, err := utils.NewAwsS3Uploader(context.Background())
    mainImageUploaded, uploadError := uploader.UploadTopoImage(context.Background(), routeName + "_" + idStr + "_main.jpg", imgConversions.Main)
    thumbnailImageUploaded, uploadError := uploader.UploadTopoImage(context.Background(), routeName  + "_" + idStr + "_thumbnail.jpg", imgConversions.Thumbnail)
    if uploadError != nil {
      return c.JSON(http.StatusBadRequest, map[string]string{"error": "Upload failed"})
    }

    images := model.RouteImages{
      Main: mainImageUploaded,
      Thumbnail: thumbnailImageUploaded,
    }

    valueJson, err := json.Marshal(images)
    if err != nil {
      return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
    }
    imagesJSON = valueJson

  } else {
    imagesJSON = nil
  }

  id, idError := strconv.ParseInt(idStr, 10, 64)
  if idError != nil || !ok {
    return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
  }

  dbpool := helpers.GetDbPool()
  route := model.NewEditRouteSectorCrag(dbpool)
  params := queries.UpdateRouteByIdParams{
    RouteID: id,
    Title: routeName,
    Grade: grade,
    RouteType: routeType,
    Images: imagesJSON,
    Description: description,
  }

  updateRoute, updateError := route.UpdateRouteById(context.Background(), params)
  if updateError != nil {
    log.Println(updateError)
    return c.JSON(http.StatusBadRequest, map[string]string{"error": "Update error"})
  }
  
  var routeImages model.RouteImages
  json.Unmarshal([]byte(updateRoute.Images), &routeImages)
  return helpers.Render(c, http.StatusOK, components.EditRoute(components.EditRouteProps{
    Name: updateRoute.Title.String, 
    Grade: updateRoute.Grade.String,
    RouteType: updateRoute.RouteType.String,
    ID: strconv.FormatInt(updateRoute.RouteID, 10),
    Image: fmt.Sprint(updateRoute.ImageMain),
    Description: updateRoute.Description.String,
  }))
}
