package comic

import (
	"http/web"
	"myproject/dapi/api/auth/session"
	"myproject/dapi/o/org/chapter"
	"myproject/dapi/o/org/comic"
	"net/http"
	"strconv"
)

type ComicServer struct {
	web.JsonServer
	*http.ServeMux
}

type Result struct {
	ID              string          `bson:"_id,omitempty" json:"id,omitempty"`      //
	MTime           int64           `bson:"mtime,omitempty" json:"mtime,omitempty"` //Modify Time
	DTime           int64           `bson:"dtime" json:"dtime"`                     //Delete Time
	Name            *string         `bson:"name,omitempty" json:"name"`
	OrtherName      *string         `bson:"orther_name,omitempty" json:"orther_name"`
	Thumbnail       *string         `bson:"thumbnail,omitempty" json:"thumbnail"`
	Summary         *string         `bson:"summary,omitempty" json:"summary"`
	Categories      *[]string       `bson:"categories,omitempty" json:"categories"`
	Author          *string         `bson:"author,omitempty" json:"author"`
	Source          *string         `bson:"source,omitempty" json:"source"`
	Alias           *string         `bson:"alias,omitempty" json:"alias"`
	MetaKeyword     *string         `bson:"metakeyword,omitempty" json:"metakeyword"`
	MetaDescription *string         `bson:"metadescription,omitempty" json:"metadescription"`
	Link            *string         `bson:"link,omitempty" json:"link"`
	State           *bool           `bson:"state,omitempty" json:"state"`
	Vote            *int64          `bson:"vote,omitempty" json:"vote"`
	Viewer          *int64          `bson:"viewer,omitempty" json:"viewer"`
	Sharer          *int64          `bson:"sharer,omitempty" json:"sharer"`
	Liker           *int64          `bson:"liker,omitempty" json:"liker"`
	Hot             *bool           `bson:"hot,omitempty" json:"hot"` // trend, most viewer
	New             *bool           `bson:"new,omitempty" json:"new"` // new comic
	Top             *bool           `bson:"top,omitempty" json:"top"` // must view
	Approve         *bool           `bson:"approve,omitempty" json:"approve"`
	Active          *bool           `bson:"active,omitempty" json:"active"`
	CTime           int64           `bson:"ctime,omitempty" json:"ctime"` // Create Time
	Chapters        chapter.Chapter `bson:"chapter" json:"chapter"`
}

func NewComicServer() *ComicServer {
	var s = &ComicServer{
		ServeMux: http.NewServeMux(),
	}

	s.HandleFunc("/create", s.HandleCreate)
	s.HandleFunc("/get", s.HandleGetByID)
	s.HandleFunc("/update", s.HandleUpdateByID)
	s.HandleFunc("/mark_delete", s.HandleMarkDelete)
	s.HandleFunc("/get_all", s.mustGetComicAll)
	//s.HandleFunc("/get_by_category", s.mustGetComicAllByCategory)
	return s
}

func (s *ComicServer) Session(w http.ResponseWriter, r *http.Request) *session.Session {

	ses, err := session.FromContext(r.Context())

	if err != nil {
		s.SendError(w, err)
	}

	return ses
}

func StrToInt(s string) int {
	i, _ := strconv.ParseInt(s, 10, 64)
	return int(i)
}

