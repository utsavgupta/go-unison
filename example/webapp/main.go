package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/julienschmidt/httprouter"
	"github.com/utsavgupta/go-unison/example/webapp/ent"
	"github.com/utsavgupta/go-unison/example/webapp/repo"
)

const (
	namespace = "demo"
)

func writeResponse(statusCode int, body interface{}, resp http.ResponseWriter) {
	bodyJSON, err := json.Marshal(body)

	if err != nil {
		strBody, ok := body.(string)

		if ok {
			resp.WriteHeader(statusCode)
			resp.Write([]byte(strBody))
			return
		}

		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte("error"))
		return
	}

	resp.Header().Set("content-type", "application/json")
	resp.WriteHeader(statusCode)
	resp.Write(bodyJSON)
}

func newHandleGetAllArtists(artistRepo repo.ArtistLoader) httprouter.Handle {
	return func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
		artists, err := artistRepo.LoadAllArtists(req.Context())

		if err != nil {
			writeResponse(http.StatusInternalServerError, "error", resp)
			return
		}

		writeResponse(http.StatusOK, ent.Artists{Items: artists}, resp)
	}
}

func newHandleGetAllAlbums(albumRepo repo.AlbumLoader) httprouter.Handle {
	return func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
		albums, err := albumRepo.LoadAllAlbums(req.Context())

		if err != nil {
			writeResponse(http.StatusInternalServerError, "error", resp)
			return
		}

		writeResponse(http.StatusOK, ent.Albums{Items: albums}, resp)
	}
}

func newHandleGetAlbumsForArtist(albumRepo repo.AlbumLoader) httprouter.Handle {
	return func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {

		artistID := params.ByName("artistID")

		if artistID == "" {
			writeResponse(http.StatusBadRequest, "artist id missing", resp)
			return
		}

		albums, err := albumRepo.LoadAllAlbumsByArtist(req.Context(), artistID)

		if err != nil {
			fmt.Println(err)
			writeResponse(http.StatusInternalServerError, "error", resp)
			return
		}

		writeResponse(http.StatusOK, ent.Albums{Items: albums}, resp)
	}
}

func main() {

	dsClient, err := datastore.NewClient(context.Background(), "*detect-project-id*")

	if err != nil {
		panic(err)
	}

	artistRepo := repo.NewArtistLoader(dsClient, namespace)
	albumRepo := repo.NewAlbumLoader(dsClient, namespace)

	r := httprouter.New()

	r.GET("/artists", newHandleGetAllArtists(artistRepo))
	r.GET("/artists/:artistID/albums", newHandleGetAlbumsForArtist(albumRepo))

	r.GET("/albums", newHandleGetAllAlbums(albumRepo))

	http.ListenAndServe(":8080", r)
}
