package router

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type recordResponse struct {
	PageSize         int                  `json:"page_size,omitempty"`
	Page             int                  `json:"page,omitempty"`
	PageCount        int                  `json:"page_count,omitempty"`
	TotalPages       int                  `json:"total_pages,omitempty"`
	TotalRecordCount int64                `json:"total_record_count,omitempty"`
	Links            *recordResponseLinks `json:"_links,omitempty"`
	Message          string               `json:"message,omitempty"`
	Error            string               `json:"error,omitempty"`
	Slug             string               `json:"slug,omitempty"`
	Record           interface{}          `json:"record,omitempty"`
	Records          interface{}          `json:"records,omitempty"`
}

// recordResponseLinks represent links that could be returned on a page
type recordResponseLinks struct {
	Self     *link `json:"self,omitempty"`
	First    *link `json:"first,omitempty"`
	Previous *link `json:"previous,omitempty"`
	Next     *link `json:"next,omitempty"`
	Last     *link `json:"last,omitempty"`
}

// link represents an address to a page
type link struct {
	Href string `json:"href,omitempty"`
}

// // notFoundResponse writes a 404 response with the given message
// func notFoundResponse(c *gin.Context, message string) {
// 	c.JSON(http.StatusNotFound, &recordResponse{Message: message})
// }

func badRequestResponse(c *gin.Context, message string, err error) {
	c.JSON(http.StatusBadRequest, &recordResponse{Message: message, Error: err.Error()})
}

func createdResponse(c *gin.Context) {
	uri := uriWithoutQueryParams(c)
	r := &recordResponse{
		Message: "resource created",
		Links: &recordResponseLinks{
			Self: &link{Href: uri},
		},
	}

	c.Header("Location", uri)
	c.JSON(http.StatusCreated, r)
}

func successResponse(c *gin.Context, obj interface{}) {
	uri := uriWithoutQueryParams(c)
	r := &recordResponse{
		Message: "resource created",
		Links: &recordResponseLinks{
			Self: &link{Href: uri},
		},
		Record: obj,
	}

	c.Header("Location", uri)
	c.JSON(http.StatusCreated, r)
}

func deletedResponse(c *gin.Context) {
	c.JSON(http.StatusOK, &recordResponse{Message: "resource deleted"})
}

// func updatedResponse(c *gin.Context, slug string) {
// 	r := &RecordResponse{
// 		Message: "resource updated",
// 		Slug:    slug,
// 		Links: &RecordResponseLinks{
// 			Self: &Link{Href: uriWithoutQueryParams(c)},
// 		},
// 	}

// 	c.JSON(http.StatusOK, r)
// }

func dbErrorResponse(c *gin.Context, err error) {
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, &recordResponse{Message: "resource not found", Error: err.Error()})
	} else {
		c.JSON(http.StatusInternalServerError, &recordResponse{Message: "datastore error", Error: err.Error()})
	}
}

func uriWithoutQueryParams(c *gin.Context) string {
	uri := c.Request.URL
	uri.RawQuery = ""

	return uri.String()
}