func (s *ComicServer) mustGetComicAll(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	pageSize := StrToInt(r.URL.Query().Get("page_size"))
	pageNumber := StrToInt(r.URL.Query().Get("page_number"))

	var res = []comic.Comic{}

	count, err := comic.GetAll(pageSize, pageNumber, sortBy, sortOrder, &res)

	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		s.SendDataSuccess(w, map[string]interface{}{
			"comics": res,
			"count":  count,
		})
	}
}
func (s *ComicServer) mustGetComicAllByCategory(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	pageSize := StrToInt(r.URL.Query().Get("page_size"))
	pageNumber := StrToInt(r.URL.Query().Get("page_number"))

	cate := r.URL.Query().Get("cate")

	var result []Result
	var arrayComic []comic.Comic
	var count = 0
	var err error
	var cates []string

	cates = append(cates, cate)
	arrayComic, err = comic.GetByCategory(pageSize, pageNumber, sortBy, sortOrder, cates)

	if err != nil {
		s.ErrorMessage(w, err.Error())
	}

	var c = chapter.Chapter{}
	var re = Result{}

	for _, i := range arrayComic {
		err = chapter.GetLatestChapterByComicID(i.ID, &c)
		if err == nil {
			re = Result{
				ID:              i.ID,
				MTime:           i.MTime,
				DTime:           i.DTime,
				Name:            i.Name,
				OrtherName:      i.OrtherName,
				Thumbnail:       i.Thumbnail,
				Summary:         i.Summary,
				Categories:      i.Categories,
				Author:          i.Author,
				Source:          i.Source,
				Alias:           i.Alias,
				MetaKeyword:     i.MetaKeyword,
				MetaDescription: i.MetaDescription,
				Link:            i.Link,
				State:           i.State,
				Vote:            i.Vote,
				Viewer:          i.Viewer,
				Sharer:          i.Sharer,
				Liker:           i.Liker,
				Hot:             i.Hot,
				New:             i.New,
				Top:             i.Top,
				Approve:         i.Approve,
				Active:          i.Active,
				CTime:           i.CTime,
				Chapters:        c,
			}
		} else {
			re = Result{
				ID:              i.ID,
				MTime:           i.MTime,
				DTime:           i.DTime,
				Name:            i.Name,
				OrtherName:      i.OrtherName,
				Thumbnail:       i.Thumbnail,
				Summary:         i.Summary,
				Categories:      i.Categories,
				Author:          i.Author,
				Source:          i.Source,
				Alias:           i.Alias,
				MetaKeyword:     i.MetaKeyword,
				MetaDescription: i.MetaDescription,
				Link:            i.Link,
				State:           i.State,
				Vote:            i.Vote,
				Viewer:          i.Viewer,
				Sharer:          i.Sharer,
				Liker:           i.Liker,
				Hot:             i.Hot,
				New:             i.New,
				Top:             i.Top,
				Approve:         i.Approve,
				Active:          i.Active,
				CTime:           i.CTime,
			}
		}

		result = append(result, re)
		count++
	}

	s.SendDataSuccess(w, map[string]interface{}{
		"comics": result,
		"count":  count,
	})
}

func (s *ComicServer) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var v = &comic.Comic{}
	s.MustDecodeBody(r, v)
	err := v.Create()
	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		s.SendDataSuccess(w, v)
	}
}

func (s *ComicServer) mustGetComic(r *http.Request) (*comic.Comic, error) {
	var id = r.URL.Query().Get("id")
	var v, err = comic.GetByID(id)
	if err != nil {
		return v, err
	} else {
		return v, nil
	}
}

func (s *ComicServer) HandleUpdateByID(w http.ResponseWriter, r *http.Request) {
	var newComic = &comic.Comic{}
	s.MustDecodeBody(r, newComic)
	v, err := s.mustGetComic(r)
	if err != nil {
		s.ErrorMessage(w, "comic_not_found")
		return
	}
	err = v.Update(newComic)
	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		result, err := comic.GetByID(v.ID)
		if err != nil {
			s.ErrorMessage(w, "comic_not_found")
			return
		}
		s.SendDataSuccess(w, result)
	}
}

func (s *ComicServer) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	v, err := s.mustGetComic(r)
	if err != nil {
		s.ErrorMessage(w, "comic_not_found")
		return
	}
	s.SendDataSuccess(w, v)
}

func (s *ComicServer) HandleMarkDelete(w http.ResponseWriter, r *http.Request) {
	v, err := s.mustGetComic(r)
	if err != nil {
		s.ErrorMessage(w, "comic_not_found")
		return
	}
	err = comic.MarkDelete(v.ID)
	if err != nil {
		s.ErrorMessage(w, err.Error())
	} else {
		s.Success(w)
	}
}
